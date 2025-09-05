package handlers

import (
	"encoding/json"
	"net/http"

	"backend/db"
	"backend/models"
)

func BoardHandler(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.DB.Query("SELECT id, title FROM columns")
	defer rows.Close()

	var board models.Board
	for rows.Next() {
		var col models.Column
		rows.Scan(&col.ID, &col.Title)

		cardRows, _ := db.DB.Query("SELECT id, title, column_id FROM cards WHERE column_id=$1", col.ID)
		defer cardRows.Close()

		var colCards []models.Card
		for cardRows.Next() {
			var c models.Card
			cardRows.Scan(&c.ID, &c.Title, &c.ColumnID)
			colCards = append(colCards, c)
		}

		board.Columns = append(board.Columns, struct {
			ID    int           `json:"id"`
			Title string        `json:"title"`
			Cards []models.Card `json:"cards"`
		}{
			ID:    col.ID,
			Title: col.Title,
			Cards: colCards,
		})
	}

	json.NewEncoder(w).Encode(board)
}
