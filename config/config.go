package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	ConfigFile = "config.json"
)

type ConfigMap map[string]string

var Config = ConfigMap{}

func init() {

	// get run path
	program, err := os.Executable()
	if err != nil {
		panic(err)
	}
	programPath := filepath.Dir(program)
	fmt.Printf("\ninit config path=%s/%s\n", programPath, ConfigFile)

	// read config
	configPath := filepath.Join(programPath, ConfigFile)
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &Config)
	if err != nil {
		log.Fatal(err)
	}
}
