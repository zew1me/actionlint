package actionlint

import (
	"path"
	"strings"
)

// localUsesSpec maps a repository-local uses reference to the workspace-relative
// form used as a metadata cache key. Self-repository references resolve against
// the defining repository, which is the project being linted, not a checkout.
// References escaping the project root with ".." are not treated as local.
// https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#example-using-an-action-in-the-same-repository-as-the-workflow
func localUsesSpec(spec string) (string, bool) {
	if strings.HasPrefix(spec, "$/") {
		local := "." + spec[1:]
		if c := path.Clean(local); c == ".." || strings.HasPrefix(c, "../") {
			return spec, false
		}
		return local, true
	}
	return spec, strings.HasPrefix(spec, "./")
}
