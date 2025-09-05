package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"backend/db"
	"backend/middleware"
	"backend/models"
)

func CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string `json:"title"`
		ColumnID int    `json:"column_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	var id int
	err := db.DB.QueryRow(
		"INSERT INTO cards (title, column_id) VALUES ($1, $2) RETURNING id",
		req.Title, req.ColumnID,
	).Scan(&id)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "Invalid column")
		return
	}

	json.NewEncoder(w).Encode(models.Card{ID: id, Title: req.Title, ColumnID: req.ColumnID})
}

func DeleteCardHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	res, err := db.DB.Exec("DELETE FROM cards WHERE id=$1", id)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "Delete error")
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		middleware.WriteJSONError(w, http.StatusNotFound, "Card not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateCardHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "Invalid card ID")
		return
	}

	var req struct {
		ColumnID int `json:"column_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	var colID int
	err = db.DB.QueryRow("SELECT id FROM columns WHERE id=$1", req.ColumnID).Scan(&colID)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "Invalid column_id")
		return
	}

	_, err = db.DB.Exec("UPDATE cards SET column_id=$1 WHERE id=$2", req.ColumnID, id)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "Update error")
		return
	}

	json.NewEncoder(w).Encode(models.Card{ID: id, ColumnID: req.ColumnID})
}
