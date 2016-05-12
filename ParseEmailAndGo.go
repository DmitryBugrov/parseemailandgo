// ParseEmailAndGo
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
	log.Init(log.LogLevelInfo, true, true, true)
	log.Print(log.LogLevelTrace, "Starting...")
	if err := c.load(); err != nil {
		log.Print(log.LogLevelError, err)
		panic(0)
	}

	err := ReciveMail()

	if err != nil {
		log.Print(log.LogLevelError, err)
	}

	//	data, err := ReadFromFile(TestFileName)
	//	if err == nil {
	//		CheckRegExp(data)
	//	} else {
	//		log.Print(log.LogLevelError, err)
	//	}
}

func Action(cmdName string) {
	cmdSlice := strings.Split(cmdName, " ")
	cmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	stdout, err := cmd.StdoutPipe()

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Print(log.LogLevelError, err)
	}
	fmt.Println("Waiting for command to finish...")
	//	log.Print(log.LogLevelInfo, stdout)

	err = cmd.Wait()
	if err != nil {
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
		fmt.Println("Processing rule: ", i+1)
		ruleIsTrueForBody := true
		ruleIsTrueForSubject := true
		if c.Rules[i].Body != "" {
			log.Print(log.LogLevelTrace, "looking for in body...", c.Rules[i].Body)
			re := regexp.MustCompile(c.Rules[i].Body)
			ruleIsTrueForBody = re.Match([]byte(body))

			//log.Print(log.LogLevelTrace, "result=", ruleIsTrueForBody)
		}
		if c.Rules[i].Subject != "" {
			log.Print(log.LogLevelTrace, "looking for in subject...", c.Rules[i].Subject)
			re := regexp.MustCompile(c.Rules[i].Subject)
			ruleIsTrueForSubject = re.Match([]byte(subject))

			//log.Print(log.LogLevelTrace, "result=", ruleIsTrueForSubject)
		}
		if ruleIsTrueForBody && ruleIsTrueForSubject {
			fmt.Println("running...", c.Rules[i].Action)
			Action(c.Rules[i].Action)

		}

		//		log.Print(log.LogLevelTrace, "subject=", subject, "body=", body)
	}
}

func ReciveMail() error {
	err := pop3.ReceiveMail(c.Address, c.User, c.Pass,
		func(number int, uid, data string, err error) (bool, error) {
			fmt.Println("getting mail...", number, uid)

			//	WriteToFile(TestFileName, data)
			CheckRegExp(data)
			return false, nil
		})
	return err
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
