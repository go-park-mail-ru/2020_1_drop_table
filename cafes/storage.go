package cafes

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"time"
)

type Cafe struct {
	CafeID      int       `json:"id"`
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Address     string    `json:"address" validate:"required"`
	Description string    `json:"description" validate:"required"`
	StaffID     int       `json:"staffID"`
	OpenTime    time.Time `json:"openTime"`
	CloseTime   time.Time `json:"closeTime"`
	Photo       string    `json:"photo"`
}

func (c *Cafe) hasPermission(StaffID int) bool {
	return c.StaffID == StaffID
}

type cafesStorage struct {
	db *sqlx.DB
}

func newCafesStorage(user string, password string, port string) (cafesStorage, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable port=%s", user, password, port)
	db, err := sqlx.Open("postgres", connStr)
	cafeStorage := cafesStorage{db}

	return cafeStorage, err
}

func (cs *cafesStorage) createTable() error {
	schema := `CREATE TABLE IF NOT EXISTS Cafe
(
    CafeID      SERIAL PRIMARY KEY,
	Name        TEXT,
	Address     TEXT,
	Description TEXT,
	StaffID     INT,
	OpenTime    TIME,
	CloseTime   TIME,
	Photo       TEXT
)
`
	_, err := cs.db.Exec(schema)
	return err
}

func (cs *cafesStorage) Append(value Cafe) (Cafe, error) {
	queryString := `insert into Cafe(
	Name, 
	Address, 
	Description, 
	StaffID, 
	OpenTime, 
	CloseTime, 
	Photo) 
	values ($1,$2,$3,$4,$5,$6,$7) 
	returning *`

	CafeDB := Cafe{}
	err := cs.db.Get(&CafeDB, queryString, value.Name, value.Address,
		value.Description, value.StaffID, value.OpenTime,
		value.CloseTime, value.Photo)
	if err != nil {
		log.Error().Msgf("error: %v, while adding cafe,  in -> %v", err, value)
		return Cafe{}, err
	}
	return CafeDB, nil
}

func (cs *cafesStorage) Get(index int) (Cafe, error) {
	queryString := `select * from Cafe where CafeID=$1`
	CafeDB := Cafe{}
	err := cs.db.Get(&CafeDB, queryString, index)

	if err != nil {
		log.Error().Msgf("error: %v, while getting cafe with index %d", err, index)
		return Cafe{}, err
	}
	return CafeDB, err
}

func (cs *cafesStorage) getOwnerCafes(staffID int) ([]Cafe, error) {
	queryString := `SELECT * FROM Cafe WHERE StaffID=$1`
	var cafes []Cafe
	err := cs.db.Select(&cafes, queryString, staffID)

	if err != nil {
		log.Error().Msgf("error: %v, while getting owner cafe, staffID: %v", err, staffID)
		return []Cafe{}, err
	}

	return cafes, nil
}

func (cs *cafesStorage) Set(i int, value Cafe) (Cafe, error) {
	queryString := `UPDATE Cafe SET 
	Name=$1, 
	Address=$2, 
	Description=$3, 
	StaffID=$4, 
	OpenTime=$5, 
	CloseTime=$6, 
	Photo=$7 
	WHERE CafeID=$8
	RETURNING *`
	cafeDB := Cafe{}
	err := cs.db.Get(&cafeDB, queryString, value.Name, value.Address, value.Description,
		value.StaffID, value.OpenTime, value.CloseTime, value.Photo, i)
	if err != nil {
		log.Error().Msgf("error: %v, while dding cafe,  in -> %v with index %d", err, value, i)
	}
	return cafeDB, err
}

func (cs *cafesStorage) Drop() error {
	_, err := cs.db.Exec("DROP TABLE IF EXISTS Cafe CASCADE")
	return err
}

func (cs *cafesStorage) Clear() {
	_ = cs.Drop()
	_ = cs.createTable()
}

var Storage, _ = newCafesStorage("postgres", "", "5431")
