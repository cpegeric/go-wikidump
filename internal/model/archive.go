package model

import (
	"database/sql"
	"fmt"
)

type Archive struct {
	FilePath  string
	IndexPath string
	ID        int64
	FileSize  int64
}

func InsertArchives(db *sql.DB, archives []*Archive) error {
	query := `
        insert or ignore 
        into Archive (FilePath,FileSize,IndexPath) 
        values 
    `
	for _, archive := range archives {
		query += fmt.Sprintf("('%v',%v,'%v'),", archive.FilePath, archive.FileSize, archive.IndexPath)
	}
	query = query[:len(query)-1] + ";"
	_, err := db.Exec(query)
	return err
}

func SelectArchives(db *sql.DB) ([]*Archive, error) {
	query := `
        select ID,FilePath,FileSize,IndexPath 
        from Archive 
        where Processed = false
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []*Archive
	for rows.Next() {
		var df Archive
		if err := rows.Scan(&df.ID, &df.FilePath, &df.FileSize, &df.IndexPath); err != nil {
			return nil, err
		}
		results = append(results, &df)
	}
	return results, nil
}

func MarkArchiveProcessed(db *sql.DB, id int64) error {
	query := `
        update Archive
        set Processed = true
        where ID = ?
    `
	_, err := db.Exec(query, id)
	return err
}
