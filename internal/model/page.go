package model

import "database/sql"

// TODO: return stream as a bufio.Scanner
// TODO: return archive as something like a scanner than can read streams.
// TODO: return page by ID.
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
