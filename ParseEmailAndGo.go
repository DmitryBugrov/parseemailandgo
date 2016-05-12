// ParseEmailAndGo
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/DmitryBugrov/log"
	"github.com/taknb2nch/go-pop3"
)

var (
	c Cfg
)

const Config_file = "./config.json"
const TestFileName = "./mail.dat"

type Cfg struct {
	Address string
	User    string
	Pass    string
	Rules   []Rule
}
type Rule struct {
	Action  string
	Subject string
	Body    string
}

func main() {
	log.Init(log.LogLevelTrace, true, true, true)
	log.Print(log.LogLevelTrace, "Starting...")
	if err := c.load(); err != nil {
		log.Print(log.LogLevelError, err)
		panic(0)
	}
	fmt.Println(c.Rules[1].Action)
	//	ReciveMail()
	data, err := ReadFromFile(TestFileName)
	if err == nil {
		CheckRegExp(data)
	} else {
		log.Print(log.LogLevelError, err)
	}
}

func WriteToFile(filename string, data string) {
	ioutil.WriteFile(filename, []byte(data), 0644)
}

func ReadFromFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	return string(data), err
}

func ParsingMail(data string) (string, string) {
	subject := "-"
	body := "--"
	re := regexp.MustCompile("Subject:.*")
	subject = string(re.Find([]byte(data)))

	re = regexp.MustCompile("(?ms)Content-Type:.*$")
	body = re.FindString(data)
	return subject, body
}

func CheckRegExp(data string) {
	subject, body := ParsingMail(data)
	for i := 0; i < len(c.Rules); i++ {
		log.Print(log.LogLevelTrace, "Processing rule: ", i+1)
		ruleIsTrue := false
		if c.Rules[i].Body != "" {
			log.Print(log.LogLevelTrace, "looking for...", c.Rules[i].Body)
			re := regexp.MustCompile(c.Rules[i].Body)
			ruleIsTrue = re.Match([]byte(body))
			log.Print(log.LogLevelTrace, "result=", ruleIsTrue)
		}
		if c.Rules[i].Subject != "" {
			log.Print(log.LogLevelTrace, "looking for...", c.Rules[i].Subject)
			re := regexp.MustCompile(c.Rules[i].Subject)
			ruleIsTrue = ruleIsTrue && re.Match([]byte(subject))
			log.Print(log.LogLevelTrace, "result=", ruleIsTrue)
		}
		if ruleIsTrue {

			log.Print(log.LogLevelTrace, "running...", c.Rules[i].Action)
		}

		//		log.Print(log.LogLevelTrace, "subject=", subject, "body=", body)
	}
}

func ReciveMail() {
	err := pop3.ReceiveMail(c.Address, c.User, c.Pass,
		func(number int, uid, data string, err error) (bool, error) {
			log.Print(log.LogLevelTrace, number, uid)

			log.Print(log.LogLevelTrace, "getting mail...")
			WriteToFile(TestFileName, data)
			CheckRegExp(data)
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
