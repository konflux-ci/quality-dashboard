export interface Alert {
    alert_message: string,
    signal: string,
}

export interface Bug {
    jira_key: string,
    summary: string,
    status: string,
    labels: string,
    priority: string,
    jira_url: string,
    triage_sli: Alert,
    response_sli: Alert,
    resolution_sli: Alert,
    global_sli: string,
    days_without_assignee: number,
    days_without_priority: number,
    days_without_resolution: number,
}

export interface GlobalSLI {
    green_sli: number,
    yellow_sli: number,
    red_sli: number,
}

export interface SLI {
    bugs: Bug[],
}

export interface Info {
    global_sli: GlobalSLI;
    response_time_sli: SLI;
    triage_time_sli: SLI;
    resolution_time_sli: SLI;
}