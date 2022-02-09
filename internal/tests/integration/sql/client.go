package sql

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jmoiron/sqlx"

	// init postgres driver.
	_ "github.com/lib/pq"
)

func env(name, def string) string {
	res := os.Getenv(name)
	if res == "" {
		res = def
	}
	return res
}

type Client struct {
	db *sqlx.DB
}

func (c *Client) Connect() error {
	db, err := sqlx.Connect("postgres", c.connectionString())
	if err == nil {
		c.db = db
	}
	return err
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) RunFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	query := string(content)

	_, err = c.db.Exec(query)
	return err
}

func (c *Client) connectionString() string {
	h := env("POSTGRES_HOST", "127.0.0.1")
	u := env("POSTGRES_USER", "cr_user")
	p := env("POSTGRES_PASSWORD", "password")
	db := env("POSTGRES_DB", "creatives")
	port := env("POSTGRES_PORT", "8086")

	return fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`, h, port, u, p, db)
}
