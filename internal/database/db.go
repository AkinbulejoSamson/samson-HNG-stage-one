package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	//execPath, err := os.Executable()
	//if err != nil {
	//	log.Fatal("Failed to get executable path: ", err)
	//}
	//dir := filepath.Dir(execPath)
	//dataDir := filepath.Join(dir, "data")
	dataDir := "data"
	dbPath := filepath.Join(dataDir, "profiles.db")

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("Failed to create database directory: ", err)
	}

	var err error

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	DB.SetMaxOpenConns(1)

	createTable := `
		CREATE TABLE IF NOT EXISTS profiles (
		    id TEXT PRIMARY KEY,
		    name TEXT UNIQUE NOT NULL,
		    gender TEXT NOT NULL,
		    gender_probability REAL NOT NULL,
		    sample_size INTEGER NOT NULL,
		    age INTEGER NOT NULL,
		    age_group TEXT NOT NULL,
		    country_id TEXT NOT NULL,
		    country_probability REAL NOT NULL,
		    created_at TEXT NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS profiles_idx ON profiles(id);
		
		CREATE INDEX IF NOT EXISTS profiles_name_idx ON profiles(name);
		
		CREATE INDEX IF NOT EXISTS profiles_gender_idx ON profiles(gender);

		CREATE INDEX IF NOT EXISTS profiles_age_idx ON profiles(age_group);

		CREATE INDEX IF NOT EXISTS profiles_country_idx ON profiles(country_id);
	`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}

	log.Println("Database initialized at : ", dbPath)
}
