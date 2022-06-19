package models

import (
	"database/sql"
	"finder_api/utils"
	"fmt"
)

// Note the struct tags to ensure idiomatic JSON (lowercase keys)
type Activity struct {
	Activity_id     int     `json:"id,omitempty"`
	Creator_user_id int     `json:"creator_user_id,omitempty"`
	Name            string  `json:"name,omitempty"`
	Description     string  `json:"description,omitempty"`
	Category        string  `json:"category,omitempty"`
	Longitude       float64 `json:"longitude,omitempty"`
	Latitude        float64 `json:"latitude,omitempty"`
}

type Activities []Activity

func (db *DB) CreateActivity(user_id string, name string, description string, category string, longitude float64, latitude float64) error {
	// Create the activity in the database
	stmt, err := db.Prepare("call join_activity_on_create(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	// Execute the statement with the data
	_, err = stmt.Exec(user_id, name, description, category, longitude, latitude)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateActivity(activity_id int, name sql.NullString, description sql.NullString, category sql.NullString, longitude sql.NullString, latitude sql.NullString) error {
	stmt, err := db.Prepare(`UPDATE activities SET
							name = COALESCE(?, name),
							description = COALESCE(?, description),
							category = COALESCE(?, category),
							longitude = COALESCE(?, longitude),
							latitude = COALESCE(?, latitude)
							WHERE id = ?;`)
	if err != nil {
		return err
	}

	// Execute the statement with the data
	_, err = stmt.Exec(name, description, category, longitude, latitude, activity_id)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetActivities(category_filter string, radius_filter float64, person_coord utils.Coord) (Activities, error) {
	var activities Activities
	var rows *sql.Rows
	var err error

	if len(category_filter) > 0 {
		queryString := `select id, name, category, latitude, longitude from activities where category = ?`
		rows, err = db.Query(queryString, category_filter)
	} else {
		queryString := `select id, name, category, latitude, longitude from activities`
		rows, err = db.Query(queryString)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through query rows
	for rows.Next() {
		activity := new(Activity)

		// Scan gets the columns one row at a time
		err := rows.Scan(&activity.Activity_id, &activity.Name, &activity.Category, &activity.Latitude, &activity.Longitude)
		if err != nil {
			return nil, err
		}

		if radius_filter != 0 {
			act_coord := utils.Coord{
				Lat: float64(activity.Latitude),
				Lon: float64(activity.Longitude),
			}

			distance := utils.Distance(person_coord, act_coord)
			fmt.Println(distance)
			if distance > radius_filter {
				continue
			}
		}
		activities = append(activities, *activity)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return activities, nil

}

func (db *DB) GetActivityByID(activity_id string) (map[string]interface{}, error) {
	// Get the User id, username and token info from the database
	queryString := `select name, description, category, latitude, longitude, users.first_name, users.last_name from activities
	left join users on users.id = activities.creator_user_id where activities.id = ?`
	stmt, err := db.Prepare(queryString)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var user User
	var activity Activity
	err = stmt.QueryRow(activity_id).Scan(&activity.Name, &activity.Description, &activity.Category, &activity.Latitude, &activity.Longitude, &user.First_name, &user.Last_name)
	if err != nil {
		return nil, err
	}

	activityDetails := map[string]interface{}{
		"activity_creator": map[string]interface{}{
			"first_name": user.First_name,
			"last_name":  user.Last_name,
		},
		"activity": map[string]interface{}{
			"name":        activity.Name,
			"description": activity.Description,
			"category":    activity.Category,
			"latitude":    activity.Latitude,
			"longitude":   activity.Longitude,
		},
	}

	return activityDetails, nil
}

func (db *DB) RetrieveUsers(activity_id string) ([]User, error) {
	var users []User

	rows, err := db.Query(`SELECT user_id, users.first_name, users.last_name 
	FROM user_activity_relation 
	LEFT JOIN users ON users.id = user_activity_relation.user_id 
	WHERE user_activity_relation.activity_id = (?) 
	ORDER BY users.first_name`, activity_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through query rows
	for rows.Next() {
		user := new(User)

		// Scan gets the columns one row at a time
		err := rows.Scan(&user.User_id, &user.First_name, &user.Last_name)
		if err != nil {
			return nil, err
		}

		// Add the message to the Messages array
		users = append(users, *user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (db *DB) DeleteActivity(activity_id string, user_id string) error {
	stmt, err := db.Prepare("delete from activities where activities.id = ? and activities.creator_user_id = ?")
	if err != nil {
		return err
	}

	// Execute the statement with the data
	_, err = stmt.Exec(activity_id, user_id)
	if err != nil {
		return err
	}

	return nil
}
