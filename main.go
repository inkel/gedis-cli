package main

import (
	"fmt"
	"github.com/inkel/gedis/client"
	"github.com/kless/term/readline"
	"os"
	"strings"
)

const (
	red    = "\033[1;33m"
	yellow = "\033[1;33m"
	blue   = "\033[1;34m"
	cyan   = "\033[0;36m"
	reset  = "\033[0m"

	cnil = cyan + "nil" + reset
)

func perror(err error) {
	fmt.Println(cerror(err))
}

func cerror(err error) string {
	return red + err.Error() + reset
}

func cint(n int64) string {
	return fmt.Sprintf("%s%d%s", yellow, n, reset)
}

func cstring(str string) string {
	return blue + str + reset
}

func tr(split []string) []interface{} {
	r := make([]interface{}, len(split))
	for i, s := range split {
		r[i] = s
	}
	return r
}

func pr(indent string, res interface{}) {
	var out string

	switch res.(type) {
	case int64:
		out = cint(res.(int64))
	case string:
		out = cstring(res.(string))
	default:
		if res != nil {
			if arr, ok := res.([]interface{}); ok {
				if len(arr) > 0 {
					for i, d := range arr {
						pr(fmt.Sprintf("%s%d) ", indent, i), d)
					}
				} else {
					fmt.Printf("%s%s(empty)%s\n", indent, cyan, reset)
				}
			} else {
				fmt.Printf("Unexpected! %#v\r\n", res)
			}
		} else {
			fmt.Println(indent + cnil)
		}
		return
	}

	fmt.Printf("%s%s\n", indent, out)
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
