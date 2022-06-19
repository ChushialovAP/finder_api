package models

import "database/sql"

// Note the struct tags to ensure idiomatic JSON (lowercase keys)
type User struct {
	User_id      int    `json:"user_id,omitempty"`
	First_name   string `json:"first_name,omitempty"`
	Last_name    string `json:"last_name,omitempty"`
	Email        string `json:"email,omitempty"`
	Phone_number string `json:"phone_number,omitempty"`
	Password     string `json:"password,omitempty"`
}

type User_detail struct {
	User
	Token
}

func (db *DB) UserSignup(first_name string, last_name string, email string, phone_number string, password []byte) error {
	// Create the user in the database
	stmt, err := db.Prepare("insert into users(first_name, last_name, email, phone_number, password) values (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	// Execute the statement with the data
	_, err = stmt.Exec(first_name, last_name, email, phone_number, password)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateUser(user_id string, first_name sql.NullString, last_name sql.NullString, email sql.NullString, phone_number sql.NullString) error {
	// Create the user in the database
	stmt, err := db.Prepare(`UPDATE users SET
							first_name = COALESCE(?, first_name),
							last_name = COALESCE(?, last_name),
							email = COALESCE(?, email),
							phone_number = COALESCE(?, phone_number)
							WHERE id = ?;`)
	if err != nil {
		return err
	}

	// Execute the statement with the data
	_, err = stmt.Exec(first_name, last_name, email, phone_number, user_id)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUserByPhoneNumber(phone_number string) (*User, error) {

	// Get the User id and password from the database
	stmt, err := db.Prepare("select id, password from users where phone_number = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	userId := 0
	accountPassword := ""

	err = stmt.QueryRow(phone_number).Scan(&userId, &accountPassword)

	if err != nil {
		return nil, err
	}

	return &User{
		User_id:  userId,
		Password: accountPassword,
	}, nil
}

func (db *DB) GetUserByAuthToken(auth_token string) (*User_detail, error) {

	// Get the User id, username and token info from the database
	queryString := `select 
                users.id,
                users.first_name,
                generated_at,
                expires_at                         
            from authentication_tokens
            left join users
            on authentication_tokens.user_id = users.id
            where auth_token = ?`
	stmt, err := db.Prepare(queryString)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var userDetails User_detail
	err = stmt.QueryRow(auth_token).Scan(&userDetails.User.User_id, &userDetails.User.First_name, &userDetails.Token.Generated_at, &userDetails.Token.Expires_at)
	if err != nil {
		return nil, err
	}

	return &userDetails, nil
}

func (db *DB) JoinActivity(user_id string, activity_id string) error {
	// Create the user/activity relation in the database
	stmt, err := db.Prepare("insert into user_activity_relation(user_id, activity_id) values (?, ?)")
	if err != nil {
		return err
	}

	// Execute the statement with the data
	_, err = stmt.Exec(user_id, activity_id)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RetriveActivities(user_id string) (Activities, error) {
	var activities Activities

	rows, err := db.Query(`SELECT activity_id, activities.name
	FROM user_activity_relation 
	LEFT JOIN activities ON activities.id = user_activity_relation.activity_id 
	WHERE user_activity_relation.user_id = (?) 
	ORDER BY activities.name`, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through query rows
	for rows.Next() {
		activity := new(Activity)

		// Scan gets the columns one row at a time
		err := rows.Scan(&activity.Activity_id, &activity.Name)
		if err != nil {
			return nil, err
		}

		// Add the message to the Messages array
		activities = append(activities, *activity)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return activities, nil
}
