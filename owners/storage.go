package owners

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"time"
)

type Owner struct {
	OwnerId  int       `json:"id"`
	Name     string    `json:"name" validate:"required,min=4,max=100"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=100"`
	EditedAt time.Time `json:"editedAt" validate:"required"`
	Photo    string    `json:"photo"`
}

type OwnerStorage struct {
	db *sqlx.DB
}

func logErr(err error, message string, where Owner) {
	log.Error().Msgf("Error: %v, %s,  in -> %v", err, message, where)
}

func (s *OwnerStorage) Append(own Owner) error {
	tx, err := s.db.Begin()
	own.Password = GetMD5Hash(own.Password)
	if err != nil {
		logErr(err, "When trying to connect to bd", own)
		return err
	}
	_, err = tx.Exec("insert into owner(OwnerId, name, email, password, editedat, photo) values ($1,$2,$3,$4,$5,$6)", own.OwnerId, own.Name, own.Email, own.Password, own.EditedAt, own.Photo)
	if err != nil {
		logErr(err, "When trying to append data", own)
		return err
	}
	err = tx.Commit()
	return err

}

func (s *OwnerStorage) AppendList(owners []Owner) error {
	tx, err := s.db.Begin()
	if err != nil {
		logErr(err, "When trying to connect to bd", owners[0])
		return err
	}
	for _, own := range owners {

		own.Password = GetMD5Hash(own.Password)
		_, err := tx.Exec("insert into owner(OwnerId, name, email, password, editedat, photo) values ($1,$2,$3,$4,$5,$6)", own.OwnerId, own.Name, own.Email, own.Password, own.EditedAt, own.Photo)
		if err != nil {
			logErr(err, "When trying to append data", own)
			return err
		}
	}
	err = tx.Commit()
	return err
}

func (s *OwnerStorage) CreateTable() error {
	schema := `CREATE TABLE IF NOT EXISTS Owner
(
    OwnerID  integer NOT NULL PRIMARY KEY,
    Name     text,
    Email    text,
    Password text,
    EditedAt timestamp,
    Photo    text
)
`

	_, err := s.db.Exec(schema)
	return err
}

func isOwnerEmpty(own *Owner) bool {
	if own.OwnerId == 0 {
		log.Info().Msgf("Owner not found")
		return true
	}
	return false
}

func (s *OwnerStorage) GetByEmailAndPassword(email string, password string) (Owner, error) {
	own := Owner{}
	password = GetMD5Hash(password)
	err := s.db.Get(&own, "select * from owner where password=$1 AND email=$2", password, email)
	if isOwnerEmpty(&own) {
		notFoundErrorMessage := fmt.Sprintf("Owner not found")
		return Owner{}, errors.New(notFoundErrorMessage)
	}
	own.EditedAt = own.EditedAt.UTC()
	return own, err
}

func (s *OwnerStorage) Get(id int) (Owner, error) {
	own := Owner{}
	err := s.db.Get(&own, "select * from owner where ownerid=$1", id)
	if err != nil {
		return Owner{}, err
	}
	if isOwnerEmpty(&own) {
		notFoundErrorMessage := fmt.Sprintf("Owner not found")
		return Owner{}, errors.New(notFoundErrorMessage)
	}
	own.EditedAt = own.EditedAt.UTC()
	return own, err
}

func (s *OwnerStorage) Set(id int, newOwner Owner) error {
	newOwner.Password = GetMD5Hash(newOwner.Password)
	_, err := s.db.Exec("UPDATE owner SET name = $1,email=$2,password=$3,editedat=$4,photo=$5 WHERE ownerid = $6", newOwner.Name, newOwner.Email, newOwner.Password, newOwner.EditedAt, newOwner.Photo, id)

	return err
}

func (s *OwnerStorage) Existed(email string, password string) (bool, Owner, error) {
	own, err := s.GetByEmailAndPassword(email, password)

	if err != nil {
		if err.Error() == "Owner not found" {
			return false, Owner{}, nil
		}
	}
	return own.OwnerId != 0, own, err
}

func NewOwnerStorage(user string, password string, port string) (OwnerStorage, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s", user, password, port) //TODO поменять порт
	db, err := sqlx.Open("postgres", connStr)
	return OwnerStorage{db}, err
}

func (s *OwnerStorage) Count() (int, error) {
	res := 0
	err := s.db.Get(&res, "SELECT COUNT(ownerid) FROM owner")
	return res, err
}

func (s *OwnerStorage) Drop() error {
	_, err := s.db.Exec("DROP TABLE IF EXISTS owner CASCADE")
	return err
}

func (s *OwnerStorage) Clear() {
	s.Drop()
	s.CreateTable()
}

func hasPermission(owner Owner, cookie string) bool {
	actualOwner, err := StorageSession.GetOwnerByCookie(cookie)
	if err != nil {
		return false
	}
	return actualOwner.OwnerId == owner.OwnerId
}
