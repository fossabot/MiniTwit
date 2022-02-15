package database

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"io/ioutil"
	"log"
	"minitwit/src/models"
	"strconv"
	"strings"
)

// const DATABASE = "/tmp/minitwit.db"
// const DATABASE = "C:/Users/hardk/source/repos/MiniTwit/minitwit.db"
//const DATABASE = "/home/turbo/ITU/DevOps/MiniTwit/tmp/minitwit.db"
const DATABASE = "C:\\Users\\JTT\\Documents\\git\\MiniTwit\\minitwit.db"

//const DATABASE = "H:/repos/MiniTwit/minitwit.db"

const PER_PAGE = 30
const DEBUG = true
const SECRET_KEY = "development key"

/****************************************
*			DATABASE RELATED			*
****************************************/
func ConnectDb() *sql.DB {
	db, err := sql.Open("sqlite3", DATABASE)
	if err != nil {
		panic(err)
	}

	return db
}

// setup
func InitDb() {
	db := ConnectDb()
	query, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}

// example Database usage
func GetUserMessages(id int) []models.Message {
	db := ConnectDb()
	query := string(`SELECT 
		message.message_id, 
		message.author_id, 
		user.username, 
		message.text, 
		message.pub_date, 
		user.email 
		FROM message, user 
		WHERE message.flagged = 0 AND 
		user.user_id = (?) AND
		user.user_id = message.author_id
		ORDER BY message.pub_date DESC 
		LIMIT 30`)
	result, err := db.Query(query, fmt.Sprint(id), fmt.Sprint(id))
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var messages []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			panic(err.Error())
		}
		messages = append(messages, msg)
	}
	return messages
}

func GetAllMessages() []models.Message {
	db := ConnectDb()
	query := string("select message.message_id , message.author_id , user.username , message.text , message.pub_date ,  user.email from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit 30")
	result, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var messages []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			panic(err.Error())
		}
		messages = append(messages, msg)
	}
	return messages
}

func AddUserToDb(username string, email string, password string) {
	db := ConnectDb()
	salt := make([]byte, 4)
	io.ReadFull(rand.Reader, salt)

	pwIteration_int, _ := strconv.Atoi("50000")
	dk := pbkdf2.Key([]byte(password), salt, pwIteration_int, 32, sha256.New)

	pw_hashed := "pbkdf2:sha256:50000$" + string(salt) + "$" + hex.EncodeToString(dk)
	query, err := db.Prepare("INSERT INTO user(username, email, pw_hash) values (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(username, email, pw_hashed)

	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()
}

func GetUserFromDb(username string) models.User {
	db := ConnectDb()
	//TODO: Prepared statements
	strs := []string{"SELECT x.* FROM 'user' x WHERE username like '", username, "'"}
	query := strings.Join(strs, "")
	row, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var user models.User
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)
	}

	return user

}
