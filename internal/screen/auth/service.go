package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
	sdk "github.com/friendly-social/golang-sdk"
)

const (
	SAVE_FOLDER = "friendly"
	SAVE_FILE   = "auth.json"
)

// AuthMsg singals that user authenticated with new credentials.
type AuthMsg struct {
	Auth *sdk.Authorization
}

type Service struct {
	client *sdk.Client
}

func NewService(client *sdk.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) loadAuth() tea.Msg {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to get user cache dir: %w", err)}
	}

	saveFile := filepath.Join(cacheDir, SAVE_FOLDER, SAVE_FILE)
	_, err = os.Stat(saveFile)
	if os.IsNotExist(err) {
		return nil
	}

	authBytes, err := os.ReadFile(saveFile)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to read authotization bytes: %w", err)}
	}

	authorization := new(sdk.Authorization)
	err = json.Unmarshal(authBytes, authorization)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to unmarshal authorization bytes: %w", err)}
	}

	return router.BroadcastMsg{Inner: AuthMsg{Auth: authorization}}
}

func (s *Service) auth(nicknameString, descriptionString, interestsString, socialString string) tea.Msg {
	nickname, err := sdk.NewNickname(nicknameString)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to create nickname: %w", err)}
	}

	description, err := sdk.NewUserDescription(descriptionString)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to create description: %w", err)}
	}

	interestsSlice := make([]sdk.Interest, 0)
	for interestStr := range strings.SplitSeq(interestsString, ",") {
		interest, err := sdk.NewInterest(strings.TrimSpace(interestStr))
		if err != nil {
			return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to create interest: %w", err)}
		}

		interestsSlice = append(interestsSlice, interest)
	}

	interests, err := sdk.NewInterests(interestsSlice)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to create interests: %w", err)}
	}

	social, err := sdk.NewSocialLink(socialString)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to create social link: %w", err)}
	}

	authorization, err := s.client.Register(context.Background(), nickname, description, interests, nil, social)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to authorize: %w", err)}
	}

	cacheFolder, err := os.UserCacheDir()
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to get user cache directory: %w", err)}
	}

	saveFolder := filepath.Join(cacheFolder, SAVE_FOLDER)
	err = os.MkdirAll(saveFolder, 0700)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to create save folder: %w", err)}
	}

	authBytes, err := json.Marshal(authorization)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to marshal authorization data: %w", err)}
	}

	saveFile := filepath.Join(saveFolder, SAVE_FILE)
	err = os.WriteFile(saveFile, authBytes, 0600)
	if err != nil {
		return screen.ErrorMsg{Value: fmt.Errorf("auth: failed to write authorization to save file: %w", err)}
	}

	return router.BroadcastMsg{Inner: AuthMsg{Auth: authorization}}
}
