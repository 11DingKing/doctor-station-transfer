package domain

import "fmt"

type ProjectState string

const (
	ProjectDraft      ProjectState = "draft"
	ProjectSubmitted  ProjectState = "submitted"
	ProjectReviewing  ProjectState = "reviewing"
	ProjectApproved   ProjectState = "approved"
	ProjectContracted ProjectState = "contracted"
	ProjectActive     ProjectState = "active"
	ProjectSuspended  ProjectState = "suspended"
	ProjectCompleted  ProjectState = "completed"
	ProjectRejected   ProjectState = "rejected"
)

type ReviewState string

const (
	ReviewPending   ReviewState = "pending"
	ReviewAccepted  ReviewState = "accepted"
	ReviewRecused   ReviewState = "recused"
	ReviewSubmitted ReviewState = "submitted"
)

type MilestoneState string

const (
	MilestonePending  MilestoneState = "pending"
	MilestoneAccepted MilestoneState = "accepted"
	MilestoneOverdue  MilestoneState = "overdue"
)

func (s ProjectState) CanMove(to ProjectState) bool {
	allowed := map[ProjectState][]ProjectState{ProjectDraft: {ProjectSubmitted}, ProjectSubmitted: {ProjectReviewing, ProjectRejected}, ProjectReviewing: {ProjectApproved, ProjectRejected}, ProjectApproved: {ProjectContracted}, ProjectContracted: {ProjectActive}, ProjectActive: {ProjectSuspended, ProjectCompleted}, ProjectSuspended: {ProjectActive, ProjectCompleted}}
	for _, x := range allowed[s] {
		if x == to {
			return true
		}
	}
	return false
}
func (s ProjectState) ValidateMove(to ProjectState) error {
	if !s.CanMove(to) {
		return fmt.Errorf("invalid project transition %s -> %s", s, to)
	}
	return nil
}
func (s ReviewState) CanMove(to ReviewState) bool {
	allowed := map[ReviewState][]ReviewState{ReviewPending: {ReviewAccepted, ReviewRecused}, ReviewAccepted: {ReviewSubmitted}}
	for _, x := range allowed[s] {
		if x == to {
			return true
		}
	}
	return false
}
func (s MilestoneState) CanMove(to MilestoneState) bool {
	return (s == MilestonePending && (to == MilestoneAccepted || to == MilestoneOverdue))
}
