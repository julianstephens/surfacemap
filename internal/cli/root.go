package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/julianstephens/go-utils/cliutil"
	"github.com/julianstephens/go-utils/logger"

	pkgerrors "github.com/julianstephens/surfacemap/pkg/errors"
)

// ErrNotImplemented is a reference to the package-level not implemented error
// Deprecated: Use pkgerrors.ErrNotImplemented directly
var ErrNotImplemented = pkgerrors.ErrNotImplemented

type Globals struct {
	VarFile  string           `help:"Path to terraform.tfvars file" type:"existingfile"`
	Root     string           `help:"Root directory of the project to analyze" type:"dir"`
	Verbose  bool             `short:"v" help:"Enable verbose output"`
	VVerbose bool             `short:"vv" help:"Enable very verbose output (debug)"`
	Quiet    bool             `short:"q" help:"Suppress all output except for errors"`
	Output   string           `short:"o" help:"Output format: text, json, markdown (default: text)" enum:"text,json,markdown" default:"text"`
	NoColor  bool             `help:"Disable colored output"`
	Config   string           `help:"Path to configuration file (default: surfacemap.yaml)" type:"existingfile" default:"surfacemap.yaml"`
	Version  VersionFlag      `short:"V" help:"Print version information and exit"`
	Logger   *logger.Logger   `kong:"-"`
	Console  *cliutil.Console `kong:"-"`
}

type CLI struct {
	Globals
	Analyze AnalyzeCmd `cmd:"" help:"Analyze the attack surface of a project"`
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	cliutil.PrintColored(fmt.Sprintf("surfacemap v%s", vars["version"]), cliutil.ColorCyan)
	app.Exit(0)
	return nil
}
