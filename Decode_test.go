package json_test

import (
	stdJSON "encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/akyoto/json"
	jsoniter "github.com/json-iterator/go"
)

var (
	jsonString = `{"title":"The Last Samurai","director":"Edward Zwick"}`
	jsonBytes  = []byte(jsonString)
	jsonReader = strings.NewReader(jsonString)
)

type Movie struct {
	Title    string `json:"title"`
	Director string `json:"director"`
}

func TestDecode(t *testing.T) {
	file, err := os.Open("testdata/movie-simple.json")

	if err != nil {
		t.Fatalf("os.Open failed: %v", err)
	}

	defer file.Close()

	movie := &Movie{}
	err = json.NewDecoder(file).Decode(movie)

	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if movie.Title != "The Last Samurai" || movie.Director != "Edward Zwick" {
		t.Fatalf("Invalid data: %v", movie)
	}
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
