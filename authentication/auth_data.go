package authentication

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	/*Required for postgres driver on sqlx package */
	_ "github.com/lib/pq"
)

/*ClientUser represents what is sent to/from the client (frontend) in a JSON format */
type ClientUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

/*ServerUser is what would be represented on backend, also reflecting our database */
type ServerUser struct {
	Userid       int    `db:"userid"`
	Username     string `db:"username"`
	PasswordHash string `db:"passwordhash"`
}

var authDB sqlx.DB

//TODO add sequence here too
var userAccountsSchema = `
CREATE SEQUENCE IF NOT EXISTS public.users_userid_seq
	AS INT
	START WITH 1
	INCREMENT BY 1
	MINVALUE 1
	NO CYCLE
	CACHE 10;
CREATE TABLE IF NOT EXISTS public.user_accounts
(
    userid bigint NOT NULL DEFAULT nextval('users_userid_seq'::regclass),
    username character varying(12) COLLATE pg_catalog."default" NOT NULL,
    firstname character varying(12) COLLATE pg_catalog."default" NOT NULL,
    lastname character varying(12) COLLATE pg_catalog."default" NOT NULL,
    passwordhash character varying(60) COLLATE pg_catalog."default" NOT NULL,
	CONSTRAINT users_pkey PRIMARY KEY (userid),
	CONSTRAINT username UNIQUE (username)
	
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.user_accounts
    OWNER to postgres;
`
var userSessionsSchema = `
CREATE TABLE public.user_session
(
    session_key text COLLATE pg_catalog."default" NOT NULL,
    userid integer NOT NULL,
    login_time timestamp with time zone NOT NULL,
    last_seen_time timestamp with time zone NOT NULL,
    CONSTRAINT user_session_pkey PRIMARY KEY (session_key)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.user_session
	OWNER to postgres;
`
var queryUserByUsername = `
SELECT userid, username, passwordhash FROM public.user_accounts WHERE username=$1
`
var queryInsertUser = `
INSERT INTO public.user_accounts (username, firstname, lastname, passwordhash)
VALUES ($1, $2, $3, $4)
`
var queryInsertSession = `
INSERT INTO public.user_session (userid, session_key, login_time, last_seen_time)
VALUES ($1, $2, $3, $4)
`
var queryDeleteUser = `
DELETE FROM public.user_accounts WHERE username=$1
`

/*Init loads in environment variables and sets a connection up to the database. */
func Init() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = initDB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func initDB() error {
	//Login to database here
	dbCreds := struct {
		Username     string
		Password     string
		DatabaseName string
	}{
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASS"),
		DatabaseName: os.Getenv("DB_NAME"),
	}
	tmpl, err := template.New("dbcreds").Parse("user={{.Username}} dbname={{.DatabaseName}} password={{.Password}} sslmode=disable")
	if err != nil {
		fmt.Println("Template failed to execute")
	}
	var dbcreds strings.Builder
	err = tmpl.Execute(&dbcreds, dbCreds)
	db, err := sqlx.Connect("postgres", `user=testuser dbname=carpoolcreator password=testpassword sslmode=disable`)
	if err != nil {
		return err
	}
	_, err = db.Exec(userAccountsSchema)
	if err != nil {
		fmt.Println(err)
		return err
	}
	authDB = *db
	return err
}

func getAuthDB() sqlx.DB {
	return authDB
}

func getServerUserByUsername(username string) (ServerUser, error) {
	var serverUser ServerUser
	var err error

	row := authDB.QueryRow(queryUserByUsername, username)
	row.Scan(
		&serverUser.Userid,
		&serverUser.Username,
		&serverUser.PasswordHash,
	)
	return serverUser, err
}

func insertUser(user ClientUser) error {
	context := authDB.MustBegin()
	_, err := context.Exec(queryInsertUser,
		user.Username,
		user.FirstName,
		user.LastName,
		hashAndSalt([]byte(user.Password)))
	context.Commit()
	return err
}

/*TODO: Probably secure this further in some way*/
func deleteUser(username string) error {
	context := authDB.MustBegin()
	_, err := context.Exec(queryDeleteUser, username)
	context.Commit()
	return err
}

func insertSessionToken(userID int, sessionKey string, currentTime string, expireTime string) error {
	context := authDB.MustBegin()
	_, err := context.Exec(queryInsertSession, userID, sessionKey, currentTime, expireTime)
	if err != nil {
		log.Println(err)
	}
	context.Commit()
	return err
}
