package database

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3" // driver import
)

var testDB = "test.db"

func TestCreate(t *testing.T) {
	os.Remove(testDB)

	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.Remove(testDB)

	if db == nil {
		t.Fatal("Open should never return a nil db")
	}

	mydb, err := Create(db, Schema{{
		Name: "test_table",
		Columns: Columns{{
			Field:      "id",
			Type:       "text",
			Constraint: "primary key not null",
		}, {
			Field: "comment",
			Type:  "text",
		}},
	}})
	if err != nil {
		t.Fatal("Table creation failed:", err)
	}
	if mydb == nil {
		t.Fatal("Create should never return a nil db")
	}
}

func TestInvalidSchema(t *testing.T) {
	os.Remove(testDB)

	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testDB)

	if db == nil {
		t.Fatal("Open should never return a nil db")
	}

	mydb, err := Create(db, Schema{{
		Name: "test_table",
		Columns: Columns{{
			Field:      "id",
			Type:       "text",
			Constraint: "primary key not null",
		}, {
			Field: "comment",
			Type:  "text",
		}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	mydb.Close()

	db, err = sql.Open("sqlite3", testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = Create(db, Schema{{
		Name: "test_table",
		Columns: Columns{{
			Field:      "id",
			Type:       "text",
			Constraint: "primary key not null",
		}},
	}})
	if err == nil {
		t.Fatal("Should get an error here")
	} else if !strings.HasPrefix(err.Error(), "Inconsistent schema") {
		t.Fatal("Expected schema error but got:", err)
	}
}

func TestValidSchema(t *testing.T) {
	os.Remove(testDB)

	db, err := sql.Open("sqlite3", testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testDB)

	if db == nil {
		t.Fatal("Open should never return a nil db")
	}

	mydb, err := Create(db, Schema{{
		Name: "test_table",
		Columns: Columns{{
			Field:      "id",
			Type:       "text",
			Constraint: "primary key not null",
		}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	mydb.Close()

	db, err = sql.Open("sqlite3", testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = Create(db, Schema{{
		Name: "test_table",
		Columns: Columns{{
			Field:      "ID",
			Type:       "TEXT",
			Constraint: "PRIMARY KEY not null",
		}},
	}})
	if err != nil {
		t.Fatal("[schema should be case-insensitive]", err)
	}
}
