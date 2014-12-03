package models

type Project struct {
	id int
	name string
	creator *User
	description string
	shortDescription string
}
