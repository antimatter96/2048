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
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome")
	fmt.Println(gameInstance.Print())
	for i := 0; i < 20; i++ {
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("You entered", string(char))
		switch char {
		case 'E':
			break
		default:
			gameInstance.Move(char)
		}
		fmt.Println(gameInstance.Print())

		char, _, _ = reader.ReadRune()

		if char == enter {
			continue
		}
	}
}
