package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/microwaves/randomizer"
)

type DatabaseConn struct {
	*sql.DB
}

func NewDatabaseConn(path string) (*DatabaseConn, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &DatabaseConn{db}, nil
}

func (db *DatabaseConn) Bootstrap() error {
	sqlStatement := `
	create table pushes (id char not null primary key, name char, token text);
	delete from pushes;
	`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return nil
}

func (db *DatabaseConn) InsertIntoPushes(name, token string) error {
	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	statement, err := transaction.Prepare(
		"insert into pushes(id, name, token) values(?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer statement.Close()

	uuid, err := randomizer.GenerateUUID()
	if err != nil {
		return err
	}
	_, err = statement.Exec(uuid, name, token)
	if err != nil {
		return err
	}
	transaction.Commit()

	return nil
}
