package gorep

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	Url      string = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-zA-Z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&\/\/=]*)`
	ImageUrl string = `(http(s?):)([/|.|\w|\s|-])*\.(?:jpe?g|gif|png)`
)

type Match struct {
	Filename string
	set      map[string]bool
	Values   []string
}

type Matches []Match

func (m *Match) initValues() {
	m.Values = []string{}
	for value, _ := range m.set {
		if strings.TrimSpace(value) != "" {
			m.Values = append(m.Values, value)
		}
	}
}

func (matches Matches) Values() []string {
	result := make([]string, 0)
	for _, m := range matches {
		m.initValues()
		result = append(result, m.Values...)
	}
	return result
}

func read(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ByteGrep(path string, content []byte, regex string) ([]Match, error) {
	if strings.HasSuffix(path, ".tar.gz") {
		return grepTarGz(content, regex)
	} else if strings.HasSuffix(path, ".gz") {
		return grepGz(path, content, regex)
	}
	return grep(path, content, regex)

}

func FileGrep(path string, regex string) ([]Match, error) {
	content, err := read(path)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(path, ".tar.gz") {
		return grepTarGz(content, regex)
	} else if strings.HasSuffix(path, ".gz") {
		return grepGz(path, content, regex)
	}

	return grep(path, content, regex)

}

func grepTarGz(content []byte, regex string) ([]Match, error) {

	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	var matches []Match
	buffer := bytes.NewBuffer(content)

	gzipReader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {

		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		if header.Typeflag == tar.TypeReg {
			content := make([]byte, header.Size)
			for {
				tmp := make([]byte, header.Size)
				i, err := tarReader.Read(tmp)
				content = append(content, tmp...)

				if i >= len(content) {
					break
				}

				if err != nil {
					break
				}

			}

			if err == nil {
				var match Match
				match.Filename = string(header.Name)

				match.set = make(map[string]bool)

				for _, value := range r.FindAll(content, -1) {
					match.set[string(value)] = true
				}
				if len(match.set) > 0 {
					match.initValues()
					matches = append(matches, match)
				}

			}

		}

	}

	return matches, nil
}

func grepGz(filename string, content []byte, regex string) ([]Match, error) {

	r, err := regexp.Compile(regex)

	if err != nil {
		return nil, err
	}

	var matches []Match
	buffer := bytes.NewBuffer(content)

	gzipReader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	decompressed, err := ioutil.ReadAll(gzipReader)

	if err != nil {
		return nil, err
	}

	var match Match
	match.Filename = string(filename)
	match.set = make(map[string]bool)

	for _, value := range r.FindAll(decompressed, -1) {
		match.set[string(value)] = true
	}

	if len(match.set) > 0 {
		match.initValues()
		matches = append(matches, match)
	}

	return matches, nil
}

func grep(filename string, content []byte, regex string) ([]Match, error) {

	r, err := regexp.Compile(regex)

	if err != nil {
		return nil, err
	}

	var matches []Match

	var match Match
	match.Filename = string(filename)
	match.set = make(map[string]bool)

	for _, value := range r.FindAll(content, -1) {
		match.set[string(value)] = true
	}
	if len(match.Values) > 0 {
		match.initValues()
		matches = append(matches, match)
	}

	return matches, nil
}
