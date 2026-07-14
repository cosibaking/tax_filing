package task

import "fmt"

const (
	StatusNotStarted          = "not_started"
	StatusPreparing           = "preparing"
	StatusPendingConfirmation = "pending_confirmation"
	StatusCompleted           = "completed"
	StatusOverdue             = "overdue"
	StatusCancelled           = "cancelled"

	EventStart    = "start"
	EventSubmit   = "submit"
	EventReject   = "reject"
	EventComplete = "complete"
	EventCancel   = "cancel"
)

func Transition(current, event string) (string, error) {
	transitions := map[string]map[string]string{
		StatusNotStarted: {
			EventStart: StatusPreparing, EventCancel: StatusCancelled,
		},
		StatusPreparing: {
			EventSubmit: StatusPendingConfirmation, EventCancel: StatusCancelled,
		},
		StatusPendingConfirmation: {
			EventReject: StatusPreparing, EventComplete: StatusCompleted,
		},
		StatusOverdue: {
			EventStart: StatusPreparing, EventSubmit: StatusPendingConfirmation,
			EventComplete: StatusCompleted, EventCancel: StatusCancelled,
		},
	}
	next, ok := transitions[current][event]
	if !ok {
		return "", fmt.Errorf("task transition %s -> %s is not allowed", current, event)
	}
	return next, nil
}
