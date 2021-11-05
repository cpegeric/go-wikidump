package model

import "database/sql"

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
