package u_test

import (
	"fmt"
	"github.com/ssgo/u"
	"testing"
)

func TestShell(t *testing.T) {

	fmt.Print("  ", u.Black("Abc"))
	fmt.Print("  ", u.Red("Abc"))
	fmt.Print("  ", u.Green("Abc"))
	fmt.Print("  ", u.Yellow("Abc"))
	fmt.Print("  ", u.Blue("Abc"))
	fmt.Print("  ", u.Magenta("Abc"))
	fmt.Print("  ", u.Cyan("Abc"))
	fmt.Print("  ", u.White("Abc"))
	fmt.Println()

	fmt.Print("  ", u.BBlack("Abc"))
	fmt.Print("  ", u.BRed("Abc"))
	fmt.Print("  ", u.BGreen("Abc"))
	fmt.Print("  ", u.BYellow("Abc"))
	fmt.Print("  ", u.BBlue("Abc"))
	fmt.Print("  ", u.BMagenta("Abc"))
	fmt.Print("  ", u.BCyan("Abc"))
	fmt.Print("  ", u.BWhite("Abc"))
	fmt.Println()

	for j := u.AttrNone; j < u.AttrCrossedOut; j++ {
		for i := u.TextBlack; i < u.TextWhite; i++ {
			fmt.Print("  ", u.Color(" Abc ", i, u.BgNone, j))
		}
		fmt.Println()
	}

	for j := u.BgBlack; j < u.BgWhite; j++ {
		for i := u.TextBlack; i < u.TextWhite; i++ {
			fmt.Print("  ", u.Color(" Abc ", i, j))
		}
		fmt.Println()
	}
}
