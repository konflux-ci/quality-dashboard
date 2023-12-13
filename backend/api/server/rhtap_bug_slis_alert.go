package server

import (
	"fmt"

	"github.com/redhat-appstudio/quality-studio/api/server/router/jira"
	"github.com/robfig/cron"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type Team struct {
	ComponentName string
	Bugs          []jira.Bug
}

func getMention(team string) string {
	switch team {
	case "quality":
		return "<!subteam^S02G2PFJ4AV|appstudio-qe-team> and <!subteam^S03E4H1JLF5|hacbs-qe-team>"
	case "spi":
		return "<!subteam^SLRSHSU1K|spi-rhtap-team>"
	case "has":
		return "<!subteam^S04MSCVRF4Z|app-has>"
	case "integration":
		return "<!subteam^S041261DDEW|rhtap-test-team>"
	case "pipeline":
		return "<!subteam^S03GF42RBE2|plnsvc-team>"
	case "gitops":
		return "<!subteam^S01AC8DU22C|gitops-team>"
	case "build":
		return "<!subteam^S014L5WTRBP|build-api-team>"
	case "release":
		return "<!subteam^S03SVBS426R|stonesoup-release-team>"
	case "o11y":
		return "<!subteam^S04S21ECL8K|rhtap-o11y-all>"
	case "ec":
		return "<!subteam^S04123TQ599|hacbs-contract-team>"
	case "hac":
		return "<!subteam^S02J1EUMMNV|hac-core-team>"
	case "sandbox":
		return "<!subteam^SKBFYSRAL|sandbox-team>"
	default:
		return ""
	}

	// missing:
	//  core
	//  docs
	//  java-rebuild
	//  performance
	//  security
	//  sre
	//  uxd
}

func getMessage(alert jira.Alert, alertType string) string {
	if alert.Signal == alertType {
		return *alert.AlertMessage + "\n\n"
	}

	return ""
}

func (s *Server) sendAlert(team, msg, color string) {
	// channel: rhtap-bug-slis-alerts
	channelID := "C062AF1RFK8"

	// get team mention
	mention := getMention(team)

	message := slack.Attachment{
		Pretext: fmt.Sprintf("Hello %s!", mention),
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

func teamExists(componentName string, teams []Team) int {
	for i, team := range teams {
		if team.ComponentName == componentName {
			return i
		}
	}

	return -1
}

func slisByTeam(bugs []jira.Bug) []Team {
	teams := make([]Team, 0)

	for _, bug := range bugs {

		idx := teamExists(bug.Component, teams)

		if idx != -1 {
			teams[idx].Bugs = append(teams[idx].Bugs, bug)
		} else {
			teams = append(teams, Team{ComponentName: bug.Component, Bugs: []jira.Bug{bug}})
		}
	}

	return teams
}

func (s *Server) sendAlerts() {
	// startDate := time.Now().AddDate(-2, 0, 0).Format(constants.DateFormat)
	// toDate := time.Now().Format(constants.DateFormat)

	// bugs, err := s.cfg.Storage.GetAllOpenRHTAPBUGS(startDate, toDate)

	bugs, err := s.cfg.Storage.GetAllOpenRHTAPBUGS()
	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to get all open RHTAP Bug SLOs", zap.Error(err))
	}

	slis := jira.GetBugSLIs(bugs)
	teams := slisByTeam(slis.Bugs)

	for _, team := range teams {
		redMsg := ""
		yellowMsg := ""

		for _, bug := range team.Bugs {
			redMsg += getMessage(*bug.TriageSLI, "red")
			redMsg += getMessage(*bug.ResponseSLI, "red")
			redMsg += getMessage(*bug.ResolutionSLI, "red")

			yellowMsg += getMessage(*bug.TriageSLI, "yellow")
			yellowMsg += getMessage(*bug.ResolutionSLI, "yellow")
		}

		if redMsg != "" {
			s.sendAlert(team.ComponentName, redMsg, "#FF0000")
		}

		if yellowMsg != "" {
			s.sendAlert(team.ComponentName, yellowMsg, "#FFFF00")
		}

	}
}

func (s *Server) SendBugSLIAlerts() {
	cron := cron.New()

	// every day at 9am
	err := cron.AddFunc("0 0 9 * * *", func() {
		s.sendAlerts()
	})
	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to add cron", zap.Error(err))
	}

	cron.Start()
}
