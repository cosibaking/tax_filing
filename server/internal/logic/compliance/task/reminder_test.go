package task

import (
	"testing"
	"time"
)

func TestPlanRemindersAtSevenThreeAndOneDay(t *testing.T) {
	due := time.Date(2026, 10, 20, 9, 0, 0, 0, time.Local)
	got := PlanReminders(10, due, []string{"in_app"})
	if len(got) != 3 {
		t.Fatalf("len=%d", len(got))
	}
	for i, days := range []int{7, 3, 1} {
		want := due.AddDate(0, 0, -days)
		if !got[i].PlannedAt.Equal(want) {
			t.Errorf("item %d=%v want=%v", i, got[i].PlannedAt, want)
		}
	}
}

func TestReminderKeyIsIdempotent(t *testing.T) {
	when := time.Date(2026, 10, 19, 9, 0, 0, 0, time.Local)
	a := Reminder{TaskID: 1, Channel: "in_app", PlannedAt: when}
	b := Reminder{TaskID: 1, Channel: "in_app", PlannedAt: when}
	if a.Key() != b.Key() {
		t.Fatalf("keys differ: %s %s", a.Key(), b.Key())
	}
}

func TestReminderFailureNeedsManualAttentionAfterThreeAttempts(t *testing.T) {
	r := Reminder{Status: ReminderPending}
	for i := 0; i < 3; i++ {
		r.MarkFailure("channel unavailable")
	}
	if r.Status != ReminderManualAttention || r.Attempts != 3 {
		t.Fatalf("reminder=%+v", r)
	}
}

func TestDisabledReminderDoesNotDispatch(t *testing.T) {
	r := Reminder{Status: ReminderPending}
	if r.CanDispatch(false, time.Now()) {
		t.Fatal("disabled reminder should not dispatch")
	}
}
