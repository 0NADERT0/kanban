package models

type User struct {
	ID       int
	Username string
	Password string
}

type Column struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Card struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	ColumnID int    `json:"column_id"`
}

type Board struct {
	Columns []struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Cards []Card `json:"cards"`
	} `json:"columns"`
}
