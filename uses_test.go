package actionlint

import (
	"path/filepath"
	"testing"
)

func TestLocalUsesSpec(t *testing.T) {
	for _, tc := range []struct {
		in, want string
		local    bool
	}{
		{"./actions/test", "./actions/test", true},
		{"$/actions/test", "./actions/test", true},
		{"$/workflows/test.yaml", "./workflows/test.yaml", true},
		{"$/", "./", true},
		{"$action", "$action", false},
		{"$", "$", false},
		{"../action", "../action", false},
		{"owner/repo@main", "owner/repo@main", false},
		{"docker://alpine:3", "docker://alpine:3", false},
	} {
		t.Run(tc.in, func(t *testing.T) {
			got, local := localUsesSpec(tc.in)
			if got != tc.want || local != tc.local {
				t.Fatalf("got (%q, %v), want (%q, %v)", got, local, tc.want, tc.local)
			}
		})
	}
}

func TestSelfRepositoryActionCacheAliases(t *testing.T) {
	for _, first := range []string{"./", "$/"} {
		t.Run(first, func(t *testing.T) {
			c := NewLocalActionsCache(&Project{root: filepath.Join("testdata", "action_metadata")}, nil)
			m, cached, err := c.FindMetadata(first + "action-yml")
			if err != nil || m == nil || cached {
				t.Fatalf("first lookup: %v, %v, %v", m, cached, err)
			}
			for _, prefix := range []string{"./", "$/"} {
				got, cached, err := c.FindMetadata(prefix + "action-yml")
				if err != nil || !cached || got != m {
					t.Fatalf("alias lookup: %v, %v, %v", got, cached, err)
				}
			}
		})
	}
}

func TestSelfRepositoryWorkflowCacheAliases(t *testing.T) {
	proj := &Project{root: filepath.Join("testdata", "reusable_workflow_metadata")}
	for _, first := range []string{"./", "$/"} {
		t.Run(first, func(t *testing.T) {
			c := NewLocalReusableWorkflowCache(proj, "", nil)
			m, err := c.FindMetadata(first + "ok.yaml")
			if err != nil || m == nil {
				t.Fatalf("first lookup: %v, %v", m, err)
			}
			for _, prefix := range []string{"./", "$/"} {
				got, err := c.FindMetadata(prefix + "ok.yaml")
				if err != nil || got != m {
					t.Fatalf("alias lookup: %v, %v", got, err)
				}
			}
			if _, err := c.FindMetadata(first + "missing.yaml"); err == nil {
				t.Fatal("missing workflow should report an error once")
			}
			for _, prefix := range []string{"./", "$/"} {
				if got, err := c.FindMetadata(prefix + "missing.yaml"); got != nil || err != nil {
					t.Fatalf("cached failure should not report again: %v, %v", got, err)
				}
			}
		})
	}

	c := NewLocalReusableWorkflowCache(proj, "", nil)
	c.WriteWorkflowCallEvent(filepath.Join(proj.root, "from_ast.yaml"), &WorkflowCallEvent{})
	m, err := c.FindMetadata("$/from_ast.yaml")
	if err != nil || m == nil {
		t.Fatalf("self-repository reference did not find AST-written metadata: %v, %v", m, err)
	}
}
