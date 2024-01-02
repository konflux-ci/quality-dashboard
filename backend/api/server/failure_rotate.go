package server

import (
	failureV1Alpha1 "github.com/redhat-appstudio/quality-studio/api/apis/failure/v1alpha1"
)

func (s *Server) UpdateFailuresByTeam() {
	teamArr, _ := s.cfg.Storage.GetAllTeamsFromDB()

	for _, team := range teamArr {
		failures, _ := s.cfg.Storage.GetAllFailures(team)

		for _, failure := range failures {

			jiraBug, err := s.cfg.Storage.GetJiraBug(failure.JiraKey)
			if err != nil {
				s.cfg.Logger.Sugar().Warnf("Failed to get jira :", err)
			}

			// if jiraStatus == "Closed" {
			// 	err = s.cfg.Storage.DeleteFailure(team.ID, failure.ID)
			// 	if err != nil {
			// 		s.cfg.Logger.Sugar().Warnf("Failed to delete failure:", err)
			// 	} else {
			// 		s.cfg.Logger.Info("Deleted closed issue from impact table", zap.String("Jira Key", failure.JiraKey))
			// 	}
			// }

			labels := ""
			if jiraBug.Labels != nil {
				labels = *jiraBug.Labels
			}

			titleFromJira := ""
			if jiraBug.Summary != "" {
				titleFromJira = jiraBug.Summary
			}

			err = s.cfg.Storage.CreateFailure(failureV1Alpha1.Failure{
				JiraKey:       failure.JiraKey,
				TitleFromJira: titleFromJira,
				JiraStatus:    failure.JiraStatus,
				ErrorMessage:  failure.ErrorMessage,
				CreatedDate:   jiraBug.CreatedAt,
				ClosedDate:    *jiraBug.ResolvedAt,
				Labels:        labels,
			}, team.ID)
			if err != nil {
				s.cfg.Logger.Sugar().Warnf("Failed to update failures:", err)
			}
		}
	}
}
