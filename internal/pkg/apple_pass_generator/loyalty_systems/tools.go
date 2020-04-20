package loyalty_systems

import "encoding/json"

func UnmarshalEmptyString(data []byte, v interface{}) error {
	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, v)
}
