package database

import (
	"errors"
	"fmt"
	"minitwit/models"
	"strings"
	"time"
)

func GormGetAllMessages() []models.Message {

	result, err := gormDb.
		Model(models.Message{}).
		Table("message").
		Order("pub_date desc").
		Limit(30).
		Where("flagged = ?", 0).
		Joins("JOIN user on message.author_id = user.user_id").
		Select("message.message_id , message.author_id , user.username , message.text , message.pub_date , user.email").
		Rows()

	if err != nil {
		panic(err)
	}

	var messages2 []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}
		}
		messages2 = append(messages2, msg)
	}

	return messages2
}

func AddMessage(userId int, message string) error {

	var messageObj = models.Message{
		AuthorId: uint(userId),
		Text:     message,
		Pubdate:  time.Now().Unix(),
		Flagged:  0,
	}

	create := gormDb.
		Select("message_id", "author_id", "text", "pub_date", "flagged").
		Table("message").Create(&messageObj)

	fmt.Println(create)

	if create.Error != nil {
		println(create.Error.Error())
		return errors.New(create.Error.Error())
	}

	return nil
}

func GetPersonalTimelineMessages(id int) []models.Message {

	//db.Find(&users, []int{1,2,3})

	subquery, err := gormDb.
		Table("follower").
		Select("whom_id").
		Where("who_id = ?", id).
		Rows()

	if err != nil {
		//TODO: Remove panic statements. it crashes the application.
		panic(err)
	}

	var follows []int
	for subquery.Next() {
		var user int
		err := subquery.Scan(&user)
		if err != nil {
			//TODO
		}
		follows = append(follows, user)
	}

	//select whom_id from follower where who_id = ?)
	result, err := gormDb.
		Model(models.Message{}).
		Table("message").
		Order("pub_date desc").
		Limit(30).
		Where("flagged = ? ", 0, "and message.author_id = user.user_id",
			"and (", "user.user_id = ? or ", id,
			"user.user_id in (?)", arrayToString(follows, ",")).
		Joins("JOIN user on message.author_id = user.user_id").
		Select("message.message_id , message.author_id , user.username , message.text , message.pub_date , user.email").
		Rows()

	if err != nil {
		//TODO: Remove panic statements. it crashes the application.
		panic(err)
	}

	var messages2 []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}
		}
		messages2 = append(messages2, msg)
	}

	return messages2
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

func GetPersonalTimelineMessages_old(id int) []models.Message {
	db := ConnectDb()
	query := string(`
		select message.message_id,message.author_id,user.username,message.text,message.pub_date, user.email
		from message, user
        where message.flagged = 0 and message.author_id = user.user_id and (
            user.user_id = ? or
            user.user_id in (select whom_id from follower
                                    where who_id = ?))
        order by message.pub_date desc limit 30`)
	result, err := db.Query(query, fmt.Sprint(id), fmt.Sprint(id))
	if err != nil {
		panic(err)
	}

	var messages []models.Message

	for result.Next() {
		var msg models.Message
		err := result.Scan(&msg.MessageId, &msg.AuthorId, &msg.Username, &msg.Text, &msg.Pubdate, &msg.Email)
		if err != nil {
			return []models.Message{}
		}
		messages = append(messages, msg)
	}
	defer result.Close()
	return messages
}
