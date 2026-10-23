package actionlint

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseCacheMode(t *testing.T) {
	for _, mode := range []string{"read", "write", "write-only", "none", "'read'", "'write'", "'write-only'", "'none'", "\"read\"", "\"write\"", "\"write-only\"", "\"none\""} {
		t.Run(mode, func(t *testing.T) {
			src := fmt.Sprintf("on: push\ncache-mode: %s\njobs:\n  test:\n    cache-mode: %s\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo test\n  call:\n    uses: owner/repo/.github/workflows/test.yaml@main\n    cache-mode: %s\n", mode, mode, mode)
			w, errs := Parse([]byte(src))
			if len(errs) != 0 {
				t.Fatal(errs)
			}
			for _, got := range []*String{w.CacheMode, w.Jobs["test"].CacheMode, w.Jobs["call"].CacheMode} {
				if got == nil || got.Value != strings.Trim(mode, "'\"") || got.Pos == nil {
					t.Fatalf("unexpected cache mode: %#v", got)
				}
			}
			if w.CacheMode.Pos.Line != 2 || w.CacheMode.Pos.Col != 13 || w.Jobs["test"].CacheMode.Pos.Line != 5 || w.Jobs["call"].WorkflowCall == nil {
				t.Fatalf("cache-mode position or reusable workflow call was lost: %#v", w)
			}
		})
	}
}

func TestParseInvalidCacheMode(t *testing.T) {
	for _, value := range []string{"READ", "read-write", "invalid", "true", "42", "null", "~", "''", "", "[]", "{}", "[read]", "{mode: read}", "' read '", "${{ 'read' }}", "read-${{ inputs.mode }}"} {
		for _, where := range []string{"workflow", "job", "call"} {
			t.Run(where+"/"+value, func(t *testing.T) {
				src := "on: push\n"
				if where == "workflow" {
					src += "cache-mode: " + value + "\n"
				}
				src += "jobs:\n  test:\n"
				if where != "workflow" {
					src += "    cache-mode: " + value + "\n"
				}
				if where == "call" {
					src += "    uses: owner/repo/.github/workflows/test.yaml@main\n"
				} else {
					src += "    runs-on: ubuntu-latest\n    steps:\n      - run: echo test\n"
				}
				_, errs := Parse([]byte(src))
				if len(errs) != 1 || errs[0].Kind != "syntax-check" {
					t.Fatalf("expected one invalid cache-mode error, got %v", errs)
				}
				line, col := 4, 17
				if where == "workflow" {
					line, col = 2, 13
				}
				if value == "" {
					col-- // YAML positions an implicit null immediately after the colon.
				}
				if errs[0].Line != line || errs[0].Column != col {
					t.Fatalf("expected error at %d:%d, got %v", line, col, errs[0])
				}
				if !strings.Contains(errs[0].Message, "cache-mode") && !strings.Contains(errs[0].Message, "scalar node") && !strings.Contains(errs[0].Message, "string should not be empty") {
					t.Fatalf("expected cache-mode value diagnostic, got %v", errs[0])
				}
			})
		}
	}
}

func TestParseCacheModeAliases(t *testing.T) {
	for _, mode := range []string{"read", "write", "write-only", "none"} {
		t.Run(mode, func(t *testing.T) {
			src := fmt.Sprintf("on: push\ncache-mode: &mode '%s'\njobs:\n  test:\n    cache-mode: *mode\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo test\n  call:\n    uses: owner/repo/.github/workflows/test.yaml@main\n    cache-mode: *mode\n", mode)
			w, errs := Parse([]byte(src))
			if len(errs) != 0 {
				t.Fatal(errs)
			}
			for _, got := range []*String{w.CacheMode, w.Jobs["test"].CacheMode, w.Jobs["call"].CacheMode} {
				if got == nil || got.Value != mode || !got.Quoted {
					t.Fatalf("cache-mode alias lost its value or quoting: %#v", got)
				}
			}
		})
	}
}

func TestCacheModeOmittedAndOverrides(t *testing.T) {
	for _, setting := range []string{"", "cache-mode: read\n", "cache-mode: write\n", "cache-mode: none\n"} {
		src := "on: push\n" + setting + "jobs:\n  default:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo test\n  override:\n    runs-on: ubuntu-latest\n    cache-mode: write-only\n    steps:\n      - run: echo test\n"
		w, errs := Parse([]byte(src))
		if len(errs) != 0 {
			t.Fatal(errs)
		}
		if (w.CacheMode == nil) != (setting == "") || w.Jobs["default"].CacheMode != nil || w.Jobs["override"].CacheMode.Value != "write-only" {
			t.Fatalf("cache-mode should preserve omissions and overrides: %#v", w)
		}
	}
}
