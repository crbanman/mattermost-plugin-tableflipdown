package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex

	configuration *configuration
	enabled       bool
}

const (
	flipTrigger string = "tableflip"
	flipASCII   string = "(╯°□°)╯︵ ┻━┻"
	downTrigger string = "tabledown"
	downASCII   string = "┬─┬ノ( º _ ºノ)"
)

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// OnActivate handles plugin deactivation
func (p *Plugin) OnActivate() error {
	p.enabled = true

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          flipTrigger,
		AutoComplete:     true,
		AutoCompleteHint: "[message]",
		AutoCompleteDesc: fmt.Sprintf("Adds %s to your message", flipASCII),
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", flipTrigger)
	}

	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          downTrigger,
		AutoComplete:     true,
		AutoCompleteHint: "[message]",
		AutoCompleteDesc: fmt.Sprintf("Adds %s to your message", downASCII),
	}); err != nil {
		return errors.Wrapf(err, "failed to register %s command", downTrigger)
	}

	return nil
}

// OnDeactivate handles plugin deactivation
func (p *Plugin) OnDeactivate() error {
	p.enabled = false
	return nil
}

// ExecuteCommand handles the core functionality of the plugin
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")
	switch trigger {
	case flipTrigger:
		return p.executeCommandTableflip(args), nil
	case downTrigger:
		return p.executeCommandTabledown(args), nil

	default:
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown command: " + args.Command),
		}, nil
	}
}

func appError(message string, err error) *model.AppError {
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}
	return model.NewAppError("TableFlipDown Plugin", message, nil, errorMessage, http.StatusBadRequest)
}

func (p *Plugin) executeCommandTableflip(args *model.CommandArgs) *model.CommandResponse {
	message := strings.TrimSpace((strings.Replace(args.Command, "/"+flipTrigger, "", 1)))
	if len(message) > 0 {
		message += " "
	}
	message += flipASCII
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         message,
	}
}

func (p *Plugin) executeCommandTabledown(args *model.CommandArgs) *model.CommandResponse {
	message := strings.TrimSpace((strings.Replace(args.Command, "/"+downTrigger, "", 1)))
	if len(message) > 0 {
		message += " "
	}
	message += downASCII
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         message,
	}
}
