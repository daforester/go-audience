package models

import "time"

type SubmissionKind string

const (
	SUBMISSION_KIND_QUESTION        SubmissionKind = "Question"
	SUBMISSION_KIND_COMPLAINT       SubmissionKind = "Complaint"
	SUBMISSION_KIND_INCIDENT_REPORT SubmissionKind = "IncidentReport"
	SUBMISSION_KIND_OTHER           SubmissionKind = "Other"
)

var SubmissionKinds = []SubmissionKind{
	SUBMISSION_KIND_QUESTION,
	SUBMISSION_KIND_COMPLAINT,
	SUBMISSION_KIND_INCIDENT_REPORT,
	SUBMISSION_KIND_OTHER,
}

func ValidSubmissionKind(kind SubmissionKind) bool {
	for _, k := range SubmissionKinds {
		if k == kind {
			return true
		}
	}
	return false
}

type SubmissionStatus string

const (
	SUBMISSION_STATUS_PENDING  SubmissionStatus = "Pending"
	SUBMISSION_STATUS_APPROVED SubmissionStatus = "Approved"
	SUBMISSION_STATUS_REJECTED SubmissionStatus = "Rejected"
	SUBMISSION_STATUS_PLAYED   SubmissionStatus = "Played"
)

type Submission struct {
	ID         string
	AuthorUUID string
	AuthorName string
	Kind       SubmissionKind
	Text       string
	Status     SubmissionStatus
	StaffVotes map[string]bool // staff UUID -> voted
	Reactions  map[string]int  // audience UUID -> +1 / -1
	Created    time.Time
	Played     *time.Time
}

func NewSubmission(id string, author *Player, kind SubmissionKind, text string) *Submission {
	return &Submission{
		ID:         id,
		AuthorUUID: author.UUID,
		AuthorName: author.Name,
		Kind:       kind,
		Text:       text,
		Status:     SUBMISSION_STATUS_PENDING,
		StaffVotes: make(map[string]bool),
		Reactions:  make(map[string]int),
		Created:    time.Now(),
	}
}

func (S *Submission) StaffVoteCount() int {
	count := 0
	for _, voted := range S.StaffVotes {
		if voted {
			count++
		}
	}
	return count
}

func (S *Submission) Ups() int {
	count := 0
	for _, r := range S.Reactions {
		if r > 0 {
			count++
		}
	}
	return count
}

func (S *Submission) Downs() int {
	count := 0
	for _, r := range S.Reactions {
		if r < 0 {
			count++
		}
	}
	return count
}

func (S *Submission) NetScore() int {
	return S.Ups() - S.Downs()
}

// SetReaction records one audience member's reaction, replacing any earlier one.
// 0 withdraws the reaction.
func (S *Submission) SetReaction(uuid string, reaction int) {
	if S.Reactions == nil {
		S.Reactions = make(map[string]int)
	}
	switch {
	case reaction > 0:
		S.Reactions[uuid] = 1
	case reaction < 0:
		S.Reactions[uuid] = -1
	default:
		delete(S.Reactions, uuid)
	}
}

// ToggleStaffVote flips one staff member's vote and returns whether they now vote for it
func (S *Submission) ToggleStaffVote(uuid string) bool {
	if S.StaffVotes == nil {
		S.StaffVotes = make(map[string]bool)
	}
	if S.StaffVotes[uuid] {
		delete(S.StaffVotes, uuid)
		return false
	}
	S.StaffVotes[uuid] = true
	return true
}
