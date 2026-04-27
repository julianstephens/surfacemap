package cli

import (
	"context"
	"fmt"

	"github.com/julianstephens/surfacemap/pkg/discovery/terraform"
)

type AnalyzeCmd struct {
}

func (c *AnalyzeCmd) Run(globals *Globals) error {
	globals.Logger.Info("starting analysis...")
	if err := globals.Console.Info("Starting analysis..."); err != nil {
		return fmt.Errorf("failed to print to console: %w", err)
	}
	parser := terraform.NewHCLParser(globals.Logger.WithField("component", "HCLParser"))
	_, err := parser.Parse(context.Background(), globals.Root)
	if err != nil {
		return fmt.Errorf("%w (root: %s)", err, globals.Root)
	}

	globals.Logger.Info("analysis complete")
	if err := globals.Console.Info("Analysis complete"); err != nil {
		return fmt.Errorf("failed to print to console: %w", err)
	}
	return nil
}
