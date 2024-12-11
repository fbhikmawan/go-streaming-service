package models

type Video struct {
	Id          	string
	Video       	string
	Title       	string
	Description 	string
	LocalPath       string
	UniqueName  	string
	S3FilesPath  	[]string

}