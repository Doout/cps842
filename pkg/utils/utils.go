package utils

import (
	"encoding/json"
	"io/ioutil"
)

func LoadJsonFromFile(t interface{}, file string) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(dat, t); err != nil {
		panic(err)
	}
}
