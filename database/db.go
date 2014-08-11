package database

import (
	"fmt"
	"strings"

	"database/sql"
)

type Schema []*Table

type Table struct {
	Name       string
	Columns    Columns
	PrimaryKey []string
}

type Columns []struct {
	Field      string
	Type       string
	Constraint string
}

type DB struct {
	*sql.DB
	schema Schema
}

type Tx struct {
	*sql.Tx
}

func Create(db *sql.DB, schema Schema) (*DB, error) {
	newDB := &DB{
		DB:     db,
		schema: schema,
	}

	err := newDB.createTablesIfNotExist()
	if err != nil {
		return nil, err
	}

	return newDB, nil
}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func (c Columns) String() string {
	flat := make([]string, len(c))
	for i, col := range c {
		s := fmt.Sprintf("%s %s", col.Field, col.Type)
		if len(col.Constraint) != 0 {
			s += " " + col.Constraint
		}
		flat[i] = s
	}
	return strings.Join(flat, ", ")
}

func (db *DB) createTablesIfNotExist() error {
	for _, table := range db.schema {
		count, err := db.QueryInt("SELECT count(1) FROM sqlite_master WHERE type=? AND name=?", "table", table.Name)
		if err != nil {
			return err
		}

		if count == 0 {
			cols := table.Columns.String()
			if len(table.PrimaryKey) != 0 {
				cols += fmt.Sprintf(", primary key (%s)", strings.Join(table.PrimaryKey, ","))
			}
			stmt := fmt.Sprintf("create table %s (%s)", table.Name, cols)
			_, err := db.Exec(stmt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *DB) QueryInt(query string, args ...interface{}) (int, error) {
	return queryInt(db.DB, query, args...)
}

func (tx *Tx) QueryInt(query string, args ...interface{}) (int, error) {
	return queryInt(tx.Tx, query, args...)
}

type queryRower interface {
	QueryRow(string, ...interface{}) *sql.Row
}

func queryInt(dbOrTx queryRower, query string, args ...interface{}) (int, error) {
	var count int
	err := dbOrTx.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *DB) QueryNonZero(query string, args ...interface{}) (bool, error) {
	count, err := db.QueryInt(query, args...)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// QM returns "(?,...,?)" with n question marks.
func QM(n int) string {
	s := make([]string, n)
	for i := range s {
		s[i] = "?"
	}
	return "(" + strings.Join(s, ",") + ")"
}
