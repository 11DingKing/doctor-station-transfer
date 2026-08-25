package domain

import (
	"errors"
	"testing"
)

func TestProjectTransitions(t *testing.T) {
	cases := []struct {
		from, to ProjectState
		ok       bool
	}{{ProjectDraft, ProjectSubmitted, true}, {ProjectSubmitted, ProjectReviewing, true}, {ProjectReviewing, ProjectApproved, true}, {ProjectApproved, ProjectContracted, true}, {ProjectContracted, ProjectActive, true}, {ProjectActive, ProjectSuspended, true}, {ProjectSuspended, ProjectActive, true}, {ProjectDraft, ProjectCompleted, false}, {ProjectCompleted, ProjectDraft, false}, {ProjectRejected, ProjectActive, false}}
	for _, c := range cases {
		if got := c.from.CanMove(c.to); got != c.ok {
			t.Errorf("%s->%s got %v", c.from, c.to, got)
		}
	}
}
func TestReviewTransitions(t *testing.T) {
	if !ReviewPending.CanMove(ReviewAccepted) {
		t.Fatal("pending should accept")
	}
	if !ReviewPending.CanMove(ReviewRecused) {
		t.Fatal("pending should recuse")
	}
	if ReviewSubmitted.CanMove(ReviewPending) {
		t.Fatal("submitted cannot rewind")
	}
}
func TestErrorChain(t *testing.T) {
	e := &CodedError{Code: "x", Err: ErrConflict}
	if !errors.Is(e, ErrConflict) {
		t.Fatal("unwrap lost")
	}
	if e.Error() == "" {
		t.Fatal("empty error")
	}
}
func TestMilestoneTransitions(t *testing.T) {
	if !MilestonePending.CanMove(MilestoneAccepted) {
		t.Fatal("complete")
	}
	if !MilestonePending.CanMove(MilestoneOverdue) {
		t.Fatal("overdue")
	}
	if MilestoneAccepted.CanMove(MilestonePending) {
		t.Fatal("accepted rewind")
	}
}
