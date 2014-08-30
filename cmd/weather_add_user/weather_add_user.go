// The weather_add_user command lets you add a new user who can log into the weather website.
package main

import (
	"bufio"
	"flag"
	"fmt"
  "github.com/bklimt/weather"
	"github.com/gcmurphy/getpass"
	"os"
	"strings"
)

func prompt(p string, hide bool, confirm bool) (string, error) {
	if hide {
		c := 0
		if confirm {
			c = 1
		}
		return getpass.GetPassWithOptions(p, c, 255)
	} else {
		fmt.Print(p)
		reader := bufio.NewReader(os.Stdin)
		s, err := reader.ReadString('\n')
		if err == nil {
			s = strings.TrimSpace(s)
		}
		return s, err
	}
}

func main() {
  config := flag.String("config", "./config.json", "File to load the database config from.")
	flag.Parse()

  weather.LoadConfig(*config)

	username, err := prompt("Username: ", false, false)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	password, err := prompt("Password: ", true, true)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

  user := weather.User{username}
  err = user.Save(password)
  if err != nil {
    fmt.Printf("Unable to save user: %v\n", err)
    return
  }
}
