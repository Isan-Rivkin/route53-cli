package expander

import (
	ui "r53/cliui"

	log "github.com/sirupsen/logrus"
)

type StopChan chan *bool

type Controller interface {
	Start() StopChan
	SubmitResult(result *CommandResult, err error)
}

type DefaultController struct {
	EventsListener chan *ui.AppEvent
	CmdExecutors   map[Command]CommandExecutor
}

func NewDefaultController() Controller {
	return &DefaultController{
		EventsListener: make(chan *ui.AppEvent),
		CmdExecutors:   map[Command]CommandExecutor{},
	}
}

func (c *DefaultController) handleEvent(event *ui.AppEvent) {
	switch event.Type {
	case ui.R53TableSelection:

	}
}

func (c *DefaultController) Start() StopChan {

	signalStop := make(chan *bool)
	go func() {
		for {
			log.Debug("controller got new event")
			event, keepOpen := <-c.EventsListener
			if !keepOpen {
				stop := true
				signalStop <- &stop
				break
			}
			c.handleEvent(event)
		}
	}()
	return signalStop
}

func (c *DefaultController) SubmitResult(result *CommandResult, err error) {
	
}
