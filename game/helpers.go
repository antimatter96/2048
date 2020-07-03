package game

type mover func(int) int
type comp func(int) bool

func inc(n int) int {
	return n + 1
}
func dec(n int) int {
	return n - 1
}

func greaterThan(limit int) func(n int) bool {
	return func(n int) bool {
		return n > limit
	}
}
func lessThan(limit int) func(n int) bool {
	return func(n int) bool {
		return n < limit
	}
}
