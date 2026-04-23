package service

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/client"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/dto"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/model"
	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/repository"
	"github.com/google/uuid"
)

type ProfileService interface {
	CreateOrRetrieveProfile(ctx context.Context, name string) (*model.Profile, bool, int, error)
	GetProfileByID(ID string) (*model.Profile, int, error)
	GetAll(query *dto.ProfileQuery) ([]*model.Profile, int, int, error)
	Delete(id string) (int, error)
}

type profileService struct {
	profileRepository repository.ProfileRepository
}

func NewProfileService(profileRepository repository.ProfileRepository) ProfileService {
	return &profileService{profileRepository: profileRepository}
}

func (s *profileService) CreateOrRetrieveProfile(ctx context.Context, name string) (*model.Profile, bool, int, error) {
	rCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	existingProfile, err := s.profileRepository.GetProfileByName(name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, false, http.StatusInternalServerError, err
	}
	if existingProfile != nil {
		return existingProfile, true, http.StatusOK, nil
	}

	var (
		mu       sync.Mutex
		once     sync.Once
		wg       sync.WaitGroup
		firstErr error
		result   = &model.Profile{Name: name}
	)

	handleError := func(err error) {
		once.Do(func() {
			firstErr = err
			cancel() // This kills the other HTTP requests immediately

		})
	}

	wg.Go(func() {
		g, err := client.FetchGenderizeRawData(rCtx, name)
		if err != nil {
			handleError(err)
		}
		mu.Lock()
		result.Gender = g.Gender
		result.GenderProbability = g.Probability
		mu.Unlock()
	})

	wg.Go(func() {
		a, err := client.FetchAgifyRawData(rCtx, name)
		if err != nil {
			handleError(err)
		}
		mu.Lock()
		result.Age = a.Age

		switch {
		case a.Age <= 12:
			result.AgeGroup = "child"
		case a.Age <= 19:
			result.AgeGroup = "teenager"
		case a.Age <= 59:
			result.AgeGroup = "adult"
		default:
			result.AgeGroup = "senior"
		}
		mu.Unlock()
	})

	wg.Go(func() {
		n, err := client.FetchNationalizeRawData(rCtx, name)
		if err != nil {
			handleError(err)
		}

		if len(n.Country) == 0 {
			handleError(fmt.Errorf("nationalize returned an invalid response"))
		}

		topCountry := slices.MaxFunc(n.Country, func(a, b dto.CountriesResponse) int {
			return cmp.Compare(a.Probability, b.Probability)
		})

		mu.Lock()
		result.CountryProbability = topCountry.Probability
		result.CountryID = topCountry.CountryID
		mu.Unlock()
	})

	wg.Wait()
	if firstErr != nil {
		return nil, false, http.StatusBadGateway, firstErr
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, false, http.StatusInternalServerError, err
	}
	result.ID = id.String()

	result.CreatedAt = time.Now()

	newProfile, err := s.profileRepository.CreateProfile(result)
	if err != nil {
		return nil, false, http.StatusInternalServerError, err
	}

	return newProfile, false, http.StatusOK, nil
}

func (s *profileService) GetProfileByID(ID string) (*model.Profile, int, error) {
	profile, err := s.profileRepository.GetProfileById(ID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return profile, http.StatusOK, nil
}

func (s *profileService) GetAll(query *dto.ProfileQuery) ([]*model.Profile, int, int, error) {
	profiles, count, err := s.profileRepository.GetProfiles(query)
	if err != nil {
		return nil, count, http.StatusInternalServerError, err
	}

	return profiles, count, http.StatusOK, nil
}

func (s *profileService) Delete(id string) (int, error) {
	if err := s.profileRepository.Delete(id); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}
