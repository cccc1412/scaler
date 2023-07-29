package scaler

import (
	"fmt"
	"log"
	"math"
	"sort"
)

type Distribution struct {
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
}

func NewDistribution(sample_num int, numBins int) *Distribution {
  return &Distribution{
    collected: false,
    bins: make([]int, numBins+1),
    collect_data: make([]int, sample_num),
    counts: make([]int, numBins),
    num_bins: numBins,
    sample_num: sample_num,
  }
}

func (d *Distribution)InitDistribution() {
  sort.Ints(d.collect_data)
	min := d.collect_data[0]
	max := d.collect_data[d.sample_num - 1]

  log.Printf("Distribution init, min = %d, max = %d", min, max)
  width := (max - min) / d.num_bins 
  for i := range d.bins {
    d.bins[i] = min + i*width
  }

  for _, num := range d.collect_data {
    d.Add(num)
  }
}

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
    fmt.Println("total:", d.total)
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
    fmt.Println("i:", i)
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
      return d.bins[i]
    }
  }

  return d.bins[len(d.bins) - 1]
}
