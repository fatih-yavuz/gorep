package gorep

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

var path string

func init() {
	path = "test_files/test.tar.gz"
}

func TestMatchInitValues(t *testing.T) {
	path = "test_files/test.tar.gz"
	m := Match{
		set: make(map[string]bool),
	}
	m.set["Doe"] = true
	m.set["John"] = true
	m.initValues()
	assert.Contains(t, m.Values, "John")
	assert.Contains(t, m.Values, "Doe")
}

func TestMatchesValues(t *testing.T) {
	path = "test_files/test.tar.gz"
	var matches = Matches{}
	m := Match{
		set: make(map[string]bool),
	}
	m.set["Doe"] = true
	m.set["John"] = true

	matches = append(matches, m)

	m = Match{
		set: make(map[string]bool),
	}
	m.set["Foo"] = true
	m.set["Bar"] = true

	matches = append(matches, m)

	values := matches.Values()
	assert.Contains(t, values, "Foo")
	assert.Contains(t, values, "Bar")
	assert.Contains(t, values, "John")
	assert.Contains(t, values, "Doe")
	assert.Equal(t, len(values), 4)

}

func TestGrepTarGz(t *testing.T) {
	path = "test_files/test.tar.gz"
	matches, err := FileGrep(path, Url)
	assert.Nil(t, err)
	assert.NotEqual(t, len(matches), 0)
}

func TestGrepGz(t *testing.T) {
	path = "test_files/test.gz"
	matches, err := FileGrep(path, Url)
	assert.Nil(t, err)
	assert.NotEqual(t, len(matches), 0)

	for _, match := range matches {
		assert.Equal(t, match.Filename, path)
	}

	{
		values := []string{}
		for _, match := range matches {
			values = append(values, match.Values...)
		}
		assert.Contains(t, values, "https://upload.wikimedia.org/wikipedia/commons/4/4b/What_Is_URL.jpg")
		assert.Contains(t, values, "https://www.tubiba.com.tr/iletisim")
	}

}

func TestImageUrl(t *testing.T) {
	path = "test_files/test.tar.gz"
	matches, err := FileGrep(path, ImageUrl)
	assert.Nil(t, err)
	assert.NotEqual(t, len(matches), 0)

	values := Matches(matches).Values()

	assert.Contains(t, values, "https://upload.wikimedia.org/wikipedia/commons/4/4b/What_Is_URL.jpg")
	assert.Contains(t, values, "http://cdn.kompass.com/_ui/desktop/theme-kompass/images/ficheEntre/contactM.png")
	assert.Contains(t, values, "http://tubiba.com/assets/img/ic_link_black_1x_web_18dp.png")

	spew.Dump(values)
}
