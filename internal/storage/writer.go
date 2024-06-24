package storage

type writer struct {
	id           string
	maxLimit     int64
	currentSize  int64
	index        int
	currentChunk map[string][]byte
	chunks       []map[string][]byte
}

func newWriter(id string, maxLimit int64, chunks ...map[string][]byte) *writer {
	index := 0
	return &writer{
		id:           id,
		maxLimit:     maxLimit,
		chunks:       chunks,
		index:        index,
		currentChunk: chunks[index],
	}
}

func (w *writer) Write(p []byte) (n int, err error) {
	if w.currentSize >= w.maxLimit {
		w.index++
		w.currentChunk = w.chunks[w.index]
		w.currentSize = 0
	}
	w.currentChunk[w.id] = append(w.currentChunk[w.id], p...)
	w.currentSize += int64(len(p))
	return len(p), nil
}
