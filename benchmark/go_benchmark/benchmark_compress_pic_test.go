package go_benchmark

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"testing"
)

func BenchmarkGzipPic(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var buf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf = bytes.NewBuffer(nil)
		zw := gzip.NewWriter(buf)
		zw.Write(PicFix)
		if err := zw.Close(); err != nil {
			b.Error(err)
		}
	}
	b.Log(buf.Len())
}

func BenchmarkZlibPic(b *testing.B) {
	b.Log("small size:", len(PicFix))
	b.ReportAllocs()
	b.ResetTimer()

	var buf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf = bytes.NewBuffer(nil)
		zw := zlib.NewWriter(buf)
		zw.Write(PicFix)
		if err := zw.Close(); err != nil {
			b.Error(err)
		}
	}
	b.Log(buf.Len())
}

func BenchmarkGzipPicDecode(b *testing.B) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)
	zw.Write(PicFix)
	if err := zw.Close(); err != nil {
		b.Error(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readBuf := bytes.NewBuffer(buf.Bytes())
		zr, _ := gzip.NewReader(readBuf)
		p := make([]byte, 200)
		zr.Read(p)
		err := zr.Close()
		if err != nil {
			b.Error(err, readBuf)
		}
	}
}

func BenchmarkZlibPicDecode(b *testing.B) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(nil)
	zw := zlib.NewWriter(buf)
	zw.Write(PicFix)
	if err := zw.Close(); err != nil {
		b.Error(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readBuf := bytes.NewBuffer(buf.Bytes())
		zr, _ := zlib.NewReader(readBuf)
		p := make([]byte, 200)
		zr.Read(p)
		err := zr.Close()
		if err != nil {
			b.Error(err, readBuf)
		}
	}
}
