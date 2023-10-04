package server

import (
	"fmt"
	"time"

	"github.com/redhat-appstudio/quality-studio/api/server/router/jira"
	"github.com/redhat-appstudio/quality-studio/pkg/constants"
	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func getMessage(alert jira.Alert, alertType string) string {
	if alert.Signal == alertType {
		return *alert.AlertMessage + "\n\n"
	}

	return ""
}

func getMention(projectKey string) string {
	mentions := map[string]string{
		"RHTAPBUGS":  "<!subteam^S02G2PFJ4AV|appstudio-qe-team>",
		"SVPI":       "<!subteam^SLRSHSU1K|spi-rhtap-team>",
		"HACBS":      "<!subteam^S03E4H1JLF5|hacbs-qe-team>",
		"GITOPSRVCE": "<!subteam^S01AC8DU22C|gitops-team>",
		"DEVHAS":     "<!subteam^S04MSCVRF4Z|app-has>",
		"STONEINTG":  "<!subteam^S041261DDEW|rhtap-test-team>",
		"SRVKP":      "<!subteam^S03GF42RBE2|plnsvc-team>",
	}

	if mention, ok := mentions[projectKey]; ok {
		return mention
	}

	return ""
}

func (s *Server) sendAlert(mention, msg, color string) {
	if msg == "" {
		return
	}

	// testing channel: cosmic-testing
	channelID := "C05S0USDKNE"

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

func (s *Server) sendAlertsByProject(project jira.Project) {
	redMsg := ""
	yellowMsg := ""
	mention := getMention(project.ProjectKey)

	if mention == "" {
		return
	}

	for _, bugSLO := range project.BugSLOs {
		redMsg += getMessage(*bugSLO.TriageSLI, "red")
		redMsg += getMessage(*bugSLO.ResponseSLI, "red")
		redMsg += getMessage(*bugSLO.ResolutionSLI, "red")

		yellowMsg += getMessage(*bugSLO.TriageSLI, "yellow")
		yellowMsg += getMessage(*bugSLO.ResolutionSLI, "yellow")
	}

	s.sendAlert(mention, redMsg, "#FF0000")
	s.sendAlert(mention, yellowMsg, "#FFFF00")
}

func (s *Server) SendBugSLOAlerts() {
	startDate := time.Now().AddDate(-2, 0, 0).Format(constants.DateFormat)
	toDate := time.Now().Format(constants.DateFormat)

	teamArr, err := s.cfg.Storage.GetAllTeamsFromDB()
	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to update cache", zap.Error(err))
	}

	bugs := make([]*db.Bugs, 0)
	for _, team := range teamArr {
		b, err := s.cfg.Storage.GetAllOpenBugSLOs(startDate, toDate, team)
		if err != nil {
			s.cfg.Logger.Sugar().Errorf("Failed to get all open bug SLOs", zap.Error(err))
		}
		bugs = append(bugs, b...)
	}

	projectSlos := jira.GetBugSLOsByProject(bugs)

	for _, projectSlo := range projectSlos {
		s.sendAlertsByProject(projectSlo)
	}
}
