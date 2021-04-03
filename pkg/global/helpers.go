package global

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func DecodeConfigFile(jsonConfigFilePath string) ([]map[string]interface{}, error) {
	jsonFile, err := ioutil.ReadFile(jsonConfigFilePath)
	if err != nil {
		return []map[string]interface{}{}, fmt.Errorf("error occurred while reading the user-provided JSON config: %+v", err)
	}

	jsonConfig := string(jsonFile)
	var inputConfigs []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonConfig), &inputConfigs); err != nil {
		return []map[string]interface{}{}, fmt.Errorf("error occurred while decoding the user-provided JSON config: %+v", err)
	}
	return inputConfigs, nil
}
