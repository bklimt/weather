package weather

import (
	"encoding/json"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"io/ioutil"
	"log"
)

var db mysql.Conn

type config struct {
	Database string
	Username string
	Password string
}

func LoadConfig(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to open config: ", err)
	}

	info := config{}

	if err = json.Unmarshal(b, &info); err != nil {
		log.Fatal("Unable to parse config: ", err)
	}

	db = mysql.New("tcp", "", "127.0.0.1:3306", info.Username, info.Password, info.Database)
	if err := db.Connect(); err != nil {
		log.Fatal("Unable to open database: ", err)
	}

	log.Printf("Opened database %s", info.Database)
}
