package actionlint

import (
	"slices"
	"testing"
)

type jobOrderRecorder struct {
	RuleBase
	ids []string
}

func (r *jobOrderRecorder) VisitJobPre(n *Job) error {
	r.ids = append(r.ids, n.ID.Value)
	return nil
}

func TestVisitorVisitsJobsInSourceOrder(t *testing.T) {
	src := `on: push
jobs:
  zebra:
    runs-on: ubuntu-latest
    steps:
      - run: echo
  apple:
    runs-on: ubuntu-latest
    steps:
      - run: echo
  mango:
    runs-on: ubuntu-latest
    steps:
      - run: echo
  banana:
    runs-on: ubuntu-latest
    steps:
      - run: echo
`
	w, errs := Parse([]byte(src))
	if len(errs) > 0 {
		t.Fatal(errs)
	}

	want := []string{"zebra", "apple", "mango", "banana"}
	// Go randomizes map iteration order so repeat to make an order-dependent bug visible
	for i := 0; i < 50; i++ {
		r := &jobOrderRecorder{}
		v := NewVisitor()
		v.AddPass(r)
		if err := v.Visit(w); err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(r.ids, want) {
			t.Fatalf("iteration %d: jobs visited in order %v, want %v", i, r.ids, want)
		}
	}
}
