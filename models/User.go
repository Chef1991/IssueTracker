package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

type User struct {
	id int
	email string
	passwordHash string
	firstName string
	lastName string
}

func (user *User) CanModify(project Project) bool {
	myLogError("meow")
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		fmt.Println("ERROR:", err)
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT writeAccess FROM AccessRights WHERE user_id=? AND project_id=?")
	if err != nil {
		fmt.Println("ERROR:", err)
		return false
	}
	defer stmt.Close()
	row := stmt.QueryRow(user.id, project.id)
	var i int8
	row.Scan(&i)
	return i == 1  		// 1 = true, 0 = false


}

func GetUserFromId(id int) (*User, error) {
	myLog("getting user from id")
	user := new(User)
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, email, password, firstName, lastName FROM Users WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&user.id, &user.email, &user.passwordHash, &user.firstName, &user.lastName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(email string, rawPassword string, firstName string, lastName string) (*User, error){
	hash := hashPassword(rawPassword)
	myLog("Creating User", email, firstName, lastName, "")
	user := new(User)
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO Users (email, password, firstName, lastName) VALUES(?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err :=stmt.Exec(email, hash, firstName, lastName)
	if err != nil {
		return nil, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	user.id = int(lastId)
	user.email = email
	user.passwordHash = hash
	user.firstName = firstName
	user.lastName = lastName

	return user, nil
}

func hashPassword(rawPassword string) string {
	salt := []byte("meow meow meow mix")
	// TODO add per user salt
	passwordBytes := []byte(rawPassword)
	defer clear(passwordBytes)
	hashBytes := pbkdf2.Key(passwordBytes, salt, 4096, sha512.Size, sha512.New)
	hash := hex.EncodeToString(hashBytes)
	return hash
}

func clear(bytes []byte) {
	for i := 0; i < len(bytes); i++ {
		bytes[i] = 0
	}
}
