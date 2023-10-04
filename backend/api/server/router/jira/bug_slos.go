package jira

import (
	"fmt"
	"math"
	"strings"

	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

// GetTriageSLI should return red when priority is not defined for more than 2 days
// GetTriageSLI should return yellow when priority is not defined for more than 1 day but less than 2 days
func GetTriageSLI(bug *db.Bugs) *Alert {
	alert := &Alert{Signal: "green"}

	if bug.Labels != nil && strings.Contains(*bug.Labels, "untriaged") {
		msg := fmt.Sprintf("Issue <%s|%s> is meeting defined Bug SLO for Triage Time. Priority should be defined between a maximum of 2 days. This issue has priority undefined for %.2f days. Please, take a time to resolve it.\n\n",
			bug.URL,
			bug.JiraKey,
			*bug.DaysWithoutResolution,
		)

		if *bug.DaysWithoutPriority > 2 {
			alert.Signal = "red"
			alert.AlertMessage = &msg
		} else if *bug.DaysWithoutPriority > 1 {
			alert.Signal = "yellow"
			alert.AlertMessage = &msg
		}
	}

	return alert
}

// GetResponseSLI should return red if there is no assignee for more than 2 days on Blocker and Critical bugs
func GetResponseSLI(bug *db.Bugs) *Alert {
	alert := &Alert{Signal: "green"}
	if bug.Priority == "Blocker" || bug.Priority == "Critical" {
		if *bug.DaysWithoutAssignee > 2 {
			msg := fmt.Sprintf("Issue <%s|%s> is meeting defined Bug SLO for Blocker or Critical Bug Response Time. Blocker or Critical bugs should be assigned between a maximum of 2 days. This issue has assignee undefined for %.2f days. Please, take a time to assign it.\n\n",
				bug.URL,
				bug.JiraKey,
				*bug.DaysWithoutAssignee,
			)
			alert.AlertMessage = &msg
			alert.Signal = "red"
		}
	}

	return alert
}

func measureBugResolutionSLI(redThreshold, yellowThreshold float64, bug *db.Bugs) *Alert {
	alert := &Alert{Signal: "green"}
	daysWithoutResolution := *bug.DaysWithoutResolution

	msg := fmt.Sprintf("Issue <%s|%s> is meeting defined Bug SLO for %s Bug Resolution Time. %s bugs should not take more than 10 days to be resolved. This issue is not resolved for %.2f days. Please, take a time to resolve it.\n\n",
		bug.URL,
		bug.JiraKey,
		bug.Priority,
		bug.Priority,
		daysWithoutResolution,
	)

	fmt.Println("*bug.DaysWithoutPriority", daysWithoutResolution)

	if daysWithoutResolution > redThreshold {
		alert.Signal = "red"
		alert.AlertMessage = &msg
	} else if daysWithoutResolution > yellowThreshold {
		alert.Signal = "yellow"
		alert.AlertMessage = &msg
	}

	return alert
}

// GetBugResolutionSLI should return:
// red
func GetResolutionSLI(bug *db.Bugs) *Alert {
	alert := &Alert{Signal: "green"}

	switch bug.Priority {
	case "Blocker":
		return measureBugResolutionSLI(10, 5, bug)
	case "Critical":
		return measureBugResolutionSLI(20, 10, bug)
	case "Major":
		return measureBugResolutionSLI(40, 20, bug)
	default:
		return alert
	}
}

func projectExists(projects []Project, project string) int {
	for idx, p := range projects {
		if p.ProjectKey == project {
			return idx
		}
	}
	return -1
}

func GetBugSLOInfo(currentInfo BugSLOInfo, currentSignal, targetSignal string, value float64) BugSLOInfo {
	if currentSignal == targetSignal {
		currentInfo.Total++
		currentInfo.Sum += value
		currentInfo.Average = currentInfo.Sum / float64(currentInfo.Total)
		currentInfo.Average = math.Round(currentInfo.Average*100) / 100
	}

	return currentInfo
}

func GetBugSLOsByProject(bugs []*db.Bugs) []Project {
	Projects := make([]Project, 0)

	for _, bug := range bugs {
		TriageSLI := GetTriageSLI(bug)
		ResponseSLI := GetResponseSLI(bug)
		ResolutionSLI := GetResolutionSLI(bug)

		if TriageSLI.Signal == "green" &&
			ResponseSLI.Signal == "green" &&
			ResolutionSLI.Signal == "green" {
			continue
		}

		new := BugSLO{
			JiraKey:               bug.JiraKey,
			JiraURL:               bug.URL,
			TriageSLI:             TriageSLI,
			ResponseSLI:           ResponseSLI,
			ResolutionSLI:         ResolutionSLI,
			DaysWithoutAssignee:   bug.DaysWithoutAssignee,
			DaysWithoutPriority:   bug.DaysWithoutPriority,
			DaysWithoutResolution: bug.DaysWithoutResolution,
		}

		idx := projectExists(Projects, *bug.ProjectKey)

		if idx != -1 {
			Projects[idx].BugSLOs = append(Projects[idx].BugSLOs, new)
			Projects[idx].RedTriageTimeBugSLOInfo = GetBugSLOInfo(Projects[idx].RedTriageTimeBugSLOInfo, new.TriageSLI.Signal, "red", *new.DaysWithoutPriority)
			Projects[idx].YellowTriageTimeBugSLOInfo = GetBugSLOInfo(Projects[idx].YellowTriageTimeBugSLOInfo, new.TriageSLI.Signal, "yellow", *new.DaysWithoutPriority)
			Projects[idx].RedResponseTimeBugSLOInfo = GetBugSLOInfo(Projects[idx].RedResponseTimeBugSLOInfo, new.ResponseSLI.Signal, "red", *new.DaysWithoutAssignee)
			Projects[idx].RedResolutionTimeBugSLOInfo = GetBugSLOInfo(Projects[idx].RedResolutionTimeBugSLOInfo, new.ResolutionSLI.Signal, "red", *new.DaysWithoutResolution)
			Projects[idx].YellowResolutionTimeBugSLOInfo = GetBugSLOInfo(Projects[idx].YellowResolutionTimeBugSLOInfo, new.ResolutionSLI.Signal, "yellow", *new.DaysWithoutResolution)
		} else {
			bugSLOInfo := BugSLOInfo{
				Total:   0,
				Sum:     0,
				Average: 0,
			}
			Projects = append(Projects, Project{
				ProjectKey:                     *bug.ProjectKey,
				RedTriageTimeBugSLOInfo:        GetBugSLOInfo(bugSLOInfo, new.TriageSLI.Signal, "red", *new.DaysWithoutPriority),
				YellowTriageTimeBugSLOInfo:     GetBugSLOInfo(bugSLOInfo, new.TriageSLI.Signal, "yellow", *new.DaysWithoutPriority),
				RedResponseTimeBugSLOInfo:      GetBugSLOInfo(bugSLOInfo, new.ResponseSLI.Signal, "red", *new.DaysWithoutAssignee),
				RedResolutionTimeBugSLOInfo:    GetBugSLOInfo(bugSLOInfo, new.ResolutionSLI.Signal, "red", *new.DaysWithoutResolution),
				YellowResolutionTimeBugSLOInfo: GetBugSLOInfo(bugSLOInfo, new.ResolutionSLI.Signal, "yellow", *new.DaysWithoutResolution),
				BugSLOs:                        []BugSLO{new},
			})
		}

	}

	return Projects
}
