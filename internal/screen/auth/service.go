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
	"github.com/friendly-social/cli/internal/system"
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

func (s Screen) initCmd() func() tea.Msg {
	return func() tea.Msg {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to get user cache dir: %w", err)}
		}

		saveFile := filepath.Join(cacheDir, SAVE_FOLDER, SAVE_FILE)
		_, err = os.Stat(saveFile)
		if os.IsNotExist(err) {
			return nil
		}

		authBytes, err := os.ReadFile(saveFile)
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to read authotization bytes: %w", err)}
		}

		authorization := new(sdk.Authorization)
		err = json.Unmarshal(authBytes, authorization)
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to unmarshal authorization bytes: %w", err)}
		}

		return router.BroadcastMsg{Inner: AuthMsg{Auth: authorization}}
	}
}

func (s Screen) submitCmd() func() tea.Msg {
	return func() tea.Msg {
		nickname, err := sdk.NewNickname(s.fields[0].Value())
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to create nickname: %w", err)}
		}

		description, err := sdk.NewUserDescription(s.fields[1].Value())
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to create description: %w", err)}
		}

		interestsSlice := make([]sdk.Interest, 0)
		for interestStr := range strings.SplitSeq(s.fields[2].Value(), ",") {
			interest, err := sdk.NewInterest(strings.TrimSpace(interestStr))
			if err != nil {
				return system.ErrorMsg{Value: fmt.Errorf("auth: failed to create interest: %w", err)}
			}

			interestsSlice = append(interestsSlice, interest)
		}

		interests, err := sdk.NewInterests(interestsSlice)
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to create interests: %w", err)}
		}

		social, err := sdk.NewSocialLink(s.fields[3].Value())
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to create social link: %w", err)}
		}

		authorization, err := s.client.Register(context.Background(), nickname, description, interests, nil, social)
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to authorize: %w", err)}
		}

		cacheFolder, err := os.UserCacheDir()
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to get user cache directory: %w", err)}
		}

		saveFolder := filepath.Join(cacheFolder, SAVE_FOLDER)
		err = os.MkdirAll(saveFolder, 0700)
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to create save folder: %w", err)}
		}

		authBytes, err := json.Marshal(authorization)
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to marshal authorization data: %w", err)}
		}

		saveFile := filepath.Join(saveFolder, SAVE_FILE)
		err = os.WriteFile(saveFile, authBytes, 0600)
		if err != nil {
			return system.ErrorMsg{Value: fmt.Errorf("auth: failed to write authorization to save file: %w", err)}
		}

		return router.BroadcastMsg{Inner: AuthMsg{Auth: authorization}}
	}
}
