package terraform

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/julianstephens/go-utils/generic"
	"github.com/julianstephens/go-utils/logger"

	"github.com/julianstephens/surfacemap/pkg/model"
)

type TerraformParser interface {
	// Parse loads and parses all .tf files in a directory
	Parse(ctx context.Context, rootPath string) (*TerraformConfig, error)

	// LoadVariables loads variable values from tfvars files
	LoadVariables(varFiles ...string) error

	// ResolveModule inlines a local module into the main config
	ResolveModule(modulePath string) (*ModuleConfig, error)
}

type TerraformConfig struct {
	Resources   []model.Resource
	DataSources []model.Resource
	Modules     []Module
	Variables   map[string]any
	Locals      map[string]any
}

type HCLParser struct {
	parser *hclparse.Parser
	logger *logger.Logger
}

func NewHCLParser(logger *logger.Logger) *HCLParser {
	return &HCLParser{
		parser: hclparse.NewParser(),
		logger: logger,
	}
}

// Parse loads and parses all .tf files in the provided root directory (recursively), extracting resources, data sources, variables, and locals into a TerraformConfig struct.
// It uses the hclparse package to parse HCL files and tfconfig to load module metadata. The function also builds an index of blocks by their file and line number to correlate with resources and data sources defined in the module.
// If any errors occur during file parsing or expression decoding, it returns a TerraformParserError with details about the failure.
func (p *HCLParser) Parse(ctx context.Context, rootPath string) (*TerraformConfig, error) {
	startTime := time.Now()
	tfDirs, err := p.getTFDirPaths(rootPath)
	if err != nil {
		return nil, &TerraformParserError{
			Err:     ErrDirWalk,
			Cause:   err,
			RootDir: rootPath,
		}
	}
	p.logger.Infof("found %d directories containing .tf files", len(tfDirs))

	if len(tfDirs) == 0 {
		return nil, &TerraformParserError{
			Err:     ErrNoHCLFiles,
			RootDir: rootPath,
		}
	}

	accumulatedConfig := &TerraformConfig{
		Resources:   []model.Resource{},
		DataSources: []model.Resource{},
		Modules:     []Module{},
		Variables:   make(map[string]any),
		Locals:      make(map[string]any),
	}

	for _, tfDir := range tfDirs {
		moduleConfig, err := p.parseAST(tfDir)
		if err != nil {
			return nil, &TerraformParserError{
				Err:       ErrParse,
				Cause:     err,
				RootDir:   rootPath,
				ModuleDir: tfDir,
			}
		}

		// Accumulate resources from this module
		accumulatedConfig.Resources = append(accumulatedConfig.Resources, moduleConfig.Resources...)
		accumulatedConfig.DataSources = append(accumulatedConfig.DataSources, moduleConfig.DataSources...)

		// Merge locals (later modules override earlier ones with same name)
		for k, v := range moduleConfig.Locals {
			accumulatedConfig.Locals[k] = v
		}

		// Track as a module
		accumulatedConfig.Modules = append(accumulatedConfig.Modules, Module{
			Name:      filepath.Base(tfDir),
			Source:    tfDir,
			Resources: moduleConfig.Resources,
			Variables: moduleConfig.Variables,
			Locals:    moduleConfig.Locals,
		})
	}

	p.logger.Infof("parsed %d resources, %d data sources, %d locals across %d modules in %s",
		len(accumulatedConfig.Resources),
		len(accumulatedConfig.DataSources),
		len(accumulatedConfig.Locals),
		len(accumulatedConfig.Modules),
		time.Since(startTime))
	return accumulatedConfig, nil
}

func (p *HCLParser) LoadVariables(varFiles ...string) error {
	return nil
}

func (p *HCLParser) ResolveModule(modulePath string) (*ModuleConfig, error) {
	return nil, nil
}

// getTFDirPaths returns a list of unique directories containing .tf files found in the provided root directory (recursive)
func (p *HCLParser) getTFDirPaths(rootPath string) (tfDirs []string, err error) {
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err2 error) error {
		if err2 != nil {
			return err2
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(d.Name(), ".tf") && !generic.Contains(tfDirs, filepath.Dir(path)) {
			tfDirs = append(tfDirs, filepath.Dir(path))
		}

		return nil
	})
	return
}

func (p *HCLParser) getTFFiles(dir string) (files []string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tf") {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	return
}

// parseAST loads and parses all .tf files in the specified module directory, returning a TerraformConfig with resources, data sources, variables, and locals defined in that module.
// It uses the hclparse package to parse HCL files and tfconfig to load module metadata. The function also builds an index of blocks by their file and line number to correlate with resources and data sources defined in the module.
// If any errors occur during file parsing or expression decoding, it returns a TerraformParserError with details about the failure.
func (p *HCLParser) parseAST(moduleDir string) (conf *TerraformConfig, err error) {
	conf = &TerraformConfig{
		Resources:   []model.Resource{},
		DataSources: []model.Resource{},
		Variables:   make(map[string]any),
		Locals:      make(map[string]any),
	}

	module, diags := tfconfig.LoadModule(moduleDir)
	if diags.HasErrors() {
		err = errors.New(diags.Error())
		return
	}

	tfFiles, err := p.getTFFiles(moduleDir)
	if err != nil {
		return
	}

	parsedFiles := make(map[string]*hcl.File)
	for _, filePath := range tfFiles {
		hclFile, diags := p.parser.ParseHCLFile(filePath)
		if diags.HasErrors() {
			p.logger.Warnf("parse error in %s: %s", filePath, diags.Error())
			continue
		}
		parsedFiles[filePath] = hclFile
	}

	// key: "file:line"
	blockIndex := make(map[string]*hclsyntax.Block)
	for filePath, hclFile := range parsedFiles {
		body, ok := hclFile.Body.(*hclsyntax.Body)
		if !ok {
			p.logger.Warnf("unexpected body type in %s: %T", filePath, hclFile.Body)
			continue
		}
		for _, block := range body.Blocks {
			if block.Type == "locals" {
				for name, attr := range block.Body.Attributes {
					exprVal, err := DecodeExpr(attr.Expr)
					if err != nil {
						p.logger.Warnf("failed to decode local %s in %s: %v", name, filePath, err)
						continue
					}
					conf.Locals[name] = exprVal.ToAny()
				}
			}
			key := fmt.Sprintf("%s:%d", filePath, block.DefRange().Start.Line)
			blockIndex[key] = block
		}
	}

	p.logger.Infof("parsed %d files, extracted %d locals, indexed %d blocks",
		len(parsedFiles), len(conf.Locals), len(blockIndex))

	for _, res := range module.ManagedResources {
		r := transformResource(res)
		key := fmt.Sprintf("%s:%d", res.Pos.Filename, res.Pos.Line)
		block := blockIndex[key]
		if block == nil {
			p.logger.Warnf("no block found for resource %s at %s", r.ID, key)
			continue
		}

		attrs, blocks, err := DecodeBlock(block.Body)
		if err != nil {
			p.logger.Errorf("failed to decode block for resource %s: %v", r.ID, err)
			continue
		}

		for name, exprVal := range attrs {
			r.Attributes[name] = exprVal.ToAny()
		}

		r.Blocks = blocks
		p.logger.Infof("extracted %d attributes and %d block types for %s",
			len(r.Attributes), len(r.Blocks), r.ID)

		conf.Resources = append(conf.Resources, r)
	}

	for _, dataSource := range module.DataResources {
		p.logger.Infof("parsing data source: %s.%s", dataSource.Type, dataSource.Name)
		r := model.Resource{
			ID:         fmt.Sprintf("data.%s.%s", dataSource.Type, dataSource.Name),
			Type:       dataSource.Type,
			Attributes: make(map[string]any),
			Blocks:     make(map[string][]map[string]any),
			Location: model.FileLocation{
				File: dataSource.Pos.Filename,
				Line: dataSource.Pos.Line,
			},
		}

		blockKey := fmt.Sprintf("%s:%d", dataSource.Pos.Filename, dataSource.Pos.Line)
		block := blockIndex[blockKey]
		if block != nil {
			attrs, blocks, err := DecodeBlock(block.Body)
			if err == nil {
				for name, exprVal := range attrs {
					r.Attributes[name] = exprVal.ToAny()
				}
				r.Blocks = blocks
			}
		}
		conf.DataSources = append(conf.DataSources, r)
	}

	return
}

func transformResource(res *tfconfig.Resource) model.Resource {
	return model.Resource{
		ID:         fmt.Sprintf("%s.%s", res.Type, res.Name),
		Type:       res.Type,
		Attributes: make(map[string]any),
		Blocks:     make(map[string][]map[string]any),
		Location: model.FileLocation{
			File: res.Pos.Filename,
			Line: res.Pos.Line,
		},
	}
}
