package configuration

import (
	"fmt"
	"time"
	data_yugabyte "yugabyte-crud/data"

	_ "github.com/lib/pq"

	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/database"
	"10.96.24.141/UDTN/integration/microservices/mcs-go/mcs-go-modules/mcs-go-core.git/database/sqlx"
)

type YugabyteConfig struct {
	Host                  string `config:"DB_ESBORACLE_HOST"`
	Port                  int    `config:"DB_ESBORACLE_PORT"`
	Username              string `config:"DB_ESBORACLE_USER"`
	Password              string `config:"DB_ESBORACLE_PASSWORD"`
	Database              string `config:"DB_ESBORACLE_DBNAME"`
	SslMode               string `config:""`
	IdleConnection        int    `config:"DB_ESBORACLE_POOL_IDLE_CONNECTION"`
	MaxConnection         int    `config:"DB_ESBORACLE_MAX_POOL_SIZE"`
	MaxLifeIdleConnection int    `config:"DB_ESBORACLE_IDLE_TIMEOUT"`  //seconds
	MaxIdleTimeConnection int    `config:"DB_ESBORACLE_MAX_LIFE_TIME"` // seconds
}

func GetYugabyteConfig() *YugabyteConfig {
	username := "postgres"
	password := ""
	host := "10.106.68.9"
	port := 5433
	sslmode := "disable"
	database := "napas"
	idleConnection := 5
	maxConnection := 5
	maxLifeTimeConnection := 30
	maxLifeIdleConnection := 30

	return &YugabyteConfig{
		Username:              username,
		Password:              password,
		Host:                  host,
		Port:                  port,
		Database:              database,
		SslMode:               sslmode,
		IdleConnection:        idleConnection,
		MaxConnection:         maxConnection,
		MaxLifeIdleConnection: maxLifeIdleConnection,
		MaxIdleTimeConnection: maxLifeTimeConnection,
	}
}

func GetYugabyteDatabase(e *YugabyteConfig) *data_yugabyte.YugabyteDatabase {

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?%s", e.Username, e.Password, e.Host, e.Port, e.Database, fmt.Sprintf("sslmode=%s", "disable"))

	fmt.Printf("DNS: %s", dsn)
	db, err := sqlx.NewSqlxGdbc("postgres", dsn,
		database.WithMaxIdleCount(e.MaxLifeIdleConnection),
		database.WithMaxOpen(e.MaxConnection),
		database.WithMaxIdleTime(time.Duration(e.MaxIdleTimeConnection)),
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &data_yugabyte.YugabyteDatabase{
		DB: db,
	}
}
