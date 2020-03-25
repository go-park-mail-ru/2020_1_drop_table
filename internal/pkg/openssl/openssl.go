package openssl

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Smime(stdin []byte, args ...string) ([]byte, error) {
	cmd := exec.Command("openssl", args...)

	in := bytes.NewReader(stdin)
	out := &bytes.Buffer{}
	errs := &bytes.Buffer{}

	cmd.Stdin, cmd.Stdout, cmd.Stderr = in, out, errs

	if err := cmd.Run(); err != nil {
		if len(errs.Bytes()) > 0 {
			return nil, fmt.Errorf("error running %s (%s):\n %v", cmd.Args, err, errs.String())
		}
		return nil, err
	}

	return out.Bytes(), nil
}
