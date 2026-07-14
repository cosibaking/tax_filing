package report

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type failingNarrator struct{}

func (failingNarrator) Narrate(context.Context, Input) (string, string, error) {
	return "", "", errors.New("model down")
}

func TestGenerateMarksMissingData(t *testing.T) {
	result := Generate(Input{PeriodKey: "2026-07"})
	if !strings.Contains(result.Content, "数据不足") {
		t.Fatalf("content=%s", result.Content)
	}
}

func TestNarratorFailureFallsBackToDeterministicReport(t *testing.T) {
	svc := NewService(nil, failingNarrator{})
	result, err := svc.Build(context.Background(), Input{PeriodKey: "2026-07"})
	if err != nil {
		t.Fatal(err)
	}
	if result.AIModel != "" || !strings.Contains(result.Content, "数据不足") {
		t.Fatalf("result=%+v", result)
	}
}

func TestPublishedReportIsImmutable(t *testing.T) {
	if err := ValidateMutation("published"); err == nil {
		t.Fatal("published report must be immutable")
	}
}
