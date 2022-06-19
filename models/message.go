package models

import (
	"time"
)

// Note the struct tags to ensure idiomatic JSON (lowercase keys)
type Message struct {
	Message_id  int    `json:"message_id,omitempty"`
	Created     string `json:"created,omitempty"`
	Text        string `json:"text,omitempty"`
	Creator_id  int    `json:"creator_id,omitempty"`
	Activity_id int    `json:"activity_id,omitempty"`
}

type MessageUserDetails struct {
	User
	Message
}

type Messages []MessageUserDetails

func (db *DB) MessageCreate(creator_id string, activity_id int, text string) error {
	// Create the message in the database
	stmt, err := db.Prepare("INSERT INTO messages(activity_id, creator_id, created, text_message) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}

	const timeLayout = "2006-01-02 15:04:05"

	dt := time.Now()

	createdAt := dt.Format(timeLayout)

	// Execute the statement with the data
	_, err = stmt.Exec(activity_id, creator_id, createdAt, text)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) MessagesRetrieve(activity_id string) (Messages, error) {
	var messages Messages

	// Get the Messages from the database
	// Note: db.Query() opens and holds a connection until rows.Close()
	rows, err := db.Query(`SELECT 
						message_id,
						users.first_name, 
						users.last_name, 
						created, 
						text_message  
						FROM messages 
						LEFT JOIN users 
						ON users.id = creator_id  
						WHERE activity_id = ? 
						ORDER BY created DESC`, activity_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through query rows

	for rows.Next() {
		message := new(MessageUserDetails)
		err = rows.Scan(&message.Message.Message_id, &message.User.First_name, &message.User.Last_name, &message.Message.Created, &message.Message.Text)
		if err != nil {
			return nil, err
		}

		const timeLayout = "2006-01-02 15:04:05"

		_, err := time.Parse(timeLayout, message.Message.Created)
		if err != nil {
			return nil, err
		}
		// Add the message to the Messages array
		messages = append(messages, *message)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return messages, nil
}
