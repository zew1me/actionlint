package actionlint

import (
	"path/filepath"
	"strings"
)

// localUsesSpec maps a repository-local uses reference to the workspace-relative
// form used as a metadata cache key. Self-repository references resolve against
// the defining repository, which is the project being linted, not a checkout.
// "$/" references escaping the project root are not treated as local.
// https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#example-using-an-action-in-the-same-repository-as-the-workflow
func localUsesSpec(spec string) (string, bool) {
	if strings.HasPrefix(spec, "$/") {
		// filepath.IsLocal also catches OS-specific escapes such as `..\` on Windows
		if rest := spec[2:]; rest != "" && !filepath.IsLocal(filepath.FromSlash(rest)) {
			return spec, false
		}
		return "." + spec[1:], true
	}
	return spec, strings.HasPrefix(spec, "./")
}

// isSelfRepositoryUsesWithRef returns true when the spec is a "$/" self-repository reference with an
// "@ref" suffix. Such a reference is invalid because it always resolves to the defining repository
// and commit, so it must not be resolved as a local path.
func isSelfRepositoryUsesWithRef(spec string) bool {
	return strings.HasPrefix(spec, "$/") && strings.ContainsRune(spec, '@')
}
