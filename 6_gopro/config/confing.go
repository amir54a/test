package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

var Conf = new(Configuration)

type Configuration struct {
	Listen struct {
		Port string `json:"port"`
	} `json:"listen"`
	Db struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"db"`
}

func ReadConfig() {

	c := flag.String("cfg", "config.json", "configuration json file.")

	flag.Parse()

	file, err := ioutil.ReadFile(*c)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &Conf)
	if err != nil {
		panic(err)
	}

}
