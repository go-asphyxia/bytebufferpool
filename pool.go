package bytebufferpool

import (
	"sort"
	"sync"
	"sync/atomic"
)

type (
	P struct {
		calls       [steps]uint64
		calibrating uint64

		defaultSize uint64
		maxSize     uint64

		pool *sync.Pool
	}

	callSize struct {
		calls uint64
		size  uint64
	}

	callSizes []callSize
)

const (
	minBitSize = 6
	steps      = 20

	minSize = 1 << minBitSize
	maxSize = 1 << (minBitSize + steps - 1)

	calibrateCallsThreshold = 42000
	maxPercentile           = 0.95
)

var (
	defaultPool = Pool()
)

func Pool() (p *P) {
	p = &P{
		pool: new(sync.Pool),
	}

	return
}

func Get() (b *B) {
	b = defaultPool.Get()
	return
}

func (p *P) Get() (b *B) {
	o := p.pool.Get()

	if o != nil {
		b = o.(*B)
		return
	}

	b = &B{
		bytes: make([]byte, 0, atomic.LoadUint64(&p.defaultSize)),
	}

	return
}

func Put(b *B) {
	defaultPool.Put(b)
}

func (p *P) Put(b *B) {
	l := len(b.bytes)
	i := index(l)

	if atomic.AddUint64(&p.calls[i], 1) > calibrateCallsThreshold {
		p.calibrate()
	}

	maxSize := int(atomic.LoadUint64(&p.maxSize))

	if maxSize == 0 || cap(b.bytes) <= maxSize {
		b.Reset()
		p.pool.Put(b)
	}
}

func (p *P) calibrate() {
	if !atomic.CompareAndSwapUint64(&p.calibrating, 0, 1) {
		return
	}

	sizes := make(callSizes, 0, steps)
	sum := uint64(0)

	for i := range p.calls {
		calls := atomic.SwapUint64(&p.calls[i], 0)
		sum += calls

		sizes = append(sizes, callSize{
			calls: calls,
			size:  minSize << i,
		})
	}

	sort.Sort(sizes)

	l := len(sizes)

	defaultSize := sizes[0].size
	maxSize := defaultSize

	maxSum := uint64(float64(sum) * maxPercentile)

	sum = 0

	for i := 0; i < l; i++ {
		if sum > maxSum {
			break
		}

		sum += sizes[i].calls
		size := sizes[i].size

		if size > maxSize {
			maxSize = size
		}
	}

	atomic.StoreUint64(&p.defaultSize, defaultSize)
	atomic.StoreUint64(&p.maxSize, maxSize)

	atomic.StoreUint64(&p.calibrating, 0)
}

func (ci callSizes) Len() (l int) {
	l = len(ci)
	return
}

func (ci callSizes) Less(i, j int) (less bool) {
	less = ci[i].calls > ci[j].calls
	return
}

func (ci callSizes) Swap(i, j int) {
	ci[i], ci[j] = ci[j], ci[i]
}

func index(n int) (index int) {
	n--
	n >>= minBitSize

	for n > 0 {
		n >>= 1
		index++
	}

	if index >= steps {
		index = steps - 1
	}

	return
}
