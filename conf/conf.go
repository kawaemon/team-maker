package conf

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/kawaemon/group-maker/g"
)

type Configuration struct {
	NameMap map[int]string `json:"names"`
	Total   int            `json:"total"`
	Women   g.Slice[int]   `json:"women"`
}

const envVarName = "GROUP_MAKER_CONF_JSON"

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

	if conf.Women.IsNil() {
		err = errors.New("field \"women\" isn't specified")
		return
	}

	return
}
