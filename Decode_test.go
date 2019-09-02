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

	// if movie.Title != "The Last Samurai" || movie.Director != "Edward Zwick" {
	// 	t.Fatalf("Invalid data: %v", movie)
	// }
}

func BenchmarkDecode(b *testing.B) {
	person := &Movie{}
	decoder := json.NewDecoder(jsonReader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		jsonReader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(person)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeStd(b *testing.B) {
	person := &Movie{}
	decoder := stdJSON.NewDecoder(jsonReader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		jsonReader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(person)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeJsoniter(b *testing.B) {
	person := &Movie{}
	decoder := jsoniter.NewDecoder(jsonReader)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		jsonReader.Seek(0, io.SeekStart) // nolint:errcheck
		err := decoder.Decode(person)

		if err != nil && err != io.EOF {
			b.Fatalf("Decode failed: %v", err)
		}
	}
}
