package conf

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

var errYamlImportDir = errors.New("yaml: can't import directory")

const yamlImportKey = "#import "

func ParseYaml(file string, out interface{}) error {
	in := readYaml(file)
	return yaml.Unmarshal(in, out)
}

func readYaml(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	if fi.IsDir() {
		panic(errYamlImportDir)
	}

	in := make([]byte, 0, fi.Size())
	for sc := bufio.NewScanner(bufio.NewReader(f)); sc.Scan(); {
		b := sc.Bytes()
		s := string(b)
		if strings.HasPrefix(s, yamlImportKey) {
			_file := strings.TrimSpace(s[len(yamlImportKey):])
			if _file == "" {
				continue
			}
			if !filepath.IsAbs(_file) {
				_file = filepath.Join(filepath.Dir(file), _file)
			}
			in = append(in, readYaml(_file)...)
		} else {
			in = append(in, b...)
			in = append(in, '\n')
		}
	}

	return in
}
