// ParseEmailAndGo
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/DmitryBugrov/log"
	"github.com/taknb2nch/go-pop3"
)

var (
	c Cfg
)

const Config_file = "./config.json"

type Cfg struct {
	Address string
	User    string
	Pass    string
}

func main() {
	log.Init(log.LogLevelTrace, true, true, true)
	log.Print(log.LogLevelTrace, "Starting...")
	if err := c.load(); err != nil {
		log.Print(log.LogLevelError, err)
		panic(0)
	}
	fmt.Println(c.Address)
	ReciveMail()
}

func ReciveMail() {
	err := pop3.ReceiveMail(c.Address, c.User, c.Pass,
		func(number int, uid, data string, err error) (bool, error) {
			log.Print(log.LogLevelTrace, "%d, %s\n", number, uid)

			log.Print(log.LogLevelWarning, "getting mail...")

			return false, nil
		})
	if err != nil {
		log.Print(log.LogLevelError, err)
	}
}

func (c *Cfg) load() error {
	log.Print(log.LogLevelTrace, "Enter to cfg.Load")
	file, err := os.Open(Config_file)
	if err != nil {
		log.Print(log.LogLevelError, "Configuration file cannot be loaded: ", Config_file)
		return err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(&c)
	if err != nil {
		log.Print(log.LogLevelError, "Unable to decode config into struct", err.Error())
	}

	return nil
}
