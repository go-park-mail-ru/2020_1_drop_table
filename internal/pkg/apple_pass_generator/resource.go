package apple_pass_generator

import (
	"fmt"
	"strings"
)

type ApplePass struct {
	design    string
	files     map[string][]byte
	envValues map[string]interface{}
}

func NewApplePass(design string, files map[string][]byte, envValues map[string]interface{}) ApplePass {
	return ApplePass{
		design:    design,
		files:     files,
		envValues: envValues,
	}
}

func (p *ApplePass) AddImage(imageName string, imageData []byte) {
	p.files[imageName] = imageData
}

func (p *ApplePass) insertValues() {
	var replaceValues []string
	for key, value := range p.envValues {
		valueName := fmt.Sprintf("<<%s>>", key)
		replaceValues = append(replaceValues, valueName)
		replaceValues = append(replaceValues, fmt.Sprintf("%v", value))
	}
	r := strings.NewReplacer(replaceValues...)
	p.design = r.Replace(p.design)
}

func (p *ApplePass) addDesignToFiles() {
	p.files["pass.json"] = []byte(p.design)
}
