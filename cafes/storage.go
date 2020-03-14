package cafes

import (
	"2020_1_drop_table/owners"
	"errors"
	"fmt"
	"sync"
	"time"
)

//ToDo make photos available
type Cafe struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Address     string    `json:"address" validate:"required"`
	Description string    `json:"description" validate:"required"`
	OwnerID     int       `json:"ownerID"`
	OpenTime    time.Time `json:"openTime"`
	CloseTime   time.Time `json:"closeTime"`
	Photo       string    `json:"photo"`
}

func (c *Cafe) hasPermission(owner owners.Owner) bool {
	return c.OwnerID == owner.OwnerId
}

type cafesStorage struct {
	sync.Mutex
	cafes []Cafe
}

func NewCafesStorage() *cafesStorage {
	return &cafesStorage{}
}

func (cs *cafesStorage) append(value Cafe) Cafe {
	value.ID = cs.count()
	cs.cafes = append(cs.cafes, value)
	return value
}

func (cs *cafesStorage) set(i int, value Cafe) Cafe {
	cs.cafes[i] = value
	return value
}

func (cs *cafesStorage) get(index int) (Cafe, error) {
	if cs.count() > index && index >= 0 {
		item := cs.cafes[index]
		return item, nil
	}
	notFoundErrorMessage := fmt.Sprintf("Cafe not fount")
	return Cafe{}, errors.New(notFoundErrorMessage)
}

func (cs *cafesStorage) count() int {
	return len(cs.cafes)
}

func (cs *cafesStorage) Append(value Cafe) (error, Cafe) {
	cs.Lock()
	defer cs.Unlock()
	value = cs.append(value)
	return nil, value
}

func (cs *cafesStorage) Count() int {
	cs.Lock()
	defer cs.Unlock()
	return cs.count()
}

func (cs *cafesStorage) Get(index int) (Cafe, error) {
	cs.Lock()
	defer cs.Unlock()
	return cs.get(index)
}

func (cs *cafesStorage) getOwnerCafes(owner owners.Owner) []Cafe {
	var ownerCafes []Cafe
	for i := 0; i < cs.Count(); i++ {
		cafe, _ := cs.Get(i)
		if cafe.OwnerID == owner.OwnerId {
			ownerCafes = append(ownerCafes, cafe)
		}
	}
	return ownerCafes
}

func (cs *cafesStorage) Set(i int, value Cafe) (Cafe, error) {
	if i > cs.Count() {
		err := errors.New(fmt.Sprintf("no user with id: %d", i))
		return Cafe{}, err
	}
	value.ID = i
	cs.Lock()
	defer cs.Unlock()
	value = cs.set(i, value)
	return value, nil
}

var Storage = NewCafesStorage()
