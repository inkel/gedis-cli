package main

import (
	"fmt"
	"github.com/inkel/gedis/client"
	"github.com/kless/term/readline"
	"os"
	"strings"
)

func perror(err error) {
	fmt.Printf("\033[0;31m%s\033[0m\r\n", err)
}

func tr(split []string) []interface{} {
	r := make([]interface{}, len(split))
	for i, s := range split {
		r[i] = s
	}
	return r
}

func pr(indent string, res interface{}) {
	var color, format string

	switch res.(type) {
	case int64:
		color = "\033[1;33m"
		format = "%d"
	case string:
		color = "\033[1;34m"
		format = "%q"
	default:
		if res != nil {
			if arr, ok := res.([]interface{}); ok {
				for i, d := range arr {
					pr(fmt.Sprintf("%s%d) ", indent, i), d)
				}
			} else {
				fmt.Printf("Unexpected! %#v\r\n", res)
			}
		} else {
			fmt.Printf("\033[0;36mnil\033[0m\r\n")
		}
		return
	}

	fmt.Printf(fmt.Sprintf("%s%s%s\033[0m\r\n", indent, color, format), res)
}

func main() {
	c, err := client.Dial("tcp", ":6379")

	if err != nil {
		perror(err)
		os.Exit(1)
	}

	ln, err := readline.NewDefaultLine(nil)

	if err != nil {
		os.Exit(2)
	}

	for {
		err = ln.Prompt()

		if err != nil {
			perror(err)
			break
		}

		line, err := ln.Read()

		if err != nil {
			if err != readline.ErrCtrlD {
				perror(err)
			}
			break
		} else {
			args := tr(strings.Split(line, " "))

			res, err := c.Send(args...)

			if err != nil {
				perror(err)
				continue
			}

			pr("", res)
		}
	}

	ln.Restore()
}
