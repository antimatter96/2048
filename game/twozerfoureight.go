package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const N = 4
const nDashed = 1 + (2 * (N + (2 << (N - 1)))) // Not even true :p

type TwoZeroFourEight struct {
	board []([]int)
	rand  *rand.Rand
}

func NewGame() *TwoZeroFourEight {
	g := &TwoZeroFourEight{}
	g.init()
	return g
}

func (g *TwoZeroFourEight) init() {
	source := rand.NewSource(time.Now().UnixNano())
	g.rand = rand.New(source)

	g.board = make([][]int, N)
	for i := 0; i < N; i++ {
		g.board[i] = make([]int, N)
	}

	g.board[N-1][0] = 2
	g.board[N-1][1] = 2
	g.board[N-2][0] = 2
}

func (g *TwoZeroFourEight) Print() string {
	var b strings.Builder
	b.Grow(350)

	b.WriteString(strings.Repeat("-", nDashed))
	b.WriteByte('\n')
	for i := 0; i < N; i++ {
		b.WriteString("|  ")
		for j := 0; j < N-1; j++ {
			b.WriteString(fmt.Sprintf("%-5d  |  ", g.board[i][j]))
		}
		b.WriteString(fmt.Sprintf("%-5d", g.board[i][N-1]))
		b.WriteString("  |\n")
		b.WriteString(strings.Repeat("-", nDashed))
		b.WriteByte('\n')
	}

	return b.String()
}

func (g *TwoZeroFourEight) boundaryCells() (int, int, error) {
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

func (g *TwoZeroFourEight) FillRandom() {
	x, y, err := g.boundaryCells()
	if err != nil {
		fmt.Println(err)
		return
	}

	g.board[x][y] = 2 << (g.rand.Intn(1) + 1)
}

func (g *TwoZeroFourEight) movesPossible() bool {
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

func (g *TwoZeroFourEight) moveHorizontal(changeI, changeJ mover, compI, compJ comp, startI, startJ int) bool {
	changed := false

	for j := startJ; compJ(j); j = changeJ(j) {
		for i := startI; compI(i); i = changeI(i) {
			if g.board[i][j] == 0 {
				continue
			}
			if g.board[i][j] == g.board[i][changeJ(j)] {
				//fmt.Println("chanign")
				g.board[i][j] = 0
				g.board[i][changeJ(j)] *= 2
				changed = true
			} else if g.board[i][changeJ(j)] == 0 {
				//fmt.Println("moving")
				g.board[i][changeJ(j)] = g.board[i][j]
				g.board[i][j] = 0
				changed = true
			}
		}
	}

	if changed {
		_ = g.moveHorizontal(changeI, changeJ, compI, compJ, startI, startJ)
	}

	return changed
}

func (g *TwoZeroFourEight) moveLeft() bool {
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

func (g *TwoZeroFourEight) moveRight() bool {
	lessThanN := lessThan(N)
	lessThanNMinusOne := lessThan(N - 1)

	//return g.moveHorizontal(inc, inc, lessThanN, lessThanNMinusOne, 0, 0)

	changed := false
	for j := 0; lessThanNMinusOne(j); j = inc(j) {
		for i := 0; lessThanN(i); i = inc(i) {
			if g.board[i][j] == 0 {
				continue
			}
			if g.board[i][j] == g.board[i][inc(j)] {
				fmt.Println("chanign")
				g.board[i][j] = 0
				g.board[i][j+1] *= 2
				changed = true
			} else if g.board[i][inc(j)] == 0 {
				fmt.Println("moving")
				g.board[i][inc(j)] = g.board[i][j]
				g.board[i][j] = 0
				changed = true
			}
		}
	}

	if changed {
		changed = changed || g.moveLeft()
	}

	return changed
}

func (g *TwoZeroFourEight) moveDown() bool {
	lessThanNMinus1 := lessThan(N - 1)
	lessThanN := lessThan(N)

	return g.moveVertical(inc, inc, lessThanNMinus1, lessThanN, 0, 0)

	// for i := 0; i < N-1; i++ {
	// 	for j := 0; j < N; j++ {
	// 		if g.board[i][j] == 0 {
	// 			continue
	// 		}
	// 		if g.board[i][j] == g.board[i+1][j] {
	// 			fmt.Println("chanign")
	// 			g.board[i][j] = 0
	// 			g.board[i+1][j] *= 2
	// 			changed = true
	// 		} else if g.board[i+1][j] == 0 {
	// 			fmt.Println("moving")
	// 			g.board[i+1][j] = g.board[i][j]
	// 			g.board[i][j] = 0
	// 		}
	// 	}
	// }
}

func (g *TwoZeroFourEight) moveUp() bool {
	greaterThan := greaterThan(0)
	lessThanN := lessThan(N)

	return g.moveVertical(dec, inc, greaterThan, lessThanN, N-1, 0)

	// for i := N - 1; i > 0; i-- {
	// 	for j := 0; j < N; j++ {
	// 		if g.board[i][j] == 0 {
	// 			continue
	// 		}
	// 		if g.board[i][j] == g.board[i-1][j] {
	// 			fmt.Println("chanign")
	// 			g.board[i][j] = 0
	// 			g.board[i-1][j] *= 2
	// 			changed = true
	// 		} else if g.board[i-1][j] == 0 {
	// 			fmt.Println("moving")
	// 			g.board[i-1][j] = g.board[i][j]
	// 			g.board[i][j] = 0
	// 		}
	// 	}
	// }
}

func (g *TwoZeroFourEight) moveVertical(changeI, changeJ mover, compI, compJ comp, startI, startJ int) bool {
	changed := false

	for i := startI; compI(i); i = changeI(i) {
		for j := startJ; compJ(j); j = changeJ(j) {
			if g.board[i][j] == 0 {
				continue
			}
			if g.board[i][j] == g.board[changeI(i)][j] {
				//fmt.Println("chanign")
				g.board[i][j] = 0
				g.board[changeI(i)][j] *= 2
				changed = true
			} else if g.board[changeI(i)][j] == 0 {
				//fmt.Println("moving")
				g.board[changeI(i)][j] = g.board[i][j]
				g.board[i][j] = 0
			}
		}
	}

	if changed {
		_ = g.moveVertical(changeI, changeJ, compI, compJ, startI, startJ)
	}

	return changed
}

func (g *TwoZeroFourEight) Move(move rune) bool {
	done := false
	switch move {
	case 'l':
		done = g.moveLeft()
	case 'L':
		done = g.moveLeft()
	case 'r':
		done = g.moveRight()
	case 'R':
		done = g.moveRight()
	case 'u':
		done = g.moveUp()
	case 'U':
		done = g.moveUp()
	case 'd':
		done = g.moveDown()
	case 'D':
		done = g.moveDown()
	default:
		done = false
	}
	if done {
		fmt.Println("Move done")
		g.FillRandom()
	}
	return done
}
