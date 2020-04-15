package update_functions

func UpdateVarPassesCount(oldVal interface{}) (interface{}, error) {
	oldValInt, ok := oldVal.(float64)

	if !ok {
		return 1, nil
	}

	return int(oldValInt) + 1, nil
}
