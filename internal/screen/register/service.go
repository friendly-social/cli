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
	saveFolder = "friendly"
	saveFile   = "user.json"
)

// AuthMsg singals that user authenticated with new credentials.
type AuthMsg struct {
	User *sdk.Authorization
}

type Service struct {
	client *sdk.Client
}

func NewService(client *sdk.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) load() (*sdk.Authorization, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("auth: failed to get user cache dir: %w", err)
	}

	saveFile := filepath.Join(cacheDir, saveFolder, saveFile)
	_, err = os.Stat(saveFile)
	if os.IsNotExist(err) {
		return nil, nil
	}

	authBytes, err := os.ReadFile(saveFile)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to read authotization bytes: %w", err)
	}

	authorization := new(sdk.Authorization)
	err = json.Unmarshal(authBytes, authorization)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to unmarshal authorization bytes: %w", err)
	}

	return authorization, nil
}

func (s *Service) auth(nicknameString, descriptionString, interestsString, socialString string) (*sdk.Authorization, error) {
	nickname, err := sdk.NewNickname(nicknameString)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to create nickname: %w", err)
	}

	description, err := sdk.NewUserDescription(descriptionString)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to create description: %w", err)
	}

	interestsSlice := make([]sdk.Interest, 0)
	for interestStr := range strings.SplitSeq(interestsString, ",") {
		interest, err := sdk.NewInterest(strings.TrimSpace(interestStr))
		if err != nil {
			return nil, fmt.Errorf("auth: failed to create interest: %w", err)
		}

		interestsSlice = append(interestsSlice, interest)
	}

	interests, err := sdk.NewInterests(interestsSlice...)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to create interests: %w", err)
	}

	socialLink, err := sdk.NewSocialLink(socialString)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to create social link: %w", err)
	}

	authorization, err := s.client.Register(context.Background(), nickname, description, interests, nil, socialLink)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to authorize: %w", err)
	}

	cacheFolder, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("auth: failed to get user cache directory: %w", err)
	}

	saveFolder := filepath.Join(cacheFolder, saveFolder)
	err = os.MkdirAll(saveFolder, 0700)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to create save folder: %w", err)
	}

	authBytes, err := json.Marshal(authorization)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to marshal authorization data: %w", err)
	}

	saveFile := filepath.Join(saveFolder, saveFile)
	err = os.WriteFile(saveFile, authBytes, 0600)
	if err != nil {
		return nil, fmt.Errorf("auth: failed to write authorization to save file: %w", err)
	}

	return authorization, nil
}
