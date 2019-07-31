package main

import (
	"fmt"
	"io"
	"os"

	"github.com/arcanericky/pushover"
	"github.com/spf13/cobra"
)

var messageCmd *cobra.Command

func outputMessageRequest(r pushover.MessageRequest) {
	pushoverURLText := "Pushover URL"
	pushoverURLTextLen := len(pushoverURLText)

	fields := []struct {
		field string
		value string
	}{
		{field: pushoverURLText, value: r.PushoverURL},
		{field: "Token", value: r.Token},
		{field: "User", value: r.User},
		{field: "Message", value: r.Message},
		{field: "Title", value: r.Title},
		{field: "URL", value: r.URL},
		{field: "URL Title", value: r.URLTitle},
		{field: "HTML", value: r.HTML},
		{field: "Monospace", value: r.Monospace},
		{field: "Sound", value: r.Sound},
		{field: "Device", value: r.Device},
		{field: "Priority", value: r.Priority},
		{field: "Timestamp", value: r.Timestamp},
	}

	for _, i := range fields {
		if len(i.value) > 0 {
			fmt.Printf("%-*s %s\n", pushoverURLTextLen, i.field+":", i.value)
		}
	}
}

func outputMessageResponse(r pushover.MessageResponse) {
	statusCodeText := "HTML Status Code:"
	maxLen := len(statusCodeText)
	fmt.Printf("%-*s %s\n", maxLen, "HTML Status:", r.HTTPStatus)
	fmt.Printf("%-*s %d\n", maxLen, statusCodeText, r.HTTPStatusCode)
	fmt.Printf("%-*s %d\n", maxLen, "API Status:", int(r.APIStatus))
	fmt.Printf("%-*s %s\n", maxLen, "Request ID:", r.Request)

	if len(r.Receipt) > 0 {
		fmt.Printf("%-*s %s\n", maxLen, "Receipt:", r.Receipt)
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

func addMessageCmd(parentCmd *cobra.Command) {
	const enable = "1"
	var token, user, title, message, url, urlTitle, sound, device, image,
		priority, timestamp, pushoverURL, htmlField, monospaceField string
	var html, monospace bool
	var imageReader io.ReadCloser
	var err error

	messageCmd = &cobra.Command{
		Use:   "message",
		Short: "Submit a message request",
		Long: `Send a Pushover message to a user or group.

Required options are:
  --token
  --user
  --message
`,
		Run: func(cmd *cobra.Command, args []string) {
			if html == true {
				htmlField = enable
			}

			if monospace == true {
				monospaceField = enable
			}

			if len(image) > 0 {
				imageReader, err = os.Open(image)
				if err != nil {
					fmt.Println("Error opening image:", err)
					return
				}
				defer imageReader.Close()
			}

			request := pushover.MessageRequest{
				PushoverURL: pushoverURL,
				Token:       token,
				User:        user,
				Message:     message,
				Title:       title,
				URL:         url,
				URLTitle:    urlTitle,
				HTML:        htmlField,
				Monospace:   monospaceField,
				Sound:       sound,
				Device:      device,
				Priority:    priority,
				Timestamp:   timestamp,
				ImageReader: imageReader,
				ImageName:   image,
			}

			fmt.Println("Request")

			outputMessageRequest(request)

			r, e := pushover.Message(request)

			fmt.Println()
			fmt.Println("Response")

			if e == nil {
				outputMessageResponse(*r)
			} else {
				fmt.Println(e)
			}
		},
	}

	// Required options
	messageCmd.Flags().StringVarP(&token, optionToken, "t", "", "Application's API token")
	messageCmd.MarkFlagRequired(optionToken)
	messageCmd.Flags().StringVarP(&user, optionUser, "u", "", "User/Group key")
	messageCmd.MarkFlagRequired(optionUser)
	messageCmd.Flags().StringVarP(&message, optionMessage, "m", "", "Notification message")
	messageCmd.MarkFlagRequired(optionMessage)

	// Optional options
	messageCmd.Flags().StringVarP(&pushoverURL, optionPushoverURL, "", "", "Pushover API URL")
	messageCmd.Flags().StringVarP(&title, optionTitle, "", "", "Message title (if empty, uses app name)")
	messageCmd.Flags().StringVarP(&url, optionURL, "", "", "Supplementary URL to show with the message")
	messageCmd.Flags().StringVarP(&urlTitle, optionURLTitle, "", "", "Title for the URL")
	messageCmd.Flags().BoolVarP(&html, optionHTML, "", false, "Enable HTML formatting")
	messageCmd.Flags().BoolVarP(&monospace, optionMonospace, "", false, "Enable monospace formatting")
	messageCmd.Flags().StringVarP(&sound, optionSound, "", "", "Name of a sound to override user's default")
	messageCmd.Flags().StringVarP(&image, optionImage, "", "", "Image attachment")
	messageCmd.Flags().StringVarP(&device, optionDevice, "", "", "Device name for message")
	messageCmd.Flags().StringVarP(&priority, optionPriority, "", "", "Message priority")
	messageCmd.Flags().StringVarP(&timestamp, optionTimestamp, "", "", "Unix timestamp for message")

	parentCmd.AddCommand(messageCmd)
}
