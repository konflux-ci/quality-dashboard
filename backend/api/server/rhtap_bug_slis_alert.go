package server

import (
	"fmt"

	"github.com/konflux-ci/quality-studio/api/server/router/jira"
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
	case "ic-appstudio-qe":
		return "<!subteam^S03PD4MV58W|ic-appstudio-qe>"
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
		return "<!subteam^S03DM1RL0TF|konfluxbld-green>"
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

func (s *Server) sendAlert(team, msg, color, channelID string) {
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
	// only for KFLUXBUGS for now
	bugs, err := s.cfg.Storage.GetAllOpenBugsForSliAlerts("KFLUXBUGS")
	if err != nil {
		s.cfg.Logger.Sugar().Errorf("Failed to get all open RHTAP Bug SLOs", zap.Error(err))
	}

	slis := jira.GetBugSLIs(bugs)
	teams := slisByTeam(slis.Bugs)

	// channel: konflux-bug-slis-alerts
	channelID := "C062AF1RFK8"

	for _, team := range teams {
		redMsg := ""
		yellowMsg := ""
		componentAssignmentMsg := ""

		for _, bug := range team.Bugs {
			redMsg += getMessage(*bug.TriageSLI, "red")
			redMsg += getMessage(*bug.ResponseSLI, "red")
			redMsg += getMessage(*bug.ResolutionSLI, "red")

			if team.ComponentName == "undefined" {
				// should mention ic-appstudio-qe
				componentAssignmentMsg += getMessage(*bug.ComponentAssignmentTriageSLI, "red")
			}

			yellowMsg += getMessage(*bug.TriageSLI, "yellow")
			yellowMsg += getMessage(*bug.ResolutionSLI, "yellow")
		}

		if redMsg != "" {
			s.sendAlert(team.ComponentName, redMsg, "#FF0000", channelID)
		}

		if componentAssignmentMsg != "" {
			s.sendAlert("ic-appstudio-qe", componentAssignmentMsg, "#FF0000", channelID)
		}

		if yellowMsg != "" {
			s.sendAlert(team.ComponentName, yellowMsg, "#FFFF00", channelID)
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
