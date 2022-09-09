package util

import "encoding/json"

func ToMapInterface(v any) (map[string]interface{}, error) {

	var mi map[string]interface{}

	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &mi)
	if err != nil {
		return nil, err
	}

	return mi, nil
}
