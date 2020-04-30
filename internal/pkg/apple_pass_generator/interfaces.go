package apple_pass_generator

import "bytes"

type Generator interface {
	CreateNewPass(pass ApplePass) (*bytes.Buffer, error)
}

type PassMeta interface {
	UpdateMeta(oldValues map[string]interface{}) (map[string]interface{}, error)
}
