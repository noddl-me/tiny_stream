// Package stream_test
package stream_test

import "stream_test/stream"

// Question1Sub1 Q1: 输入 employees，返回 年龄 >22岁 的所有员工，年龄总和
func Question1Sub1(employees []*Employee) int64 {
	return stream.OfSlice(employees).
		Filter(func(e stream.T) bool {
			age := e.(*Employee).Age
			return age != nil && *age > 22
		}).
		Reduce(func(acc stream.R, e stream.T, idx int, sLen int) stream.R {
			return acc.(int64) + int64(*e.(*Employee).Age)
		}, int64(0)).(int64)
}

// Question1Sub2 Q2: - 输入 employees，返回 id 最小的十个员工，按 id 升序排序
func Question1Sub2(employees []*Employee) []*Employee {
	return stream.OfSlice(employees).
		Sort(func(left stream.T, right stream.T) int {
			return int(left.(*Employee).Id - right.(*Employee).Id)
		}).Limit(10).
		Reduce(func(acc stream.R, e stream.T, idx int, sLen int) stream.R {
			return append(acc.([]*Employee), e.(*Employee))
		}, []*Employee{}).([]*Employee)
}

// Question1Sub3 Q3: - 输入 employees，对于没有手机号为0的数据，随机填写一个
func Question1Sub3(employees []*Employee) []*Employee {
	return stream.OfSlice(employees).
		Map(func(e stream.T) stream.R {
			ele := e.(*Employee)
			if ele.Phone == nil || *ele.Phone == "" || *ele.Phone == "0" {
				phone := "10086"
				ele = &Employee{
					Id:       ele.Id,
					Age:      ele.Age,
					Name:     ele.Name,
					Position: ele.Position,
					Phone:    &phone,
				}
			}
			return ele
		}).
		Reduce(func(acc stream.R, e stream.T, idx int, sLen int) stream.R {
			return append(acc.([]*Employee), e.(*Employee))
		}, make([]*Employee, 0, len(employees))).([]*Employee)
}

// Question1Sub4 Q4: - 输入 employees ，返回一个map[int][]int，其中 key 为 员工年龄 Age，value 为该年龄段员工ID
func Question1Sub4(employees []*Employee) map[int][]int64 {
	return stream.OfSlice(employees).
		Reduce(func(acc stream.R, e stream.T, idx int, sLen int) stream.R {
			ele := e.(*Employee)
			m := acc.(map[int][]int64)
			m[*ele.Age] = append(m[*ele.Age], ele.Id)
			return acc
		}, make(map[int][]int64)).(map[int][]int64)
}
