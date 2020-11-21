package conf

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Configuration struct {
	NameMap      map[int]string `json:"names"`
	Total        int            `json:"total"`
	WomanNumbers []int          `json:"womanNumbers"`
}

const envVarName = "TEAM_MAKER_CONF_JSON"

func FromEnv() (conf Configuration, err error) {
	data := os.Getenv(envVarName)

	if data == "" {
		err = fmt.Errorf("config environment varibale isn't set, Please set it to %s", envVarName)
		return
	}

	dec, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		err = fmt.Errorf("couldn't decode config environment(%s) variable from base64: %w", envVarName, err)
		return
	}

	err = json.Unmarshal(dec, &conf)

	if err != nil {
		err = fmt.Errorf("couldn't decode config environment(%s) variable from json: %w", envVarName, err)
		return
	}

	if conf.NameMap == nil {
		err = errors.New("field \"names\" isn't specified")
		return
	}

	if conf.WomanNumbers == nil {
		err = errors.New("field \"womanNumbers\" isn't specified")
		return
	}

    return
}
