package entities

import "strings"

type base struct {
	class, name, num string
	energy           float64
}

func newBase(class string, in Info) (base, error) {
	name, num := strings.TrimSpace(in.Name), strings.TrimSpace(in.Num)

	switch {
	case name == "":
		return base{}, ErrEmptyName
	case num == "":
		return base{}, ErrEmptyNum
	case in.EnergyKWh < 0:
		return base{}, ErrInvalidEnergy
	}

	return base{class: class, name: name, num: num, energy: in.EnergyKWh}, nil
}

func (b base) Name() string            { return b.name }
func (b base) InventoryNum() string    { return b.num }
func (b base) Class() string           { return b.class }
func (b base) DailyEnergyKWh() float64 { return b.energy }

type vehicle struct {
	base
	simpl int
}

func (v vehicle) SimplForBeginners() int { return v.simpl }

func NewVehicle(class string, in Info) (Item, error) {
	b, err := newBase(class, in)
	if err != nil {
		return nil, err
	}
	if in.Simpl < MinSimpl || in.Simpl > MaxSimpl {
		return nil, ErrInvalidSimpl
	}

	return vehicle{base: b, simpl: in.Simpl}, nil
}

type gear struct{ base }

func NewGear(class string, in Info) (Item, error) {
	b, err := newBase(class, in)
	if err != nil {
		return nil, err
	}

	return gear{b}, nil
}
