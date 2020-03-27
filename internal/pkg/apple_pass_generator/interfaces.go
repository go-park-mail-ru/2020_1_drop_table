package apple_pass_generator

import "bytes"

type Generator interface {
	CreateNewPass(pass ApplePass) (*bytes.Buffer, error)
}
