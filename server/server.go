package server

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
)

// bufferSize is the size in bytes of a buffer used by the HTTP server to read in the body of a request
const bufferSize = 2048

// BasicHandler handles requests made to the server without using any object pool based optimizations
func BasicHandler(w http.ResponseWriter, req *http.Request) {
	buffer := bytes.NewBuffer(make([]byte, bufferSize))
	if _, err := buffer.ReadFrom(req.Body); err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

var p = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, bufferSize))
	},
}

// ObjectPoolHandler handles requests to the server using an object pool based optimization
func ObjectPoolHandler(w http.ResponseWriter, req *http.Request) {
	buffer := p.Get().(*bytes.Buffer)
	defer func() {
		buffer.Reset()
		p.Put(buffer)
	}()
	if _, err := buffer.ReadFrom(req.Body); err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

// SynchronizedBuffer is a thread-safe buffer
type SynchronizedBuffer struct {
	buffer *bytes.Buffer // buffer is the backing slice of bytes
	lock   sync.Mutex    // lock is used for synchronizing write access to the buffer
	isUsed bool          // isUsed indicates whether the buffer is currently used by some goroutine.
}

// Acquire tries to gain ownership of the buffer
func (b *SynchronizedBuffer) Acquire() {
	b.lock.Lock()
}

// Release relinquishes ownership of the buffer
func (b *SynchronizedBuffer) Release() {
	b.lock.Unlock()
}

// BufferPool is an object pool holding buffers with an upper bound on the total size of
// the pool at any given time
type BufferPool struct {
	size               uint64               // size is the size of a single buffer in the pool
	length             uint64               // length is the number of buffers in the pool
	nextAllocatedIndex int32                // nextAllocatedIndex is the index of the next buffer to allocate to clients
	buffers            []SynchronizedBuffer // buffers contains all buffers allocated by the pool
}

// NewBufferPool creates a new buffer pool with a given size for each single buffer allocated
// by the pool, the desired initial number of buffers in the pool and an upper limit on the number
// of buffers that can be owned by the pool at any given time.
func NewBufferPool(singleBufferSize, length uint64) *BufferPool {
	return &BufferPool{
		size:   singleBufferSize,
		length: length,
		buffers: func() []SynchronizedBuffer {
			buffers := make([]SynchronizedBuffer, length)
			for i := uint64(0); i < length; i++ {
				buffers[i].buffer = bytes.NewBuffer(make([]byte, singleBufferSize))
			}
			return buffers
		}(),
	}
}

// Get repeatedly tries to gain ownership of a buffer in the pool until it succeeds
func (r *BufferPool) Get() *SynchronizedBuffer {
	currIndex := atomic.AddInt32(&r.nextAllocatedIndex, 1) % int32(r.length)
	r.buffers[currIndex].Acquire()
	return &r.buffers[currIndex]
}

// Put returns ownership of a buffer to the pool
func (r *BufferPool) Put(b *SynchronizedBuffer) {
	b.Release()
}

var bp = NewBufferPool(bufferSize, 10)

// BoundedPoolHandler handles requests to the server using an object pool based optimization
// with an upper bound on the object pool's size
func BoundedPoolHandler(w http.ResponseWriter, req *http.Request) {
	syncedBuffer := bp.Get()
	buffer := syncedBuffer.buffer
	defer func() {
		buffer.Reset()
		bp.Put(syncedBuffer)
	}()
	if _, err := buffer.ReadFrom(req.Body); err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
