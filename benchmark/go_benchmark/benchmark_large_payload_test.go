package go_benchmark

import (
	"encoding/json"
	"testing"

	"github.com/buger/jsonparser"
	jsoniter "github.com/json-iterator/go"
)

func BenchmarkJsonParserLarge(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		count := 0
		jsonparser.ArrayEach(largeFixture, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			count++
		}, "topics", "topics")
	}
}

func BenchmarkJsoniterLarge(b *testing.B) {
	iter := jsoniter.ParseBytes(jsoniter.ConfigDefault, largeFixture)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter.ResetBytes(largeFixture)
		count := 0
		for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
			if "topics" != field {
				iter.Skip()
				continue
			}
			for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
				if "topics" != field {
					iter.Skip()
					continue
				}
				for iter.ReadArray() {
					iter.Skip()
					count++
				}
				break
			}
			break
		}
	}
}

func BenchmarkEncodingJsonLarge(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		payload := &LargePayload{}
		json.Unmarshal(largeFixture, payload)
	}
}

func BenchmarkEncodingJsoniterLarge(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		payload := &LargePayload{}
		jsoniter.Unmarshal(largeFixture, payload)
	}
}
