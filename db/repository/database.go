package repository

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

type database struct {
	db       *sqlx.DB
	dbPath   string
	dbDriver string
	sqlPath  string
}

type prepareFunc func(query string) (*sqlx.NamedStmt, error)

func (d *database) disconnect() {
	d.db.Close()
}

func (d *database) rawSQL(filename string) (string, error) {
	sql, err := ioutil.ReadFile(filepath.Join(d.sqlPath, filename))
	if err != nil {
		log.Printf("error: Error retrieving sql file. filename=%s\n", filename)
		return "", err
	}

	return string(sql), nil
}

// prepareStmt reads a sql file and preapare named query with it
func (d *database) prepareStmt(filename string) (*sqlx.NamedStmt, error) {
	sql, err := d.rawSQL(filename)
	if err != nil {
		return nil, err
	}

	return d.db.PrepareNamed(sql)
}

// queryRow queries a single row
func (d *database) queryRow(filename string, data interface{}) (*sqlx.Row, error) {
	args := data
	if args == nil {
		args = map[string]interface{}{}
	}

	stmt, err := d.prepareStmt(filename)
	if err != nil {
		log.Printf("error: failed to prepare sql. filename=%s. err=%s\n", filename, err)
		return nil, err
	}

	return stmt.QueryRowx(args), nil
}

// queryRows queries multiple rows
func (d *database) queryRows(filename string, data interface{}) (*sqlx.Rows, error) {
	args := data
	if args == nil {
		args = map[string]interface{}{}
	}

	stmt, err := d.prepareStmt(filename)
	if err != nil {
		log.Printf("error: failed to prepare sql statement. err=%v\n", err)
		return nil, err
	}

	rows, err := stmt.Queryx(args)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// command executes sql statement except select
func (d *database) command(filename string, data interface{}) error {
	args := data
	if args == nil {
		args = ""
	}

	stmt, err := d.prepareStmt(filename)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(args); err != nil {
		log.Printf("error: Failed to execute SQL. sql=%s err: %v\n", filename, err)
		return err
	}

	return nil
}

func (d *database) txCommand(tx *sqlx.Tx, filename string, data interface{}) error {
	sql, err := d.rawSQL(filename)
	if err != nil {
		return err
	}

	log.Printf("debug: action=txCommand. sql=%s", sql)
	if _, err := tx.NamedExec(sql, data); err != nil {
		log.Printf("debug: action=txCommand. sql=%s", sql)
		return err
	}

	log.Println("debug: action=txCommand. filename=", filename)
	return nil
}
