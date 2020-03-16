package owners

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"time"
)

type Staff struct {
	StaffID  int       `json:"id"`
	Name     string    `json:"name" validate:"required,min=4,max=100"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=100"`
	EditedAt time.Time `json:"editedAt" validate:"required"`
	Photo    string    `json:"photo"`
	IsOwner  bool      `json:"isOwner"`
}

type StaffStorage struct {
	db *sqlx.DB
}

func logErr(err error, message string, where Staff) {
	log.Error().Msgf("error: %v, %s,  in -> %v", err, message, where)
}

func (s *StaffStorage) Append(st Staff) (Staff, error) {
	_, found, err := s.GetByEmailAndPassword(st.Email, st.Password)
	if found {
		return Staff{}, errors.New("user already in base")
	}
	if err != nil {
		logErr(err, "error when trying to check is existed", st)
		return Staff{}, err
	}
	st.Password = GetMD5Hash(st.Password)
	dbOwner := Staff{}
	err = s.db.Get(&dbOwner, "insert into staff(name, email, password, editedat, photo,isowner) values ($1,$2,$3,$4,$5,$6) returning *", st.Name, st.Email, st.Password, st.EditedAt, st.Photo, st.IsOwner)
	if err != nil {
		logErr(err, "when trying to append data", st)
		return Staff{}, err
	}
	return dbOwner, err

}

func (s *StaffStorage) CreateTable() error {
	schema := `CREATE TABLE IF NOT EXISTS Staff
(
    StaffID  Bigserial PRIMARY KEY,
    Name     text,
    Email    text,
    Password text,
    EditedAt timestamp,
    Photo    text,
    IsOwner  boolean
)
`

	_, err := s.db.Exec(schema)
	return err
}

func isStaffEmpty(own *Staff) bool {
	if own.StaffID == 0 && own.Name == "" {
		log.Info().Msgf("staff not found")
		return true
	}
	return false
}

func (s *StaffStorage) GetByEmailAndPassword(email string, password string) (staff Staff, found bool, e error) {
	own := Staff{}
	password = GetMD5Hash(password)
	err := s.db.Get(&own, "select * from Staff where password=$1 AND email=$2", password, email)
	if isStaffEmpty(&own) {
		return Staff{}, false, nil
	}
	own.EditedAt = own.EditedAt.UTC()
	return own, true, err
}

func (s *StaffStorage) Get(id int) (Staff, error) {
	sta := Staff{}
	err := s.db.Get(&sta, "select * from Staff where StaffID=$1", id)
	if err != nil {
		return Staff{}, err
	}
	if isStaffEmpty(&sta) {
		notFoundErrorMessage := fmt.Sprintf("staff not found")
		return Staff{}, errors.New(notFoundErrorMessage)
	}
	sta.EditedAt = sta.EditedAt.UTC()
	return sta, err
}

func (s *StaffStorage) Set(id int, newStaff Staff) (Staff, error) {
	newStaff.Password = GetMD5Hash(newStaff.Password)
	_, err := s.db.Exec("UPDATE Staff SET name=$1,email=$2,password=$3,editedat=$4,photo=$5 WHERE staffid = $6", newStaff.Name, newStaff.Email, newStaff.Password, newStaff.EditedAt, newStaff.Photo, id)

	return newStaff, err
}

func NewStaffStorage(user string, password string, port string) (StaffStorage, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s", user, password, port)
	db, err := sqlx.Open("postgres", connStr)
	ownStorage := StaffStorage{db}
	return ownStorage, err
}

func (s *StaffStorage) Count() (int, error) {
	res := 0
	err := s.db.Get(&res, "SELECT COUNT(StaffID) FROM Staff")
	return res, err
}

func (s *StaffStorage) Drop() error {
	_, err := s.db.Exec("DROP TABLE IF EXISTS Staff CASCADE")
	return err
}

func (s *StaffStorage) Clear() {
	_ = s.Drop()
	_ = s.CreateTable()
}
