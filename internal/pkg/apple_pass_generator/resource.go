package apple_pass_generator

type applePass struct {
	design string
	files  map[string][]byte
}

func NewApplePass(design string, files map[string][]byte) applePass {
	return applePass{
		design: design,
		files:  files,
	}
}

func NewApplePassWithDesign(design string) applePass {
	return applePass{
		design: design,
		files:  map[string][]byte{},
	}
}

func (p *applePass) AddImage(imageName string, imageData []byte) {
	p.files[imageName] = imageData
}

func (p *applePass) addDesignToFiles() {
	p.files["pass.json"] = []byte(p.design)
}
