package model

import (
	"database/sql"
	"strings"
)

// TODO: return group of pages by ID list.
// TODO: extract all templates.

func SelectPage(db *sql.DB, pageID int64) (int64, error) {
	query := `
        select StreamID 
        from Page 
        where ID = ?
    `
	var streamID int64
	err := db.QueryRow(query, pageID).Scan(&streamID)
	if err != nil {
		return 0, err
	}
	return streamID, nil
}

func SelectPages(db *sql.DB, pageIDs []int64) ([]int64, error) {
	query := `
        select StreamID
        from Page
        Where ID in (?` + strings.Repeat(",?", len(pageIDs)-1) + ")"
	args := make([]interface{}, len(pageIDs))
	for i, id := range pageIDs {
		args[i] = id
	}
	rows, err := db.Query(query, args...)
	results := make([]int64, len(pageIDs))
	if err != nil {
		return nil, err
	}
	i := 0
	for rows.Next() {
		err := rows.Scan(&results[i])
		if err != nil {
			return nil, err
		}
		i++
	}
	return results, nil
}
