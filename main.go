package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-template/rules"
)

// PluginConfig is the configuration for the plugin
type PluginConfig struct {
	AttributeNames []string `hclext:"attribute_names"`
}

// RuleSet is the custom ruleset with plugin configuration
type RuleSet struct {
	*tflint.BuiltinRuleSet
	config *PluginConfig
}

// ConfigSchema returns the plugin config schema
func (r *RuleSet) ConfigSchema() *hclext.BodySchema {
	return hclext.ImpliedBodySchema(&PluginConfig{})
}

// ApplyConfig applies the plugin configuration
func (r *RuleSet) ApplyConfig(body *hclext.BodyContent) error {
	config := &PluginConfig{}
	if err := hclext.DecodeBody(body, nil, config); err != nil {
		return err
	}
	r.config = config
	
	// Apply config to the module attribute comment checker rule
	for _, rule := range r.BuiltinRuleSet.Rules {
		if checker, ok := rule.(*rules.ModuleAttributeCommentCheckerRule); ok {
			checker.SetConfig(config.AttributeNames)
		}
	}
	
	return r.BuiltinRuleSet.ApplyConfig(body)
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &RuleSet{
			BuiltinRuleSet: &tflint.BuiltinRuleSet{
				Name:    "template",
				Version: "0.1.0",
				Rules: []tflint.Rule{
					rules.NewAwsInstanceExampleTypeRule(),
					rules.NewAwsS3BucketExampleLifecycleRule(),
					rules.NewGoogleComputeSSLPolicyRule(),
					rules.NewTerraformBackendTypeRule(),
					rules.NewModuleAttributeCommentCheckerRule(),
				},
			},
		},
	})
}
