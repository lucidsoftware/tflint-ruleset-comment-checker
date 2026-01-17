package main

import (
	"github.com/lucidsoftware/tflint-ruleset-comment-checker/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "comment-checker",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				&rules.ModuleAttributeCommentsRule{},
				// TODO: add rule for resources as well?
				// TODO: require comments on other types of things
			},
		},
	})
}
