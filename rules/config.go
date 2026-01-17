package rules

// Configuration for a rule about comments
// Contains a list of "input" block containin information
// about inputs that should commented
type InputCommentRuleConfig struct {
	Attributes []CommentRule `hclext:"attribute,block"`
}

type CommentRule struct {
	Name    string `hclext:"name"`
	Message string `hclext:"message,optional"`
	// TODO: add support for blocks? not needed for modules
}
