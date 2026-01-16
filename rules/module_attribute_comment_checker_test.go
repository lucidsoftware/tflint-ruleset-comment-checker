package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_ModuleAttributeCommentCheckerRule(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "attribute with comment - no issue",
			Content: `
module "example" {
  source = "./modules/example"
  # This is a comment for instance_type
  instance_type = "t2.micro"
}`,
			Config: `
rule "module_attribute_comment_checker" {
  enabled = true
  attribute_names = ["instance_type"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "attribute without comment - issue found",
			Content: `
module "example" {
  source = "./modules/example"
  instance_type = "t2.micro"
}`,
			Config: `
rule "module_attribute_comment_checker" {
  enabled = true
  attribute_names = ["instance_type"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewModuleAttributeCommentCheckerRule(),
					Message: `attribute "instance_type" in module "example" should have a comment`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 4, Column: 29},
					},
				},
			},
		},
		{
			Name: "multiple attributes - mixed comments",
			Content: `
module "example" {
  source = "./modules/example"
  # Comment for instance_type
  instance_type = "t2.micro"
  count = 5
  # Comment for name
  name = "test"
}`,
			Config: `
rule "module_attribute_comment_checker" {
  enabled = true
  attribute_names = ["instance_type", "count", "name"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewModuleAttributeCommentCheckerRule(),
					Message: `attribute "count" in module "example" should have a comment`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 6, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 12},
					},
				},
			},
		},
		{
			Name: "module without configured attribute - no issue",
			Content: `
module "example" {
  source = "./modules/example"
  other_attr = "value"
}`,
			Config: `
rule "module_attribute_comment_checker" {
  enabled = true
  attribute_names = ["instance_type"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "multiple modules - some with comments",
			Content: `
module "with_comment" {
  source = "./modules/example"
  # This has a comment
  instance_type = "t2.micro"
}

module "without_comment" {
  source = "./modules/example"
  instance_type = "t2.small"
}`,
			Config: `
rule "module_attribute_comment_checker" {
  enabled = true
  attribute_names = ["instance_type"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewModuleAttributeCommentCheckerRule(),
					Message: `attribute "instance_type" in module "without_comment" should have a comment`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 10, Column: 3},
						End:      hcl.Pos{Line: 10, Column: 29},
					},
				},
			},
		},
		{
			Name: "comment with double slashes",
			Content: `
module "example" {
  source = "./modules/example"
  // This is a double-slash comment
  instance_type = "t2.micro"
}`,
			Config: `
rule "module_attribute_comment_checker" {
  enabled = true
  attribute_names = ["instance_type"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "empty attribute names - no issues",
			Content: `
module "example" {
  source = "./modules/example"
  instance_type = "t2.micro"
}`,
			Config: `
rule "module_attribute_comment_checker" {
  enabled = true
  attribute_names = []
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewModuleAttributeCommentCheckerRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content, ".tflint.hcl": test.Config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
