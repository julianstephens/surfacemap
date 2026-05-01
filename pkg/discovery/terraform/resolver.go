package terraform

import (
	"path/filepath"
	"strings"

	"github.com/julianstephens/go-utils/logger"

	pkgerrors "github.com/julianstephens/surfacemap/pkg/errors"
)

// ModuleResolver resolves module sources to filesystem paths
type ModuleResolver interface {
	// Resolve converts a module source to an absolute path
	// Parameters:
	//	fromPath: The path of the module that contains the source reference
	//	source: The module source string (e.g., "../.." or "registry/...")
	// Returns:
	//	string: The resolved absolute path to the module, or empty if remote/unsupported
	//	error: Any errors encountered during resolution
	Resolve(fromPath, source string) (string, error)

	// IsLocal determines if a source is a local path
	IsLocal(source string) bool

	// IsRemote determines if a source is a remote reference
	IsRemote(source string) bool
}

// PathResolver implements ModuleResolver
type PathResolver struct {
	logger *logger.Logger
}

// NewPathResolver creates a new module resolver
func NewPathResolver(logger *logger.Logger) *PathResolver {
	return &PathResolver{
		logger: logger.WithField("component", "path_resolver"),
	}
}

// Resolve converts a module source to an absolute path
func (r *PathResolver) Resolve(fromPath, source string) (string, error) {
	if source == "" {
		return "", nil
	}

	r.logger.Debugf("resolving source '%s' from '%s'", source, fromPath)

	if r.IsLocal(source) {
		r.logger.Debugf("source '%s' is identified as local", source)
		return r.NormalizePath(filepath.Join(fromPath, source))
	}

	r.logger.Warnf("remote module sources are not supported in this version: '%s'", source)
	return "", nil
}

// IsLocal determines if a source is a local filesystem path
func (r *PathResolver) IsLocal(source string) bool {
	if (strings.HasPrefix(source, ".") || strings.HasPrefix(source, "/")) && (!strings.HasPrefix(source, "terraform-") && !strings.Contains(source, "://")) {
		return true
	}
	return false
}

// IsRemote determines if a source is a remote reference
func (r *PathResolver) IsRemote(source string) bool {
	registryFormat := strings.Count(source, "/") == 2 && !strings.Contains(source, "://") && !strings.HasPrefix(source, ".") && !strings.HasPrefix(source, "/")
	if strings.Contains(source, "://") || registryFormat {
		return true
	}
	return false
}

// NormalizePath normalizes a path for comparison
func (r *PathResolver) NormalizePath(path string) (string, error) {
	r.logger.Debugf("normalizing path: %s", path)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", pkgerrors.NewFileError(
			pkgerrors.ErrFileRead,
			"resolve absolute path",
			path,
			err,
		)
	}
	cleanPath := filepath.Clean(absPath)
	if cleanPath != absPath {
		r.logger.Debugf("cleaned path: %s -> %s", absPath, cleanPath)
	}
	return cleanPath, nil
}
