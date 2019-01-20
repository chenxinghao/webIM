package cmd

import (
	"github.com/chenxinghao/gotools/Util/RandomUtils"
	"strconv"
)

type CmdDice struct {
}

func (this *CmdDice) CmdHandler() (string, error) {
	randomNumber := RandomUtils.RandomNumber{}
	num, err := randomNumber.CryptoRandInt(1, 7)
	return strconv.Itoa(num), err
}

func (this *CmdDice) CmdHandlerWithoutResult() error {
	return nil
}
