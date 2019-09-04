package json_test

import (
	"bytes"
	stdJSON "encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/akyoto/assert"
	"github.com/akyoto/json"
	jsoniter "github.com/json-iterator/go"
)

var (
	jsonBytes  []byte
	jsonReader io.ReadSeeker
)

func init() {
	jsonBytes, _ = ioutil.ReadFile("testdata/movie-numbers.json")
	jsonReader = bytes.NewReader(jsonBytes)
}

type Movie struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
	Duration int    `json:"duration"`
	Budget   int    `json:"budget"`
}

func TestDecodeStrings(t *testing.T) {
	movie, err := load("testdata/movie-strings.json")
	assert.Nil(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, movie.Title, "The Last Samurai")
	assert.Equal(t, movie.Director, "Edward Zwick")
}

func TestDecodeNumbers(t *testing.T) {
	movie, err := load("testdata/movie-numbers.json")
	assert.Nil(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, movie.Title, "The Last Samurai")
	assert.Equal(t, movie.Director, "Edward Zwick")
	assert.Equal(t, movie.Year, 2003)
	assert.Equal(t, movie.Duration, 160)
	assert.Equal(t, movie.Budget, 140000000)
}

// load loads a single JSON file as movie data.
func load(path string) (*Movie, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	movie := &Movie{}
	err = json.NewDecoder(file).Decode(movie)

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func BenchmarkDecode(b *testing.B) {
	movie := &Movie{}
	decoder := json.NewDecoder(jsonReader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		jsonReader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeStd(b *testing.B) {
	movie := &Movie{}
	decoder := stdJSON.NewDecoder(jsonReader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		jsonReader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeJsoniter(b *testing.B) {
	movie := &Movie{}
	decoder := jsoniter.NewDecoder(jsonReader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		jsonReader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkFullDecode(b *testing.B) {
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		jsonReader.Seek(0, io.SeekStart) // nolint:errcheck
		decoder := json.NewDecoder(jsonReader)
		err := decoder.Decode(movie)
		decoder.Close()

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkFullDecodeStd(b *testing.B) {
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		jsonReader.Seek(0, io.SeekStart) // nolint:errcheck
		err := stdJSON.NewDecoder(jsonReader).Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkFullDecodeJsoniter(b *testing.B) {
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		jsonReader.Seek(0, io.SeekStart) // nolint:errcheck
		err := jsoniter.NewDecoder(jsonReader).Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBytes, movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkUnmarshalStd(b *testing.B) {
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := stdJSON.Unmarshal(jsonBytes, movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkUnmarshalJsoniter(b *testing.B) {
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(jsonBytes, movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}
