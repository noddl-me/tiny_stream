package stream_test

// Employee the struct of employee
type Employee struct {
	Id       int64
	Name     *string
	Age      *int
	Phone    *string
	Position *PositionInfo
}

// PositionInfo is the struct of info
type PositionInfo struct {
	Province *string
	Country  *string
	City     *string
}
