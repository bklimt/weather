package weather

import (
	"encoding/json"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"io/ioutil"
	"log"
)

var db mysql.Conn

var config struct {
	Database string
	Username string
	Password string
}

// Reads in the database config for the server and connects to the database.
// Kills the server if unsuccessful.
func LoadConfig(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to open config: ", err)
	}

	if err = json.Unmarshal(b, &config); err != nil {
		log.Fatal("Unable to parse config: ", err)
	}

	LoadDb()
}

func LoadDb() {
	db = mysql.New("tcp", "", "127.0.0.1:3306", config.Username, config.Password, config.Database)
	if err := db.Connect(); err != nil {
		log.Fatal("Unable to open database: ", err)
	}

	log.Printf("Opened database %s", config.Database)
}
