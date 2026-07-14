package ruleengine

import (
	"context"
	"testing"
)

type versionMemoryStore struct {
	items map[uint64]RuleVersion
}

func (m *versionMemoryStore) Get(_ context.Context, id uint64) (*RuleVersion, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, nil
	}
	return &item, nil
}

func (m *versionMemoryStore) Save(_ context.Context, item RuleVersion) error {
	m.items[item.ID] = item
	return nil
}

func (m *versionMemoryStore) RetirePublished(_ context.Context, ruleID, exceptID uint64) error {
	for id, item := range m.items {
		if item.RuleID == ruleID && id != exceptID && item.Status == RuleStatusPublished {
			item.Status = RuleStatusRetired
			m.items[id] = item
		}
	}
	return nil
}

func validRuleVersion() RuleVersion {
	return RuleVersion{
		ID: 1, RuleID: 10, Version: 1, Status: RuleStatusDraft, CreatedBy: 100,
		Expression: Expression{All: []Node{{Condition: &Condition{
			Field: "profile.region", Op: "eq", Value: "CN-BJ",
		}}}},
	}
}

func TestRuleVersionTransitionsToReview(t *testing.T) {
	store := &versionMemoryStore{items: map[uint64]RuleVersion{1: validRuleVersion()}}
	service := NewVersionService(store)

	got, err := service.Transition(context.Background(), TransitionInput{VersionID: 1, Action: ActionSubmitReview, OperatorID: 100})
	if err != nil || got.Status != RuleStatusReviewing {
		t.Fatalf("got=%+v err=%v", got, err)
	}
}

func TestRuleVersionRequiresDifferentReviewer(t *testing.T) {
	version := validRuleVersion()
	version.Status = RuleStatusReviewing
	store := &versionMemoryStore{items: map[uint64]RuleVersion{1: version}}
	service := NewVersionService(store)

	if _, err := service.Transition(context.Background(), TransitionInput{VersionID: 1, Action: ActionPublish, OperatorID: 100}); err == nil {
		t.Fatal("expected creator publishing own version to fail")
	}
}

func TestPublishingRetiresPreviousVersion(t *testing.T) {
	old := validRuleVersion()
	old.ID = 1
	old.Status = RuleStatusPublished
	current := validRuleVersion()
	current.ID = 2
	current.Version = 2
	current.Status = RuleStatusReviewing
	store := &versionMemoryStore{items: map[uint64]RuleVersion{1: old, 2: current}}
	service := NewVersionService(store)

	got, err := service.Transition(context.Background(), TransitionInput{VersionID: 2, Action: ActionPublish, OperatorID: 200})
	if err != nil || got.Status != RuleStatusPublished || store.items[1].Status != RuleStatusRetired {
		t.Fatalf("current=%+v old=%+v err=%v", got, store.items[1], err)
	}
}

func TestPublishedVersionCannotReturnToDraft(t *testing.T) {
	version := validRuleVersion()
	version.Status = RuleStatusPublished
	store := &versionMemoryStore{items: map[uint64]RuleVersion{1: version}}
	service := NewVersionService(store)

	if _, err := service.Transition(context.Background(), TransitionInput{VersionID: 1, Action: ActionSubmitReview, OperatorID: 200}); err == nil {
		t.Fatal("expected immutable published version transition to fail")
	}
}

func TestSimulateReturnsMatchWithoutSaving(t *testing.T) {
	store := &versionMemoryStore{items: map[uint64]RuleVersion{1: validRuleVersion()}}
	service := NewVersionService(store)
	matched, err := service.Simulate(context.Background(), 1, map[string]any{
		"profile": map[string]any{"region": "CN-BJ"},
	})
	if err != nil || !matched || len(store.items) != 1 {
		t.Fatalf("matched=%v items=%d err=%v", matched, len(store.items), err)
	}
}
