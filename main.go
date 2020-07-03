package main

import (
	"fmt"

	"github.com/antimatter96/2048/game"
)

func main() {
	x := &game.TwoZeroFourEight{}
	x.Init()
	x.Print()
	fmt.Println("")
	x.FillRandom()
	x.Print()
	fmt.Println("")
	x.MoveLeft()
	x.Print()
	fmt.Println("")

	x.FillRandom()
	x.Print()
	fmt.Println("")
	x.MoveLeft()
	x.Print()
	fmt.Println("")

	x.FillRandom()
	x.Print()
	fmt.Println("")
	x.MoveLeft()
	x.Print()
	fmt.Println("")

	x.FillRandom()
	x.Print()
	fmt.Println("")
	x.MoveLeft()
	x.Print()
	fmt.Println("")

}
