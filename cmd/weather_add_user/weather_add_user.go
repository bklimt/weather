package main

import (
	"bufio"
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/gcmurphy/getpass"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
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
	dbname := flag.String("db", "weather", "Name of the database to connect to.")
	flag.Parse()

	user, err := prompt("Database username: ", false, false)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	pass, err := prompt("Database password: ", true, false)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, *dbname)
	if err := db.Connect(); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

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

	salt := make([]byte, 20)
	// There's no need to seed with rand.Read.
	_, err = rand.Read(salt)
	if err != nil {
		fmt.Printf("Unable to generate salt: %v\n", err)
		return
	}
	// fmt.Printf("salt: %v\n", salt)

	b := append(salt, password...)
	bcpw, err := bcrypt.GenerateFromPassword(b, 0)
	if err != nil {
		fmt.Printf("Unable to bcrypt password: %v\n", err)
		return
	}

	stmt, err := db.Prepare("INSERT INTO user (name, salt, bcrypt) values (?, ?, ?)")
	if err != nil {
		fmt.Printf("Unable to prepare statement: %v\n", err)
		return
	}

	stmt.Bind(username, salt, bcpw)
	_, err = stmt.Run()
	if err != nil {
		fmt.Printf("Unable to insert row into database: %v\n", err)
		return
	}
}
