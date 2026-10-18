package actionlint

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseConcurrencyQueue(t *testing.T) {
	for _, tc := range []struct {
		name, queue, cancel, wantError string
	}{
		{name: "default"},
		{name: "default_cancel", cancel: "true"},
		{name: "single", queue: "single"},
		{name: "single_cancel", queue: "single", cancel: "true"},
		{name: "max", queue: "max"},
		{name: "max_no_cancel", queue: "max", cancel: "false"},
		{name: "quoted_max", queue: `'max'`},
		{name: "conflict", queue: "max", cancel: "true", wantError: `"queue: max" cannot be combined with "cancel-in-progress: true"`},
		{name: "conflict_upper_bool", queue: "max", cancel: "TRUE", wantError: `"queue: max" cannot be combined with "cancel-in-progress: true"`},
		{name: "conflict_title_bool", queue: "max", cancel: "True", wantError: `"queue: max" cannot be combined with "cancel-in-progress: true"`},
		{name: "upper_false", queue: "max", cancel: "FALSE"},
		{name: "invalid", queue: "unlimited", wantError: `invalid value "unlimited"`},
		{name: "case_sensitive", queue: "Max", wantError: `invalid value "Max"`},
		{name: "whitespace", queue: `' max '`, wantError: `invalid value " max "`},
		{name: "integer", queue: "100", wantError: `invalid value "100"`},
		{name: "boolean", queue: "true", wantError: `invalid value "true"`},
		{name: "null", queue: "null", wantError: `invalid value "null"`},
		{name: "empty", queue: `''`, wantError: "string should not be empty"},
		{name: "sequence", queue: "[max]", wantError: "expected scalar node"},
		{name: "mapping", queue: "{value: max}", wantError: "expected scalar node"},
		{name: "queue_expression", queue: "${{ inputs.queue }}", cancel: "true"},
		{name: "cancel_expression", queue: "max", cancel: "${{ inputs.cancel }}"},
		{name: "both_expressions", queue: "${{ inputs.queue }}", cancel: "${{ inputs.cancel }}"},
		{name: "template", queue: "${{ inputs.prefix }}ax", cancel: "true"},
		// Parsing does not evaluate expressions, even constant ones.
		{name: "constant_queue_expression", queue: "${{ 'max' }}", cancel: "true"},
		{name: "constant_cancel_expression", queue: "max", cancel: "${{ true }}"},
	} {
		for _, jobLevel := range []bool{false, true} {
			for _, cancelFirst := range []bool{false, true} {
				t.Run(fmt.Sprintf("%s/job=%t/cancelFirst=%t", tc.name, jobLevel, cancelFirst), func(t *testing.T) {
					section := "group: test\n"
					queue, cancel := "", ""
					if tc.queue != "" {
						queue = "queue: " + tc.queue + "\n"
					}
					if tc.cancel != "" {
						cancel = "cancel-in-progress: " + tc.cancel + "\n"
					}
					if cancelFirst {
						section += cancel + queue
					} else {
						section += queue + cancel
					}
					indent := "  "
					prefix := "on: push\nconcurrency:\n"
					suffix := "jobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo test\n"
					if jobLevel {
						indent = "      "
						prefix = "on: push\njobs:\n  test:\n    runs-on: ubuntu-latest\n    concurrency:\n"
						suffix = "    steps:\n      - run: echo test\n"
					}
					src := prefix + indent + strings.ReplaceAll(strings.TrimSuffix(section, "\n"), "\n", "\n"+indent) + "\n" + suffix
					w, errs := Parse([]byte(src))
					if tc.wantError != "" {
						if len(errs) != 1 || !strings.Contains(errs[0].Message, tc.wantError) {
							t.Fatalf("expected one error containing %q, got %v", tc.wantError, errs)
						}
						return
					}
					if len(errs) != 0 {
						t.Fatal(errs)
					}
					c := w.Concurrency
					if jobLevel {
						c = w.Jobs["test"].Concurrency
					}
					if tc.queue == "" {
						if c.Queue != nil {
							t.Fatal("omitted queue must remain nil")
						}
					} else if c.Queue == nil || c.Queue.Value != strings.Trim(tc.queue, "'") || c.Queue.Quoted != strings.HasPrefix(tc.queue, "'") || c.Queue.Pos.Col != len(indent)+8 {
						t.Fatalf("unexpected queue AST: %#v", c.Queue)
					}
				})
			}
		}
	}
}

func TestParseConcurrencyShorthand(t *testing.T) {
	w, errs := Parse([]byte("on: push\nconcurrency: test\njobs:\n  test:\n    runs-on: ubuntu-latest\n    concurrency: job-test\n    steps:\n      - run: echo test\n"))
	if len(errs) != 0 {
		t.Fatal(errs)
	}
	for _, c := range []*Concurrency{w.Concurrency, w.Jobs["test"].Concurrency} {
		if c.Group == nil || c.Queue != nil || c.CancelInProgress != nil {
			t.Fatalf("unexpected shorthand AST: %#v", c)
		}
	}
}
