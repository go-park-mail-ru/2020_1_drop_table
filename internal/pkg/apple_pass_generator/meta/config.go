package meta

import (
	"2020_1_drop_table/internal/pkg/apple_pass_generator/meta/update_functions"
)

type updateFunc func(oldVal interface{}) (interface{}, error)

var onUpdateFunctions map[string]updateFunc

var EmptyMeta map[string]interface{}

func init() {
	onUpdateFunctions = map[string]updateFunc{}

	onUpdateFunctions["PassesCount"] = update_functions.UpdateVarPassesCount

	EmptyMeta = map[string]interface{}{}
	for key := range onUpdateFunctions {
		EmptyMeta[key] = ""
	}
}
