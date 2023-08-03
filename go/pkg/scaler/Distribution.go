package scaler

import (
	"fmt"
	"log"
	"math"
	"sort"
  "errors"
)


// RingBuffer represents a ring buffer.
type RingBuffer struct {
	buffer []int
	size   int
	head   int
	tail   int
}

// NewRingBuffer creates a new ring buffer with the specified size.
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]int, size),
		size:   size,
		head:   -1,
		tail:   -1,
	}
}

// IsEmpty checks if the ring buffer is empty.
func (rb *RingBuffer) IsEmpty() bool {
	return rb.head == -1
}

// IsFull checks if the ring buffer is full.
func (rb *RingBuffer) IsFull() bool {
	return (rb.tail+1)%rb.size == rb.head
}

// Enqueue adds a value to the ring buffer.
func (rb *RingBuffer) Enqueue(value int) error {
	if rb.IsFull() {
		return errors.New("Ring buffer is full")
	}

	if rb.IsEmpty() {
		rb.head = 0
	}

	rb.tail = (rb.tail + 1) % rb.size
	rb.buffer[rb.tail] = value

	return nil
}

// Dequeue removes and returns a value from the ring buffer.
func (rb *RingBuffer) Dequeue() (int, error) {
	if rb.IsEmpty() {
		return 0, errors.New("Ring buffer is empty")
	}

	value := rb.buffer[rb.head]

	if rb.head == rb.tail {
		rb.head = -1
		rb.tail = -1
	} else {
		rb.head = (rb.head + 1) % rb.size
	}

	return value, nil
}

// ToArray converts the ring buffer to an array.
func (rb *RingBuffer) ToArray() []int {
	if rb.IsEmpty() {
		return []int{}
	}

	if rb.head <= rb.tail {
		return rb.buffer[rb.head : rb.tail+1]
	}

	return append(rb.buffer[rb.head:], rb.buffer[:rb.tail+1]...)
}
func intToFloatArray(intArray []int) []float64 {
	floatArray := make([]float64, len(intArray))

	for i, val := range intArray {
		floatArray[i] = float64(val)
	}

	return floatArray
}
type Distribution struct {
    assgin_ts RingBuffer
    idle_ts RingBuffer  
    collected bool
    sample_num int 
    cur_sample int
    num_bins int
    bins   []int // 每个桶的边界
    counts []int     // 每个桶中的样本数量
    mean   float64   // 均值
    m2     float64   // 二阶中心矩
    // data   []int
    collect_data []int
    total  int
    num_reuse int
}

func NewDistribution(sample_num int, numBins int) *Distribution {
  return &Distribution{
    assgin_ts: *NewRingBuffer(100),
    idle_ts: *NewRingBuffer(100),
    collected: false,
    bins: make([]int, numBins+1),
    collect_data: make([]int, sample_num),
    counts: make([]int, numBins),
    num_bins: numBins,
    sample_num: sample_num,
  }
}

func (d *Distribution) PrintDistribution() {
  debug_str := "Distribution init:"
  for i := range d.counts {
    debug_str += fmt.Sprint(d.counts[i]) + ", " 
  }
  log.Printf(debug_str)
}

func (d *Distribution)InitDistribution() {
  sort.Ints(d.collect_data)
	min := d.collect_data[0]
	max := d.collect_data[d.sample_num - 1]

  log.Printf("Distribution init, min = %d, max = %d, cv = %f", min, max, d.CV())
  debug_str := "Distribution init"
  width := (max - min) / d.num_bins 
  for i := range d.bins {
    d.bins[i] = min + i*width
    debug_str += fmt.Sprint(d.bins[i]) + ", "
  }
  log.Printf(debug_str)
  for _, num := range d.collect_data {
    d.Add(num)
  }
  d.PrintDistribution()
}

func (d *Distribution) PrintReuseRate() {
  log.Printf("reuseRate : %f", float64(d.num_reuse) / float64(d.total))
}


func (d *Distribution) AddAssignTs(x int) {
  d.assgin_ts.Enqueue(x)
}

func (d *Distribution) AddIdleTs(x int) {
  d.idle_ts.Enqueue(x)
}

// func (d *Distribution) IsPreload() bool {
//   last_assign_ts, _ := d.assgin_ts.Last()
//   last_idle_ts, _ := d.idle_ts.Last()
//   if(last_assign_ts < last_idle_ts) {
//     return true
//   }
//   return false
// }

func (d *Distribution) Add(x int) {
  if(!d.collected) {
    d.collect_data[d.cur_sample] = x
    d.cur_sample++
    if(d.cur_sample == d.sample_num) {
      d.collected = true
      d.InitDistribution()
    }
    return
  }
    // 找到 x 所在的桶
  if(d.collected) {
    // d.data = append(d.data, x)
    d.total ++
    // fmt.Println("total:", d.total)
    i := sort.SearchInts(d.bins, x)
    if i == len(d.bins) {
        i--
    }
    if i > 0 && x < d.bins[i-1]+(d.bins[i]-d.bins[i-1])/2 {
        i--
    }
    if i > 0 {
      i--
    }
    // 更新统计信息
    // fmt.Println("i:", i)
    d.counts[i]++
    delta := float64(x) - d.mean
    d.mean += delta / float64(d.total)
    delta2 := float64(x) - d.mean
    d.m2 += delta * delta2
  }
}

func (d *Distribution) CV() float64 {
    if d.total == 0 {
        return math.NaN()
    }
    return math.Sqrt(d.m2 / (float64(d.total - 1))) / d.mean
}

func (d *Distribution) GetQuantiles(p float64) int {
  log.Printf("Distribution, min = %d, max = %d", d.bins[0], d.bins[len(d.bins) - 1])
  sum := 0
  for i := range d.counts {
    sum += d.counts[i]
    if(float64(sum) / float64(d.total) > p) {
      return (d.bins[i] + d.bins[i+1]) / 2
    }
  }

  return (d.bins[len(d.bins) - 2] + d.bins[len(d.bins) - 1]) / 2
}
