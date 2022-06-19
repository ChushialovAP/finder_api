package models

// Then underscore in front of an import does something I can't remember
import (
	"database/sql"
	"finder_api/utils"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Datastore interface {
	// work with messages
	MessageCreate(creator_id string, activity_id int, text string) error
	MessagesRetrieve(activity_id string) (Messages, error)
	// create user in database
	UserSignup(first_name string, last_name string, email string, phone_number string, password []byte) error
	// update user info
	UpdateUser(user_id string, first_name sql.NullString, last_name sql.NullString, email sql.NullString, phone_number sql.NullString) error
	// get user data via phone number and auth token
	GetUserByPhoneNumber(phone_number string) (*User, error)
	GetUserByAuthToken(auth_token string) (*User_detail, error)
	// generate new token using number
	GenerateToken(phone_number string, password string) (map[string]interface{}, error)
	// validate if token is not expired
	ValidateToken(authToken string) (map[string]interface{}, error)
	// create activity in database
	CreateActivity(user_id string, name string, description string, category string, longitude float64, latitude float64) error
	// update activity info
	UpdateActivity(activity_id int, name sql.NullString, description sql.NullString, category sql.NullString, longitude sql.NullString, latitude sql.NullString) error
	// get activity by id
	GetActivityByID(activity_id string) (map[string]interface{}, error)
	// delete activity by id
	DeleteActivity(activity_id string, user_id string) error
	// join activity
	JoinActivity(user_id string, activity_id string) error
	// retrieve all users from activity
	RetrieveUsers(activity_id string) ([]User, error)
	// retrieve all user's activities
	RetriveActivities(user_id string) (Activities, error)
	// get activities according to filters
	GetActivities(category_filter string, radius_filter float64, person_coord utils.Coord) (Activities, error)
}

type DB struct {
	*sql.DB
}

func New() *DB {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/api_db")
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return &DB{db}
}
