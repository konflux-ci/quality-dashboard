package server

import (
	failureV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/failure/v1alpha1"
)

func (s *Server) UpdateFailuresByTeam() {
	teamArr, _ := s.cfg.Storage.GetAllTeamsFromDB()

	for _, team := range teamArr {
		failures, _ := s.cfg.Storage.GetAllFailures(team)

		for _, failure := range failures {
			jiraStatus, err := s.cfg.Storage.GetJiraStatus(failure.JiraKey)
			if err != nil {
				s.cfg.Logger.Sugar().Warnf("Failed to get jira status:", err)
			}

			err = s.cfg.Storage.CreateFailure(failureV1Alpha1.Failure{
				JiraKey:      failure.JiraKey,
				JiraStatus:   jiraStatus,
				ErrorMessage: failure.ErrorMessage,
			}, team.ID)
			if err != nil {
				s.cfg.Logger.Sugar().Warnf("Failed to update failures:", err)
			}
		}
	}
}
