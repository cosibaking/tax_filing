package profile

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
)

func TestMergeProfileJSONProducesValidJSON(t *testing.T) {
	merged := mergeProfileJSON([]byte(`{"region":"CN-BJ"}`), []byte(`{"employeeBand":"0"}`))
	var decoded map[string]any
	if err := json.Unmarshal([]byte(merged), &decoded); err != nil {
		t.Fatalf("invalid merged JSON: %v", err)
	}
}

type memoryRepository struct {
	items map[uint64][]Profile
}

func (m *memoryRepository) Active(_ context.Context, opcID uint64) (*Profile, error) {
	items := m.items[opcID]
	if len(items) == 0 {
		return nil, nil
	}
	item := items[len(items)-1]
	return &item, nil
}

func (m *memoryRepository) SaveVersion(_ context.Context, item Profile) (*Profile, error) {
	item.ID = uint64(len(m.items[item.OpcEntityID]) + 1)
	m.items[item.OpcEntityID] = append(m.items[item.OpcEntityID], item)
	return &item, nil
}

type fixedOwner struct {
	memberID uint64
}

func (f fixedOwner) MemberOwnsOpc(_ context.Context, memberID, _ uint64) (bool, error) {
	return memberID == f.memberID, nil
}

func validInput() SaveInput {
	return SaveInput{
		MemberID: 1,
		OpcID:    10,
		Source:   SourceUser,
		Data: Data{
			Region:               "CN-BJ",
			EntityType:           "one_person_limited_company",
			TaxpayerType:         "small_scale",
			VATPeriod:            "quarterly",
			EmployeeCount:        0,
			InvoiceEnabled:       true,
			HasRevenue:           false,
			HasPublicBankAccount: true,
			Complexity:           "low",
		},
	}
}

func TestSaveCreatesFirstActiveVersion(t *testing.T) {
	repo := &memoryRepository{items: map[uint64][]Profile{}}
	service := NewService(repo, fixedOwner{memberID: 1})

	got, changed, err := service.Save(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if !changed || got.Version != 1 || got.Status != StatusActive {
		t.Fatalf("unexpected first version: changed=%v profile=%+v", changed, got)
	}
	if got.Labels["employeeBand"] != "0" || got.Labels["region"] != "CN-BJ" {
		t.Fatalf("unexpected labels: %#v", got.Labels)
	}
}

func TestSaveDefaultsEmptySourceToUser(t *testing.T) {
	repo := &memoryRepository{items: map[uint64][]Profile{}}
	service := NewService(repo, fixedOwner{memberID: 1})
	in := validInput()
	in.Source = ""

	got, _, err := service.Save(context.Background(), in)
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if got.Source != SourceUser {
		t.Fatalf("source=%q want=%q", got.Source, SourceUser)
	}
}

func TestSaveSameDataDoesNotCreateVersion(t *testing.T) {
	repo := &memoryRepository{items: map[uint64][]Profile{}}
	service := NewService(repo, fixedOwner{memberID: 1})
	first, _, _ := service.Save(context.Background(), validInput())

	second, changed, err := service.Save(context.Background(), validInput())
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if changed || second.ID != first.ID || len(repo.items[10]) != 1 {
		t.Fatalf("same data created a version: changed=%v items=%d", changed, len(repo.items[10]))
	}
}

func TestSaveChangedDataCreatesNextVersion(t *testing.T) {
	repo := &memoryRepository{items: map[uint64][]Profile{}}
	service := NewService(repo, fixedOwner{memberID: 1})
	_, _, _ = service.Save(context.Background(), validInput())

	in := validInput()
	in.Data.EmployeeCount = 2
	got, changed, err := service.Save(context.Background(), in)
	if err != nil || !changed || got.Version != 2 {
		t.Fatalf("expected version 2, got=%+v changed=%v err=%v", got, changed, err)
	}
	if reflect.DeepEqual(repo.items[10][0].Data, got.Data) {
		t.Fatal("new version did not preserve changed data")
	}
}

func TestAdvisorChangeRequiresReason(t *testing.T) {
	repo := &memoryRepository{items: map[uint64][]Profile{}}
	service := NewService(repo, fixedOwner{memberID: 1})
	in := validInput()
	in.Source = SourceAdvisor

	if _, _, err := service.Save(context.Background(), in); err == nil {
		t.Fatal("expected advisor change without reason to fail")
	}
}

func TestSaveRejectsOpcOwnedByAnotherMember(t *testing.T) {
	repo := &memoryRepository{items: map[uint64][]Profile{}}
	service := NewService(repo, fixedOwner{memberID: 99})

	if _, _, err := service.Save(context.Background(), validInput()); err == nil {
		t.Fatal("expected ownership validation error")
	}
}
