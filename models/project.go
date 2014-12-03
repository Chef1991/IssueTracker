package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"errors"
)

type Project struct {
	id int
	name string
	creator *User
	description string
	shortDescription string
}

// prints pretty
func (project Project) String() string {
	return fmt.Sprintf("Project:\n\tid:\t\t\t\t\t%d\n\tname:\t\t\t\t%s\n\tdescription:\t\t%s\n\tshortDescription:\t%s\n\tcreator:\t\t\t%v",
						project.id, project.name, project.description, project.shortDescription, *project.creator);
}

//TODO check to see if their is a record already if their is update else make a new one
func (project *Project) accessRecordExists(user *User) bool {
	user.id = user.id
	panic(errors.New("Not implemented"))
	return false
}
func (project *Project) AllowReadWrite(user *User) error{
	errorFunc := func(err error, transaction *sql.Tx) error {
		if transaction != nil {
			transaction.Rollback()
		}
		return err
	}
	db, err := sql.Open("mysql", connectionStr)
	if err != nil { return errorFunc(err, nil)}
	defer db.Close()

	transaction, err := db.Begin()
	if err != nil {return errorFunc(err, transaction)}

	err = project.allowReadWrite(user, transaction)
	if err != nil {return errorFunc(err, transaction)}

	transaction.Commit()
	return nil
}

func (project *Project) AllowRead(user *User) error {
	errorFunc := func(err error, transaction *sql.Tx)  error {
		if transaction != nil {
			transaction.Rollback()
		}
		return err
	}
	db, err := sql.Open("mysql", connectionStr)
	if err != nil { return errorFunc(err, nil)}
	defer db.Close()

	transaction, err := db.Begin()
	if err != nil {return errorFunc(err, transaction)}

	err = project.allowRead(user, transaction)
	if err != nil {return errorFunc(err, transaction)}

	transaction.Commit()
	return nil
}

// TODO: Get rid of this, it is useless, you cant have write access without read access
func (project *Project) allowWrite(user *User, transaction *sql.Tx) error {
	stmt, err := transaction.Prepare("INSERT INTO AccessRights (user_id, project_id, readAccess, writeAccess) VALUES(?,?,?,?)")
	if err != nil { return err }

	_, err = stmt.Exec(user.id, project.id, 0, 1)
	return err
}

func (project *Project) allowRead(user *User, transaction *sql.Tx) error {
	stmt, err := transaction.Prepare("INSERT INTO AccessRights (user_id, project_id, readAccess, writeAccess) VALUES(?,?,?,?)")
	if err != nil { return err }

	_, err = stmt.Exec(user.id, project.id, 1, 0)
	return err
}

func (project *Project) allowReadWrite(user *User, transaction *sql.Tx) error {

	stmt, err := transaction.Prepare("INSERT INTO AccessRights (user_id, project_id, readAccess, writeAccess) VALUES(?,?,?,?)")
	if err != nil { return err }

	_, err = stmt.Exec(user.id, project.id, 1, 1)
	return err
}

// creates a project from the database with the specified id, Project is nil if not found, error is nil if found
func GetProjectFromId(id int) (*Project, error) {
	project := new(Project)
	errorFunc := func(err error) (*Project, error) {
		return nil, err
	}
	myLog("getting user from id")
	db, err := sql.Open("mysql", connectionStr)
	if err != nil { return errorFunc(err)}
	//TODO: use a transaction here
	defer db.Close()

	stmt, err := db.Prepare("SELECT name, creatorId, description, shortDescription FROM Project WHERE id=?")
	if err != nil {
		fmt.Println(1)
		return errorFunc(err)
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)


	creatorId := -1
	row.Scan(&project.name, &creatorId, &project.description, &project.shortDescription)
	if (creatorId == -1) {
		err = errors.New("Project with specified id does not exist")
		return errorFunc(err)
	}
	project.id = id
	project.creator, err = GetUserFromId(creatorId)
	if err != nil { return errorFunc(err)}
	return project, nil
}

func CreateProject(name string, creator *User, description string, shortDescription string) (*Project, error) {
	// create the project
	// set read/write access for creator
	errorFunc := func(err error, transaction *sql.Tx) (*Project, error) {
		if transaction != nil {
			transaction.Rollback()
		}
		return nil, err
	}
	db, err := sql.Open("mysql", connectionStr)
	if err != nil { return errorFunc(err, nil)}
	defer db.Close()

	transaction, err := db.Begin()
	if err != nil {return errorFunc(err, transaction)}

	stmt, err := transaction.Prepare("INSERT INTO Project (name, creatorId, description, shortDescription) VALUES(?, ?, ?, ?)")
	if err != nil { return errorFunc(err, transaction)}
	res, err := stmt.Exec(name, creator.id, description, shortDescription)
	if err != nil { return errorFunc(err, transaction)}

	project := new(Project)
	projectId, err := res.LastInsertId()
	if err != nil { return errorFunc(err, transaction)}
	project.id = int(projectId)
	project.name = name
	project.creator = creator
	project.description = description
	project.shortDescription = shortDescription

	err = project.allowReadWrite(creator, transaction)
	if err != nil { return errorFunc(err, transaction)}

	transaction.Commit()
	return project, nil
}
