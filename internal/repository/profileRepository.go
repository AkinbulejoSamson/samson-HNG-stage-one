package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/dto"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/model"
)

type ProfileRepository interface {
	CreateProfile(p *model.Profile) (*model.Profile, error)
	GetProfileByName(name string) (*model.Profile, error)
	GetProfileById(id string) (*model.Profile, error)
	GetProfiles(q *dto.ProfileQuery) ([]*model.Profile, int, error)
	Delete(id string) error
}

type profileRepository struct {
	db *sql.DB
}

func NewProfileRepository(DB *sql.DB) ProfileRepository {
	return &profileRepository{
		db: DB,
	}
}

func (r *profileRepository) CreateProfile(p *model.Profile) (*model.Profile, error) {
	_, err := r.db.Exec(`
		INSERT INTO profiles
		(id, name, gender, gender_probability, age, age_group, country_id, country_probability, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Gender, p.GenderProbability,
		p.Age, p.AgeGroup, p.CountryID, p.CountryProbability,
		p.CreatedAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *profileRepository) GetProfileByName(name string) (*model.Profile, error) {
	return r.scanProfile(r.db.QueryRow(`
		SELECT id, name, gender, gender_probability, age, age_group, country_id, country_probability, created_at
		FROM profiles WHERE LOWER(name) = LOWER(?)`, name))
}

func (r *profileRepository) GetProfileById(id string) (*model.Profile, error) {
	return r.scanProfile(r.db.QueryRow(`
		SELECT id, name, gender, gender_probability, age, age_group, country_id, country_probability, created_at
		FROM profiles WHERE id = ?`, id))
}

func (r *profileRepository) GetProfiles(q *dto.ProfileQuery) ([]*model.Profile, int, error) {
	//query := `SELECT * FROM profiles WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM profiles`

	filter := " WHERE 1=1"
	args := []any{}

	if q.Gender != "" {
		filter += " AND LOWER(gender) = LOWER(?)"
		args = append(args, q.Gender)
	}
	if q.CountryID != "" {
		filter += " AND LOWER(country_id) = LOWER(?)"
		args = append(args, q.CountryID)
	}
	if q.CountryName != "" {
		filter += " AND LOWER(country_name) = LOWER(?)"
		args = append(args, q.CountryName)
	}
	if q.AgeGroup != "" {
		filter += " AND LOWER(age_group) = LOWER(?)"
		args = append(args, q.AgeGroup)
	}
	if q.MinAge > 0 {
		filter += " AND age >= ?"
		args = append(args, q.MinAge)
	}
	if q.MaxAge > 0 {
		filter += " AND age <= ?"
		args = append(args, q.MaxAge)
	}
	if q.MinGenderProbability > 0 {
		filter += " AND gender_probability >= ?"
		args = append(args, q.MinGenderProbability)
	}
	if q.MinCountryProbability > 0 {
		filter += " AND country_probability >= ?"
		args = append(args, q.MinCountryProbability)
	}

	var total int
	err := r.db.QueryRow(countQuery+filter, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	sortBy := "created_at"
	allowedSortFields := map[string]bool{
		"age": true, "created_at": true, "gender_probability": true,
	}
	if allowedSortFields[q.SortBy] {
		sortBy = q.SortBy
	}
	order := "ASC"
	if strings.ToUpper(q.OrderBy) == "DESC" {
		order = "DESC"
	}

	query := "SELECT id, name, gender, gender_probability, age, age_group, country_id, country_name, country_probability, created_at FROM profiles" +
		filter +
		fmt.Sprintf(" ORDER BY %s %s", sortBy, order) +
		" LIMIT ? OFFSET ?"

	pagination := append(args, q.Limit, (q.Page-1)*q.Limit)

	rows, err := r.db.Query(query, pagination...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	profiles := []*model.Profile{}

	for rows.Next() {
		var p model.Profile
		var createdAt string
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Gender,
			&p.GenderProbability,
			&p.Age,
			&p.AgeGroup,
			&p.CountryID,
			&p.CountryName,
			&p.CountryProbability,
			&createdAt,
		); err != nil {
			return nil, 0, err
		}

		p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		profiles = append(profiles, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}

func (r *profileRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM profiles WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil
}

// scanProfile is a private helper to scan a row into a Profile struct
func (r *profileRepository) scanProfile(row *sql.Row) (*model.Profile, error) {
	var p model.Profile
	var createdAt string
	err := row.Scan(
		&p.ID, &p.Name, &p.Gender, &p.GenderProbability,
		&p.Age, &p.AgeGroup, &p.CountryID,
		&p.CountryProbability, &createdAt,
	)
	if err != nil {
		return nil, err
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &p, nil
}
