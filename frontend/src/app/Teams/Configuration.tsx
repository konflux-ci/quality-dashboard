/* eslint-disable @typescript-eslint/explicit-module-boundary-types */

export interface Configuration {
    team_name: string;
    jira_config: JiraConfig;
    bug_slos_config: string;
}

export interface JiraConfig {
    bugs_collect_query: string;
}

export const generateJiraConfig = (query) => {
    const cfg: JiraConfig = {
        bugs_collect_query: query
    };
    const config = JSON.stringify(cfg);

    return config
}
