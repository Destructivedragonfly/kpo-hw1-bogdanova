package entities

import "errors"

const (
	MinSimpl = 1
	MaxSimpl = 10
)

var (
	ErrNilItem         = errors.New("no entity found")
	ErrEmptyName       = errors.New("null naming")
	ErrEmptyNum        = errors.New("null inventary number")
	ErrInvalidSimpl    = errors.New("simplicity is not in [1,10]")
	ErrInvalidEnergy   = errors.New("energy consumption is wrong")
	ErrDuplicateNum    = errors.New("same entity")
	ErrRejected        = errors.New("item rejected")
	ErrEnergyOverLimit = errors.New("energy consumption over the limit")
	ErrNeedExam        = errors.New("exam is obligatory")
)

type Item interface {
	Name() string
	InventoryNum() string
	Class() string
}

type BeginnerRated interface {
	Item
	SimplForBeginners() int
}

type EnergyConsumer interface {
	DailyEnergyKWh() float64
}

type Info struct {
	Name      string
	Num       string
	Simpl     int
	EnergyKWh float64
}
