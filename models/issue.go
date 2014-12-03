package models

type Issue struct {
	id int
	project *Project
	creator *User
	title string
	description string
	shortDescription string
	//issueType int
	//priority int
}
