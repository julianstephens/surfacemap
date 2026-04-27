package extractors

import (
	"errors"

	"github.com/julianstephens/surfacemap/pkg/model"
)

var (
	ErrMissingFunctionName = errors.New("missing function_name attribute in Lambda resource")
	ErrMissingRoleArn      = errors.New("missing role attribute in Lambda resource")
)

type LambdaExtractor struct {
}

func (le *LambdaExtractor) Extract(resource model.Resource) (*model.LambdaFunction, error) {
	if resource.Type != "aws_lambda_function" {
		return nil, ErrUnsupportedResourceType(resource.Type)
	}

	funcName, ok := resource.Attributes["function_name"].(string)
	if !ok {
		return nil, ErrMissingFunctionName
	}
	roleArn, ok := resource.Attributes["role"].(string)
	if !ok {
		return nil, ErrMissingRoleArn
	}
	handler, ok := resource.Attributes["handler"].(string)
	if !ok {
		handler = ""
	}
	runtime, ok := resource.Attributes["runtime"].(string)
	if !ok {
		runtime = ""
	}
	sourceCodePath, ok := resource.Attributes["filename"].(string)
	if !ok {
		sourceCodePath = ""
	}
	architectures, ok := resource.Attributes["architectures"].([]string)
	if !ok {
		architectures = nil
	}
	environment, ok := resource.Blocks["environment"]
	var envVars map[string]string
	if ok && len(environment) > 0 {
		envBlock := environment[0]
		if vars, ok := envBlock["variables"].(map[string]string); ok {
			envVars = vars
		}
	}
	vpcConfig, ok := resource.Blocks["vpc_config"]
	var vpcCfg *model.VPCConfig
	if ok && len(vpcConfig) > 0 {
		vpcBlock := vpcConfig[0]
		subnetIds, _ := vpcBlock["subnet_ids"].([]string)
		securityGroupIds, _ := vpcBlock["security_group_ids"].([]string)
		vpcCfg = &model.VPCConfig{
			SubnetIds:        &subnetIds,
			SecurityGroupIds: &securityGroupIds,
		}
	}

	return &model.LambdaFunction{
		Resource:       resource,
		FunctionName:   funcName,
		RoleArn:        roleArn,
		Handler:        &handler,
		Runtime:        &runtime,
		SourceCodePath: &sourceCodePath,
		Architectures:  &architectures,
		Environment: &struct {
			Variables map[string]string "json:\"variables\""
		}{
			Variables: envVars,
		},
		VPCConfig: vpcCfg,
	}, nil
}

func ExtractLambdaFunction(resource model.Resource) (*model.LambdaFunction, error) {
	extractor := &LambdaExtractor{}
	return extractor.Extract(resource)
}
