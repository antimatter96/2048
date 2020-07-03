package game

type game interface {
	Move(byte)
	End()
	Lost()
	Won()
	Print() string
}
