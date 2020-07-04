package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/antimatter96/2048/game"
	"go.uber.org/zap"
)

const enter rune = 10

func main() {

	logger, _ := zap.NewProduction()

	gameInstance := game.NewGame(logger)
	fmt.Println(gameInstance.Print())
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < 20; i++ {
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Got a", char)
		switch char {
		case 'E':
			break
		default:
			fmt.Println(gameInstance.Move(char))
		}
		fmt.Println(gameInstance.Print())

		char, _, _ = reader.ReadRune()

		if char == enter {
			continue
		}
	}
}
