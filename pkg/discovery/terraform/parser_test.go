package terraform

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/julianstephens/go-utils/logger"
)

func TestGetTFDirPaths(t *testing.T) {
	// Create temporary test directory structure
	tmpDir := t.TempDir()

	// Create directories with .tf files
	dir1 := filepath.Join(tmpDir, "module1")
	dir2 := filepath.Join(tmpDir, "module2")
	dir3 := filepath.Join(tmpDir, "nested", "module3")
	dirEmpty := filepath.Join(tmpDir, "empty")

	for _, dir := range []string{dir1, dir2, dir3, dirEmpty} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}

	// Create .tf files
	testFiles := []string{
		filepath.Join(dir1, "main.tf"),
		filepath.Join(dir1, "variables.tf"),
		filepath.Join(dir2, "resources.tf"),
		filepath.Join(dir3, "data.tf"),
	}

	for _, f := range testFiles {
		if err := os.WriteFile(f, []byte("# test"), 0644); err != nil {
			t.Fatalf("failed to write file %s: %v", f, err)
		}
	}

	// Also create a non-.tf file to ensure it's ignored
	nonTFFile := filepath.Join(dir1, "README.md")
	if err := os.WriteFile(nonTFFile, []byte("# README"), 0644); err != nil {
		t.Fatalf("failed to write non-tf file: %v", err)
	}

	parser := NewHCLParser(logger.NewNoop())

	dirs, err := parser.getTFDirPaths(tmpDir)
	if err != nil {
		t.Fatalf("getTFDirPaths() error = %v", err)
	}

	// Should find 3 directories (dir1, dir2, dir3), not dirEmpty
	if len(dirs) != 3 {
		t.Errorf("getTFDirPaths() found %d dirs, want 3", len(dirs))
	}

	// Verify each directory is in the list
	expectedDirs := []string{dir1, dir2, dir3}
	for _, expected := range expectedDirs {
		found := false
		for _, actual := range dirs {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("getTFDirPaths() missing expected dir %s", expected)
		}
	}

	// Verify empty dir is NOT in the list
	for _, dir := range dirs {
		if dir == dirEmpty {
			t.Errorf("getTFDirPaths() should not include empty dir %s", dirEmpty)
		}
	}
}

func TestGetTFDirPaths_Empty(t *testing.T) {
	tmpDir := t.TempDir()

	parser := NewHCLParser(logger.NewNoop())

	dirs, err := parser.getTFDirPaths(tmpDir)
	if err != nil {
		t.Fatalf("getTFDirPaths() error = %v", err)
	}

	if len(dirs) != 0 {
		t.Errorf("getTFDirPaths() on empty dir = %d dirs, want 0", len(dirs))
	}
}

func TestGetTFDirPaths_NoDuplicates(t *testing.T) {
	tmpDir := t.TempDir()
	dir1 := filepath.Join(tmpDir, "module")

	if err := os.MkdirAll(dir1, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	// Create multiple .tf files in the same directory
	files := []string{"main.tf", "variables.tf", "outputs.tf"}
	for _, f := range files {
		path := filepath.Join(dir1, f)
		if err := os.WriteFile(path, []byte("# test"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	parser := NewHCLParser(logger.NewNoop())

	dirs, err := parser.getTFDirPaths(tmpDir)
	if err != nil {
		t.Fatalf("getTFDirPaths() error = %v", err)
	}

	// Should find only 1 unique directory despite multiple .tf files
	if len(dirs) != 1 {
		t.Errorf("getTFDirPaths() found %d dirs, want 1", len(dirs))
	}
}

func TestGetTFFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	tfFiles := []string{"main.tf", "variables.tf", "outputs.tf"}
	nonTFFiles := []string{"README.md", "config.json", ".gitignore"}

	for _, f := range tfFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("# test"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	for _, f := range nonTFFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	// Create a subdirectory with a .tf file (should not be included)
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "nested.tf"), []byte("# test"), 0644); err != nil {
		t.Fatalf("failed to write nested file: %v", err)
	}

	parser := NewHCLParser(logger.NewNoop())

	files, err := parser.getTFFiles(tmpDir)
	if err != nil {
		t.Fatalf("getTFFiles() error = %v", err)
	}

	// Should find exactly 3 .tf files
	if len(files) != 3 {
		t.Errorf("getTFFiles() found %d files, want 3", len(files))
	}

	// Verify all returned files are .tf files
	for _, f := range files {
		if !strings.HasSuffix(f, ".tf") {
			t.Errorf("getTFFiles() returned non-.tf file: %s", f)
		}
	}

	// Verify all expected .tf files are present
	for _, expected := range tfFiles {
		expectedPath := filepath.Join(tmpDir, expected)
		found := false
		for _, actual := range files {
			if actual == expectedPath {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("getTFFiles() missing expected file %s", expected)
		}
	}
}

func TestGetTFFiles_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	parser := NewHCLParser(logger.NewNoop())

	files, err := parser.getTFFiles(tmpDir)
	if err != nil {
		t.Fatalf("getTFFiles() error = %v", err)
	}

	if len(files) != 0 {
		t.Errorf("getTFFiles() on empty dir = %d files, want 0", len(files))
	}
}

func TestGetTFFiles_NonexistentDir(t *testing.T) {
	parser := NewHCLParser(logger.NewNoop())

	_, err := parser.getTFFiles("/nonexistent/directory")
	if err == nil {
		t.Error("getTFFiles() on nonexistent dir should return error")
	}
}

func TestTransformResource(t *testing.T) {
	tfconfigRes := &tfconfig.Resource{
		Type: "aws_lambda_function",
		Name: "example",
		Pos: tfconfig.SourcePos{
			Filename: "main.tf",
			Line:     10,
		},
	}

	result := transformResource(tfconfigRes)

	// Check ID format
	expectedID := "aws_lambda_function.example"
	if result.ID != expectedID {
		t.Errorf("transformResource() ID = %s, want %s", result.ID, expectedID)
	}

	// Check Type
	if result.Type != "aws_lambda_function" {
		t.Errorf("transformResource() Type = %s, want %s", result.Type, "aws_lambda_function")
	}

	// Check Attributes initialized
	if result.Attributes == nil {
		t.Error("transformResource() Attributes is nil, want initialized map")
	}

	// Check Blocks initialized
	if result.Blocks == nil {
		t.Error("transformResource() Blocks is nil, want initialized map")
	}

	// Check Location
	if result.Location.File != "main.tf" {
		t.Errorf("transformResource() Location.File = %s, want main.tf", result.Location.File)
	}
	if result.Location.Line != 10 {
		t.Errorf("transformResource() Location.Line = %d, want 10", result.Location.Line)
	}
}

func TestParse_NoTFFiles(t *testing.T) {
	tmpDir := t.TempDir()

	parser := NewHCLParser(logger.NewNoop())

	_, err := parser.Parse(context.Background(), tmpDir)
	if err == nil {
		t.Error("Parse() on dir without .tf files should return error")
	}

	// Should be ErrNoHCLFiles
	if parserErr, ok := err.(*TerraformParserError); ok {
		if parserErr.Err != ErrNoHCLFiles {
			t.Errorf("Parse() error type = %v, want ErrNoHCLFiles", parserErr.Err)
		}
	} else {
		t.Errorf("Parse() error should be *TerraformParserError, got %T", err)
	}
}

func TestParse_WithValidTerraform(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple valid Terraform file
	tfContent := `
resource "aws_s3_bucket" "example" {
  bucket = "my-test-bucket"
  
  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  owners = ["099720109477"]
}

locals {
  bucket_prefix = "test-prefix"
  environment   = "development"
}
`

	tfFile := filepath.Join(tmpDir, "main.tf")
	if err := os.WriteFile(tfFile, []byte(tfContent), 0644); err != nil {
		t.Fatalf("failed to write test .tf file: %v", err)
	}

	parser := NewHCLParser(logger.NewNoop())

	config, err := parser.Parse(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check that we got a config
	if config == nil {
		t.Fatal("Parse() returned nil config")
	}

	// Should have 1 resource
	if len(config.Resources) != 1 {
		t.Errorf("Parse() found %d resources, want 1", len(config.Resources))
	} else {
		// Check resource details
		res := config.Resources[0]
		if res.Type != "aws_s3_bucket" {
			t.Errorf("resource type = %s, want aws_s3_bucket", res.Type)
		}
		if res.ID != "aws_s3_bucket.example" {
			t.Errorf("resource ID = %s, want aws_s3_bucket.example", res.ID)
		}
		// Check bucket attribute
		if bucket, ok := res.Attributes["bucket"]; !ok {
			t.Error("resource missing 'bucket' attribute")
		} else if bucket != "my-test-bucket" {
			t.Errorf("bucket attribute = %v, want my-test-bucket", bucket)
		}
	}

	// Should have 1 data source
	if len(config.DataSources) != 1 {
		t.Errorf("Parse() found %d data sources, want 1", len(config.DataSources))
	} else {
		// Check data source details
		ds := config.DataSources[0]
		if ds.Type != "aws_ami" {
			t.Errorf("data source type = %s, want aws_ami", ds.Type)
		}
		if ds.ID != "data.aws_ami.ubuntu" {
			t.Errorf("data source ID = %s, want data.aws_ami.ubuntu", ds.ID)
		}
	}

	// Should have 2 locals
	if len(config.Locals) != 2 {
		t.Errorf("Parse() found %d locals, want 2", len(config.Locals))
	} else {
		// Check specific locals
		if val, ok := config.Locals["bucket_prefix"]; !ok {
			t.Error("missing local 'bucket_prefix'")
		} else if val != "test-prefix" {
			t.Errorf("local bucket_prefix = %v, want test-prefix", val)
		}

		if val, ok := config.Locals["environment"]; !ok {
			t.Error("missing local 'environment'")
		} else if val != "development" {
			t.Errorf("local environment = %v, want development", val)
		}
	}

	// Should have 1 module
	if len(config.Modules) != 1 {
		t.Errorf("Parse() found %d modules, want 1", len(config.Modules))
	}
}

func TestParse_MultipleModules(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two module directories
	module1 := filepath.Join(tmpDir, "module1")
	module2 := filepath.Join(tmpDir, "module2")

	for _, dir := range []string{module1, module2} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
	}

	// Module 1 with one resource
	tf1 := `
resource "aws_s3_bucket" "bucket1" {
  bucket = "bucket-one"
}
`
	if err := os.WriteFile(filepath.Join(module1, "main.tf"), []byte(tf1), 0644); err != nil {
		t.Fatalf("failed to write module1 tf: %v", err)
	}

	// Module 2 with one resource
	tf2 := `
resource "aws_s3_bucket" "bucket2" {
  bucket = "bucket-two"
}
`
	if err := os.WriteFile(filepath.Join(module2, "main.tf"), []byte(tf2), 0644); err != nil {
		t.Fatalf("failed to write module2 tf: %v", err)
	}

	parser := NewHCLParser(logger.NewNoop())

	config, err := parser.Parse(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should have 2 resources (accumulated from both modules)
	if len(config.Resources) != 2 {
		t.Errorf("Parse() found %d resources, want 2", len(config.Resources))
	}

	// Should have 2 modules
	if len(config.Modules) != 2 {
		t.Errorf("Parse() found %d modules, want 2", len(config.Modules))
	}

	// Verify each module has its own resource
	for _, module := range config.Modules {
		if len(module.Resources) != 1 {
			t.Errorf("module %s has %d resources, want 1", module.Name, len(module.Resources))
		}
	}
}

func TestParse_WithNestedBlocks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Terraform file with nested blocks
	tfContent := `
resource "aws_lambda_function" "example" {
  function_name = "test-function"
  runtime       = "python3.9"

  environment {
    variables = {
      KEY1 = "value1"
      KEY2 = "value2"
    }
  }

  logging_config {
    log_level = "INFO"
  }
}
`

	tfFile := filepath.Join(tmpDir, "main.tf")
	if err := os.WriteFile(tfFile, []byte(tfContent), 0644); err != nil {
		t.Fatalf("failed to write test .tf file: %v", err)
	}

	parser := NewHCLParser(logger.NewNoop())

	config, err := parser.Parse(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(config.Resources) != 1 {
		t.Fatalf("Parse() found %d resources, want 1", len(config.Resources))
	}

	res := config.Resources[0]

	// Check attributes
	if res.Attributes["function_name"] != "test-function" {
		t.Errorf("function_name = %v, want test-function", res.Attributes["function_name"])
	}

	// Check blocks
	if res.Blocks == nil {
		t.Fatal("resource Blocks is nil")
	}

	// Should have environment block
	if envBlocks, ok := res.Blocks["environment"]; !ok {
		t.Error("missing environment block")
	} else if len(envBlocks) != 1 {
		t.Errorf("found %d environment blocks, want 1", len(envBlocks))
	}

	// Should have logging_config block
	if logBlocks, ok := res.Blocks["logging_config"]; !ok {
		t.Error("missing logging_config block")
	} else if len(logBlocks) != 1 {
		t.Errorf("found %d logging_config blocks, want 1", len(logBlocks))
	}
}

func TestParse_WithVariableReferences(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Terraform file with variable references
	tfContent := `
variable "bucket_name" {
  type = string
}

resource "aws_s3_bucket" "example" {
  bucket = var.bucket_name
  
  tags = {
    Environment = var.environment
  }
}

locals {
  full_name = "${var.prefix}-${var.bucket_name}"
}
`

	tfFile := filepath.Join(tmpDir, "main.tf")
	if err := os.WriteFile(tfFile, []byte(tfContent), 0644); err != nil {
		t.Fatalf("failed to write test .tf file: %v", err)
	}

	parser := NewHCLParser(logger.NewNoop())

	config, err := parser.Parse(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(config.Resources) != 1 {
		t.Fatalf("Parse() found %d resources, want 1", len(config.Resources))
	}

	res := config.Resources[0]

	// Check that variable reference is preserved as symbolic value
	if bucket, ok := res.Attributes["bucket"]; !ok {
		t.Error("missing bucket attribute")
	} else {
		bucketStr, ok := bucket.(string)
		if !ok {
			t.Errorf("bucket should be string, got %T", bucket)
		} else if !strings.Contains(bucketStr, "var.bucket_name") {
			t.Errorf("bucket = %s, should contain var.bucket_name", bucketStr)
		}
	}

	// Check locals with template
	if len(config.Locals) != 1 {
		t.Errorf("found %d locals, want 1", len(config.Locals))
	} else {
		if fullName, ok := config.Locals["full_name"]; !ok {
			t.Error("missing local 'full_name'")
		} else {
			// Should be a template string with variable references
			fullNameStr, ok := fullName.(string)
			if !ok {
				t.Errorf("full_name should be string, got %T", fullName)
			} else if !strings.Contains(fullNameStr, "var.prefix") || !strings.Contains(fullNameStr, "var.bucket_name") {
				t.Errorf("full_name = %s, should contain variable references", fullNameStr)
			}
		}
	}
}
