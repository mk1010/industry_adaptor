package go_benchmark

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"testing"
)

func BenchmarkGzipSmall(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var buf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf = bytes.NewBuffer(nil)
		zw := gzip.NewWriter(buf)
		zw.Write(smallFixture)
		if err := zw.Close(); err != nil {
			b.Error(err)
		}
	}
	b.Log(buf.Len())
}

func BenchmarkZlibSmall(b *testing.B) {
	b.Log("small size:", len(smallFixture))
	b.ReportAllocs()
	b.ResetTimer()

	var buf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf = bytes.NewBuffer(nil)
		zw := zlib.NewWriter(buf)
		zw.Write(smallFixture)
		if err := zw.Close(); err != nil {
			b.Error(err)
		}
	}
	b.Log(buf.Len())
}

func BenchmarkGzipSmallDecode(b *testing.B) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)
	zw.Write(smallFixture)
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

func BenchmarkZlibSmallDecode(b *testing.B) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(nil)
	zw := zlib.NewWriter(buf)
	zw.Write(smallFixture)
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

func BenchmarkGzipMedium(b *testing.B) {
	b.Log("medium size:", len(mediumFixture))
	b.ReportAllocs()
	b.ResetTimer()
	var buf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf = bytes.NewBuffer(nil)
		zw := gzip.NewWriter(buf)
		zw.Write(mediumFixture)
		if err := zw.Close(); err != nil {
			b.Error(err)
		}
	}
	b.Log(buf.Len())
}

func BenchmarkZlibMedium(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var buf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf = bytes.NewBuffer(nil)
		zw := zlib.NewWriter(buf)
		zw.Write(mediumFixture)
		if err := zw.Close(); err != nil {
			b.Error(err)
		}
	}
	b.Log(buf.Len())
}

func BenchmarkGzipMediumDecode(b *testing.B) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)
	zw.Write(mediumFixture)
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

func BenchmarkZlibMediumDecode(b *testing.B) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(nil)
	zw := zlib.NewWriter(buf)
	zw.Write(mediumFixture)
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

func BenchmarkGzipLarge(b *testing.B) {
	b.Log("large size:", len(largeFixture))
	b.ReportAllocs()
	b.ResetTimer()
	var buf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf = bytes.NewBuffer(nil)
		zw := gzip.NewWriter(buf)
		zw.Write(largeFixture)
		if err := zw.Close(); err != nil {
			b.Error(err)
		}
	}
	b.Log(buf.Len())
}

func BenchmarkZlibLarge(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var buf *bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf = bytes.NewBuffer(nil)
		zw := zlib.NewWriter(buf)
		zw.Write(largeFixture)
		if err := zw.Close(); err != nil {
			b.Error(err)
		}
	}
	b.Log(buf.Len())
}

func BenchmarkGzipLargeDecode(b *testing.B) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)
	zw.Write(largeFixture)
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

func BenchmarkZlibLargeDecode(b *testing.B) {
	var buf *bytes.Buffer
	buf = bytes.NewBuffer(nil)
	zw := zlib.NewWriter(buf)
	zw.Write(largeFixture)
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
