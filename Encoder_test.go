package json_test

import (
	"io/ioutil"
	"testing"

	stdJSON "encoding/json"

	"github.com/akyoto/json"
	jsoniter "github.com/json-iterator/go"
)

var object = &Movie{
	Title:    "The Last Samurai",
	Director: "Edward Zwick",
}

func TestEncode(t *testing.T) {
	err := json.Encode(object, ioutil.Discard)

	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
}

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := json.Encode(object, ioutil.Discard)

		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
	}
}

func BenchmarkEncodeJsoniter(b *testing.B) {
	b.ReportAllocs()
	encoder := jsoniter.NewEncoder(ioutil.Discard)

	for i := 0; i < b.N; i++ {
		err := encoder.Encode(object)

		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
	}
}

func BenchmarkEncodeStd(b *testing.B) {
	b.ReportAllocs()
	encoder := stdJSON.NewEncoder(ioutil.Discard)

	for i := 0; i < b.N; i++ {
		err := encoder.Encode(object)

		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
	}
}
