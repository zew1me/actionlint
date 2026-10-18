package actionlint

import "strings"

// localUsesSpec maps a repository-local uses reference to the workspace-relative
// form used as a metadata cache key. Self-repository references resolve against
// the defining repository, which is the project being linted, not a checkout.
// https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#example-using-an-action-in-the-same-repository-as-the-workflow
func localUsesSpec(spec string) (string, bool) {
	if strings.HasPrefix(spec, "$/") {
		return "." + spec[1:], true
	}
	return spec, strings.HasPrefix(spec, "./")
}
