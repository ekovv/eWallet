package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Config struct {
	Host  string `json:"host"`
	DB    string `json:"dsn"`
	CFile string
}

type F struct {
	host  *string
	db    *string
	cFile *string
}

var f F

const addr = ":8080"

func init() {
	f.host = flag.String("a", addr, "-a=")
	f.db = flag.String("d", "", "-d=db")
	f.cFile = flag.String("c", "", "-c=")

}

func New() (c Config) {
	flag.Parse()
	if envHost := os.Getenv("HOST"); envHost != "" {
		f.host = &envHost
	}
	if envDB := os.Getenv("DB_CONNECTION_STRING"); envDB != "" {
		f.db = &envDB
	}
	c.Host = *f.host
	c.DB = *f.db
	c.CFile = *f.cFile
	file, err := os.Open(c.CFile)
	if err != nil {
		return
	}
	defer file.Close()

	all, err := io.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(all, &c)
	if err != nil {
		return
	}
	return c

}
