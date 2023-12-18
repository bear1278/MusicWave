package repository

import (
	"database/sql"
	"fmt"
	music "github.com/bear1278/MusicWave"
)

type GenreMysql struct {
	db *sql.DB
}

func NewGenreMysql(db *sql.DB) *GenreMysql {
	return &GenreMysql{db: db}
}

func (g *GenreMysql) InsertGenre(artistID string, genres []music.Genre) error {
	var genreID int64
	query2 := fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", genresTable)
	query3 := fmt.Sprintf("INSERT INTO %s VALUES (?,?)", artistGenreTable)
	query4 := fmt.Sprintf("SELECT * FROM %s WHERE ID_Artist=? AND ID_Genre=?", artistGenreTable)
	for _, genre := range genres {
		transaction, err := g.db.Begin()
		if err != nil {
			return err
		}
		id, alreadyExist, err := g.GenreAlreadyExist(genre.Name)
		if err != nil {
			transaction.Rollback()
			return err
		}
		if !alreadyExist {
			dbResult, err := transaction.Exec(query2, genre.Name)
			if err != nil {
				transaction.Rollback()
				return err
			}
			genreID, err = dbResult.LastInsertId()
			if err != nil {
				transaction.Rollback()
				return err
			}
			_, err = transaction.Exec(query3, artistID, genreID)
			if err != nil {
				transaction.Rollback()
				return err
			}
			err = transaction.Commit()
			if err != nil {
				transaction.Rollback()
				return err
			}
		} else {
			genreID = id
			rows, err := transaction.Query(query4, artistID, genreID)
			if err != nil {
				transaction.Rollback()
				return err
			}
			if !rows.Next() {
				_, err = transaction.Exec(query3, artistID, genreID)
				if err != nil {
					transaction.Rollback()
					return err
				}
			}
			err = transaction.Commit()
			if err != nil {
				transaction.Rollback()
				return err
			}
		}
	}
	return nil
}

func (g *GenreMysql) GenreAlreadyExist(name string) (int64, bool, error) {
	var id int64
	query := fmt.Sprintf("SELECT id FROM %s where name=?", genresTable)
	rows, err := g.db.Query(query, name)
	if err != nil {
		return 0, false, err
	}
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, false, err
		}
	}
	if id != 0 {
		return id, true, err
	} else {
		return id, false, err
	}
}

func (g *GenreMysql) SelectUserGenre(userID int64) ([]music.Genre, error) {
	var genre music.Genre
	var genres []music.Genre
	query := fmt.Sprintf("SELECT id,name FROM %s as G "+
		"INNER JOIN %s AS UG ON UG.ID_Genre=G.id "+
		"WHERE UG.ID_USER=?", genresTable, userGenreTable)
	rows, err := g.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&genre.Id, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
		genre = music.Genre{}
	}
	return genres, err
}
