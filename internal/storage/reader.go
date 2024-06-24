package storage

import "io"

type reader struct {
	id           string
	currentChunk map[string][]byte
	chunks       []map[string][]byte
	index        int64
	chunkIndex   int
}

func newReader(id string, chunks ...map[string][]byte) *reader {
	chunkIndex := 0
	return &reader{
		id:           id,
		chunks:       chunks,
		currentChunk: chunks[chunkIndex],
		chunkIndex:   chunkIndex,
	}
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.index >= int64(len(r.currentChunk[r.id])) {
		r.chunkIndex++
		if r.chunkIndex >= len(r.chunks) {
			err = io.EOF
			return
		}
		r.currentChunk, r.index = r.chunks[r.chunkIndex], 0
	}
	n = copy(p, r.currentChunk[r.id][r.index:])
	r.index += int64(n)

	return n, nil
}
