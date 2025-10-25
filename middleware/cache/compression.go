package cache

import (
	"sync"

	"github.com/klauspost/compress/zstd"
)

var (
	writerPool = sync.Pool{
		New: func() any {
			writer, err := zstd.NewWriter(
				nil,
				zstd.WithEncoderCRC(true),
				zstd.WithEncoderLevel(zstd.SpeedDefault),
			)
			if err != nil {
				return nil
			}
			return writer
		},
	}
	readerPool = sync.Pool{
		New: func() any {
			reader, err := zstd.NewReader(nil)
			if err != nil {
				return nil
			}
			return reader
		},
	}
)

// ZstdDecode decompresses zstd-encoded bytes.
func ZstdDecode(value []byte) ([]byte, error) {
	reader := readerPool.Get().(*zstd.Decoder)
	defer readerPool.Put(reader)
	return reader.DecodeAll(value, nil)
}

// ZstdEncode compresses bytes with zstd.
func ZstdEncode(value []byte) []byte {
	writer := writerPool.Get().(*zstd.Encoder)
	defer writerPool.Put(writer)
	return writer.EncodeAll(value, nil)
}
