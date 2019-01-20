package services

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	if i, ok := Cmd.CmdMap["Dice"]; ok {
		fmt.Println(i())
	}
}
