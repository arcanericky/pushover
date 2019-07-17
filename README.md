# Pushover

Go package for the [Pushover](https://pushover.net/) API.

[![Build Status](https://travis-ci.com/arcanericky/pushover.svg?branch=master)](https://travis-ci.com/arcanericky/pushover)
[![codecov](https://codecov.io/gh/arcanericky/pushover/branch/master/graph/badge.svg)](https://codecov.io/gh/arcanericky/pushover)
[![GoDoc](https://img.shields.io/badge/docs-GoDoc-brightgreen.svg)](https://godoc.org/github.com/arcanericky/pushover)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

## About

This package enables your Go application to send requests to the [Pushover](https://pushover.net/) service through the [Pushover REST API](https://pushover.net/api). Implementing a notification message in your Go code is as simple as `pushover.Message(pushover.MessageRequest{User: "user", Token: "token", Message: "message"})`. It's just a straightforward function call - no fancy methods attached to functions, just populate a structure and [Go](https://golang.org/).

Note that Pushover has many APIs available, but currently this package only supports:
-  [Messages](https://pushover.net/api#messages)
-  [User/Group Validation](https://pushover.net/api#validate)

The Pushover service is a great way to send notifications to your device for any purpose. The device application is [free for 7 days](https://pushover.net/faq#overview-fees), after which you must purchase it for a one-time price of $4.99 per platform. It comes with a [7,500 message per month limit](https://pushover.net/faq#overview-limits) with the [ability to pay for more messages](https://pushover.net/faq#overview-usage).


## Using the Package

Obtaining a Pushover account and adding an application for using this library are not covered in this README. Pushover API and user tokens are required.

To use this Pushover package, just import it and make a single call to `pushover.Message`.

```
$ cat > main.go << EOF
package main

import (
	"context"
	"fmt"
	"os"
	"github.com/arcanericky/pushover"
)

func main() {
	var r *pushover.MessageResponse
	var e error
	if r, e = pushover.Message(pushover.MessageRequest{
		Token: os.Args[1], User: os.Args[2], Message: os.Args[3]},
	); e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(r)
}
EOF
$ GO111MOD=on go mod init demo
$ GO111MOD=on go build
$ ./demo api-token user-token "Test Message"
```

## Using the Utility

A simple application to demonstrate and test the Pushover package is included with this repository in [Released executables](https://github.com/arcanericky/pushover/releases) and is useful on its own. While using Pushover via [`curl`](https://curl.haxx.se/) is simple enough, this utility makes it even easier.

```
$ pushover message --token token --user user --message message
```

For help, use the `--help` option.

```
$ pushover --help
Pushover CLI version 1.0.0

Submit various requests to the Pushover API. Currently only
message (notification) and validate are supported.

See the README at https://github.com/arcanericky/pushover for
more information. For details on Pushover, see
https://pushover.net/.

Usage:
  pushover [command]

Available Commands:
  help        Help about any command
  message     Submit a message request
  validate    Submit a validate request

Flags:
  -h, --help      help for pushover
      --version   version for pushover

Use "pushover [command] --help" for more information about a command.
```

And for help with the various subcommands, issue the subcommand followed by `--help`.
```
$ pushover message --help
Send a Pushover message to a user or group.

Required options are:
  --token
  --user
  --message

Usage:
  pushover message [flags]

Flags:
      --device string        Device name for message
  -h, --help                 help for message
      --html                 Enable HTML formatting
      --image string         Image attachment
  -m, --message string       Notification message
      --monospace            Enable monospace formatting
      --priority string      Message priority
      --pushoverurl string   Pushover API URL
      --sound string         Name of a sound to override user's default
      --timestamp string     Unix timestamp for message
      --title string         Message title (if empty, uses app name)
  -t, --token string         Application's API token
      --url string           Supplementary URL to show with the message
      --urltitle string      Title for the URL
  -u, --user string          User/Group key
```

## Contributing

Contributions and corrections are welcome. If you're adding a feature, please submit an issue so it can be discussed and ensure the work isn't being duplicated. Unit test coverage is required for all pull requests.

Some features that are not implemented but would be welcome:
  
- Implement other Pushover APIs
  - [Receipt and Callback](https://pushover.net/api/receipts)
  - [Subscription](https://pushover.net/api/subscriptions)
  - [Groups](https://pushover.net/api/groups)
  - [Glances](https://pushover.net/api/groups)
  - [Licensing](https://pushover.net/api/licensing)
  - [Open Client](https://pushover.net/api/client)
  - [Limits](https://pushover.net/api#limits)
- Use of environment variables for API token in the CLI

## Inspiration

During my development workday I often find I need to monitor state of systems, programs, networks, etc. Eventually I concluded the Pushover service was best for me and I've been using it in some monitoring scripts using `curl`. One day I found a need to do this in one of my Go utilities and that's when this code was started.