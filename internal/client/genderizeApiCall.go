package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AkinbulejoSamson/samson-HNG-stage-one/internal/dto"
)

func FetchGenderizeRawData(ctx context.Context, name string) (*dto.GenderizeRawData, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(" genderize api returned status: %d", resp.Status)
	}

	var result dto.GenderizeRawData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
