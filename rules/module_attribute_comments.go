package rules

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// ModuleAttributeCommentsRule checks whether module attributes have comments
type ModuleAttributeCommentsRule struct {
	tflint.DefaultRule
}

// Name returns the rule name
func (r *ModuleAttributeCommentsRule) Name() string {
	return "module_attribute_comments"
}

// Enabled returns whether the rule is enabled by default
func (r *ModuleAttributeCommentsRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
// We use ERROR because if you are using this ruleset, then you probably want to
// enforce that there are comments.
// TODO: Make this configurable later?
func (r *ModuleAttributeCommentsRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *ModuleAttributeCommentsRule) Link() string {
	return "https://github.com/lucidsoftware/tflint-ruleset-comment-checker/blob/main/README.md#module_attribute_comments"
}

// Check checks whether configured attributes in module calls have comments
func (r *ModuleAttributeCommentsRule) Check(runner tflint.Runner) error {
	config := &InputCommentRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	// If no inputs configured, nothing to check
	if len(config.Attributes) == 0 {
		logger.Warn("module_input_comments is enabled, but no inputs were configured")
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
		// I'm not sure if we need expand mode to be none or not
	}, nil)
	if err != nil {
		return err
	}

	// Check each module block
	for _, module := range modules.Blocks {
		// Check each configured attribute
		for _, input := range config.Attributes {
			if attr, exists := module.Body.Attributes[input.Name]; exists {
				file, err := runner.GetFile(attr.Range.Filename)
				if err != nil {
					return err
				}
				// Check if there's a comment immediately preceding this attribute
				if !hasCommentBefore(attr, file) {
					err := runner.EmitIssue(
						r,
						fmt.Sprintf("%q in module %q should have a comment. %s", input.Name, module.Labels[0], input.Message),
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

	// Get the byte at the beginning of the attribute definition
	start := attr.Range.Start.Byte

	// Figure out the range of the previous line before the atttribute occurance
	prevLineEnd := bytes.LastIndexByte(file.Bytes[:start], '\n')
	// We add 1, so that the start is AFTER the newline
	prevLineStart := bytes.LastIndexByte(file.Bytes[:prevLineEnd], '\n') + 1

	prevLine := file.Bytes[prevLineStart:prevLineEnd]

	// Strip off any whitespace
	prevLine = bytes.TrimSpace(prevLine)

	return (bytes.HasPrefix(prevLine, []byte("#")) && len(prevLine) > 1) ||
		(bytes.HasPrefix(prevLine, []byte("//")) && len(prevLine) > 2)
}
