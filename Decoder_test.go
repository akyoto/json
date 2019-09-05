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

type Movie struct {
	Title       string            `json:"title"`
	Director    string            `json:"director"`
	Plot        string            `json:"plot"`
	Year        int               `json:"year"`
	Duration    int               `json:"duration"`
	Budget      int               `json:"budget"`
	Rating      float64           `json:"rating"`
	AspectRatio float64           `json:"aspectRatio"`
	Profit      float64           `json:"profit"`
	Starring    []string          `json:"starring"`
	Producers   []string          `json:"producers"`
	Mappings    map[string]string `json:"mappings"`
}

func TestDecodeStrings(t *testing.T) {
	movie, err := load("testdata/movie-strings.json")
	assert.Nil(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, movie.Title, "The Last Samurai")
	assert.Equal(t, movie.Director, "Edward Zwick")
	assert.Equal(t, len(movie.Plot), 682)
}

func TestDecodeIntegers(t *testing.T) {
	movie, err := load("testdata/movie-integers.json")
	assert.Nil(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, movie.Year, 2003)
	assert.Equal(t, movie.Duration, 160)
	assert.Equal(t, movie.Budget, 140000000)
}

func TestDecodeFloats(t *testing.T) {
	movie, err := load("testdata/movie-floats.json")
	assert.Nil(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, movie.Rating, 7.7)
	assert.Equal(t, movie.AspectRatio, 2.35)
	assert.Equal(t, movie.Profit, 454.627263)
}

func TestDecodeArrays(t *testing.T) {
	movie, err := load("testdata/movie-arrays.json")
	assert.Nil(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, len(movie.Starring), 8)
	assert.Equal(t, movie.Starring[0], "Tom Cruise")
	assert.Equal(t, movie.Starring[1], "Timothy Spall")
	assert.Equal(t, movie.Starring[2], "Ken Watanabe")
	assert.Equal(t, movie.Starring[3], "Billy Connolly")
	assert.Equal(t, movie.Starring[4], "Tony Goldwyn")
	assert.Equal(t, movie.Starring[5], "Hiroyuki Sanada")
	assert.Equal(t, movie.Starring[6], "Koyuki")
	assert.Equal(t, movie.Starring[7], "Shin Koyamada")
	assert.Equal(t, movie.Producers[0], "Marshall Herskovitz")
	assert.Equal(t, movie.Producers[1], "Edward Zwick")
	assert.Equal(t, movie.Producers[2], "Tom Cruise")
	assert.Equal(t, movie.Producers[3], "Paula Wagner")
	assert.Equal(t, movie.Producers[4], "Scott Kroopf")
	assert.Equal(t, movie.Producers[5], "Tom Engelman")
}

func TestDecodeAll(t *testing.T) {
	movie, err := load("testdata/movie.json")
	assert.Nil(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, movie.Title, "The Last Samurai")
	assert.Equal(t, movie.Director, "Edward Zwick")
	assert.Equal(t, len(movie.Plot), 682)
	assert.Equal(t, movie.Year, 2003)
	assert.Equal(t, movie.Duration, 160)
	assert.Equal(t, movie.Budget, 140000000)
	assert.Equal(t, movie.Rating, 7.7)
	assert.Equal(t, movie.AspectRatio, 2.35)
	assert.Equal(t, movie.Profit, 454.627263)
	assert.Equal(t, len(movie.Starring), 8)
	assert.Equal(t, movie.Starring[0], "Tom Cruise")
	assert.Equal(t, movie.Starring[1], "Timothy Spall")
	assert.Equal(t, movie.Starring[2], "Ken Watanabe")
	assert.Equal(t, movie.Starring[3], "Billy Connolly")
	assert.Equal(t, movie.Starring[4], "Tony Goldwyn")
	assert.Equal(t, movie.Starring[5], "Hiroyuki Sanada")
	assert.Equal(t, movie.Starring[6], "Koyuki")
	assert.Equal(t, movie.Starring[7], "Shin Koyamada")
	assert.Equal(t, movie.Producers[0], "Marshall Herskovitz")
	assert.Equal(t, movie.Producers[1], "Edward Zwick")
	assert.Equal(t, movie.Producers[2], "Tom Cruise")
	assert.Equal(t, movie.Producers[3], "Paula Wagner")
	assert.Equal(t, movie.Producers[4], "Scott Kroopf")
	assert.Equal(t, movie.Producers[5], "Tom Engelman")
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

func BenchmarkDecodeIntegers(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-integers.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := json.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeIntegersJsoniter(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-integers.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := jsoniter.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeIntegersStd(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-integers.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := stdJSON.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeFloats(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-floats.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := json.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeFloatsJsoniter(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-floats.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := jsoniter.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeFloatsStd(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-floats.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := stdJSON.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeStrings(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-strings.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := json.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeStringsJsoniter(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-strings.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := jsoniter.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeStringsStd(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-strings.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := stdJSON.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeArrays(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-arrays.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := json.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeArraysJsoniter(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-arrays.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := jsoniter.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeArraysStd(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie-arrays.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := stdJSON.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeAll(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := json.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeAllJsoniter(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := jsoniter.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeAllStd(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	decoder := stdJSON.NewDecoder(reader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkWithNewDecoder(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		decoder := json.NewDecoder(reader)
		err := decoder.Decode(movie)
		decoder.Close()

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkWithNewDecoderJsoniter(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := jsoniter.NewDecoder(reader).Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkWithNewDecoderStd(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie.json")
	reader := bytes.NewReader(data)
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart) // nolint:errcheck
		err := stdJSON.NewDecoder(reader).Decode(movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}
