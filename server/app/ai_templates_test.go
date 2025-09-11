package app

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mattermost/focalboard/server/model"
	"github.com/stretchr/testify/require"
)

func TestGenerateTemplateFallback(t *testing.T) {
	th, tearDown := SetupTestHelper(t)
	defer tearDown()

	const teamID = "team1"
	const userID = "user1"
	prompt := "Test template"

	inserted := &model.Board{ID: "new-board", TeamID: teamID}
	th.Store.EXPECT().InsertBoard(gomock.AssignableToTypeOf(&model.Board{}), userID).Return(inserted, nil)
	th.Store.EXPECT().GetMembersForBoard(inserted.ID).Return([]*model.BoardMember{}, nil)

	board, err := th.App.GenerateTemplate(teamID, userID, prompt)
	require.NoError(t, err)
	require.Equal(t, inserted.ID, board.ID)
	require.Equal(t, teamID, board.TeamID)
}
