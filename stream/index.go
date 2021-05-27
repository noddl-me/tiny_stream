// Package stream is the package of stream
package stream

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type sortable struct {
	List []T
	Cmp  Comparator
}

func (a *sortable) Len() int {
	return len(a.List)
}

func (a *sortable) Less(i, j int) bool {
	return a.Cmp(a.List[i], a.List[j]) < 0
}

func (a *sortable) Swap(i, j int) {
	a.List[i], a.List[j] = a.List[j], a.List[i]
}

type stream struct {
	data []T     // data is the source of stream
	opts []stage // opts is the operations of stream
	para uint32  // para is the atomic field, which flag the stream execute parallel

	// prod []R     // prod is the product of stream
	// prev Stream  // prev is the stream state, the stream is head if this field is nil
	// wrap wrapper // wrap is the wrapper of the next stream, which will recurves call in terminal
}

func (s *stream) getStageMachine() *stageMachine {
	ret := &stageMachine{
		s: s,
	}

	stages := make([][]stage, 0)
	for _, sg := range s.opts {
		if sLen := len(stages); sLen > 0 {
			if stages[sLen-1][0].stageFlag == sg.stageFlag && sg.stageFlag == stageStateless {
				// only stateless actions can be group
				stages[sLen-1] = append(stages[sLen-1], sg)
			} else {
				stages = append(stages, make([]stage, 1))
				stages[sLen][0] = sg
			}
		} else {
			stages = append(stages, make([]stage, 1))
			stages[0][0] = sg
		}
		if sg.prepare != nil {
			sg.prepare()
		}
	}
	ret.stages = stages
	return ret
}

func (s *stream) terminate(sg stage) {
	opts := s.opts
	// each add stage was a make(), so the opts.cap always equal to opts.len
	// that append will alloc a new slice
	s.opts = append(s.opts, sg)
	machine := s.getStageMachine()
	machine.run()
	s.opts = opts // make last effect (terminate op) unavailable

	// s.opts = s.opts[:len(s.opts)-1] // make last effect (terminate op) unavailable
}

func (s *stream) addStage(action func(ele T, i int) (R, bool), flag int, prepare func()) *stream {
	ret := &stream{
		data: s.data,
		opts: make([]stage, len(s.opts)+1),
		para: s.para,
	}
	copy(ret.opts, s.opts)
	ret.opts[len(s.opts)] = stage{
		prepare:   prepare,
		action:    action,
		stageFlag: flag,
	}
	return ret
}

// Concat concat with stream
func (s *stream) Concat(other Stream) Stream {
	return &stream{
		data: append(s.data, other.(*stream).data...), // concat the source, but it operation will not effects outside
		opts: append(s.opts, other.(*stream).opts...), // then concat the opts
		para: s.para + other.(*stream).para,
	}
}

// Filter filters out if elements match the condition
func (s *stream) Filter(f Predicate) Stream {
	return s.addStage(
		func(ele T, i int) (R, bool) {
			if !f(ele) {
				ele = nil
			}
			return ele, false
		}, stageStateless, nil)
}

// Limit limits elements
func (s *stream) Limit(limit int) Stream {
	count := 0
	return s.addStage(
		func(ele T, i int) (R, bool) {
			if count >= limit {
				return nil, true
			}
			count++
			return ele, false
		}, stageStateless, func() {
			count = 0
		})
}

// Map maps elements with function
func (s *stream) Map(f Function) Stream {
	return s.addStage(
		func(ele T, i int) (R, bool) {
			return f(ele), false
		}, stageStateless, nil)
}

// Skip skips elements
func (s *stream) Skip(num int) Stream {
	return s.addStage(
		func(ele T, i int) (R, bool) {
			if i < num {
				return nil, false
			}
			return ele, false
		}, stageStateless, nil)
}

// Slice return the stream with[start, start+count)
func (s *stream) Slice(start, count int) Stream {
	return s.Skip(start).Limit(count)
}

// Fill fill the stream with T
func (s *stream) Fill(e T) Stream {
	return s.addStage(
		func(ele T, i int) (R, bool) {
			return e, false
		}, stageStateless, nil)
}

// Pop pop the last element
func (s *stream) Pop() Stream {
	return s.addStage(
		func(ele T, i int) (R, bool) {
			prod := ele.([]T)
			return prod[:i-1], false
		}, stageStateful, nil)
}

// Push insert the element at last
func (s *stream) Push(e T) Stream {
	return s.addStage(
		func(ele T, i int) (R, bool) {
			prod := ele.([]T)
			return append(prod, e), false
		}, stageStateful, nil)
}

// Reverse reverse the stream
func (s *stream) Reverse() Stream {
	return s.addStage(
		func(ele T, pLen int) (R, bool) {
			prod := ele.([]T)
			last := pLen - 1
			for i := 0; i < pLen/2; i++ {
				prod[i], prod[last-i] = prod[last-i], prod[i]
			}
			return prod, false
		}, stageStateful, nil)
}

// Shift remove the first element
func (s *stream) Shift() Stream {
	return s.addStage(
		func(ele T, i int) (R, bool) {
			prod := ele.([]T)
			return prod[1:], false
		}, stageStateful, nil)
}

// Unique de-duplicates elements
func (s *stream) Unique(f IntFunction) Stream {
	return s.addStage(
		func(ele T, sLen int) (R, bool) {
			prod := ele.([]T)
			set := make(map[int]bool)
			i := 0
			for _, e := range prod {
				h := f(e)
				if _, ok := set[h]; !ok {
					set[h] = true
					prod[i] = e
					i++
				}
			}
			return prod[:i], false
		}, stageStateful, nil)
}

// Unshift insert the element at front
func (s *stream) Unshift(e T) Stream {
	return s.addStage(
		func(ele T, i int) (R, bool) {
			prod := ele.([]T)
			return append([]T{e}, prod...), false
		}, stageStateful, nil)
}

// Sort sorts elements
func (s *stream) Sort(f Comparator) Stream {
	return s.addStage(
		func(ele T, i int) (R, bool) {
			sort.Sort(&sortable{
				List: ele.([]T),
				Cmp:  f,
			})
			return ele, false
		}, stageStateful, nil)
}

// AllMatch test if all elements match the condition
func (s *stream) AllMatch(f Predicate) (ret bool) {
	s.terminate(stage{
		action: func(ele T, i int) (R, bool) {
			ret = f(ele)
			return nil, !ret
		},
		stageFlag: stageShortcut,
	})
	return
}

// AnyMatch test if any element matches the condition
func (s *stream) AnyMatch(f Predicate) (ret bool) {
	s.terminate(stage{
		action: func(ele T, i int) (R, bool) {
			ret = f(ele)
			return nil, ret
		},
		stageFlag: stageShortcut,
	})
	return
}

// Count return the count of stream
func (s *stream) Count() (ret int) {
	s.terminate(stage{
		action: func(ele T, i int) (R, bool) {
			ret = i
			return nil, true
		},
		stageFlag: stageNonShortcut,
	})
	return
}

// FindFirst return the first element that matches the condition
func (s *stream) FindFirst(f Predicate) (ret T) {
	s.terminate(stage{
		action: func(ele T, i int) (R, bool) {
			if f(ele) {
				ret = ele
				return nil, true
			}
			return nil, false
		},
		stageFlag: stageShortcut,
	})
	return
}

// ForEach traversal the stream
func (s *stream) ForEach(f Consumer) {
	s.terminate(stage{
		action: func(ele T, i int) (R, bool) {
			for _, e := range ele.([]T) {
				f(e)
			}
			return nil, true
		},
		stageFlag: stageNonShortcut,
	})
}

// Join join all element with splitter
func (s *stream) Join(split string) string {
	return s.Reduce(func(acc R, e T, i int, sLen int) R {
		sb := acc.(*strings.Builder)
		sb.WriteString(fmt.Sprintf("%v", e))
		if i < sLen-1 {
			sb.WriteString(split)
		}
		return sb
	}, &strings.Builder{}).(*strings.Builder).String()
}

// Reduce return nil if no element. calculate result by (R, T) -> R from init element, panic if reduction is nil
func (s *stream) Reduce(accumulator func(R, T, int, int) R, initValue R) (ret R) {
	ret = initValue
	s.terminate(stage{
		action: func(ele T, sLen int) (R, bool) {
			for i, e := range ele.([]T) {
				ret = accumulator(ret, e, i, sLen)
			}
			return nil, true
		},
		stageFlag: stageNonShortcut,
	})
	return
}

// ToSlice reduce the stream to slice
func (s *stream) ToSlice() (ret []T) {
	s.terminate(stage{
		action: func(ele T, i int) (R, bool) {
			ret = ele.([]T)
			return nil, true
		},
		stageFlag: stageNonShortcut,
	})
	return
}

// Of creates a Stream from slice
func Of(elements ...T) Stream {
	return &stream{
		data: elements,
		opts: make([]stage, 0),
		para: 0,
	}
}

// OfSlice construct a stream form slice
func OfSlice(s T) Stream {
	// if reflect.TypeOf(s).Kind() != reflect.Slice {
	// 	panic("arg is not slice")
	// }
	v := reflect.ValueOf(s)
	it := make([]T, v.Len())
	for i := range it {
		it[i] = v.Index(i).Interface()
	}
	return &stream{
		data: it,
		opts: make([]stage, 0),
		para: 0,
	}
}

// Repeat is constructor of repeat element
func Repeat(ele T, times int) Stream {
	s := make([]T, times)
	for i := range s {
		s[i] = ele
	}
	return Of(s...)
}
