package apple_pass_generator

type ApplePass struct {
	design string
	files  map[string][]byte
}

func NewApplePass(design string, files map[string][]byte) ApplePass {
	return ApplePass{
		design: design,
		files:  files,
	}
}

func (p *ApplePass) AddImage(imageName string, imageData []byte) {
	p.files[imageName] = imageData
}

func (p *ApplePass) addDesignToFiles() {
	p.files["pass.json"] = []byte(p.design)
}
