package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Config struct {
	Host  string `json:"host"`
	CFile string
}

type F struct {
	host  *string
	cFile *string
}

var f F

const addr = "localhost:8080"

func init() {
	f.host = flag.String("a", addr, "-a=")
	f.cFile = flag.String("c", "", "config file")

}

func New() (c Config) {
	flag.Parse()
	if envHost := os.Getenv("HOST"); envHost != "" {
		f.host = &envHost
	}
	c.Host = *f.host
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
