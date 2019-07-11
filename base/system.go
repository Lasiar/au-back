package base

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

const pathConfigFile string = "./config.json"

// Config config running applications
type Config struct {
	ConnStr string `json:"connect_string"`
	Port    string `json:"port"`
}

var (
	_config     *Config
	_onceConfig sync.Once
)

// GetConfig получение объекта конфига
func GetConfig() *Config {
	_onceConfig.Do(func() {
		_config = new(Config)
		file, err := os.Open(pathConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		_config.load(file)
	})
	return _config
}

func (c *Config) load(r io.Reader) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		log.Fatal("Read Config file: ", err)
	}

	if c.ConnStr == "" {
		log.Fatal("Can`t read connection string: ", c.ConnStr)
	}

}
