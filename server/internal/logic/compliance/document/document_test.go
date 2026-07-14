package document

import "testing"

func TestProcessingStatusPreservesUploadOnOCRFailure(t *testing.T) {
	got := ResolveProcessingStatus(false, 0)
	if got != StatusManualReview {
		t.Fatalf("status=%s want=%s", got, StatusManualReview)
	}
}

func TestProcessingStatusRequiresConfirmationForLowConfidence(t *testing.T) {
	if got := ResolveProcessingStatus(true, 0.69); got != StatusPendingConfirmation {
		t.Fatalf("status=%s", got)
	}
	if got := ResolveProcessingStatus(true, 0.95); got != StatusPendingConfirmation {
		t.Fatalf("AI result must still require confirmation, status=%s", got)
	}
}

func TestDuplicateHashIsWarningNotDeletion(t *testing.T) {
	items := []BusinessDocument{{ID: 1, OpcEntityID: 10, FileHash: "same", Deleted: false}}
	duplicate := FindDuplicate(items, 10, "same")
	if duplicate == nil || duplicate.Deleted {
		t.Fatalf("duplicate=%+v", duplicate)
	}
}

func TestConfirmRejectsDifferentMember(t *testing.T) {
	item := BusinessDocument{ID: 1, MemberID: 2, ProcessStatus: StatusPendingConfirmation}
	if err := Confirm(&item, 1, "bank_statement", map[string]any{}); err == nil {
		t.Fatal("expected cross-member confirmation to fail")
	}
}

func TestConfirmMarksDocumentConfirmed(t *testing.T) {
	item := BusinessDocument{ID: 1, MemberID: 1, ProcessStatus: StatusPendingConfirmation}
	fields := map[string]any{"amount": 100.50}
	if err := Confirm(&item, 1, "bank_statement", fields); err != nil {
		t.Fatalf("Confirm returned error: %v", err)
	}
	if item.ProcessStatus != StatusConfirmed || item.DocumentType != "bank_statement" || item.ConfirmedAt == 0 {
		t.Fatalf("item=%+v", item)
	}
}
