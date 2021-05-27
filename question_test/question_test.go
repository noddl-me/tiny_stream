package question_test

import (
	"testing"

	"stream_test"

	"github.com/Pallinder/go-randomdata"
	"github.com/sumory/idgen"
)

// ----------------------------------------
var employees []*stream_test.Employee
var str string
var strList []string
var idGen *idgen.IdWorker

// ----------------------------------------

func BenchmarkQuestion1Sub1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stream_test.Question1Sub1(employees)
	}
}

func BenchmarkQuestion1Sub2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stream_test.Question1Sub2(employees)
	}
}

func TestQuestion1Sub1(t *testing.T) {
	answer := stream_test.Question1Sub1(employees)
	t.Log(answer)
}

func TestQuestion1Sub2(t *testing.T) {
	answer := stream_test.Question1Sub2(employees)
	t.Log(answer)
}

func TestQuestion1Sub3(t *testing.T) {
	answer := stream_test.Question1Sub3(employees)
	t.Log(answer)
}

func TestQuestion1Sub4(t *testing.T) {
	answer := stream_test.Question1Sub4(employees)
	t.Log(answer)
}

func TestQuestion2Sub1(t *testing.T) {
	answer := stream_test.Question2Sub1(str)
	t.Log(answer)
}

func BenchmarkQuestion2Sub1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stream_test.Question2Sub1(str)
	}
}

func TestQuestion2Sub2(t *testing.T) {
	answer := stream_test.Question2Sub2(strList)
	t.Log(answer)
}

func BenchmarkQuestion2Sub2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stream_test.Question2Sub2(strList)
	}
}

func TestQuestion3Sub1(t *testing.T) {
	answer := stream_test.Question3Sub1(str, 20)
	t.Log(answer)
}

func init() {
	_, idGen = idgen.NewIdWorker(int64(randomdata.Number(18, 100)))

	for i := 0; i < 10000; i++ {
		employees = append(employees, RandomInstance())
	}

	str = randomdata.RandStringRunes(100)
	for i := 0; i < 10000; i++ {
		strList = append(strList, randomdata.RandStringRunes(2000))
	}
}

func RandomInstance() *stream_test.Employee {
	_, id := idGen.NextId()
	country := randomdata.Country(randomdata.FullCountry)
	province := randomdata.ProvinceForCountry(country)
	city := randomdata.City()
	age := randomdata.Number(18, 100)
	name := randomdata.SillyName()
	phone := ""
	if randomdata.Boolean() {
		phone = randomdata.PhoneNumber()
	}
	return &stream_test.Employee{
		Id:    id,
		Name:  &name,
		Age:   &age,
		Phone: &phone,
		Position: &stream_test.PositionInfo{
			Province: &province,
			Country:  &country,
			City:     &city,
		},
	}
}
