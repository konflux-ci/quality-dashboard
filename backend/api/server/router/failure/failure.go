package failure

import (
	"context"
	"encoding/json"
	"net/http"

	failureV1Alpha1 "github.com/konflux-ci/quality-studio/api/apis/failure/v1alpha1"
	"github.com/konflux-ci/quality-studio/api/types"
	"github.com/konflux-ci/quality-studio/pkg/utils/httputils"
	"go.uber.org/zap"
)

func (f *failureRouter) createFailure(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var fr failureV1Alpha1.Failure
	if err := json.NewDecoder(r.Body).Decode(&fr); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Error reading team/error_message/jira_key value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	team, err := f.Storage.GetTeamByName(fr.TeamName)
	if err != nil {
		f.Logger.Error("Failed to fetch team. Make sure the team exists", zap.String("team", fr.TeamName), zap.Error(err))

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	jiraBug, err := f.Storage.GetJiraBug(fr.JiraKey)
	if err != nil {
		f.Logger.Sugar().Warnf("Failed to get jira :", err)
	}

	labels := ""
	if jiraBug.Labels != nil {
		labels = *jiraBug.Labels
	}

	titleFromJira := ""
	if jiraBug.Summary != "" {
		titleFromJira = jiraBug.Summary
	}

	err = f.Storage.CreateFailure(failureV1Alpha1.Failure{
		JiraKey:       fr.JiraKey,
		TitleFromJira: titleFromJira,
		JiraStatus:    jiraBug.Status,
		ErrorMessage:  fr.ErrorMessage,
		CreatedDate:   jiraBug.CreatedAt,
		ClosedDate:    *jiraBug.ResolvedAt,
		Labels:        labels,
	}, team.ID)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Failed to save failure data in database.",
			StatusCode: http.StatusBadRequest,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message:    "Successfully created failure in quality-studio",
		StatusCode: http.StatusCreated,
	})
}

func (f *failureRouter) getFailures(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	teamName := r.URL.Query()["team_name"]
	startDate := r.URL.Query()["start_date"]
	endDate := r.URL.Query()["end_date"]

	if len(teamName) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "team_name value not present in query",
			StatusCode: 400,
		})
	} else if len(startDate) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "start_date value not present in query",
			StatusCode: 400,
		})
	} else if len(endDate) == 0 {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "end_date value not present in query",
			StatusCode: 400,
		})
	}

	team, err := f.Storage.GetTeamByName(teamName[0])
	if err != nil {
		f.Logger.Error("Failed to get team")

		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	}

	failures, err := f.Storage.GetFailuresByDate(team, startDate[0], endDate[0])
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to get failures by team.",
			StatusCode: 400,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, failures)
}

func (f *failureRouter) deleteFailure(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var fr failureV1Alpha1.Failure
	if err := json.NewDecoder(r.Body).Decode(&fr); err != nil {
		return httputils.WriteJSON(w, http.StatusInternalServerError, &types.ErrorResponse{
			Message:    "Error reading team/error_message/jira_key value from body",
			StatusCode: http.StatusBadRequest,
		})
	}

	err := f.Storage.DeleteFailure(fr.TeamID, fr.JiraID)
	if err != nil {
		return httputils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{
			Message:    "Failed to delete failure",
			StatusCode: 400,
		})
	}

	return httputils.WriteJSON(w, http.StatusOK, types.SuccessResponse{
		Message: "Failure deleted",
	})
}
