package jira

import (
	"fmt"
	"math"
	"strings"

	"github.com/redhat-appstudio/quality-studio/pkg/storage/ent/db"
)

// GetPriorityTriageSLI should return red when priority is not defined for more than 2 days on untriaged bugs
// GetPriorityTriageSLI should return yellow when priority is not defined for more than 1 day but less than 2 days on untriaged bugs
func GetPriorityTriageSLI(bug *db.Bugs) *Alert {
	alert := &Alert{Signal: "green"}

	if bug.Labels != nil && bug.DaysWithoutPriority != nil && strings.Contains(*bug.Labels, "untriaged") {
		msg := fmt.Sprintf("Priority should be defined between a maximum of 2 days. This issue has priority undefined for %.2f days. Please, take a time to prioritize it.\n\n",
			*bug.DaysWithoutPriority,
		)

		if *bug.DaysWithoutPriority > 2 {
			alert.Signal = "red"
			msg = fmt.Sprintf("Issue <%s|%s> is not meeting defined Bug SLO for Priority Triage Time. %s",
				bug.URL,
				bug.JiraKey,
				msg,
			)
			alert.AlertMessage = &msg
		} else if *bug.DaysWithoutPriority > 1 {
			alert.Signal = "yellow"
			msg = fmt.Sprintf("Issue <%s|%s> is almost not meeting defined Bug SLO for Priority Triage Time. %s",
				bug.URL,
				bug.JiraKey,
				msg,
			)
			alert.AlertMessage = &msg
		}
	}

	return alert
}

// GetResponseSLI should return red if there is no assignee for more than 2 days on Blocker and Critical bugs
func GetResponseSLI(bug *db.Bugs) *Alert {
	alert := &Alert{Signal: "green"}

	if bug.DaysWithoutAssignee != nil && (bug.Priority == "Blocker" || bug.Priority == "Critical") {
		if *bug.DaysWithoutAssignee > 2 {
			msg := fmt.Sprintf("Issue <%s|%s> is not meeting defined Bug SLO for Blocker and Critical Bug Response Time. Blocker and Critical bugs should be assigned between a maximum of 2 days. This issue has assignee undefined for %.2f days. Please, take a time to assign it.\n\n",
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

// GetComponentAssignmentTriageSLI should return red if there is no component assigned for more than 1 day
func GetComponentAssignmentTriageSLI(bug *db.Bugs) *Alert {
	alert := &Alert{Signal: "green"}

	if bug.DaysWithoutComponent != nil {
		if *bug.DaysWithoutComponent > 1 {
			msg := fmt.Sprintf("Issue <%s|%s> is not meeting defined Bug SLO for Component Assignment Triage Time. Bugs should have component assigned between a maximum of 1 day. This issue has component undefined for %.2f days. Please, take a time to assign a component to it.\n\n",
				bug.URL,
				bug.JiraKey,
				*bug.DaysWithoutComponent,
			)
			alert.AlertMessage = &msg
			alert.Signal = "red"
		}
	}

	return alert
}

func measureBugResolutionSLI(redThreshold, yellowThreshold float64, bug *db.Bugs) *Alert {
	alert := &Alert{Signal: "green"}

	if bug.DaysWithoutResolution != nil {
		daysWithoutResolution := *bug.DaysWithoutResolution

		msg := fmt.Sprintf("%s bugs should not take more than %g days to be resolved. This issue is not resolved for %.2f days. Please, take a time to resolve it.\n\n",
			bug.Priority,
			redThreshold,
			daysWithoutResolution,
		)

		if daysWithoutResolution > redThreshold {
			alert.Signal = "red"
			msg = fmt.Sprintf("Issue <%s|%s> is not meeting defined Bug SLO for %s Bug Resolution Time. %s",
				bug.URL,
				bug.JiraKey,
				bug.Priority,
				msg,
			)
			alert.AlertMessage = &msg
		} else if daysWithoutResolution > yellowThreshold {
			alert.Signal = "yellow"
			msg = fmt.Sprintf("Issue <%s|%s> is almost not meeting defined Bug SLO for %s Bug Resolution Time. %s",
				bug.URL,
				bug.JiraKey,
				bug.Priority,
				msg,
			)
			alert.AlertMessage = &msg
		}
	}

	return alert
}

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

func GetMetric(currentInfo Metric, currentSignal, targetSignal string, value float64) Metric {
	if currentSignal == targetSignal {
		currentInfo.Total++
		currentInfo.Sum += value
		currentInfo.Average = currentInfo.Sum / float64(currentInfo.Total)
		currentInfo.Average = math.Round(currentInfo.Average*100) / 100
	}

	return currentInfo
}

func GetGlobalSLI(triageSLI, responseSLI, resolutionSLI, componentSLI string) string {
	if triageSLI == "red" || responseSLI == "red" || resolutionSLI == "red" || componentSLI == "red" {
		return "red"
	} else if triageSLI == "yellow" || responseSLI == "yellow" || resolutionSLI == "yellow" || componentSLI == "yellow" {
		return "yellow"
	}

	return "green"
}

func GetSLI(info GlobalSLI, targetSLI string) GlobalSLI {
	switch targetSLI {
	case "red":
		info.RedSLI++
	case "yellow":
		info.YellowSLI++
	default:
		info.GreenSLI++
	}
	return info
}

func GetBugSLIs(bugs []*db.Bugs) BugSlisInfo {
	info := BugSlisInfo{
		GlobalSLI:                    GlobalSLI{GreenSLI: 0, YellowSLI: 0, RedSLI: 0},
		TriageTimeSLI:                SLI{Bugs: []Bug{}, Red: Metric{}, Yellow: Metric{}},
		ResponseTimeSLI:              SLI{Bugs: []Bug{}, Red: Metric{}, Yellow: Metric{}},
		ResolutionTimeSLI:            SLI{Bugs: []Bug{}, Red: Metric{}, Yellow: Metric{}},
		ComponentAssignmentTriageSLI: SLI{Bugs: []Bug{}, Red: Metric{}, Yellow: Metric{}},
		Bugs:                         []Bug{},
	}

	for _, bug := range bugs {
		TriageSLI := GetPriorityTriageSLI(bug)
		ComponentAssignmentTriageSLI := GetComponentAssignmentTriageSLI(bug)
		ResponseSLI := GetResponseSLI(bug)
		ResolutionSLI := GetResolutionSLI(bug)
		GlobalSLI := GetGlobalSLI(TriageSLI.Signal, ResponseSLI.Signal, ResolutionSLI.Signal, ComponentAssignmentTriageSLI.Signal)

		new := Bug{
			JiraKey:                      bug.JiraKey,
			JiraURL:                      bug.URL,
			Priority:                     bug.Priority,
			Status:                       bug.Status,
			Summary:                      bug.Summary,
			TriageSLI:                    TriageSLI,
			ResponseSLI:                  ResponseSLI,
			ResolutionSLI:                ResolutionSLI,
			ComponentAssignmentTriageSLI: ComponentAssignmentTriageSLI,
			GlobalSLI:                    GlobalSLI,
			DaysWithoutAssignee:          bug.DaysWithoutAssignee,
			DaysWithoutPriority:          bug.DaysWithoutPriority,
			DaysWithoutResolution:        bug.DaysWithoutResolution,
			DaysWithoutComponent:         bug.DaysWithoutComponent,
		}

		if bug.Labels != nil {
			new.Labels = *bug.Labels
		} else {
			new.Labels = ""
		}

		if bug.Component != nil {
			new.Component = *bug.Component
		} else {
			new.Component = ""
		}

		if GlobalSLI == "green" {
			info.GlobalSLI.GreenSLI++
			continue
		}

		info.Bugs = append(info.Bugs, new)

		if TriageSLI.Signal != "green" {
			info.TriageTimeSLI.Bugs = append(info.TriageTimeSLI.Bugs, new)
		}

		if ResponseSLI.Signal != "green" {
			info.ResponseTimeSLI.Bugs = append(info.ResponseTimeSLI.Bugs, new)
		}

		if ResolutionSLI.Signal != "green" {
			info.ResolutionTimeSLI.Bugs = append(info.ResolutionTimeSLI.Bugs, new)
		}

		if ComponentAssignmentTriageSLI.Signal != "green" {
			info.ComponentAssignmentTriageSLI.Bugs = append(info.ComponentAssignmentTriageSLI.Bugs, new)
		}

		info.GlobalSLI = GetSLI(info.GlobalSLI, GlobalSLI)
	}

	return info
}
