package errorMessage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var em map[string]ErrorMessage

func Init(path string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(pwd + path)
	if err != nil {
		return err
	}

	em = make(map[string]ErrorMessage)

	data := make(map[string]interface{})

	err = json.Unmarshal(raw, &data)
	if err != nil {
		return err
	}

	if err := set(data, ""); err != nil {
		return err
	}

	return nil
}

const IsMessage = "IS_MESSAGE"

func set(data map[string]interface{}, prefix string) error {
	if prefix != "" {
		prefix = strings.ToLower(prefix) + "."
	}
	for k, v := range data {
		key := fmt.Sprintf("%s%s", prefix, strings.ToLower(k))
		if k == "code" || k == "message" {
			return errors.New(IsMessage)
		} else {
			err := set(v.(map[string]interface{}), key)
			if err != nil {
				if err.Error() == IsMessage {
					raw, err := json.Marshal(v)
					if err != nil {
						return err
					}
					var value ErrorMessage
					err = json.Unmarshal(raw, &value)
					if err != nil {
						return err
					}
					em[key] = value
				}
			}
		}
	}
	return nil
}

func Get(key string) ErrorMessage {
	if value, ok := em[key]; ok {
		return value
	}
	m := make(map[string]string)
	m["en"] = "System process error, please try again."
	return ErrorMessage{
		Code:    9999,
		Message: m,
	}
}
