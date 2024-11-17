package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Hackaton-UDEVS/auth/internal/config"
	logger "github.com/Hackaton-UDEVS/auth/internal/logger"
	"github.com/Hackaton-UDEVS/auth/internal/storage"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

type Storage struct {
	Db    *sql.DB
	Rd    *redis.Client
	Useri storage.UserStorageI
}

func ConnectPostgres() (*Storage, error) {
	logs, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	conf := config.Load()
	dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.DBHOST, conf.DBPORT, conf.DBUSER, conf.DBPASSWORD, conf.DBNAME)
	fmt.Println(dns)
	db, err := sql.Open("postgres", dns)
	if err != nil {
		logs.Error("Error while connecting postgres")
	}
	//err = db.Ping()
	//if err != nil {
	//	logs.Error("Error while pinging postgres")
	//}
	logs.Info("Successfully connected to postgres")

	rd := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	user := NewUserRepo(db, rd)
	return &Storage{
		Db:    db,
		Rd:    rd,
		Useri: user,
	}, nil
}

func (stg *Storage) User() *storage.UserStorageI {
	if stg.Useri == nil {
		stg.Useri = NewUserRepo(stg.Db, stg.Rd)
	}
	return &stg.Useri
}
