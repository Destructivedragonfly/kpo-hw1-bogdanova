package entities

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrBadRegistration   = errors.New("empty registration")
	ErrUnknownClass      = errors.New("unknown class")
	ErrAlreadyRegistered = errors.New("same registration")
)

type Creator func(Info) (Item, error)
type Option struct{ ID, Title string }
type class struct {
	id, title  string
	isConsumer bool
}

type entry struct {
	Option
	create Creator
}
type Factory struct{ entries []entry }

func NewFactory() *Factory { return &Factory{} }

func (factory *Factory) Register(id, title string, create Creator) error {
	if id == "" || title == "" || create == nil {
		return ErrBadRegistration
	}
	if factory.find(id) >= 0 {
		return fmt.Errorf("%w: %q", ErrAlreadyRegistered, id)
	}
	factory.entries = append(factory.entries, entry{Option{id, title}, create})

	return nil
}

func (factory *Factory) Create(id string, in Info) (Item, error) {
	i := factory.find(id)

	if i < 0 {
		return nil, fmt.Errorf("%w: %q", ErrUnknownClass, id)
	}

	return factory.entries[i].create(in)
}

func (factory *Factory) Options() []Option {
	res := make([]Option, len(factory.entries))

	for i, ent := range factory.entries {
		res[i] = ent.Option
	}

	return res
}

func (factory *Factory) find(id string) int {
	for i, e := range factory.entries {
		if e.ID == id {
			return i
		}
	}

	return -1
}

func NewVehicleFactory() *Factory {
	return build([]class{
		{"bicycle", "Bicycle", false},
		{"scooter", "Scooter", true},
		{"ebike", "Ebike", true},
	}, NewVehicle, false)
}

func NewGearFactory() *Factory {
	return build([]class{
		{"charging_station", "Charging station", true},
		{"charging_cabinet", "Charging cabinet", true},
		{"helmet", "Helmet", false},
		{"dock_station", "Dock station", true},
	}, NewGear, true)
}

func build(classes []class, construct func(string, Info) (Item, error), defaultName bool) *Factory {
	factory := NewFactory()

	for _, cl := range classes {
		err := factory.Register(cl.id, cl.title, func(in Info) (Item, error) {

			if !cl.isConsumer {
				in.EnergyKWh = 0
			} else if in.EnergyKWh <= 0 {
				return nil, ErrInvalidEnergy
			}
			if defaultName && strings.TrimSpace(in.Name) == "" {
				in.Name = cl.title
			}

			return construct(cl.title, in)
		})

		if err != nil {
			panic(err)
		}
	}

	return factory
}
