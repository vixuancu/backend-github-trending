package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq" // Postgres driver
)

type Sql struct {
	Db       *sqlx.DB
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// kết nối đến cơ sở dữ liệu postgres
func (s *Sql) Connect() error {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.Username, s.Password, s.Database)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	s.Db = db
	fmt.Println("Connected to database successfully")
	return nil
}
func (s *Sql) Close() {
	s.Db.Close()
}
