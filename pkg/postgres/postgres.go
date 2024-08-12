package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

// for /registrarion
func InsertDb(login string, password string) error {
	// Connect
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), `INSERT INTO accounts (login, password_users) VALUES ($1, $2)`, login, password)
	if err != nil {
		return err
	}
	return nil
}

func CheckPassword(login string) string {
	//Connect
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	conn, err := pgx.Connect(context.Background(), urlToDataBase)
	if err != nil {
		log.Println(err.Error())
	}
	defer conn.Close(context.Background())

	// Get data
	var passwordForCompare string = ""
	//command := fmt.Sprintf(`SELECT (password_users) FROM %s WHERE login = $1)`, Cfg.PGnameTable)
	//err = conn.QueryRow(context.Background(), `SELECT (password_users) FROM accounts WHERE login = $1`, login).Scan(&passwordForCompare)
	command := fmt.Sprintf(`SELECT password_users FROM %s WHERE login = $1`, Cfg.PGnameTable)
	err = conn.QueryRow(context.Background(), command, login).Scan(&passwordForCompare)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println("password from db:", passwordForCompare)
	return passwordForCompare // if exist => return password for compare otherwise return ""
}

func init() {
	file, err := os.Open("config.cfg")
	if err != nil {
		fmt.Println("Error open .cfg", err)
		panic("Can't open the file \"setting.cfg\"")
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	readSetting := make([]byte, fileInfo.Size())
	_, err = file.Read(readSetting)
	if err != nil {
		panic("can't read file")
	}
	// fmt.Println(string(readSetting))  работает

	err = json.Unmarshal(readSetting, &Cfg)
	if err != nil {
		panic("json err")
	}
}

type setting struct { // должен повторять структуру json
	PGaddress   string
	PGpassword  string
	PGuser      string
	PGdbname    string
	PGPort      string
	PGnameTable string
}

var (
	Cfg setting // for use in main for open db
)
