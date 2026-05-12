package register

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sdk "github.com/friendly-social/golang-sdk"
)

const (
	saveFile   = "user.json"
	saveFolder = "friendly"
)

// Service provides registration logic.
type Service struct {
	client *sdk.Client
}

// NewService creates Service from sdk.Client.
func NewService(client *sdk.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) load() (*sdk.Authorization, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("register: failed to get user cache dir: %w", err)
	}

	saveFile := filepath.Join(cacheDir, saveFolder, saveFile)
	_, err = os.Stat(saveFile)
	if os.IsNotExist(err) {
		return nil, nil
	}

	userBytes, err := os.ReadFile(saveFile)
	if err != nil {
		return nil, fmt.Errorf("register: failed to read user bytes: %w", err)
	}

	user := new(sdk.Authorization)
	err = json.Unmarshal(userBytes, user)
	if err != nil {
		return nil, fmt.Errorf("register: failed to unmarshal user bytes: %w", err)
	}

	return user, nil
}

func (s *Service) register(nicknameString, descriptionString, interestsString, socialString string) (*sdk.Authorization, error) {
	nickname, err := sdk.NewNickname(nicknameString)
	if err != nil {
		return nil, fmt.Errorf("register: failed to create nickname: %w", err)
	}

	description, err := sdk.NewUserDescription(descriptionString)
	if err != nil {
		return nil, fmt.Errorf("register: failed to create description: %w", err)
	}

	interestsSlice := make([]sdk.Interest, 0)
	for interestStr := range strings.SplitSeq(interestsString, ",") {
		interest, err := sdk.NewInterest(strings.TrimSpace(interestStr))
		if err != nil {
			return nil, fmt.Errorf("register: failed to create interest: %w", err)
		}

		interestsSlice = append(interestsSlice, interest)
	}

	interests, err := sdk.NewInterests(interestsSlice...)
	if err != nil {
		return nil, fmt.Errorf("register: failed to create interests: %w", err)
	}

	socialLink, err := sdk.NewSocialLink(socialString)
	if err != nil {
		return nil, fmt.Errorf("register: failed to create social link: %w", err)
	}

	user, err := s.client.Register(context.Background(), nickname, description, interests, nil, socialLink)
	if err != nil {
		return nil, fmt.Errorf("register: failed to register: %w", err)
	}

	cacheFolder, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("register: failed to get user cache directory: %w", err)
	}

	saveFolder := filepath.Join(cacheFolder, saveFolder)
	err = os.MkdirAll(saveFolder, 0700)
	if err != nil {
		return nil, fmt.Errorf("register: failed to create save folder: %w", err)
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("register: failed to marshal user data: %w", err)
	}

	saveFile := filepath.Join(saveFolder, saveFile)
	err = os.WriteFile(saveFile, userBytes, 0600)
	if err != nil {
		return nil, fmt.Errorf("register: failed to write user data to save file: %w", err)
	}

	return user, nil
}
