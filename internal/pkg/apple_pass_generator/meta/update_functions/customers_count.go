package update_functions

func UpdateVarPassesCount(oldVal interface{}) (interface{}, error) {
	if oldVal == nil {
		oldVal = 0
	} else if oldValInt, ok := oldVal.(int); ok {
		oldVal = float64(oldValInt)
	}

	oldValInt, ok := oldVal.(float64)

	if !ok {
		return 1, nil
	}

	return int(oldValInt) + 1, nil
}
