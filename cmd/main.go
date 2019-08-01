package main

import (
	"github.com/spf13/cobra"
)

const (
	optionDevice      = "device"
	optionExpire      = "expire"
	optionHTML        = "html"
	optionImage       = "image"
	optionMessage     = "message"
	optionMonospace   = "monospace"
	optionPriority    = "priority"
	optionPushoverURL = "pushoverurl"
	optionRetry       = "retry"
	optionSound       = "sound"
	optionTimestamp   = "timestamp"
	optionTitle       = "title"
	optionToken       = "token"
	optionURL         = "url"
	optionURLTitle    = "urltitle"
	optionUser        = "user"
)

var versionText string

func main() {
	rootCmd := &cobra.Command{
		Use:   "pushover",
		Short: "Pushover",
		Long: `Pushover CLI version ` + versionText + `

Submit various requests to the Pushover API. Currently only
message (notification) and validate are supported.

See the README at https://github.com/arcanericky/pushover for
more information. For details on Pushover, see
https://pushover.net/.`,
		Version: versionText,
	}

	addMessageCmd(rootCmd)
	addValidateCmd(rootCmd)

	rootCmd.Execute()
}
