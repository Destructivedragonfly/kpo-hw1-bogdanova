package infrastructure

import (
	"kpo_hw1/entities"
	"slices"
)

type CurStorage struct {
	items []entities.Item
	nums  map[string]bool
}

type Storage interface {
	Has(number string) bool
	AddNew(item entities.Item) error
	Clone() []entities.Item
}

func NewCurStorage() *CurStorage {
	return &CurStorage{nums: make(map[string]bool)}
}

func (store *CurStorage) AddNew(item entities.Item) error {
	if item == nil {
		return entities.ErrNilItem
	}
	if store.nums[item.InventoryNum()] {
		return entities.ErrDuplicateNum
	}

	store.nums[item.InventoryNum()] = true
	store.items = append(store.items, item)

	return nil
}

func (store *CurStorage) Has(num string) bool { return store.nums[num] }

func (store *CurStorage) Clone() []entities.Item {
	return slices.Clone(store.items)
}
