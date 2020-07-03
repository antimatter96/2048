package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/antimatter96/2048/game"
)

func main() {
	game := game.NewGame()
	fmt.Println(game.Print())
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < 20; i++ {
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Got an ", char)
		switch char {
		case 'E':
			break
		default:
			fmt.Println(game.Move(char))
			break
		}
		fmt.Println(game.Print())

		char, _, err = reader.ReadRune()

		if char == 10 {
			continue
		}
	}
}
