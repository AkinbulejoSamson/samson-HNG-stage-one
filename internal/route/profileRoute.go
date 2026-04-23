package route

import (
	"net/http"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/handler"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/service"
)

func SetupProfileRoutes(profileService service.ProfileService) *http.ServeMux {
	profileMux := http.NewServeMux()
	profileHandler := handler.NewProfileHandler(profileService)

	// POST /api/profiles - Create or retrieve profile
	profileMux.HandleFunc("POST /api/profiles", profileHandler.CreateOrRetrieveProfile)

	// GET /api/profiles - Get all profiles with optional filters
	profileMux.HandleFunc("GET /api/profiles", profileHandler.GetAll)

	// GET /api/profiles/search - Get all profiles with natural language processing
	profileMux.HandleFunc("GET /api/profiles/search", profileHandler.Search)

	// GET /api/profiles/{id} - Get profile by ID
	profileMux.HandleFunc("GET /api/profiles/{id}", profileHandler.GetProfileByID)

	// DELETE /api/profiles/{id} - Delete profile by ID
	profileMux.HandleFunc("DELETE /api/profiles/{id}", profileHandler.Delete)

	return profileMux
}
