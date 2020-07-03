package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type game interface {
	Move()
	Init()
	End()
	Lost()
	Won()
}

const N = 4

type TwoZeroFourEight struct {
	board []([]int)

	rand *rand.Rand
}

func (g *TwoZeroFourEight) FillRandom() {
	x, y, err := g.BoundaryCells()
	if err != nil {
		fmt.Println(err)
		return
	}

	g.board[x][y] = 2 << (g.rand.Intn(1) + 1)
}

func (game *TwoZeroFourEight) Init() {
	source := rand.NewSource(time.Now().UnixNano())
	game.rand = rand.New(source)

	game.board = make([][]int, N)
	for i := 0; i < N; i++ {
		game.board[i] = make([]int, N)
	}

	game.board[N-1][0] = 2
	game.board[N-1][1] = 2
	game.board[N-2][0] = 2
}

func (g *TwoZeroFourEight) Print() {
	fmt.Println(strings.Repeat("-", 1+2<<(N)+N))
	for i := 0; i < N; i++ {
		fmt.Print("|  ")
		for j := 0; j < N-1; j++ {
			fmt.Printf("%-4d  |  ", g.board[i][j])
		}
		fmt.Printf("%-4d", g.board[i][N-1])
		fmt.Print("  |\n")
		fmt.Println(strings.Repeat("-", 1+2<<(N)+N))
	}
}

func (g *TwoZeroFourEight) BoundaryCells() (int, int, error) {
	mp := make(map[string]bool)

	for i := 0; i < N; i++ {
		if g.board[i][0] == 0 && g.board[i][N-1] == 0 {
			continue
		}
		if g.board[i][0] == 0 {
			for j := (N - 1); j > -1; j-- {
				if g.board[i][j] == 0 {
					mp[fmt.Sprintf("%d,%d", i, j)] = true
					break
				}
			}
		} else {
			for j := 0; j < N; j++ {
				if g.board[i][j] == 0 {
					mp[fmt.Sprintf("%d,%d", i, j)] = true
					break
				}
			}
		}
	}

	for j := 0; j < N; j++ {
		if g.board[0][j] == 0 && g.board[N-1][j] == 0 {
			continue
		}
		if g.board[0][j] == 0 {
			for i := (N - 1); i > -1; i-- {
				if g.board[i][j] == 0 {
					mp[fmt.Sprintf("%d,%d", i, j)] = true
					break
				}
			}
		} else {
			for i := 0; i < N; i++ {
				if g.board[i][j] == 0 {
					mp[fmt.Sprintf("%d,%d", i, j)] = true
					break
				}
			}
		}
	}

	var arr []string

	for k := range mp {
		arr = append(arr, k)
	}

	if len(arr) == 0 {
		fmt.Println(">>>>>")
		g.Print()
		fmt.Println("<<<<<")
		return 0, 0, fmt.Errorf("%v", "FULL BOARD")
	}

	picked := g.rand.Intn(len(arr))

	xy := strings.Split(arr[picked], ",")
	x, _ := strconv.Atoi(xy[0])
	y, _ := strconv.Atoi(xy[1])

	return x, y, nil
}

func (g *TwoZeroFourEight) MovesPossible() bool {
	for i := 0; i < N; i++ {
		for j := 0; j < N-1; j++ {
			if g.board[i][j] == 0 {
				continue
			}
			if g.board[i][j] == g.board[i][j+1] {
				return true
			}
		}
	}

	for j := 0; j < N; j++ {

		for i := 0; i < N-1; i++ {
			if g.board[i][j] == 0 {
				continue
			}
			if g.board[i][j] == g.board[i+1][j] {
				return true
			}
		}
	}

	return false
}

func inc(n int) int {
	return n + 1
}
func dec(n int) int {
	return n - 1
}

type mover func(int) int
type comp func(int) bool

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

func (g *TwoZeroFourEight) moveHorizontal(changeI, changeJ mover, compI, compJ comp, startI, startJ int) bool {

	changed := false

	for j := startJ; compJ(j); j = changeJ(j) {
		for i := startI; compI(i); i = changeI(i) {
			if g.board[i][j] == 0 {
				continue
			}
			if g.board[i][j] == g.board[i][changeJ(j)] {
				fmt.Println("chanign")
				g.board[i][j] = 0
				g.board[i][changeJ(j)] *= 2
				changed = true
			} else if g.board[i][changeJ(j)] == 0 {
				fmt.Println("moving")
				g.board[i][changeJ(j)] = g.board[i][j]
				g.board[i][j] = 0
			}
		}
	}

	if changed {
		changed = changed || g.moveHorizontal(changeI, changeJ, compI, compJ, startI, startJ)
	}

	return changed
}

func (g *TwoZeroFourEight) MoveLeft() bool {
	greaterThanZero := greaterThan(0)
	lessThanN := lessThan(N)

	return g.moveHorizontal(inc, dec, lessThanN, greaterThanZero, 0, N-1)

	// for j := N - 1; greaterThanZero(j); j = dec(j) {
	// 	for i := 0; lessThanN(i); i = inc(i) {
	// 		if g.board[i][j] == 0 {
	// 			continue
	// 		}
	// 		if g.board[i][j] == g.board[i][dec(j)] {
	// 			fmt.Println("chanign")
	// 			g.board[i][j] = 0
	// 			g.board[i][dec(j)] *= 2
	// 			changed = true
	// 		} else if g.board[i][dec(j)] == 0 {
	// 			fmt.Println("moving")
	// 			g.board[i][dec(j)] = g.board[i][j]
	// 			g.board[i][j] = 0
	// 		}
	// 	}
	// }
}

func (g *TwoZeroFourEight) MoveRight() bool {
	lessThanN := lessThan(N)
	lessThanNMinusOne := lessThan(N - 1)

	return g.moveHorizontal(inc, inc, lessThanN, lessThanNMinusOne, 0, 0)

	// for j := 0; lessThanNMinusOne(j); j = inc(j) {
	// 	for i := 0; lessThanN(i); i = inc(i) {
	// 		if g.board[i][j] == 0 {
	// 			continue
	// 		}
	// 		if g.board[i][j] == g.board[i][inc(j)] {
	// 			fmt.Println("chanign")
	// 			g.board[i][j] = 0
	// 			g.board[i][j+1] *= 2
	// 			changed = true
	// 		} else if g.board[i][inc(j)] == 0 {
	// 			fmt.Println("moving")
	// 			g.board[i][inc(j)] = g.board[i][j]
	// 			g.board[i][j] = 0
	// 		}
	// 	}
	// }
}

// func (g *TwoZeroFourEight) MoveLeft() bool {
// 	changed := false
// 	for j := N - 1; j > 0; j-- {
// 		for i := 0; i < N; i++ {
// 			if g.board[i][j] == 0 {
// 				continue
// 			}
// 			if g.board[i][j] == g.board[i][j-1] {
// 				fmt.Println("chanign")
// 				g.board[i][j] = 0
// 				g.board[i][j-1] *= 2
// 				changed = true
// 			}
// 		}
// 	}

// 	return changed
// }

// func (g *TwoZeroFourEight) MoveLeft() bool {
// 	changed := false
// 	for j := N - 1; j > 0; j-- {
// 		for i := 0; i < N; i++ {
// 			if g.board[i][j] == 0 {
// 				continue
// 			}
// 			if g.board[i][j] == g.board[i][j-1] {
// 				fmt.Println("chanign")
// 				g.board[i][j] = 0
// 				g.board[i][j-1] *= 2
// 				changed = true
// 			}
// 		}
// 	}

// 	return changed
// }
