package auth

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/app"
	sdk "github.com/friendly-social/golang-sdk"
)

func (m Model) authCmd() tea.Cmd {
	return func() tea.Msg {
		nickname, err := sdk.NewNickname(m.inputs[0].Raw().Value())
		if err != nil {
			return app.ErrorMsg{Error: err}
		}

		description, err := sdk.NewUserDescription(m.inputs[1].Raw().Value())
		if err != nil {
			return app.ErrorMsg{Error: err}
		}

		social, err := sdk.NewSocialLink(m.inputs[3].Raw().Value())
		if err != nil {
			return app.ErrorMsg{Error: err}
		}

		interests := make([]sdk.Interest, 0)
		for interestStr := range strings.SplitSeq(m.inputs[2].Raw().Value(), ",") {
			interest, err := sdk.NewInterest(strings.TrimSpace(interestStr))
			if err != nil {
				return app.ErrorMsg{Error: err}
			}
			interests = append(interests, interest)
		}

		auth, err := m.client.Generate(context.Background(), nickname, description, interests, nil, social)
		if err != nil {
			return app.ErrorMsg{Error: err}
		}

		configDir, err := os.UserConfigDir()
		if err != nil {
			return app.ErrorMsg{Error: err}
		}

		saveFolderPath := filepath.Join(configDir, "friendly")
		err = os.MkdirAll(saveFolderPath, 0700)
		if err != nil {
			return app.ErrorMsg{Error: err}
		}

		authBytes, err := json.Marshal(auth)
		if err != nil {
			return app.ErrorMsg{Error: err}
		}

		saveFilePath := filepath.Join(saveFolderPath, "auth.json")
		err = os.WriteFile(saveFilePath, authBytes, 0600)
		if err != nil {
			return app.ErrorMsg{Error: err}
		}

		return AuthorizedMsg{Auth: auth}
	}
}

func (m Model) initCmd() tea.Msg {
	filename, err := os.UserConfigDir()
	if err != nil {
		return app.ErrorMsg{Error: err}
	}

	filename = filepath.Join(filename, "friendly", "auth.json")
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}

	authBytes, err := os.ReadFile(filename)
	if err != nil {
		return app.ErrorMsg{Error: err}
	}

	auth := new(sdk.Authorization)
	err = json.Unmarshal(authBytes, auth)
	if err != nil {
		return app.ErrorMsg{Error: err}
	}

	return AuthorizedMsg{Auth: auth}
}
