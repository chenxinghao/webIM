package services

import (
	"errors"
	"webIM/cmd"
)

var Cmd *CmdHandler

type CmdHandler struct {
	CmdMap map[string]func() (string, error)
}

func init() {
	Cmd = &CmdHandler{make(map[string]func() (string, error))}
	Cmd.CmdMap["Dice"] = (&cmd.CmdDice{}).CmdHandler
}

func (this *CmdHandler) Exec(cmd string) (string, error) {
	if i, ok := this.CmdMap[cmd]; ok {
		return i()
	}
	return "failed", errors.New("can't find cmd")
}
