package entities

import (
	"errors"
	"testing"
)

func TestNewVehicle(t *testing.T) {
	v, err := NewVehicle("Scooter", Info{Name: " Pro ", Num: " SC-1 ", Simpl: 7, EnergyKWh: 1.5})

	if err != nil {
		t.Fatal(err)
	}
	if v.Name() != "Pro" || v.InventoryNum() != "SC-1" || v.Class() != "Scooter" {
		t.Errorf("unexpected vehicle: %q %q %q", v.Name(), v.InventoryNum(), v.Class())
	}
	if r, ok := v.(BeginnerRated); !ok || r.SimplForBeginners() != 7 {
		t.Error("vehicle must expose simplicity 7")
	}
}

func TestNewVehicleInvalid(t *testing.T) {
	cases := []struct {
		in   Info
		want error
	}{
		{Info{Name: "", Num: "1", Simpl: 5}, ErrEmptyName},
		{Info{Name: "n", Num: "", Simpl: 5}, ErrEmptyNum},
		{Info{Name: "n", Num: "1", Simpl: 0}, ErrInvalidSimpl},
		{Info{Name: "n", Num: "1", Simpl: 11}, ErrInvalidSimpl},
		{Info{Name: "n", Num: "1", Simpl: 5, EnergyKWh: -1}, ErrInvalidEnergy},
	}

	for _, c := range cases {
		if _, err := NewVehicle("Scooter", c.in); !errors.Is(err, c.want) {
			t.Errorf("%+v: got %v, want %v", c.in, err, c.want)
		}
	}
}

func TestNewGear(t *testing.T) {
	g, err := NewGear("Helmet", Info{Name: "Giro", Num: "H-1"})

	if err != nil {
		t.Fatal(err)
	}
	if g.InventoryNum() != "H-1" {
		t.Errorf("got number %q, want H-1", g.InventoryNum())
	}
	if _, ok := g.(BeginnerRated); ok {
		t.Error("gear must not have a beginner rating")
	}
	if _, err := NewGear("Helmet", Info{Name: "Giro"}); !errors.Is(err, ErrEmptyNum) {
		t.Errorf("got %v, want ErrEmptyNum", err)
	}
}
