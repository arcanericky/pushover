package main

import (
	"fmt"

	"github.com/arcanericky/pushover"
	"github.com/spf13/cobra"
)

var validateCmd *cobra.Command

func outputValidateRequest(r pushover.ValidateRequest) {
	pushoverURLText := "Pushover URL"
	pushoverURLTextLen := len(pushoverURLText)

	fields := []struct {
		field string
		value string
	}{
		{field: pushoverURLText, value: r.PushoverURL},
		{field: "Token", value: r.Token},
		{field: "User", value: r.User},
	}

	for _, i := range fields {
		if len(i.value) > 0 {
			fmt.Printf("%-*s %s\n", pushoverURLTextLen, i.field+":", i.value)
		}
	}
}

func outputValidateResponse(r pushover.ValidateResponse) {
	statusCodeText := "HTML Status Code:"
	maxLen := len(statusCodeText)
	fmt.Printf("%-*s %s\n", maxLen, "HTML Status:", r.HTTPStatus)
	fmt.Printf("%-*s %d\n", maxLen, statusCodeText, r.HTTPStatusCode)
	fmt.Printf("%-*s %d\n", maxLen, "API Status:", int(r.APIStatus))
	fmt.Printf("%-*s %s\n", maxLen, "Request ID:", r.Request)
	fmt.Printf("%-*s %d\n", maxLen, "Group:", r.Group)

	fmt.Println("Licenses:")
	for _, v := range r.Licenses {
		fmt.Println(" ", v)
	}

	fmt.Println("Devices:")
	for _, v := range r.Devices {
		fmt.Println(" ", v)
	}

	if len(r.ErrorParameters) > 0 {
		maxLen := 0
		for k := range r.ErrorParameters {
			curLen := len(k)
			if curLen > maxLen {
				maxLen = curLen
			}
		}
		maxLen++

		fmt.Println("Parameter Errors:")
		for k, v := range r.ErrorParameters {
			fmt.Printf("  %-*s %s\n", maxLen, k+":", v)
		}
	}

	if len(r.Errors) > 0 {
		fmt.Println("Errors:")
		for _, v := range r.Errors {
			fmt.Println(" ", v)
		}
	}

	fmt.Println("Response Body:", r.ResponseBody)
}

func addValidateCmd(parentCmd *cobra.Command) {
	var token, user, device, pushoverURL string

	validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Submit a validate request",
		Long: `Validate a Pushover user or group, and optionally a
device name.

Required options are:
  --token
  --user
`,
		Run: func(cmd *cobra.Command, args []string) {
			request := pushover.ValidateRequest{
				PushoverURL: pushoverURL,
				Token:       token,
				User:        user,
				Device:      device,
			}

			fmt.Println("Request")

			outputValidateRequest(request)

			r, e := pushover.Validate(request)

			fmt.Println()
			fmt.Println("Response")

			if e == nil {
				outputValidateResponse(*r)
			} else {
				fmt.Println(e)
			}
		},
	}

	// Required options
	validateCmd.Flags().StringVarP(&token, optionToken, "t", "", "Application's API token")
	validateCmd.MarkFlagRequired(optionToken)
	validateCmd.Flags().StringVarP(&user, optionUser, "u", "", "User/Group key")
	validateCmd.MarkFlagRequired(optionUser)

	// Optional options
	validateCmd.Flags().StringVarP(&pushoverURL, optionPushoverURL, "", "", "Pushover API URL")
	validateCmd.Flags().StringVarP(&device, optionDevice, "", "", "Device name to validate")

	parentCmd.AddCommand(validateCmd)
}
