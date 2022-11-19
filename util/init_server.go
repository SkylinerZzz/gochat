package util

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	Config    map[string]string // server configs
	RedisPool *redis.Pool       // redis connection pool
	RedisConn *redis.Conn       // redis connection
	DB        *gorm.DB          // mysql handler
)

func init() {
	loadConfig("../config")
	initRedis()
	initDB()
}

func loadConfig(dir string) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)

	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to load config file, err = %s", err)
	}
	err = v.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("failed to resolve config file, err = %s", err)
	}
}

func initRedis() {
	addr := fmt.Sprintf("%s:%s", Config["redis_host"], Config["redis_port"])
	// init RedisPool
	RedisPool = &redis.Pool{
		MaxIdle:     1024,
		MaxActive:   16,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			addr := fmt.Sprintf("%s:%s", Config["redis_host"], Config["redis_port"])
			return redis.Dial("tcp", addr, redis.DialPassword(Config["redis_password"]))
		},
	}
	// init RedisConn
	conn, err := redis.Dial("tcp", addr, redis.DialPassword(Config["redis_password"]))
	if err != nil {
		log.Fatalf("failed to connect to redis, err = %s", err)
	}
	RedisConn = &conn
}

func initDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", Config["mysql_username"],
		Config["mysql_password"], Config["mysql_host"], Config["mysql_port"], Config["mysql_database"])
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to mysql, err = %s", err)
	}
}
