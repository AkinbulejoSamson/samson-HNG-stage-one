package repository

import (
	"database/sql"
	"time"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/model"
)

type ProfileRepository interface {
	CreateProfile(p *model.Profile) (*model.Profile, error)
	GetProfileByName(name string) (*model.Profile, error)
	GetProfileById(id string) (*model.Profile, error)
	GetProfiles(gender, countryID, ageGroup string) ([]*model.ProfileListItem, int, error)
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
		(id, name, gender, gender_probability, sample_size, age, age_group, country_id, country_probability, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Gender, p.GenderProbability, p.SampleSize,
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
		SELECT id, name, gender, gender_probability, sample_size, age, age_group, country_id, country_probability, created_at
		FROM profiles WHERE LOWER(name) = LOWER(?)`, name))
}

func (r *profileRepository) GetProfileById(id string) (*model.Profile, error) {
	return r.scanProfile(r.db.QueryRow(`
		SELECT id, name, gender, gender_probability, sample_size, age, age_group, country_id, country_probability, created_at
		FROM profiles WHERE id = ?`, id))
}

func (r *profileRepository) GetProfiles(gender, countryID, ageGroup string) ([]*model.ProfileListItem, int, error) {
	query := `SELECT id, name, gender, age, age_group, country_id FROM profiles WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM profiles WHERE 1=1`

	filter := ""
	args := []any{}

	if gender != "" {
		filter += " AND LOWER(gender) = LOWER(?)"
		args = append(args, gender)
	}
	if countryID != "" {
		filter += " AND LOWER(country_id) = LOWER(?)"
		args = append(args, countryID)
	}
	if ageGroup != "" {
		filter += " AND LOWER(age_group) = LOWER(?)"
		args = append(args, ageGroup)
	}

	var total int
	err := r.db.QueryRow(countQuery+filter, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(query+filter, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	profiles := make([]*model.ProfileListItem, 0, total)

	for rows.Next() {
		p := new(model.ProfileListItem)
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Gender,
			&p.Age,
			&p.AgeGroup,
			&p.CountryID,
		); err != nil {
			return nil, 0, err
		}
		profiles = append(profiles, p)
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
		&p.SampleSize, &p.Age, &p.AgeGroup, &p.CountryID,
		&p.CountryProbability, &createdAt,
	)
	if err != nil {
		return nil, err
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return &p, nil
}
