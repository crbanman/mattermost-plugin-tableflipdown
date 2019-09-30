package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

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

	p.API.RegisterCommand(&model.Command{
		Trigger:          flipTrigger,
		AutoComplete:     true,
		AutoCompleteHint: "[message]",
		AutoCompleteDesc: "Adds " + flipASCII + " to your message",
	})

	return p.API.RegisterCommand(&model.Command{
		Trigger:          downTrigger,
		AutoComplete:     true,
		AutoCompleteHint: "[message]",
		AutoCompleteDesc: "Adds " + downASCII + " to your message",
	})
}

// OnDeactivate handles plugin deactivation
func (p *Plugin) OnDeactivate() error {
	p.enabled = false
	return nil
}

// ExecuteCommand handles the core functionality of the plugin
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	if !p.enabled {
		return nil, appError("Cannot execute command while the plugin is disabled.", nil)
	}

	if p.API == nil {
		return nil, appError("Cannot access the plugin API.", nil)
	}

	slashCommand := "/"
	ascii := ""

	if strings.HasPrefix(args.Command, "/"+flipTrigger) {
		slashCommand += flipTrigger
		ascii = flipASCII
	} else if strings.HasPrefix(args.Command, "/"+downTrigger) {
		slashCommand += downTrigger
		ascii = downASCII
	} else {
		return nil, appError("Expected trigger "+flipTrigger+" or "+downTrigger+", but got "+args.Command, nil)
	}

	message := strings.TrimSpace((strings.Replace(args.Command, slashCommand, "", 1)))

	if len(message) == 0 {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
			Text:         fmt.Sprintf(ascii),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         fmt.Sprintf(message + " " + ascii),
	}, nil
}

func appError(message string, err error) *model.AppError {
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}
	return model.NewAppError("Acro Plugin", message, nil, errorMessage, http.StatusBadRequest)
}
