package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/julianstephens/go-utils/cliutil"

	"github.com/julianstephens/surfacemap/internal"
	"github.com/julianstephens/surfacemap/internal/cli"
)

type Globals struct {
	VarFile string      `help:"Path to terraform.tfvars file" type:"existingfile"`
	Root    string      `help:"Root directory of the project to analyze" type:"dir"`
	Verbose bool        `short:"v" help:"Enable verbose output"`
	Quiet   bool        `short:"q" help:"Suppress all output except for errors"`
	Output  string      `short:"o" help:"Output format: text, json, markdown (default: text)" enum:"text,json,markdown" default:"text"`
	NoColor bool        `help:"Disable colored output"`
	Config  string      `help:"Path to configuration file (default: surfacemap.yaml)" type:"existingfile"`
	Version VersionFlag `short:"V" help:"Print version information and exit"`
}

type CLI struct {
	Globals
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	cliutil.PrintColored(fmt.Sprintf("surfacemap v%s", vars["version"]), cliutil.ColorCyan)
	app.Exit(0)
	return nil
}

func main() {
	app := &CLI{}

	kongCtx := kong.Parse(app,
		kong.Name("surfacemap"),
		kong.Description("A tool for visualizing and analyzing your attack surface."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": internal.Version,
		},
	)

	err := kongCtx.Run()
	if err != nil {
		if errors.Is(err, cli.ErrNotImplemented) {
			os.Exit(2)
		}
		kongCtx.FatalIfErrorf(err)
	}
}
