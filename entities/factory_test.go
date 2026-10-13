package entities

import (
	"errors"
	"testing"
)

func TestFactoryRegisterAndCreate(t *testing.T) {
	f := NewFactory()
	create := func(in Info) (Item, error) { return NewGear("Test", in) }

	if err := f.Register("a", "Type A", create); err != nil {
		t.Fatal(err)
	}
	if err := f.Register("a", "Again", create); !errors.Is(err, ErrAlreadyRegistered) {
		t.Errorf("duplicate id: got %v, want ErrAlreadyRegistered", err)
	}
	if err := f.Register("", "No id", create); !errors.Is(err, ErrBadRegistration) {
		t.Errorf("empty id: got %v, want ErrBadRegistration", err)
	}
	if _, err := f.Create("a", Info{Name: "n", Num: "1"}); err != nil {
		t.Errorf("create registered type: %v", err)
	}
	if _, err := f.Create("zzz", Info{}); !errors.Is(err, ErrUnknownClass) {
		t.Errorf("unknown type: got %v, want ErrUnknownClass", err)
	}
	if len(f.Options()) != 1 {
		t.Errorf("got %d options, want 1", len(f.Options()))
	}
}

func TestVehicleFactoryEnergyRules(t *testing.T) {
	f := NewVehicleFactory()

	_, err := f.Create("scooter", Info{Name: "n", Num: "1", Simpl: 5})
	if !errors.Is(err, ErrInvalidEnergy) {
		t.Errorf("scooter without energy: got %v, want ErrInvalidEnergy", err)
	}

	bike, err := f.Create("bicycle", Info{Name: "n", Num: "2", Simpl: 5, EnergyKWh: 5})
	if err != nil {
		t.Fatal(err)
	}

	if bike.(EnergyConsumer).DailyEnergyKWh() != 0 {
		t.Error("bicycle must ignore energy")
	}
}

func TestGearFactoryDefaultName(t *testing.T) {
	item, err := NewGearFactory().Create("charging_cabinet", Info{Num: "CC-1", EnergyKWh: 6})

	if err != nil {
		t.Fatal(err)
	}
	if item.Name() != "Charging cabinet" {
		t.Errorf("got name %q, want the kind title", item.Name())
	}
}
