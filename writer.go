package disruptor

type Writer struct {
	written  *Cursor
	upstream Barrier
	capacity int64
	previous int64
	gate     int64
	pad1     int64
	pad2     int64
}

func NewWriter(written *Cursor, upstream Barrier, capacity int64) *Writer {
	assertPowerOfTwo(capacity)

	return &Writer{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: InitialSequenceValue,
		gate:     InitialSequenceValue,
	}
}

func assertPowerOfTwo(value int64) {
	if value > 0 && (value&(value-1)) != 0 {
		// Wikipedia entry: http://bit.ly/1krhaSB
		panic("The ring capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}
}

func (this *Writer) Reserve() int64 {
	next := this.previous + 1

	for next-this.capacity > this.gate {
		this.gate = this.upstream.Read(0)
	}

	this.previous = next
	return next
}

func (this *Writer) Commit(sequence int64) {
	this.written.sequence = sequence
}
