package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	// Adress that the webserver will run under
	APIAdress string

	// SQL Information
	// Default SQLProtocol is mysql
	SQLUser     string
	SQLPassword string
	SQLIP       string
	SQLPort     string
	SQLDatabase string

	// How often the server will check clients, in seconds
	Interval int
}

func DefaultConfig() *Config {
	return &Config{
		APIAdress:   "127.0.0.1",
		SQLUser:     "root",
		SQLPassword: "",
		SQLIP:       "127.0.0.1",
		SQLPort:     "3306",
		SQLDatabase: "MSTT-Monitor",
		Interval:    1,
	}
}

func NewConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	conf := DefaultConfig()
	if err != nil {
		_, err = os.Create(file)
		if err != nil {
			return nil, err
		}
		err = SaveFile(file, conf)
		return conf, err
	}
	err = ReadFile(f, &conf)
	return conf, err
}

func ReadFile(file *os.File, v interface{}) error {
	s, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(s, &v)
}

func SaveFile(file string, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(file, b, 0777)
	return nil
}
