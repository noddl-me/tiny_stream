package stream

type (
	// T is a empty interface, that is `any` type.
	// since Go is not support generics now(but will coming soon),
	// so we use T to represent any type
	T interface{}
	// R is another `any` type used to distinguish T
	R interface{}
	// U is another `any` type used to distinguish T and R
	U interface{}
	// Function represents a conversion ability, which accepts one argument and produces a result
	Function func(e T) R
	// IntFunction is a Function, which result type is int
	IntFunction func(e T) int
	// Predicate is a Function, which produces a bool value. usually used to test a value whether satisfied condition
	Predicate func(e T) bool
	// UnaryOperator is a Function, which argument and result are the same type
	UnaryOperator func(e T) T
	// Consumer accepts one argument and not produces any result
	Consumer func(e T)
	// Supplier returns a result. each time invoked it can returns a new or distinct result
	Supplier func() T
	// BiFunction like Function, but is accepts two arguments and produces a result
	BiFunction func(t T, u U) R
	// BinaryOperator is a BiFunction which input and result are the same type
	BinaryOperator func(e1 T, e2 T) T
	// Comparator is a BiFunction, which two input arguments are the type, and returns a int.
	// if left is greater then right, it returns a positive number;
	// if left is less then right, it returns a negative number; if the two input are equal, it returns 0
	Comparator func(left T, right T) int
	// Pair is a pair of two element
	Pair struct {
		First  T // First is first element
		Second R // Second is second element
	}
)

// Stream is the main interface of stream
type Stream interface {
	// Intermediate operations
	// Stateless operation

	// Concat concat with stream
	Concat(Stream) Stream
	// Filter filters out if elements match the condition
	Filter(Predicate) Stream
	// Limit limits elements
	Limit(int) Stream
	// Map maps elements with function
	Map(Function) Stream
	// Skip skips elements
	Skip(int) Stream
	// Slice return the stream with[start, start+count)
	Slice(start, count int) Stream

	// Stateful operation

	// Fill fill the stream with T
	Fill(T) Stream
	// Pop pop the last element
	Pop() Stream
	// Push insert the element at last
	Push(T) Stream
	// Reverse reverse the stream
	Reverse() Stream
	// Shift remove the first element
	Shift() Stream
	// Unique de-duplicates elements
	Unique(IntFunction) Stream
	// Unshift insert the element at front
	Unshift(T) Stream
	// Sort sorts elements
	Sort(Comparator) Stream

	// Terminate operation
	// Non-short-circuiting

	// Count return the count of stream
	Count() int
	// ForEach traversal the stream
	ForEach(Consumer)
	// Join join all element with splitter
	Join(string) string
	// Reduce return nil if no element. calculate result by (R, T) -> R from init element, panic if reduction is nil
	Reduce(accumulator func(acc R, e T, idx int, sLen int) R, initValue R) R
	// ToSlice reduce the stream to slice
	ToSlice() []T

	// Short-circuiting

	// AllMatch test if all elements match the condition
	AllMatch(Predicate) bool
	// AnyMatch test if any element matches the condition
	AnyMatch(Predicate) bool
	// FindFirst return the first element that matches the condition
	FindFirst(Predicate) T
}
