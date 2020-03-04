package owners

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

//ToDo make photos available
type Owner struct {
	ID       int       `json:"id"`
	Name     string    `json:"name" validate:"required,min=4,max=100"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=100"`
	EditedAt time.Time `json:"editedAt" validate:"required"`
	Photo    string    `json:"photo"`
}

type owners struct {
	sync.Mutex
	owners []Owner
}

func NewOwnersStorage() *owners {
	return &owners{}
}

func (ds *owners) append(value Owner) Owner {
	value.ID = ds.count()
	ds.owners = append(ds.owners, value)
	return value
}

func (ds *owners) set(i int, value Owner) Owner {
	ds.owners[i] = value
	return value
}

func (ds *owners) get(index int) (Owner, error) {
	if ds.count() > index && index >= 0 {
		item := ds.owners[index]
		return item, nil
	}
	notFoundErrorMessage := fmt.Sprintf("Owner not fount")
	return Owner{}, errors.New(notFoundErrorMessage)
}

func (ds *owners) count() int {
	return len(ds.owners)
}

func (ds *owners) isRegistered(email, password string) (int, Owner) {
	password = GetMD5Hash(password)
	for i := 0; i < ds.count(); i++ {
		owner, _ := ds.Get(i)
		if owner.Email == email && owner.Password == password {
			return 2, owner
		} else if owner.Email == email {
			return 1, Owner{}
		}
	}
	return -1, Owner{}
}

func (ds *owners) Append(value Owner) (error, Owner) {
	if n, _ := ds.isRegistered(value.Email, ""); n != -1 {
		err := errors.New("user with this email already existed")
		return err, Owner{}
	}
	value.Password = GetMD5Hash(value.Password)
	ds.Lock()
	defer ds.Unlock()
	value = ds.append(value)
	return nil, value
}

func (ds *owners) Set(i int, value Owner) (Owner, error) {
	if i > ds.Count() {
		err := errors.New(fmt.Sprintf("no user with id: %d", i))
		return Owner{}, err
	}
	value.ID = i

	ds.Lock()
	defer ds.Unlock()
	value = ds.set(i, value)
	return value, nil
}

func (ds *owners) Get(index int) (Owner, error) {
	ds.Lock()
	defer ds.Unlock()
	return ds.get(index)
}

func (ds *owners) Count() int {
	ds.Lock()
	defer ds.Unlock()
	return ds.count()
}

func (ds *owners) Existed(email string, password string) (bool, Owner) {
	code, owner := ds.isRegistered(email, password)
	return code == 2, owner
}

func hasPermission(owner Owner, cookie string) bool {
	actualOwner, err := StorageSession.GetOwnerByCookie(cookie)
	if err != nil {
		return false
	}
	return actualOwner.ID == owner.ID
}
