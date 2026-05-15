package internal

type State int

const (
	StateNone State = iota
	StateAwaitingCaseLink
	StateAwaitingCaseTitle
	StateAwaitingCaseDescription
	StateAwaitingReviewLink
	StateAwaitingReviewAuthor
	StateAwaitingReviewText
	StateAwaitingEdit
	StateAwaitingRename
)

type UserState struct {
	State State
	Data  map[string]string
}
