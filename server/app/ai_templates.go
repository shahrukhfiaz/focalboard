package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mattermost/focalboard/server/model"
)

type aiTemplateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type aiTemplateResponse struct {
	Board  *model.Board   `json:"board"`
	Blocks []*model.Block `json:"blocks"`
}

// GenerateTemplate creates a board template based on an AI prompt. It will
// attempt to contact the configured AI service to generate the template. If the
// service is not configured or fails, a simple placeholder board is created
// instead.
func (a *App) GenerateTemplate(teamID, userID, prompt string) (*model.Board, error) {
	// Attempt to contact external AI service if configured
	if a.config != nil && a.config.AIURL != "" {
		reqBody := aiTemplateRequest{Model: a.config.AIModel, Prompt: prompt}
		data, err := json.Marshal(reqBody)
		if err == nil {
			req, err := http.NewRequest(http.MethodPost, a.config.AIURL, bytes.NewReader(data))
			if err == nil {
				req.Header.Set("Content-Type", "application/json")
				if a.config.AIAPIKey != "" {
					req.Header.Set("Authorization", "Bearer "+a.config.AIAPIKey)
				}
				client := &http.Client{Timeout: 30 * time.Second}
				resp, err := client.Do(req)
				if err == nil && resp != nil {
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						var aiResp aiTemplateResponse
						if json.NewDecoder(resp.Body).Decode(&aiResp) == nil && aiResp.Board != nil {
							aiResp.Board.TeamID = teamID
							aiResp.Board.Type = model.BoardTypeOpen
							aiResp.Board.IsTemplate = true
							aiResp.Board.TemplateVersion = 1
							board, err := a.CreateBoard(aiResp.Board, userID, false)
							if err != nil {
								return nil, err
							}
							if len(aiResp.Blocks) > 0 {
								for _, block := range aiResp.Blocks {
									block.BoardID = board.ID
								}
								if _, err = a.InsertBlocks(aiResp.Blocks, userID); err != nil {
									return nil, err
								}
							}
							return board, nil
						}
					}
				}
			}
		}
	}

	// Fallback: simple placeholder board
	title := fmt.Sprintf("AI Template: %s", prompt)
	if a.config != nil && a.config.AIModel != "" {
		title = fmt.Sprintf("AI Template (%s): %s", a.config.AIModel, prompt)
	}

	board := &model.Board{
		TeamID:          teamID,
		Title:           title,
		Type:            model.BoardTypeOpen,
		IsTemplate:      true,
		TemplateVersion: 1,
	}

	return a.CreateBoard(board, userID, false)
}
