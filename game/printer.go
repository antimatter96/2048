package game

// Copied from https://en.wikipedia.org/wiki/ANSI_escape_code and
// https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux

const colorOff string = "\033[0m" // Text Reset

var paddedText = []string{
	"      ",
	"   2  ",
	"   4  ",
	"   8  ",
	"  16  ",
	"  32  ",
	"  64  ",
	"  128 ",
	"  256 ",
	"  512 ",
	" 1024 ",
	" 2048 ",
}

var colorCodes = []string{
	"0",
	"226",
	"190",
	"220",
	"184",
	"214",
	"178",
	"208",
	"172",
	"202",
	"166",
	"196",
}
