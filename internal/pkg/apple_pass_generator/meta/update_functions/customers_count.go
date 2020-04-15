package update_functions

func UpdateVarPassesCount(oldVal interface{}) (interface{}, error) {
	oldValInt, ok := oldVal.(int)

	if !ok {
		return -1, ErrNotInt
	}

	return oldValInt + 1, nil
}
