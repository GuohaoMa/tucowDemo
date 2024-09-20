package model

type Edge struct {
	Id           int
	Identity     string `xml:"id"`
	FromId       int
	FromIdentity string `xml:"from"`
	ToId         int
	ToIdentity   string  `xml:"to"`
	Cost         float64 `xml:"cost"`
}
