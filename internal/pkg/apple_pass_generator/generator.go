package apple_pass_generator

import (
	"2020_1_drop_table/internal/pkg/hasher"
	"2020_1_drop_table/internal/pkg/openssl"
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
)

type generator struct {
	wwdr        string
	certificate string
	key         string
	password    string
}

func NewGenerator(wwdr, certificate, key, password string) generator {
	return generator{
		wwdr:        wwdr,
		certificate: certificate,
		key:         key,
		password:    password,
	}
}

func (g *generator) createManifest(pass applePass) ([]byte, error) {
	files := map[string]string{}
	files["pass.json"] = hasher.GetSha1([]byte(pass.design))
	for key, data := range pass.files {
		files[key] = hasher.GetSha1(data)
	}

	return json.Marshal(files)
}

func (g *generator) createSignature(manifest []byte) ([]byte, error) {
	smimeCMD := []string{
		"smime",
		"-binary",
		"-sign",
		"-certfile",
		g.wwdr,
		"-signer",
		g.certificate,
		"-inkey",
		g.key,
		"-outform",
		"DER",
		"-passin",
		fmt.Sprintf("pass:%s", g.password)}

	signature, err := openssl.Smime(manifest, smimeCMD...)
	if err != nil {
		return []byte{}, err
	}

	return signature, nil
}

func (g *generator) createZip(files map[string][]byte) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipPass := zip.NewWriter(buf)

	for key, data := range files {
		f, err := zipPass.Create(key)
		if err != nil {
			return nil, err
		}
		_, err = f.Write(data)
		if err != nil {
			return nil, err
		}
	}

	err := zipPass.Close()
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (g *generator) CreateNewPass(pass applePass) (*bytes.Buffer, error) {
	manifest, err := g.createManifest(pass)
	if err != nil {
		return nil, err
	}

	signature, err := g.createSignature(manifest)
	if err != nil {
		return nil, err
	}

	pass.addDesignToFiles()

	pass.files["signature"] = signature
	pass.files["manifest.json"] = manifest

	return g.createZip(pass.files)
}
