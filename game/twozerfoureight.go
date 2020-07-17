package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

const N = 4
const nDashed = 1 + (2 * (N + (2 << (N - 1)))) // Not even true :p
const minStringSize = 350                      // Is 342 to be exqct :p

type TwoZeroFourEight struct {
	board  [][]int
	rand   *rand.Rand
	logger *zap.Logger
}

func NewGame(logger *zap.Logger) *TwoZeroFourEight {
	g := &TwoZeroFourEight{}
	g.init(logger)
	return g
}

func (g *TwoZeroFourEight) init(logger *zap.Logger) {
	source := rand.NewSource(time.Now().UnixNano())
	g.rand = rand.New(source)

	g.board = make([][]int, N)
	for i := 0; i < N; i++ {
		g.board[i] = make([]int, N)
	}

	g.board[N-1][0] = 2 << g.rand.Intn(2)
	g.board[N-1][1] = 2 << g.rand.Intn(2)
	g.board[N-2][0] = 2 << g.rand.Intn(2)

	g.logger = logger
}

func (g *TwoZeroFourEight) Print() string {
	var b strings.Builder
	b.Grow(minStringSize)

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
		return 0, 0, fmt.Errorf("%v\n%s", "FULL BOARD", g.Print())
	}

	picked := g.rand.Intn(len(arr))

	xy := strings.Split(arr[picked], ",")
	x, _ := strconv.Atoi(xy[0])
	y, _ := strconv.Atoi(xy[1])

	return x, y, nil
}

func (g *TwoZeroFourEight) fillRandom() {
	x, y, err := g.boundaryCells()
	if err != nil {
		g.logger.Error("Error while filling", zap.Error(err))
		return
	}

	toAdd := 2 << g.rand.Intn(2)
	g.board[x][y] = toAdd
	g.logger.Debug("Added field",
		zap.Int("added", toAdd),
		zap.Int("x", x),
		zap.Int("y", y),
	)
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
	g.logger.Debug("Call to moveHorizontal")

	changed := false
	moved := false

	for j := startJ; compJ(j); j = changeJ(j) {
		for i := startI; compI(i); i = changeI(i) {
			if g.board[i][j] == 0 {
				continue
			}
			if g.board[i][j] == g.board[i][changeJ(j)] {
				g.logger.Debug("changing",
					zap.Int("x", i),
					zap.Int("y", changeJ(j)),
					zap.Int("value", g.board[i][j]),
					zap.Int("new value", g.board[i][j]*2),
				)
				g.board[i][j] = 0
				g.board[i][changeJ(j)] *= 2
				changed = true
			} else if g.board[i][changeJ(j)] == 0 {
				g.logger.Debug("moving",
					zap.Int("x", i),
					zap.Int("y", changeJ(j)),
					zap.Int("new value", g.board[i][j]),
				)
				g.board[i][changeJ(j)] = g.board[i][j]
				g.board[i][j] = 0
				moved = true
			}
		}
	}

	if changed || moved {
		temp := g.moveHorizontal(changeI, changeJ, compI, compJ, startI, startJ)
		changed = changed || temp
	}

	return changed
}

func (g *TwoZeroFourEight) moveLeft() bool {
	greaterThanZero := greaterThan(0)
	lessThanN := lessThan(N)

	return g.moveHorizontal(inc, dec, lessThanN, greaterThanZero, 0, N-1)
}

func (g *TwoZeroFourEight) moveRight() bool {
	lessThanN := lessThan(N)
	lessThanNMinusOne := lessThan(N - 1)

	return g.moveHorizontal(inc, inc, lessThanN, lessThanNMinusOne, 0, 0)
}

func (g *TwoZeroFourEight) moveVertical(changeI, changeJ mover, compI, compJ comp, startI, startJ int) bool {
	g.logger.Debug("Call to moveVertical")

	changed := false
	moved := false

	for i := startI; compI(i); i = changeI(i) {
		for j := startJ; compJ(j); j = changeJ(j) {
			if g.board[i][j] == 0 {
				continue
			}
			if g.board[i][j] == g.board[changeI(i)][j] {
				g.logger.Debug("changing",
					zap.Int("x", changeI(i)),
					zap.Int("y", j),
					zap.Int("value", g.board[i][j]),
					zap.Int("new value", g.board[i][j]*2),
				)
				g.board[i][j] = 0
				g.board[changeI(i)][j] *= 2
				changed = true
			} else if g.board[changeI(i)][j] == 0 {
				g.logger.Debug("moving",
					zap.Int("x", changeI(i)),
					zap.Int("y", j),
					zap.Int("new value", g.board[i][j]),
				)
				g.board[changeI(i)][j] = g.board[i][j]
				g.board[i][j] = 0
				moved = true
			}
		}
	}

	if changed || moved {
		temp := g.moveVertical(changeI, changeJ, compI, compJ, startI, startJ)
		changed = changed || temp
	}

	return changed
}

func (g *TwoZeroFourEight) moveDown() bool {
	lessThanNMinus1 := lessThan(N - 1)
	lessThanN := lessThan(N)

	return g.moveVertical(inc, inc, lessThanNMinus1, lessThanN, 0, 0)
}

func (g *TwoZeroFourEight) moveUp() bool {
	greaterThan := greaterThan(0)
	lessThanN := lessThan(N)

	return g.moveVertical(dec, inc, greaterThan, lessThanN, N-1, 0)
}

// Move is used by the controller
func (g *TwoZeroFourEight) Move(move rune) bool {
	done := false
	switch move {
	case 'A':
		done = g.moveLeft()
	case 'a':
		done = g.moveLeft()
	case 'D':
		done = g.moveRight()
	case 'd':
		done = g.moveRight()
	case 'W':
		done = g.moveUp()
	case 'w':
		done = g.moveUp()
	case 'S':
		done = g.moveDown()
	case 's':
		done = g.moveDown()
	default:
		done = false
	}
	if done {
		g.logger.Debug("Move done")
		g.fillRandom()
	}
	return done
}
