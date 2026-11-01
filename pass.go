package actionlint

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"time"
)

// Pass is an interface to traverse a workflow syntax tree
type Pass interface {
	// VisitStep is callback when visiting Step node. It returns internal error when it cannot continue the process
	VisitStep(node *Step) error
	// VisitJobPre is callback when visiting Job node before visiting its children. It returns internal error when it cannot continue the process
	VisitJobPre(node *Job) error
	// VisitJobPost is callback when visiting Job node after visiting its children. It returns internal error when it cannot continue the process
	VisitJobPost(node *Job) error
	// VisitWorkflowPre is callback when visiting Workflow node before visiting its children. It returns internal error when it cannot continue the process
	VisitWorkflowPre(node *Workflow) error
	// VisitWorkflowPost is callback when visiting Workflow node after visiting its children. It returns internal error when it cannot continue the process
	VisitWorkflowPost(node *Workflow) error
}

// Visitor visits syntax tree from root in depth-first order
type Visitor struct {
	passes []Pass
	dbg    io.Writer
}

// NewVisitor creates Visitor instance
func NewVisitor() *Visitor {
	return &Visitor{}
}

// AddPass adds new pass which is called on traversing a syntax tree
func (v *Visitor) AddPass(p Pass) {
	v.passes = append(v.passes, p)
}

// EnableDebug enables debug output when non-nil io.Writer value is given. All debug outputs from
// visitor will be written to the writer.
func (v *Visitor) EnableDebug(w io.Writer) {
	v.dbg = w
}

func (v *Visitor) reportElapsedTime(what string, start time.Time) {
	fmt.Fprintf(v.dbg, "[Visitor] %s took %vms\n", what, time.Since(start).Milliseconds())
}

// jobsInSourceOrder returns the jobs in the order they appear in the workflow file. Go randomizes map
// iteration order. Rules share state across jobs (e.g. the reusable workflow metadata cache), so
// visiting jobs in random order makes the reported errors and the rules reporting them vary between
// runs for the same input.
func jobsInSourceOrder(jobs map[string]*Job) []*Job {
	sorted := make([]*Job, 0, len(jobs))
	for _, j := range jobs {
		sorted = append(sorted, j)
	}
	slices.SortFunc(sorted, func(a, b *Job) int {
		if a.ID == nil || b.ID == nil || a.ID.Pos == nil || b.ID.Pos == nil {
			return 0 // Nodes not built by the parser have no positions
		}
		return cmp.Or(
			cmp.Compare(a.ID.Pos.Line, b.ID.Pos.Line),
			cmp.Compare(a.ID.Pos.Col, b.ID.Pos.Col),
		)
	})
	return sorted
}

// Visit visits given syntax tree in depth-first order
func (v *Visitor) Visit(n *Workflow) error {
	var t time.Time
	if v.dbg != nil {
		t = time.Now()
	}

	for _, p := range v.passes {
		if err := p.VisitWorkflowPre(n); err != nil {
			return err
		}
	}

	if v.dbg != nil {
		v.reportElapsedTime("VisitWorkflowPre", t)
		t = time.Now()
	}

	for _, j := range jobsInSourceOrder(n.Jobs) {
		if err := v.visitJob(j); err != nil {
			return err
		}
	}

	if v.dbg != nil {
		v.reportElapsedTime(fmt.Sprintf("Visiting %d jobs", len(n.Jobs)), t)
		t = time.Now()
	}

	for _, p := range v.passes {
		if err := p.VisitWorkflowPost(n); err != nil {
			return err
		}
	}

	if v.dbg != nil {
		v.reportElapsedTime("VisitWorkflowPost", t)
	}

	return nil
}

func (v *Visitor) visitJob(n *Job) error {
	var t time.Time
	if v.dbg != nil {
		t = time.Now()
	}

	for _, p := range v.passes {
		if err := p.VisitJobPre(n); err != nil {
			return err
		}
	}

	if v.dbg != nil {
		v.reportElapsedTime(fmt.Sprintf("VisitWorkflowJobPre at job %q", n.ID.Value), t)
		t = time.Now()
	}

	for _, s := range n.Steps {
		if err := v.visitStep(s); err != nil {
			return err
		}
	}

	if v.dbg != nil {
		v.reportElapsedTime(fmt.Sprintf("Visiting %d steps at job %q", len(n.Steps), n.ID.Value), t)
		t = time.Now()
	}

	for _, p := range v.passes {
		if err := p.VisitJobPost(n); err != nil {
			return err
		}
	}

	if v.dbg != nil {
		v.reportElapsedTime(fmt.Sprintf("VisitWorkflowJobPost at job %q", n.ID.Value), t)
	}

	return nil
}

func (v *Visitor) visitStep(n *Step) error {
	var t time.Time
	if v.dbg != nil {
		t = time.Now()
	}

	for _, p := range v.passes {
		if err := p.VisitStep(n); err != nil {
			return err
		}
	}

	if v.dbg != nil {
		v.reportElapsedTime(fmt.Sprintf("VisitStep at %s", n.Pos), t)
	}

	return nil
}
