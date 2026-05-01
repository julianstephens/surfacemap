package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/julianstephens/go-utils/cliutil"
	"github.com/julianstephens/go-utils/generic"
	"github.com/julianstephens/go-utils/helpers"
	"github.com/julianstephens/go-utils/logger"

	"github.com/julianstephens/surfacemap/internal"
	"github.com/julianstephens/surfacemap/internal/cli"
	pkgerrors "github.com/julianstephens/surfacemap/pkg/errors"
)

func main() {
	app := &cli.CLI{
		Globals: cli.Globals{
			Version: cli.VersionFlag(internal.Version),
		},
	}

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

	if err := helpers.Ensure(app.Root, true); err != nil {
		kongCtx.FatalIfErrorf(pkgerrors.NewFileError(
			pkgerrors.ErrInvalidConfiguration,
			"validate",
			app.Root,
			nil,
		))
	}

	if app.Quiet {
		app.Logger = logger.NewNoop()
	} else {
		app.Logger = logger.New().WithField("component", "root")
		logLevel := "error"
		if app.VVerbose {
			logLevel = "debug"
		} else if app.Verbose {
			logLevel = "info"
		}
		if err := app.Logger.SetLogLevel(logLevel); err != nil {
			panic(pkgerrors.NewFileError(
				pkgerrors.ErrInternal,
				"initialize",
				"logger",
				err,
			))
		}
	}

	cliutil.SetDefaultConsole(cliutil.NewConsole(generic.If[cliutil.Formatter](app.NoColor, cliutil.NewPlainTextFormatter(), cliutil.NewColoredFormatter()), os.Stdout))
	app.Console = cliutil.DefaultConsole()

	err := kongCtx.Run(&app.Globals)
	if err != nil {
		if pkgerrors.HasSentinel(err, pkgerrors.ErrNotImplemented) {
			os.Exit(2)
		}
		kongCtx.FatalIfErrorf(err)
	}
}
