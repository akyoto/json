package json_test

import (
	stdJSON "encoding/json"
	"io"
	"io/ioutil"
	"testing"

	"github.com/akyoto/json"
	jsoniter "github.com/json-iterator/go"
)

func BenchmarkUnmarshal(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie.json")
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(data, movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Unmarshal failed: %v", err)
		}
	}
}

func BenchmarkUnmarshalJsoniter(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie.json")
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := jsoniter.Unmarshal(data, movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Unmarshal failed: %v", err)
		}
	}
}

func BenchmarkUnmarshalStd(b *testing.B) {
	data, _ := ioutil.ReadFile("testdata/movie.json")
	movie := &Movie{}
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := stdJSON.Unmarshal(data, movie)

		if err != nil && err != io.EOF {
			b.Fatalf("Unmarshal failed: %v", err)
		}
	}
}
