package infrastructure

import (
	"cmp"
	"kpo_hw1/entities"
	"slices"
)

type BeginnerRule func(entities.BeginnerRated) bool

func CheckIsSimpl(min int) BeginnerRule {
	return func(entity entities.BeginnerRated) bool { return entity.SimplForBeginners() >= min }
}

type Reception struct {
	storage   Storage
	rule      BeginnerRule
	inspector Inspector
}

func New(storage Storage, inspector Inspector, rule BeginnerRule) *Reception {
	return &Reception{storage: storage, inspector: inspector, rule: rule}
}

func (rec *Reception) CheckHas(num string) bool { return rec.storage.Has(num) }
func (rec *Reception) Register(item entities.Item) error {
	if _, isVehicle := item.(entities.BeginnerRated); isVehicle {
		return entities.ErrNeedExam
	}

	return rec.storage.AddNew(item)
}

func (rec *Reception) Exam(item entities.Item, result bool) (ExamResult, error) {
	if item == nil {
		return ExamResult{}, entities.ErrNilItem
	}
	if rec.storage.Has(item.InventoryNum()) {
		return ExamResult{}, entities.ErrDuplicateNum
	}

	final := rec.inspector.Inspect(item, result)
	if !final.IsAccepted {
		return final, nil
	}
	if err := rec.storage.AddNew(item); err != nil {
		return ExamResult{}, err
	}

	return final, nil
}

type Report struct {
	Total         int
	EnergyKWh     float64
	AmountInClass map[string]int
}

func (rec *Reception) Report() Report {
	items := rec.storage.Clone()
	report := Report{Total: len(items), AmountInClass: map[string]int{}}

	for _, it := range items {
		report.AmountInClass[it.Class()]++
		if consumer, ok := it.(entities.EnergyConsumer); ok {
			report.EnergyKWh += consumer.DailyEnergyKWh()
		}
	}

	return report
}

func (rec *Reception) Inventory() []entities.Item {
	items := rec.storage.Clone()
	slices.SortFunc(items, func(a, b entities.Item) int {
		return cmp.Compare(a.InventoryNum(), b.InventoryNum())
	})

	return items
}

func (rec *Reception) Beginners() []entities.BeginnerRated {
	var res []entities.BeginnerRated
	for _, it := range rec.storage.Clone() {
		if entity, ok := it.(entities.BeginnerRated); ok && rec.rule(entity) {
			res = append(res, entity)
		}
	}

	slices.SortFunc(res, func(a, b entities.BeginnerRated) int {
		return cmp.Compare(a.InventoryNum(), b.InventoryNum())
	})

	return res
}
