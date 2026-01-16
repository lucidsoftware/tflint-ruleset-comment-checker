package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// ModuleAttributeCommentCheckerRule checks whether module attributes have comments
type ModuleAttributeCommentCheckerRule struct {
	tflint.DefaultRule
	attributeNames []string
}

// NewModuleAttributeCommentCheckerRule returns a new rule
func NewModuleAttributeCommentCheckerRule() *ModuleAttributeCommentCheckerRule {
	return &ModuleAttributeCommentCheckerRule{}
}

// Name returns the rule name
func (r *ModuleAttributeCommentCheckerRule) Name() string {
	return "module_attribute_comment_checker"
}

// Enabled returns whether the rule is enabled by default
func (r *ModuleAttributeCommentCheckerRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *ModuleAttributeCommentCheckerRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *ModuleAttributeCommentCheckerRule) Link() string {
	return ""
}

// SetConfig sets the attribute names to check from plugin configuration
func (r *ModuleAttributeCommentCheckerRule) SetConfig(attributeNames []string) {
	r.attributeNames = attributeNames
}

// Check checks whether configured attributes in module calls have comments
func (r *ModuleAttributeCommentCheckerRule) Check(runner tflint.Runner) error {
	// If no attributes configured, nothing to check
	if len(r.attributeNames) == 0 {
		return nil
	}

	// Get all module blocks
	modules, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "module",
				Body: &hclext.BodySchema{
					Mode: hclext.SchemaJustAttributesMode,
				},
				LabelNames: []string{"name"},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	// Get all files to access raw content for comment checking
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	// Check each module block
	for _, module := range modules.Blocks {
		// Check each configured attribute
		for _, attrName := range r.attributeNames {
			if attr, exists := module.Body.Attributes[attrName]; exists {
				// Check if there's a comment immediately preceding this attribute
				if !hasCommentBefore(attr, files[attr.Range.Filename]) {
					err := runner.EmitIssue(
						r,
						fmt.Sprintf("attribute %q in module %q should have a comment", attrName, module.Labels[0]),
						attr.Range,
					)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// hasCommentBefore checks if there's a comment on the line immediately before the attribute
func hasCommentBefore(attr *hclext.Attribute, file *hcl.File) bool {
	if file == nil {
		return false
	}

	// Get the line number of the attribute
	attrLine := attr.Range.Start.Line

	// Get the file bytes
	lines := strings.Split(string(file.Bytes), "\n")

	// Check if there's a previous line and if attrLine is within bounds
	// attrLine is 1-indexed, so we need attrLine-2 to be valid (>=0 and < len(lines))
	if attrLine <= 1 || attrLine-1 >= len(lines) {
		return false
	}

	// Get the previous line (attrLine is 1-indexed)
	prevLine := strings.TrimSpace(lines[attrLine-2])

	// Check if the previous line is a comment
	return strings.HasPrefix(prevLine, "#") || strings.HasPrefix(prevLine, "//")
}
