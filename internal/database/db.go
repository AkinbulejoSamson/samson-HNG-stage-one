package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/model"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
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
		    age INTEGER NOT NULL,
		    age_group TEXT NOT NULL,
		    country_id TEXT NOT NULL,
		    country_name TEXT NOT NULL,
		    country_probability REAL NOT NULL,
		    created_at TEXT NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS profiles_idx ON profiles(id);
		
		CREATE INDEX IF NOT EXISTS profiles_name_idx ON profiles(name);
		
		CREATE INDEX IF NOT EXISTS profiles_gender_idx ON profiles(gender);

		CREATE INDEX IF NOT EXISTS profiles_gender_probability_idx ON profiles(gender_probability);

		CREATE INDEX IF NOT EXISTS profiles_age_idx ON profiles(age);

		CREATE INDEX IF NOT EXISTS profiles_age_group_idx ON profiles(age_group);

		CREATE INDEX IF NOT EXISTS profiles_country_idx ON profiles(country_id);

		CREATE INDEX IF NOT EXISTS profiles_country_probability_idx ON profiles(country_probability);
	`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}

	wd, _ := os.Getwd()
	fmt.Println("Looking for file in:", wd)
	file, err := os.ReadFile("internal/database/seed_profiles.json")
	if err != nil {
		log.Fatal("Failed to read seed_profiles.json: ", err)
	}

	var profiles []*model.Profile
	if err := json.Unmarshal(file, &profiles); err != nil {
		log.Fatal("Failed to unmarshal seed_profiles.json: ", err)
	}

	inserted := 0
	skipped := 0

	for _, p := range profiles {
		id, _ := uuid.NewV7()
		_, err := DB.Exec(`
			INSERT OR IGNORE INTO profiles
			(id, name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id.String(),
			p.Name,
			p.Gender,
			p.GenderProbability,
			p.Age,
			p.AgeGroup,
			p.CountryID,
			p.CountryName,
			p.CountryProbability,
			time.Now().UTC().Format(time.RFC3339),
		)
		if err != nil {
			log.Printf("Failed to insert %s in database: %v", p.Name, err)
			skipped++
			continue
		}
		inserted++
	}

	log.Println("Database initialized at : ", dbPath)
}
