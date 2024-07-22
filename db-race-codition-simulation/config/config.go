package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type config struct {
	MysqlHost                  string `envconfig:"MYSQL_HOST" default:"localhost" prompt:"Enter mysql host"`
	MysqlPort                  int    `envconfig:"MYSQL_PORT" default:"3306" prompt:"Enter mysql port"`
	MysqlDBName                string `envconfig:"MYSQL_DB_NAME" default:"test_race" prompt:"Enter database name"`
	MysqlUsername              string `envconfig:"MYSQL_USERNAME" default:"" prompt:"Enter mysql username"`
	MysqlPassword              string `envconfig:"MYSQL_PASSWORD" default:"" prompt:"Enter mysql password" secret:"true"`
	MysqlLogMode               int    `envconfig:"MYSQL_LOG_MODE" default:"1" prompt:"Enter gorm log mode 1-4"`
	MysqlParseTime             bool   `envconfig:"MYSQL_PARSE_TIME" default:"true" prompt:"Parse mysql time to local"`
	MysqlCharset               string `envconfig:"MYSQL_CHARSET" default:"utf8mb4" prompt:"Enter mysql database charset"`
	MysqlLoc                   string `envconfig:"MYSQL_LOC" default:"Local" prompt:"Enter mysql local time"`
	MysqlMaxLifetimeConnection int    `envconfig:"MYSQL_MAX_LIFETIME_CONNECTION" default:"10" prompt:"Enter mysql maximum amount of time a connection may be reused, in minute"`
	MysqlMaxOpenConnection     int    `envconfig:"MYSQL_MAX_OPEN_CONNECTION" default:"50" prompt:"Enter mysql maximum number of open connections to the database"`
	MysqlMaxIdleConnection     int    `envconfig:"MYSQL_MAX_IDLE_CONNECTION" default:"10" prompt:"Enter mysql maximum number of connections in the idle connection pool"`

	RedisHost     string `envconfig:"REDIS_HOST" default:"localhost" prompt:"Enter redis host"`
	RedisPort     int    `envconfig:"REDIS_PORT" default:"6379" prompt:"Enter redis port"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:"" prompt:"Enter redis password" secret:"true"`
}

func InitDB() (*gorm.DB, *redis.Client) {
	err := godotenv.Overload()
	if err != nil {
		if err.Error() == "open .env: no such file or directory" {
			GenerateDotEnv()
			return InitDB()
		}
		log.Fatal(err)
	}
	var cfg config
	envconfig.MustProcess("", &cfg)
	return cfg.ConnectDB(), cfg.ConnectRedis()
}

// connect to db
func (cfg *config) ConnectDB() *gorm.DB {
	// construct connection string
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%+v&loc=%s",
		cfg.MysqlUsername,
		cfg.MysqlPassword,
		cfg.MysqlHost,
		cfg.MysqlPort,
		cfg.MysqlDBName,
		cfg.MysqlCharset,
		cfg.MysqlParseTime,
		cfg.MysqlLoc)

	// open mysql connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(cfg.MysqlLogMode)),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set configuration pooling connection
	mysqlDb, _ := db.DB()
	mysqlDb.SetMaxOpenConns(cfg.MysqlMaxOpenConnection)
	mysqlDb.SetConnMaxLifetime(time.Duration(cfg.MysqlMaxLifetimeConnection) * time.Minute)
	mysqlDb.SetMaxIdleConns(cfg.MysqlMaxIdleConnection)

	return db
}

func (cfg *config) ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
	})
}
