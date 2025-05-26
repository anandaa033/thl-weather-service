package gormDB

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
)

func New(host string, port int, username, password, dbname string) (*Client, error) {

	c := Client{
		config: config{
			username: username,
			password: password,
			host:     host,
			port:     port,
		},
	}

	dbQuery := url.Values{}
	dbQuery.Add("dbname", dbname)

	dbUrl := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(username, password),
		Host:     fmt.Sprintf("%s:%d", c.config.host, c.config.port),
		RawQuery: dbQuery.Encode(),
	}

	db, err := gorm.Open(postgres.Open(dbUrl.String()), &gorm.Config{})
	if err != nil {
		// panic("failed to connect database")
		return nil, err
	}

	c.client = db
	return &c, nil
}

type Client struct {
	config config
	client *gorm.DB
}

type config struct {
	username string
	password string
	host     string
	port     int
}

func (c *Client) GetDB() *gorm.DB {
	return c.client
}
