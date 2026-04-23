package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/dto"
)

func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(dto.ErrorResponse{
		Status:  "error",
		Message: message,
	})
	if err != nil {
		return
	}
}

func WriteJSONSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	if data == nil {
		WriteJSONError(w, http.StatusInternalServerError, "no data to return")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func IsAlpha(name string) bool {
	return regexp.MustCompile("^[A-Za-z ]+$").MatchString(name)
}

func ParseProfileQuery(r *http.Request) (*dto.ProfileQuery, error) {
	q := &dto.ProfileQuery{}

	// String filters
	q.Gender = strings.ToLower(strings.TrimSpace(r.URL.Query().Get("gender")))
	q.AgeGroup = strings.ToLower(strings.TrimSpace(r.URL.Query().Get("age_group")))
	q.CountryID = strings.ToLower(strings.TrimSpace(r.URL.Query().Get("country_id")))

	// Int filters
	if v := r.URL.Query().Get("min_age"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("min_age must be an integer")
		}
		q.MinAge = val
	}

	if v := r.URL.Query().Get("max_age"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("max_age must be an integer")
		}
		q.MaxAge = val
	}

	// Float filters
	if v := r.URL.Query().Get("min_gender_probability"); v != "" {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("min_gender_probability must be a float")
		}
		q.MinGenderProbability = val
	}

	if v := r.URL.Query().Get("min_country_probability"); v != "" {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("min_country_probability must be a float")
		}
		q.MinCountryProbability = val
	}

	// Sorting
	allowedSortFields := map[string]bool{
		"age": true, "created_at": true, "gender_probability": true,
	}
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy != "" && !allowedSortFields[sortBy] {
		return nil, fmt.Errorf("sort_by must be one of: age, gender_probability, created_at")

	}
	q.SortBy = sortBy

	order := strings.ToLower(r.URL.Query().Get("order"))
	if order != "" && order != "asc" && order != "desc" {
		return nil, fmt.Errorf("order must be either 'asc' or 'desc'")
	}
	q.OrderBy = order

	paginate(r, q)

	log.Println(q)

	return q, nil
}

func ParseNaturalLanguage(r *http.Request, q string) (*dto.ProfileQuery, error) {
	if strings.TrimSpace(q) == "" {
		return nil, fmt.Errorf("Unable to interpret query")
	}

	query := &dto.ProfileQuery{
		Page:  1,
		Limit: 10,
	}

	lower := strings.ToLower(q)
	matched := false

	// Gender
	if strings.Contains(lower, "femaile") || strings.Contains(lower, "females") {
		query.Gender = "female"
		matched = true
	} else if strings.Contains(lower, "male") || strings.Contains(lower, "males") {
		query.Gender = "male"
		matched = true
	}

	// Age Group
	if strings.Contains(lower, "young") {
		query.MinAge = 16
		query.MaxAge = 24
		matched = true
	} else if strings.Contains(lower, "child") || strings.Contains(lower, "children") {
		query.AgeGroup = "child"
		matched = true
	} else if strings.Contains(lower, "teenager") || strings.Contains(lower, "teenagers") {
		query.AgeGroup = "teenager"
		matched = true
	} else if strings.Contains(lower, "adult") || strings.Contains(lower, "adults") {
		query.AgeGroup = "adult"
		matched = true
	} else if strings.Contains(lower, "senior") || strings.Contains(lower, "seniors") {
		query.AgeGroup = "senior"
		matched = true
	}

	// "above X" → min_age
	if idx := strings.Index(lower, "above "); idx != -1 {
		var age int
		fmt.Sscanf(lower[idx+6:], "%d", &age)
		if age > 0 {
			query.MinAge = age
			matched = true
		}
	}

	// "below X" → max_age
	if idx := strings.Index(lower, "below "); idx != -1 {
		var age int
		fmt.Sscanf(lower[idx+6:], "%d", &age)
		if age > 0 {
			query.MaxAge = age
			matched = true
		}
	}

	// "under X" → max_age
	if idx := strings.Index(lower, "under "); idx != -1 {
		var age int
		fmt.Sscanf(lower[idx+6:], "%d", &age)
		if age > 0 {
			query.MaxAge = age
			matched = true
		}
	}

	// "from country" → country_name
	if idx := strings.Index(lower, "from "); idx != -1 {
		remaining := lower[idx+5:]
		parts := strings.Fields(remaining)
		if len(parts) > 0 {
			query.CountryName = parts[0]
			matched = true
		}
	}

	paginate(r, query)

	if !matched {
		return nil, fmt.Errorf("unable to interpret query")
	}

	return query, nil
}

func paginate(r *http.Request, q *dto.ProfileQuery) {
	// Pagination
	if v := r.URL.Query().Get("page"); v != "" {
		val, err := strconv.Atoi(v)
		if err != nil && val > 0 {
			q.Page = val
		}
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	v := r.URL.Query().Get("limit")
	q.Limit = 10
	if v != "" {
		val, err := strconv.Atoi(v)
		if err == nil {
			q.Limit = val
		}
	}
	if q.Limit < 1 {
		q.Limit = 1
	} else if q.Limit > 50 {
		q.Limit = 50
	}
}
