package cmd

type CmdInterface interface {
	CmdHandler() (string, error)
	CmdHandlerWithoutResult() error
}
