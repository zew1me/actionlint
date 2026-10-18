package actionlint

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestPermissionsNewScopes(t *testing.T) {
	for _, scope := range []struct {
		name   string
		levels []string
	}{
		{"code-quality", []string{"read", "write", "none"}},
		{"copilot-requests", []string{"write", "none"}},
		{"vulnerability-alerts", []string{"read", "none"}},
	} {
		for _, level := range []string{"read", "write", "none", "READ", "invalid"} {
			for _, where := range []string{"workflow", "job", "call"} {
				t.Run(scope.name+"/"+level+"/"+where, func(t *testing.T) {
					permission := fmt.Sprintf("permissions:\n%s%s: %s\n", map[string]string{"workflow": "  ", "job": "      ", "call": "      "}[where], scope.name, level)
					src := "on: push\n"
					if where == "workflow" {
						src += permission
					}
					src += "jobs:\n  test:\n"
					if where != "workflow" {
						src += "    " + permission
					}
					if where == "call" {
						src += "    uses: owner/repo/.github/workflows/test.yml@main\n"
					} else {
						src += "    runs-on: ubuntu-latest\n    steps:\n      - run: echo test\n"
					}
					l, err := NewLinter(io.Discard, &LinterOptions{})
					if err != nil {
						t.Fatal(err)
					}
					errs, err := l.Lint("test.yaml", []byte(src), nil)
					if err != nil {
						t.Fatal(err)
					}
					valid := false
					for _, allowed := range scope.levels {
						valid = valid || level == allowed
					}
					if valid {
						if len(errs) != 0 {
							t.Fatal(errs)
						}
					} else if len(errs) != 1 || errs[0].Kind != "permissions" || !strings.Contains(errs[0].Message, "invalid as permission of scope") {
						t.Fatalf("expected an invalid permission level error, got %v", errs)
					}
				})
			}
		}
	}
}
