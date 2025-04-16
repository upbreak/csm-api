package entity

type RestDel struct {
	DateName string `json:"dateName"`
	Locdate  int64  `json:"locdate"`
}

type RestDels struct {
	Item []RestDel `json:"item"`
}

type RestDate struct {
	Reason   string `json:"reason"`
	RestDate int64  `json:"rest_date"`
}
type RestDates []RestDate
