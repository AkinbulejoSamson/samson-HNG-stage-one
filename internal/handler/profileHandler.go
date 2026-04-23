package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/dto"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/handler/helpers"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/service"
)

type ProfileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// POST /api/profiles
func (h *ProfileHandler) CreateOrRetrieveProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Decode body
	var body dto.CreateProfileReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.WriteJSONError(w, http.StatusUnprocessableEntity, "Invalid type")
		return
	}
	if body.Name == nil || strings.TrimSpace(*body.Name) == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Missing or empty name")
		return
	}

	name := strings.TrimSpace(*body.Name)

	profile, existed, statusCode, err := h.profileService.CreateOrRetrieveProfile(r.Context(), name)
	if err != nil {
		helpers.WriteJSONError(w, statusCode, err.Error())
		return
	}

	if existed {
		helpers.WriteJSONSuccess(w, http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "Profile already exists",
			"data":    profile,
		})
		return
	}

	helpers.WriteJSONSuccess(w, statusCode, map[string]interface{}{
		"status": "success",
		"data":   profile,
	})
}

// GET /api/profiles/{id}
func (h *ProfileHandler) GetProfileByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := r.PathValue("id")
	if id == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Missing id")
		return
	}

	profile, statusCode, err := h.profileService.GetProfileByID(id)
	if err != nil {
		helpers.WriteJSONError(w, statusCode, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, statusCode, map[string]interface{}{
		"status": "success",
		"data":   profile,
	})
}

// GET /api/profiles
func (h *ProfileHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	query, err := helpers.ParseProfileQuery(r)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	profiles, count, statusCode, err := h.profileService.GetAll(query)
	if err != nil {
		helpers.WriteJSONError(w, statusCode, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"page":   query.Page,
		"limit":  query.Limit,
		"total":  count,
		"data":   profiles,
	})
}

// GET /api/profiles/ssearch
func (h *ProfileHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Missing or empty query")
		return
	}

	query, err := helpers.ParseNaturalLanguage(r, q)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	profiles, count, statusCode, err := h.profileService.GetAll(query)
	if err != nil {
		helpers.WriteJSONError(w, statusCode, err.Error())
		return
	}

	helpers.WriteJSONSuccess(w, http.StatusOK, map[string]interface{}{
		"status": "success",
		"page":   query.Page,
		"limit":  query.Limit,
		"total":  count,
		"data":   profiles,
	})
}

// DELETE /api/profiles/{id}
func (h *ProfileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := r.PathValue("id")
	if id == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "Missing id")
		return
	}

	statusCode, err := h.profileService.Delete(id)
	if err != nil {
		helpers.WriteJSONError(w, statusCode, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
