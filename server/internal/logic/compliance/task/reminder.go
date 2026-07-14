package task

import (
	"fmt"
	"time"
)

const (
	ReminderPending         = "pending"
	ReminderSent            = "sent"
	ReminderManualAttention = "manual_attention"
)

type Reminder struct {
	ID        uint64    `json:"id"`
	TaskID    uint64    `json:"taskId"`
	Channel   string    `json:"channel"`
	PlannedAt time.Time `json:"plannedAt"`
	Status    string    `json:"status"`
	Attempts  int       `json:"attempts"`
	LastError string    `json:"lastError"`
	SentAt    time.Time `json:"sentAt"`
}

func PlanReminders(taskID uint64, due time.Time, channels []string) []Reminder {
	items := make([]Reminder, 0, len(channels)*3)
	for _, channel := range channels {
		for _, days := range []int{7, 3, 1} {
			items = append(items, Reminder{TaskID: taskID, Channel: channel, PlannedAt: due.AddDate(0, 0, -days), Status: ReminderPending})
		}
	}
	return items
}

func (r Reminder) Key() string {
	return fmt.Sprintf("%d:%s:%d", r.TaskID, r.Channel, r.PlannedAt.Unix())
}
func (r Reminder) CanDispatch(enabled bool, now time.Time) bool {
	return enabled && r.Status == ReminderPending && !now.Before(r.PlannedAt)
}
func (r *Reminder) MarkFailure(message string) {
	r.Attempts++
	r.LastError = message
	if r.Attempts >= 3 {
		r.Status = ReminderManualAttention
	}
}
func (r *Reminder) MarkSent(now time.Time) { r.Status = ReminderSent; r.SentAt = now; r.LastError = "" }
