package authentication

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type RegisterUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type ClientUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerUser struct {
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
var queryUserByUsername = `
SELECT username, passwordhash FROM public.user_accounts WHERE username=$1
`
var queryInsertUser = `
INSERT INTO public.user_accounts (username, firstname, lastname, passwordhash)
VALUES ($1, $2, $3, $4)
`
var insertUserTest = `
INSERT INTO public.user_accounts (username, firstname, lastname, passwordhash)
VALUES ('testusername', 'testfname', 'testlastname', 'testpasshash')
`

func Init() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Auth loaded env")
	err = initDB()
	if err != nil {
		fmt.Println(err)
	}
	context := authDB.MustBegin()
	context.Exec(insertUserTest)
	context.Commit()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Auth loaded DB")
	return err
}

func initDB() error {
	//Login to database here
	db, err := sqlx.Connect("insert your database here!! use env variables maybe!!", "hey here's another string to fit this param!!")
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
		&serverUser.Username,
		&serverUser.PasswordHash,
	)
	return serverUser, err
}

func insertUser(user RegisterUser) error {
	context := authDB.MustBegin()
	_, err := context.Exec(queryInsertUser,
		user.Username,
		user.FirstName,
		user.LastName,
		hashAndSalt([]byte(user.Password)))
	context.Commit()
	return err
}
