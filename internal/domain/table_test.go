package domain

import "testing"

func TestAllIllegalProjectMoves(t *testing.T) {
	states := []ProjectState{ProjectDraft, ProjectSubmitted, ProjectReviewing, ProjectApproved, ProjectContracted, ProjectActive, ProjectSuspended, ProjectCompleted, ProjectRejected}
	for _, from := range states {
		for _, to := range states {
			if from == to {
				if from.CanMove(to) {
					t.Errorf("self move %s", from)
				}
				continue
			}
			allowed := from.CanMove(to)
			if allowed {
				if err := from.ValidateMove(to); err != nil {
					t.Fatalf("%s %s %v", from, to, err)
				}
			} else {
				if err := from.ValidateMove(to); err == nil {
					t.Errorf("illegal %s %s", from, to)
				}
			}
		}
	}
}
func TestTransitionPath(t *testing.T) {
	p := ProjectDraft
	path := []ProjectState{ProjectSubmitted, ProjectReviewing, ProjectApproved, ProjectContracted, ProjectActive, ProjectCompleted}
	for _, next := range path {
		if !p.CanMove(next) {
			t.Fatalf("path %s to %s", p, next)
		}
		p = next
	}
	if p != ProjectCompleted {
		t.Fatal(p)
	}
}
func TestSuspensionRoundTrip(t *testing.T) {
	if !ProjectActive.CanMove(ProjectSuspended) {
		t.Fatal("suspend")
	}
	if !ProjectSuspended.CanMove(ProjectActive) {
		t.Fatal("resume")
	}
	if ProjectSuspended.CanMove(ProjectSubmitted) {
		t.Fatal("submitted")
	}
}
func TestReviewMatrix(t *testing.T) {
	valid := map[ReviewState][]ReviewState{ReviewPending: {ReviewAccepted, ReviewRecused}, ReviewAccepted: {ReviewSubmitted}, ReviewRecused: {}, ReviewSubmitted: {}}
	for from, tos := range valid {
		for _, to := range []ReviewState{ReviewPending, ReviewAccepted, ReviewRecused, ReviewSubmitted} {
			want := false
			for _, x := range tos {
				if x == to {
					want = true
				}
			}
			if from.CanMove(to) != want {
				t.Errorf("%s to %s", from, to)
			}
		}
	}
}
func TestMilestoneMatrix(t *testing.T) {
	for _, from := range []MilestoneState{MilestonePending, MilestoneAccepted, MilestoneOverdue} {
		for _, to := range []MilestoneState{MilestonePending, MilestoneAccepted, MilestoneOverdue} {
			want := from == MilestonePending && (to == MilestoneAccepted || to == MilestoneOverdue)
			if from.CanMove(to) != want {
				t.Errorf("%s %s", from, to)
			}
		}
	}
}
func TestRoleValues(t *testing.T) {
	vals := []Role{RoleAdmin, RoleResearcher, RoleReviewer, RoleAuditor}
	seen := map[Role]bool{}
	for _, v := range vals {
		if seen[v] {
			t.Fatal(v)
		}
		seen[v] = true
		if string(v) == "" {
			t.Fatal("empty")
		}
	}
}
