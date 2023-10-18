package server

import (
	"context"
	"fmt"
	"time"

	"github.com/redhat-appstudio/quality-studio/api/server/router/jira"
	"github.com/redhat-appstudio/quality-studio/pkg/constants"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func getMessage(alert jira.Alert, alertType string) string {
	if alert.Signal == alertType {
		return *alert.AlertMessage + "\n\n"
	}

	return ""
}

func (s *Server) sendAlert(msg, color string) {
	// channel: rhtap-bug-slis-alert
	channelID := "C061N8AL2SW"

	// mention appstudio-qe-team and hacbs-qe-team
	mention := "<!subteam^S02G2PFJ4AV|appstudio-qe-team> and <!subteam^S03E4H1JLF5|hacbs-qe-team>"

	message := slack.Attachment{
		Pretext: fmt.Sprintf("Hello, %s!", mention),
		Text:    msg,
		Color:   color,
	}
	_, _, err := s.cfg.Slack.PostMessage(
		channelID,
		slack.MsgOptionAttachments(message),
	)

	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to post message on slack", zap.Error(err))
	}
}

func (s *Server) sendAlerts(slis jira.BugSlisInfo) {
	redMsg := ""
	yellowMsg := ""
	bugs := slis.ResolutionTimeSLI.Bugs
	bugs = append(bugs, slis.ResponseTimeSLI.Bugs...)
	bugs = append(bugs, slis.TriageTimeSLI.Bugs...)

	for _, bug := range bugs {
		redMsg += getMessage(*bug.TriageSLI, "red")
		redMsg += getMessage(*bug.ResponseSLI, "red")
		redMsg += getMessage(*bug.ResolutionSLI, "red")

		yellowMsg += getMessage(*bug.TriageSLI, "yellow")
		yellowMsg += getMessage(*bug.ResolutionSLI, "yellow")
	}

	if redMsg != "" {
		s.sendAlert(redMsg, "#FF0000")
	}

	if yellowMsg != "" {
		s.sendAlert(yellowMsg, "#FFFF00")
	}
}

func (s *Server) SendBugSLIAlerts() {
	startDate := time.Now().AddDate(-2, 0, 0).Format(constants.DateFormat)
	toDate := time.Now().Format(constants.DateFormat)

	bugs, err := s.cfg.Storage.GetAllOpenRHTAPBUGS(startDate, toDate)
	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to get all open RHTAP Bug SLOs", zap.Error(err))
	}

	ctx := context.TODO()
	slis := jira.GetBugSLIs(bugs)

	go func() {
		s.sendAlerts(slis)

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Hour * 24):
				s.sendAlerts(slis)
			}
		}
	}()
}
