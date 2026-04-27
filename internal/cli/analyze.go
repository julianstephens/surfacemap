package cli

import (
	"context"
	"fmt"

	"github.com/julianstephens/surfacemap/pkg/discovery/terraform"
	"github.com/julianstephens/surfacemap/pkg/discovery/terraform/extractors"
	"github.com/julianstephens/surfacemap/pkg/model"
)

type AnalyzeCmd struct {
}

func (c *AnalyzeCmd) Run(globals *Globals) error {
	globals.Logger.Info("starting analysis...")
	if err := globals.Console.Info("Starting analysis..."); err != nil {
		return fmt.Errorf("failed to print to console: %w", err)
	}
	parser := terraform.NewHCLParser(globals.Logger.WithField("component", "HCLParser"))
	conf, err := parser.Parse(context.Background(), globals.Root)
	if err != nil {
		return fmt.Errorf("%w (root: %s)", err, globals.Root)
	}

	resourceCollection := &model.ResourceCollection{
		Lambdas: make([]*model.LambdaFunction, 0),
	}

	for _, res := range conf.Resources {
		globals.Logger.WithField("resource_id", res.ID).Info("discovered resource")
		extracted, err := extractors.ExtractResource(res)
		if err != nil {
			globals.Logger.WithField("resource_id", res.ID).Errorf("failed to extract resource: %v", err)
			continue
		}
		if extracted == nil {
			globals.Logger.WithField("resource_id", res.ID).Info("skipping unsupported resource type")
			continue
		}
		switch r := extracted.(type) {
		case *model.LambdaFunction:
			resourceCollection.Lambdas = append(resourceCollection.Lambdas, r)
			globals.Logger.WithField("resource_id", res.ID).Info("extracted Lambda function")
		default:
			globals.Logger.WithField("resource_id", res.ID).Warn("extracted resource of unknown type")
		}
	}

	globals.Logger.Info("analysis complete")
	if err := globals.Console.Info("Analysis complete"); err != nil {
		return fmt.Errorf("failed to print to console: %w", err)
	}
	return nil
}
