package model

type Node struct {
	Id       int
	Identity string `xml:"id"`
	Name     string `xml:"name"`
}
