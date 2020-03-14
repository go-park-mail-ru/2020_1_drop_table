package owners

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"time"
)

type Stuff struct {
	StuffID  int       `json:"id"`
	Name     string    `json:"name" validate:"required,min=4,max=100"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=100"`
	EditedAt time.Time `json:"editedAt" validate:"required"`
	Photo    string    `json:"photo"`
	IsOwner  bool      `json:"isowner"`
}

type StuffStorage struct {
	db *sqlx.DB
}

func logErr(err error, message string, where Stuff) {
	log.Error().Msgf("Error: %v, %s,  in -> %v", err, message, where)
}

func (s *StuffStorage) Append(st Stuff) (Stuff, error) {
	isExisted, _, err := s.Existed(st.Email, st.Password)
	if isExisted {
		return Stuff{}, errors.New("User already in base")
	}
	st.Password = GetMD5Hash(st.Password)
	if err != nil {
		logErr(err, "error when trying to check is existed", st)
		return Stuff{}, err
	}
	dbOwner := Stuff{}
	err = s.db.Get(&dbOwner, "insert into stuff(name, email, password, editedat, photo,isowner) values ($1,$2,$3,$4,$5,$6) returning *", st.Name, st.Email, st.Password, st.EditedAt, st.Photo, st.IsOwner)
	if err != nil {
		logErr(err, "When trying to append data", st)
		return Stuff{}, err
	}
	return dbOwner, err

}

func (s *StuffStorage) CreateTable() error {
	schema := `CREATE TABLE IF NOT EXISTS stuff
(
    StuffID  Bigserial PRIMARY KEY,
    Name     text,
    Email    text,
    Password text,
    EditedAt timestamp,
    Photo    text,
    IsOwner bool
)
`

	_, err := s.db.Exec(schema)
	return err
}

func isStuffEmpty(st *Stuff) bool {
	if st.StuffID == 0 && st.Name == "" {
		log.Info().Msgf("Stuff not found")
		return true
	}
	return false
}

func (s *StuffStorage) GetByEmailAndPassword(email string, password string) (Stuff, error) {
	st := Stuff{}
	password = GetMD5Hash(password)
	err := s.db.Get(&st, "select * from stuff where password=$1 AND email=$2", password, email)
	if isStuffEmpty(&st) {
		notFoundErrorMessage := fmt.Sprintf("Stuff not found")
		return Stuff{}, errors.New(notFoundErrorMessage)
	}
	st.EditedAt = st.EditedAt.UTC()
	return st, err
}

func (s *StuffStorage) Get(id int) (Stuff, error) {
	st := Stuff{}
	err := s.db.Get(&st, "select * from stuff where stuffid=$1", id)
	if err != nil {
		return Stuff{}, err
	}
	if isStuffEmpty(&st) {
		notFoundErrorMessage := fmt.Sprintf("Stuff not found")
		return Stuff{}, errors.New(notFoundErrorMessage)
	}
	st.EditedAt = st.EditedAt.UTC()
	return st, err
}

func (s *StuffStorage) Set(id int, newOwner Stuff) (Stuff, error) {
	newOwner.Password = GetMD5Hash(newOwner.Password)
	_, err := s.db.Exec("UPDATE stuff SET name = $1,email=$2,password=$3,editedat=$4,photo=$5 WHERE stuffid = $6", newOwner.Name, newOwner.Email, newOwner.Password, newOwner.EditedAt, newOwner.Photo, id)

	return newOwner, err
}

func (s *StuffStorage) Existed(email string, password string) (bool, Stuff, error) {
	st, err := s.GetByEmailAndPassword(email, password)

	if err != nil {
		if err.Error() == "Stuff not found" {
			return false, Stuff{}, nil
		}
	}
	isEmpty := isStuffEmpty(&st)
	return !isEmpty, st, err
}

func NewStuffStorage(user string, password string, port string) (StuffStorage, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s", user, password, port)
	db, err := sqlx.Open("postgres", connStr)
	ownStorage := StuffStorage{db}
	return ownStorage, err
}

func (s *StuffStorage) Count() (int, error) {
	res := 0
	err := s.db.Get(&res, "SELECT COUNT(stuffid) FROM stuff")
	return res, err
}

func (s *StuffStorage) Drop() error {
	_, err := s.db.Exec("DROP TABLE IF EXISTS stuff CASCADE")
	return err
}

func (s *StuffStorage) Clear() {
	_ = s.Drop()
	_ = s.CreateTable()
}

func hasPermission(stuff Stuff, cookie string) bool {
	actualOwner, err := StorageSession.GetOwnerByCookie(cookie)
	if err != nil {
		return false
	}
	return actualOwner.StuffID == stuff.StuffID
}
