package infrastructure

import (
	"kpo_hw1/entities"
)

type ExamResult struct {
	IsAccepted bool
	Problem    string
}

type ServiceCenter struct {
	MaxEnergyKWh float64
}

type Inspector interface {
	Inspect(item entities.Item, result bool) ExamResult
}

func (center ServiceCenter) Inspect(item entities.Item, result bool) ExamResult {
	if item == nil {
		return ExamResult{Problem: entities.ErrNilItem.Error()}
	}
	if !result {
		return ExamResult{Problem: entities.ErrRejected.Error()}
	}

	energyCons, ok := item.(entities.EnergyConsumer)
	if ok && center.MaxEnergyKWh > 0 && energyCons.DailyEnergyKWh() > center.MaxEnergyKWh {
		return ExamResult{Problem: entities.ErrEnergyOverLimit.Error()}
	}

	return ExamResult{IsAccepted: true}
}
