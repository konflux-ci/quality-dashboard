/* eslint-disable @typescript-eslint/explicit-module-boundary-types */

export interface Configuration {
    team_name: string;
    jira_config: JiraConfig;
    bug_slos_config: string;
}

export interface JiraConfig {
    bugs_collect_query: string;
    ci_impact_query: string;
}

export const generateJiraConfig = (bugsCollectQuery, ciImpactQuery) => {
    const cfg: JiraConfig = {
        bugs_collect_query: bugsCollectQuery,
        ci_impact_query: ciImpactQuery
    };
    const config = JSON.stringify(cfg);

    return config
}
