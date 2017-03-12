package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/hjson/hjson-go"
)

type Config struct {
	APIAdress string

	SQLProtocol string
	SQLUser     string
	SQLPassword string
	SQLIP       string
	SQLPort     string
	SQLDatabase string

	Interval float64 // Tydligen kräver hjson att nummer alltid är float64.
}

func (c *Config) create(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	f.Write([]byte(`{
	// Adress that the webserver will run under
	APIAdress: ":8080"

	// SQL Information
	// Default SQLProtocol is mysql
	SQLProtocol: "mysql"
	SQLUser: "example"
	SQLPassword: ""
	SQLIP: "localhost"
	SQLPort: "3306"
	SQLDatabase: "example-database"

 	// How often the server will check clients, in seconds
	Interval: 1
}`))
	f.Close()
	return nil
}

func (c *Config) setField(name string, value interface{}) error {
	structValue := reflect.ValueOf(c).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in config", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("canno set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func (c *Config) fillFields(m map[string]interface{}) error {
	for k, v := range m {
		err := c.setField(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) Load() {
	f := "config.hjson"

	b, err := ioutil.ReadFile(f)
	if err != nil {
		err = c.create(f)
		if err != nil {
			log.Panic(err, "Can't create config")
		}

		b, err = ioutil.ReadFile(f)
		if err != nil {
			log.Panic(err, "Can't open config")
		}
	}

	var m map[string]interface{}
	err = hjson.Unmarshal(b, &m)
	if err != nil {
		log.Panic(err, "Can't parse config")
	}

	err = c.fillFields(m)
	if err != nil {
		log.Panic(err, "Can't parse config map")
	}
}
