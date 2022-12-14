package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var OwnerID string

func init() {

	//	if ./config.json exists, then:
	//	else if ./config.json does not exist, then:
	//	else if ./config.json is Schrödingers pet, then:
	if _, err := os.Stat("./config.json"); err == nil {
		// log to console: config file found
		log.Println("Configuration file found, loading configuration...")

	} else if errors.Is(err, os.ErrNotExist) {
		//	create a map with contents to include in the JSON file.
		config := map[string]interface{}{
			"token": "tokenhere",
			"owner": "owneridhere",
		}
		//	Pretty JSON
		data, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		// 	create config.json, and write data to it
		_ = ioutil.WriteFile("config.json", data, 0644)
		log.Fatal("Configuration file not found! Creating new one. Please insert bot token into file before restarting")
		return
	} else {
		log.Fatal(err)
	}

	jsonFile, err := jsonParse()

	if err != nil {
		log.Fatal(err)
		return
	}

	disToken := jsonFile["token"].(string)
	OwnerID = jsonFile["owner"].(string)

	InitDB()
	initDiscord(disToken)

}
func main() {

	fmt.Println("selfbot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	s.Close()

}

func jsonParse() (map[string]interface{}, error) {
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result, err
}
