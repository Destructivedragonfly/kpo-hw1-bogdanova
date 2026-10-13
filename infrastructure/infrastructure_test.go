package infrastructure

import (
	"errors"
	"kpo_hw1/entities"
	"testing"
)

func vehicle(t *testing.T, num string, simpl int, kwh float64) entities.Item {
	t.Helper()
	v, err := entities.NewVehicle("Scooter", entities.Info{Name: "n", Num: num, Simpl: simpl, EnergyKWh: kwh})

	if err != nil {
		t.Fatal(err)
	}

	return v
}

func gear(t *testing.T, num string, kwh float64) entities.Item {
	t.Helper()
	g, err := entities.NewGear("Gear", entities.Info{Name: "g", Num: num, EnergyKWh: kwh})

	if err != nil {
		t.Fatal(err)
	}

	return g
}

func newReception() *Reception {
	return New(NewCurStorage(), ServiceCenter{MaxEnergyKWh: 5}, CheckIsSimpl(6))
}

func TestInspect(t *testing.T) {
	sc := ServiceCenter{MaxEnergyKWh: 5}

	if !sc.Inspect(vehicle(t, "1", 7, 3), true).IsAccepted {
		t.Error("good vehicle must be accepted")
	}
	if sc.Inspect(vehicle(t, "2", 7, 3), false).IsAccepted {
		t.Error("mechanic said no: must be rejected")
	}
	if sc.Inspect(vehicle(t, "3", 7, 9), true).IsAccepted {
		t.Error("energy over the limit: must be rejected")
	}
}

func TestStorageRejectsDuplicate(t *testing.T) {
	s := NewCurStorage()
	_ = s.AddNew(gear(t, "G-1", 0))

	if err := s.AddNew(gear(t, "G-1", 0)); !errors.Is(err, entities.ErrDuplicateNum) {
		t.Errorf("got %v, want ErrDuplicateNum", err)
	}
}

func TestExam(t *testing.T) {
	rec := newReception()

	if res, _ := rec.Exam(vehicle(t, "A", 7, 2), true); !res.IsAccepted || !rec.CheckHas("A") {
		t.Error("accepted vehicle must be stored")
	}
	if res, _ := rec.Exam(vehicle(t, "B", 7, 2), false); res.IsAccepted || rec.CheckHas("B") {
		t.Error("rejected vehicle must not be stored")
	}
	if _, err := rec.Exam(vehicle(t, "A", 7, 2), true); !errors.Is(err, entities.ErrDuplicateNum) {
		t.Errorf("duplicate number: got %v, want ErrDuplicateNum", err)
	}
}

func TestRegisterGear(t *testing.T) {
	rec := newReception()

	if err := rec.Register(gear(t, "H-1", 0)); err != nil || !rec.CheckHas("H-1") {
		t.Errorf("gear must be registered, err = %v", err)
	}
}

func TestReport(t *testing.T) {
	rec := newReception()
	_, _ = rec.Exam(vehicle(t, "SC-1", 7, 1.5), true)
	_ = rec.Register(gear(t, "C-1", 6.5))

	r := rec.Report()
	if r.Total != 2 || r.EnergyKWh != 8 || r.AmountInClass["Scooter"] != 1 {
		t.Errorf("unexpected report: %+v", r)
	}
}

func TestBeginners(t *testing.T) {
	rec := newReception()
	_, _ = rec.Exam(vehicle(t, "A", 6, 1), true)
	_, _ = rec.Exam(vehicle(t, "B", 5, 1), true)
	_ = rec.Register(gear(t, "H-1", 0))

	got := rec.Beginners()
	if len(got) != 1 || got[0].InventoryNum() != "A" {
		t.Errorf("got %d beginner items, want only A", len(got))
	}
}

func TestInventorySorted(t *testing.T) {
	rec := newReception()
	_ = rec.Register(gear(t, "B", 0))
	_ = rec.Register(gear(t, "A", 0))

	if items := rec.Inventory(); items[0].InventoryNum() != "A" {
		t.Errorf("first item is %q, want A", items[0].InventoryNum())
	}
}
