package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// N : N * N grid
const N = 4
const nDashed = (6 * N) + N + 1 + (4 * N) // 6 : width, 4 : 2 + 2 (the spaces around each pipe)
const minStringSize = nDashed * ((2 * N) - 1)

// TwoZeroFourEight holds the state of the board,
// the random seed and the colors mapping
type TwoZeroFourEight struct {
	board  [][]int
	rand   *rand.Rand
	logger *zap.Logger

	colorCodes map[int]string
	paddedText map[int]string
}

// NewGame initializes a new 2048 game
func NewGame(logger *zap.Logger) *TwoZeroFourEight {
	g := &TwoZeroFourEight{logger: logger}
	g.init()

	return g
}

func (g *TwoZeroFourEight) init() {
	// Create random source
	source := rand.NewSource(time.Now().UnixNano())
	g.rand = rand.New(source)

	// Initialize the 2D array
	g.board = make([][]int, N)
	for i := 0; i < N; i++ {
		g.board[i] = make([]int, N)
	}

	// Fill bottom 3 cells with random numbers
	g.board[N-1][0] = 2 << g.rand.Intn(2)
	g.board[N-1][1] = 2 << g.rand.Intn(2)
	g.board[N-2][0] = 2 << g.rand.Intn(2)

	// Create a number to its text representation mapping
	g.paddedText = make(map[int]string)
	g.paddedText[0] = paddedText[0]
	for i, j := 2, 1; i < 2048+1; i, j = i*2, j+1 {
		g.paddedText[i] = paddedText[j]
	}

	// Create a number to its color representation mapping
	// 8-bit colors
	g.colorCodes = make(map[int]string)
	g.colorCodes[0] = "\033[38;5;" + colorCodes[0] + "m"
	for i, j := 2, 1; i < 2048+1; i, j = i*2, j+1 {
		g.colorCodes[i] = "\033[38;5;" + colorCodes[j] + "m"
	}
}

func (g *TwoZeroFourEight) getString(i int) string {
	return fmt.Sprintf("%s%s%s", g.colorCodes[i], g.paddedText[i], colorOff)
}

// Print returns a formatted visual rep of the board
func (g *TwoZeroFourEight) Print() string {
	var b strings.Builder
	b.Grow(minStringSize)

	b.WriteString(strings.Repeat("-", nDashed))
	b.WriteByte('\n')
	for i := 0; i < N; i++ {
		b.WriteString("|  ")
		for j := 0; j < N-1; j++ {
			b.WriteString(fmt.Sprintf("%s  |  ", g.getString(g.board[i][j])))
		}
		b.WriteString(fmt.Sprintf("%s  |\n", g.getString(g.board[i][N-1])))
		b.WriteString(strings.Repeat("-", nDashed))
		b.WriteByte('\n')
	}

	return b.String()
}

func (g *TwoZeroFourEight) emptyCells() (int, int, error) {
	mp := make(map[string]bool)

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if g.board[i][j] == 0 {
				mp[fmt.Sprintf("%d,%d", i, j)] = true
				break
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
	x, y, err := g.emptyCells()
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

func (g *TwoZeroFourEight) moveHorizontal(changeI, changeJ mover, compI, compJ comp, startI, startJ int, combine bool) bool {
	g.logger.Debug("Call to moveHorizontal")

	changed := false

	for j := startJ; compJ(j); j = changeJ(j) {
		for i := startI; compI(i) && compJ(j); i = changeI(i) {
			if g.board[i][j] == 0 {
				continue
			}
			if combine && g.board[i][j] == g.board[i][changeJ(j)] {
				g.logger.Debug("changing",
					zap.Int("x", i),
					zap.Int("y", changeJ(j)),
					zap.Int("value", g.board[i][j]),
					zap.Int("new value", g.board[i][j]*2),
				)
				g.board[i][j] = 0
				g.board[i][changeJ(j)] *= 2
				j = changeJ(j)
				changed = true
			} else if g.board[i][changeJ(j)] == 0 {
				g.logger.Debug("moving",
					zap.Int("x", i),
					zap.Int("y", changeJ(j)),
					zap.Int("new value", g.board[i][j]),
				)
				g.board[i][changeJ(j)] = g.board[i][j]
				g.board[i][j] = 0
				changed = true
			}
		}
	}

	if changed {
		temp := g.moveHorizontal(changeI, changeJ, compI, compJ, startI, startJ, false)
		changed = changed || temp
	}

	return changed
}

func (g *TwoZeroFourEight) moveLeft() bool {
	greaterThanZero := greaterThan(0)
	lessThanN := lessThan(N)

	return g.moveHorizontal(inc, dec, lessThanN, greaterThanZero, 0, N-1, true)
}

func (g *TwoZeroFourEight) moveRight() bool {
	lessThanN := lessThan(N)
	lessThanNMinusOne := lessThan(N - 1)

	return g.moveHorizontal(inc, inc, lessThanN, lessThanNMinusOne, 0, 0, true)
}

func (g *TwoZeroFourEight) moveVertical(changeI, changeJ mover, compI, compJ comp, startI, startJ int, combine bool) bool {
	g.logger.Debug("Call to moveVertical")

	changed := false

	for i := startI; compI(i); i = changeI(i) {
		for j := startJ; compJ(j) && compI(i); j = changeJ(j) {
			if g.board[i][j] == 0 {
				continue
			}
			if combine && g.board[i][j] == g.board[changeI(i)][j] {
				g.logger.Debug("changing",
					zap.Int("x", changeI(i)),
					zap.Int("y", j),
					zap.Int("value", g.board[i][j]),
					zap.Int("new value", g.board[i][j]*2),
				)
				g.board[i][j] = 0
				g.board[changeI(i)][j] *= 2
				i = changeI(i)
				changed = true
			} else if g.board[changeI(i)][j] == 0 {
				g.logger.Debug("moving",
					zap.Int("x", changeI(i)),
					zap.Int("y", j),
					zap.Int("new value", g.board[i][j]),
				)
				g.board[changeI(i)][j] = g.board[i][j]
				g.board[i][j] = 0
				changed = true
			}
		}
	}

	if changed {
		temp := g.moveVertical(changeI, changeJ, compI, compJ, startI, startJ, false)
		changed = changed || temp
	}

	return changed
}

func (g *TwoZeroFourEight) moveDown() bool {
	lessThanNMinus1 := lessThan(N - 1)
	lessThanN := lessThan(N)

	return g.moveVertical(inc, inc, lessThanNMinus1, lessThanN, 0, 0, true)
}

func (g *TwoZeroFourEight) moveUp() bool {
	greaterThan := greaterThan(0)
	lessThanN := lessThan(N)

	return g.moveVertical(dec, inc, greaterThan, lessThanN, N-1, 0, true)
}

// Move is used by the controller
func (g *TwoZeroFourEight) Move(move rune) bool {
	var done bool
	switch move {
	case 'A':
		fallthrough
	case 'a':
		done = g.moveLeft()
	case 'D':
		fallthrough
	case 'd':
		done = g.moveRight()
	case 'W':
		fallthrough
	case 'w':
		done = g.moveUp()
	case 'S':
		fallthrough
	case 's':
		done = g.moveDown()
	}

	if done {
		g.logger.Debug("Move done")
		g.fillRandom()
	}

	return done
}
