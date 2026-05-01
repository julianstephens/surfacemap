package cli

import (
	"context"
	"fmt"

	"github.com/julianstephens/surfacemap/pkg/discovery/terraform"
	"github.com/julianstephens/surfacemap/pkg/discovery/terraform/extractors"
	pkgerrors "github.com/julianstephens/surfacemap/pkg/errors"
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

	// Accumulate extraction errors
	extractionErrors := pkgerrors.NewMultiError("extracting resources")
	successCount := 0
	skippedCount := 0

	for _, res := range conf.Resources {
		globals.Logger.WithField("resource_id", res.ID).Info("discovered resource")
		extracted, err := extractors.ExtractResource(res)
		if err != nil {
			globals.Logger.WithField("resource_id", res.ID).Errorf("failed to extract resource: %v", err)
			extractionErrors.Add(err)
			continue
		}
		if extracted == nil {
			globals.Logger.WithField("resource_id", res.ID).Info("skipping unsupported resource type")
			skippedCount++
			continue
		}
		switch r := extracted.(type) {
		case *model.LambdaFunction:
			resourceCollection.Lambdas = append(resourceCollection.Lambdas, r)
			globals.Logger.WithField("resource_id", res.ID).Info("extracted Lambda function")
			successCount++
		default:
			globals.Logger.WithField("resource_id", res.ID).Warn("extracted resource of unknown type")
			skippedCount++
		}
	}

	// Log summary
	globals.Logger.Infof("extraction summary: %d successful, %d skipped, %d failed", successCount, skippedCount, len(extractionErrors.Errors))
	if extractionErrors.HasErrors() {
		globals.Logger.Warnf("encountered %d extraction errors:", len(extractionErrors.Errors))
		for i, err := range extractionErrors.Errors {
			globals.Logger.Warnf("  %d. %v", i+1, err)
		}
	}

	globals.Logger.Info("analysis complete")
	if err := globals.Console.Info("Analysis complete"); err != nil {
		return fmt.Errorf("failed to print to console: %w", err)
	}
	return nil
}
