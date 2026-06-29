package repository

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/niljimeno/seamail/models"
)

type DB struct {
	Path string
	Pool *sql.DB
}

var Database DB

func (db *DB) Connect() error {
	shouldMigrate := false
	if _, err := os.Stat(db.Path); err != nil {
		shouldMigrate = true
		defer db.Migrate()
	}

	var err error
	db.Pool, err = sql.Open("sqlite3", db.Path)
	if err != nil {
		return err
	}

	if shouldMigrate {
		err := db.Migrate()
		if err != nil {
			os.Remove(db.Path)
			return err
		}
	}

	return nil
}

func (db *DB) Migrate() error {
	var err error

	_, err = db.Pool.Exec("create table if not exists users (id integer primary key, name text unique not null, password TEXT)", nil)
	if err != nil {
		return err
	}

	log.Println("Here i go")
	_, err = db.Pool.Exec("insert into users (name, password) values (?, ?)", "nil", "pass")
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec("create table if not exists mail (id integer primary key, read integer not null default 0)", nil)
	if err != nil {
		return err
	}

	_, err = db.Pool.Exec("create table if not exists connections (id integer primary key, mail_id integer, user_id integer)", nil)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateMail(recievers []string) (int64, error) {
	result, err := db.Pool.Exec("insert into mail default values", nil)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, reciever := range recievers {
		var uid int64

		err = db.Pool.QueryRow("select id from users where name = ?", reciever).Scan(&uid)
		if err != nil {
			log.Println("User not found")
			return 0, err
		}

		_, err = db.Pool.Exec("insert into connections (mail_id, user_id) values (?, ?)", id, uid)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (db *DB) RemoveMail(id int64, reciever string) error {
	tx, err := db.Pool.Begin()
	defer tx.Rollback()

	_, err = tx.Exec("delete from mail where id = ?", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from connections where mail_id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) IsRead(id int64) (bool, error) {
	var read bool

	err := db.Pool.QueryRow("select read from mail where id = ?", id).Scan(&read)
	if err != nil {
		return false, err
	}

	return read, nil
}

func (db *DB) MarkAsRead(id int64) error {
	_, err := db.Pool.Exec("update mail set read = 1 where id = ?", id)
	return err
}

func (db *DB) ListMail() ([]models.Mail, error) {
	rows, err := db.Pool.Query("select * from mail", nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.Mail{}

	for rows.Next() {
		var mail models.Mail
		err := rows.Scan(&mail.Id, &mail.Read)
		if err != nil {
			log.Println("here", err)
			return nil, err
		}
		result = append(result, mail)
	}

	return result, nil
}
