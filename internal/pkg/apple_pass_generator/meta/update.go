package meta

import (
	"errors"
	"fmt"
)

type Meta struct{}

func (m *Meta) UpdateMeta(oldValues map[string]interface{}) (map[string]interface{}, error) {
	newValues := map[string]interface{}{}
	for key, element := range oldValues {
		f, ok := onUpdateFunctions[key]
		if !ok {
			message := fmt.Sprintf("not found update func for var <<%s>>", key)
			return nil, errors.New(message)
		}

		newVal, err := f(element)
		if err != nil {
			return nil, err
		}
		newValues[key] = newVal
	}

	return newValues, nil
}
