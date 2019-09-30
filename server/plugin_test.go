package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/mattermost/mattermost-server/plugin/plugintest/mock"
)

// TestTableFlip - Test that the standalone tableflip command works.
func TestTableFlip(t *testing.T) {
	resp, err := runTestPluginCommand(t, "/tableflip")

	assert.NotNil(t, resp)
	assert.Nil(t, err)

	// Negative tests.
	assert.False(t, strings.Contains(resp.Text, "┬─┬ノ( º _ ºノ)"))

	// Positive tests.
	assert.True(t, strings.Contains(resp.Text, "(╯°□°)╯︵ ┻━┻"))
}

// TestTableFlipWithText - Test that the tableflip command with text works.
func TestTableFlipWithText(t *testing.T) {
	resp, err := runTestPluginCommand(t, "/tableflip This is horrible")

	assert.NotNil(t, resp)
	assert.Nil(t, err)

	// Negative tests.
	assert.False(t, strings.Contains(resp.Text, "┬─┬ノ( º _ ºノ)"))
	assert.False(t, strings.Contains(resp.Text, "(╯°□°)╯︵ ┻━┻ This is horrible"))

	// Positive tests.
	assert.True(t, strings.Contains(resp.Text, "This is horrible (╯°□°)╯︵ ┻━┻"))
}

// TestTableDown - Test that the standalone tabledown command works.
func TestTableDown(t *testing.T) {
	resp, err := runTestPluginCommand(t, "/tabledown")

	assert.NotNil(t, resp)
	assert.Nil(t, err)

	// Negative tests.
	assert.False(t, strings.Contains(resp.Text, "(╯°□°)╯︵ ┻━┻"))

	// Positive tests.
	assert.True(t, strings.Contains(resp.Text, "┬─┬ノ( º _ ºノ)"))
}

// TestTableDownWithText - Test that the tabledown command with text works.
func TestTableDownWithText(t *testing.T) {
	resp, err := runTestPluginCommand(t, "/tabledown Oh, actually it's fine")

	assert.NotNil(t, resp)
	assert.Nil(t, err)

	// Negative tests.
	assert.False(t, strings.Contains(resp.Text, "(╯°□°)╯︵ ┻━┻"))
	assert.False(t, strings.Contains(resp.Text, "┬─┬ノ( º _ ºノ) Oh, actually it's fine"))

	// Positive tests.
	assert.True(t, strings.Contains(resp.Text, "Oh, actually it's fine ┬─┬ノ( º _ ºノ)"))
}

// -----------------------------------------------------------------------------
// Utilities
// -----------------------------------------------------------------------------

func runTestPluginCommand(t *testing.T, cmd string) (*model.CommandResponse, *model.AppError) {
	p := initTestPlugin(t)
	assert.Nil(t, p.OnActivate())

	var command *model.CommandArgs
	command = &model.CommandArgs{
		Command: cmd,
	}

	return p.ExecuteCommand(&plugin.Context{}, command)
}

func initTestPlugin(t *testing.T) *Plugin {
	api := &plugintest.API{}
	api.On("RegisterCommand", mock.Anything).Return(nil)
	api.On("UnregisterCommand", mock.Anything, mock.Anything).Return(nil)
	api.On("GetUser", mock.Anything).Return(&model.User{
		Id:        "userid",
		Nickname:  "User",
		Username:  "hunter2",
		FirstName: "User",
		LastName:  "McUserface",
	}, (*model.AppError)(nil))

	p := Plugin{}
	p.SetAPI(api)

	return &p
}