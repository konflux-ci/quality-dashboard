/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import axios, { AxiosResponse } from 'axios';
import _ from 'lodash';
import { JobsStatistics, JobMetric } from '@app/utils/sharedComponents';
import { getRepoNameFormatted, sortGlobalSLI, teamIsNotEmpty } from '@app/utils/utils';
import { formatDateTime } from '@app/Reports/utils';

type ApiResponse = {
  code: number;
  data: any;
};

type RepositoriesApiResponse = {
  code: number;
  data: any;
  all: any;
};

const API_URL = process.env.REACT_APP_API_SERVER_URL || 'http://localhost:9898';

async function getVersion() {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/server/info';
  const uri = API_URL + subPath;

  try {
    await axios
      .get(uri)
      .then((res: AxiosResponse) => {
        if (res != undefined) {
          result.code = res.status;
          result.data = res.data;
        }
      })
      .catch((err) => {
        result.code = err.response.status;
        result.data = err.response.data;
      });
  } catch (error) {
    result.code = 400;
  }

  return result;
}

async function getJiras() {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/jira/bugs/all';
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getJirasResolutionTime(priority: string, team: string, rangeDateTime: Date[]) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/jira/bugs/metrics/resolution';
  const uri = API_URL + subPath;
  const start_date = formatDateTime(rangeDateTime[0]);
  const end_date = formatDateTime(rangeDateTime[1]);

  await axios
    .post(uri, {
      priority: priority,
      team_name: team,
      start: start_date,
      end: end_date,
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function listJiraProjects() {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/jira/project/list';
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getJirasOpen(priority: string, team: string, rangeDateTime: Date[]) {
  const start_date = formatDateTime(rangeDateTime[0]);
  const end_date = formatDateTime(rangeDateTime[1]);

  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/jira/bugs/metrics/open';
  const uri = API_URL + subPath;
  await axios
    .post(uri, {
      priority: priority,
      team_name: team,
      start: start_date,
      end: end_date,
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getRepositories(perPage = 50, team: string) {
  const REPOS_IN_PAGE = perPage;
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/repositories/list';
  const uri = API_URL + subPath;

  if (!teamIsNotEmpty(team)) return result;

  await axios
    .get(uri, {
      headers: {},
      params: {
        team_name: team,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {

        result.code = res.status;
        result.all = res.data;
        if (res.data.length >= REPOS_IN_PAGE) {
          result.data = _.chunk(res.data, REPOS_IN_PAGE);
        } else {
          result.data = _.chunk(res.data, res.data.length);
        }
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getAllRepositoriesWithOrgs(team: string, openshift: boolean, rangeDateTime: Date[]) {
  let repoAndOrgs = [];
  const start_date = formatDateTime(rangeDateTime[0]);
  const end_date = formatDateTime(rangeDateTime[1]);

  if (!teamIsNotEmpty(team)) return repoAndOrgs;

  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/repositories/list?team_name=' +
    team +
    '&openshift_ci=' +
    (openshift ? 'true' : 'false') +
    '&start_date=' +
    start_date +
    '&end_date=' +
    end_date
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  if (result.code != 200) {
    throw 'Error fetching data from server.';
  } else {
    result.data.sort((a, b) => (a.repository_name < b.repository_name ? -1 : 1));
    repoAndOrgs = result.data.map((row) => {
      return {
        repoName: row.repository_name,
        repoNameFormatted: getRepoNameFormatted(row.repository_name),
        organization: row.git_organization,
        description: row.description,
        url: row.git_url,
        coverage: row.code_coverage,
        prs: row.prs,
        workflows: row.workflows,
      };
    });
  }
  return repoAndOrgs;
}

async function getWorkflowByRepositoryName(repositoryName: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/workflows/get';
  const uri = API_URL + subPath;
  await axios
    .get(uri, {
      headers: {},
      params: {
        repository_name: repositoryName,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  return result;
}

async function createRepository(data = {}) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/repositories/create';
  const uri = API_URL + subPath;
  await axios
    .request({
      method: 'POST',
      url: uri,
      data: { ...data },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
      return result;
    });
  return result;
}

async function getLatestProwJob(repoName: string, repoOrg: string, jobType: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/prow/results/latest/get?repository_name=' +
    repoName +
    '&git_organization=' +
    repoOrg +
    '&job_type=' +
    jobType

  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  if (result.code != 200) {
    throw 'Error fetching data from server. ';
  }

  return result.data;
}

async function getProwJobStatistics(repoName: string, repoOrg: string, jobType: string, jobName: string, rangeDateTime: Date[]) {
  const start_date = formatDateTime(rangeDateTime[0]);
  const end_date = formatDateTime(rangeDateTime[1]);
  const result: ApiResponse = { code: 0, data: {} };
  const uri = API_URL + '/api/quality/prow/metrics/get?repository_name=' +
    repoName +
    '&git_organization=' +
    repoOrg +
    '&job_type=' +
    jobType +
    '&job_name=' +
    jobName +
    '&start_date=' +
    start_date +
    '&end_date=' +
    end_date;

  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  if (result.code != 200) {
    throw 'Error fetching data from server. ';
  }

  const statistics: JobsStatistics = result.data;

  return statistics;
}

async function getTeams() {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/teams/list/all';
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function createTeam(data = {}) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/teams/create';
  const uri = API_URL + subPath;
  await axios
    .request({
      method: 'POST',
      url: uri,
      data: { ...data },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getJobTypes(repoName: string, repoOrg: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/repositories/getJobTypesFromRepo?repository_name=' +
    repoName +
    '&git_organization=' +
    repoOrg
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  if (result.code != 200) {
    throw 'Error fetching data from server. ';
  }

  return result.data.sort((a, b) => (a < b ? -1 : 1));
}

async function getJobNamesAndTypes(repoName: string, repoOrg: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/prow/jobs/types?repository_name=' +
    repoName +
    '&git_organization=' +
    repoOrg
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  if (result.code != 200) {
    throw 'Error fetching data from server. ';
  }

  if (!result.data) {
    throw 'No jobs detected in OpenShift CI';
  }

  return result.data.sort((a, b) => (a < b ? -1 : 1));
}

// deleteInApi deletes data in the given subPath
// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
async function deleteInApi(data = {}, subPath: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const uri = API_URL + subPath;
  await axios
    .delete(uri, {
      headers: {},
      data: data,
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

// updateTeam updates a team in the database
async function updateTeam(data = {}) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/teams/put';
  const uri = API_URL + subPath;
  await axios
    .request({
      method: 'PUT',
      url: uri,
      data: { ...data },
      timeout: 120000,
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
      return result;
    });
  return result;
}

// checkDbConnection checks if the database is available
async function checkDbConnection() {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/database/ok';
  const uri = API_URL + subPath;

  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  return result;
}

async function getPullRequests(repoName: string, repoOrg: string, rangeDateTime: Date[]) {
  const start_date = formatDateTime(rangeDateTime[0]);
  const end_date = formatDateTime(rangeDateTime[1]);

  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/prs/get?repository_name=' +
    repoName +
    '&git_organization=' +
    repoOrg +
    '&start_date=' +
    start_date +
    '&end_date=' +
    end_date

  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  if (result.code != 200) {
    throw 'Error fetching data from server. ';
  }

  return result.data;
}

async function getProwJobs(repoName: string, repoOrg: string, rangeDateTime: Date[]) {
  const start_date = formatDateTime(rangeDateTime[0]);
  const end_date = formatDateTime(rangeDateTime[1]);
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/prow/results/get?repository_name=' + repoName +
    '&git_organization=' + repoOrg +
    '&start_date=' + start_date +
    '&end_date=' + end_date;
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  if (result.code != 200) {
    throw 'Error fetching data from server. ';
  }

  return result.data
}

async function listBugsAffectingCI(team: string, start_date: string, end_date: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/jira/bugs/e2e?team_name=' + team +
    '&start_date=' + start_date +
    '&end_date=' + end_date;
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getFailures(team: string, rangeDateTime: Date[]) {
  const start_date = formatDateTime(rangeDateTime[0]);
  const end_date = formatDateTime(rangeDateTime[1]);
  const result: ApiResponse = { code: 0, data: {} };
  const uri = API_URL +
    '/api/quality/failures/get?team_name=' + team +
    '&start_date=' +
    start_date +
    '&end_date=' +
    end_date;

  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.statusCode;
      result.data = err.data;
    });

  return result;
}

async function createFailure(team: string, jiraKey: string, errorMessage: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/failures/create';
  const uri = API_URL + subPath;
  await axios
    .request({
      method: 'POST',
      url: uri,
      data: {
        team: team,
        jira_key: jiraKey,
        error_message: errorMessage,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  return result;
}

async function bugExists(jiraKey: string, teamName: string) {
  const uri = API_URL + '/api/quality/jira/bugs/exist?team_name=' + teamName + '&jira_key=' + jiraKey
  const result: ApiResponse = { code: 0, data: {} };

  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.status;
      result.data = err.data;
    });

  return result;
}

// async function getBugSLIs(team: string, rangeDateTime: Date[]) {
async function getBugSLIs(team: string) {
  const result: ApiResponse = { code: 0, data: {} };
  // const start_date = formatDateTime(rangeDateTime[0]);
  // const end_date = formatDateTime(rangeDateTime[1]);

  const subPath = '/api/quality/jira/slis/list?team_name=' + team
  // '&start_date=' +
  // start_date +
  // '&end_date=' +
  // end_date

  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });


  sortGlobalSLI(result.data.resolution_time_sli.bugs)
  sortGlobalSLI(result.data.response_time_sli.bugs)
  sortGlobalSLI(result.data.triage_time_sli.bugs)

  return result;
}

async function getRepositoriesWithJobs(team: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/prow/repositories/list';
  const uri = API_URL + subPath;

  if (!teamIsNotEmpty(team)) return result;

  await axios
    .get(uri, {
      headers: {},
      params: {
        team_name: team,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getFlakyData(team: string, job: string, repo: string, rangeDateTime: Date[], gitOrg: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/suites/occurrences';
  const uri = API_URL + subPath;
  const startDate = rangeDateTime[0].toISOString();
  const endDate = rangeDateTime[1].toISOString();

  await axios
    .get(uri, {
      headers: {},
      params: {
        repository_name: repo,
        job_name: job,
        start_date: startDate,
        end_date: endDate,
        git_org: gitOrg,
        team_name: team
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getGlobalImpactData(team: string, job: string, repo: string, rangeDateTime: Date[], gitOrg: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/suites/flaky/trends';
  const uri = API_URL + subPath;
  const startDate = rangeDateTime[0].toISOString();
  const endDate = rangeDateTime[1].toISOString();

  await axios
    .get(uri, {
      headers: {},
      params: {
        repository_name: repo,
        job_name: job,
        start_date: startDate,
        end_date: endDate,
        git_org: gitOrg,
        team_name: team
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  return result;
}

async function listTeamRepos(team: string) {
  let repos = [];

  if (!teamIsNotEmpty(team)) return repos;

  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/prow/repositories/list?team_name=' + team
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  if (result.code != 200) {
    throw 'Error fetching data from server.';
  } else {
    repos = result.data.sort((a, b) => (a.repository_name < b.repository_name ? -1 : 1));
  }
  return repos;
}

async function getProwJobMetricsDaily(repoName: string, repoOrg: string, jobType: string, jobName: string, rangeDateTime: Date[]) {
  const start_date = formatDateTime(rangeDateTime[0]);
  const end_date = formatDateTime(rangeDateTime[1]);
  const result: ApiResponse = { code: 0, data: {} };
  const uri = API_URL + '/api/quality/prow/metrics/daily?repository_name=' +
    repoName +
    '&git_organization=' +
    repoOrg +
    '&job_type=' +
    jobType +
    '&job_name=' +
    jobName +
    '&start_date=' +
    start_date +
    '&end_date=' +
    end_date;

  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  if (result.code != 200) {
    throw 'Error fetching data from server. ';
  }

  const statistics: JobMetric[] = result.data;

  return statistics;
}

async function createUser(userEmail: string, userConfig: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/users/create';
  const uri = API_URL + subPath;
  await axios
    .request({
      method: 'POST',
      url: uri,
      data: {
        user_email: userEmail,
        user_config: userConfig,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  return result;
}

// useful for debug in Console
async function listUsers() {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/users/get/all';
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getUser(userEmail: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/users/get/user';
  const uri = API_URL + subPath;

  await axios
    .get(uri, {
      headers: {},
      params: {
        user_email: userEmail,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function checkGithubRepositoryUrl(repoOrg: string, repoName: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/repositories/verify?repository_name=' +
    repoName +
    '&git_organization=' +
    repoOrg
  const uri = API_URL + subPath;

  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  return result;
}

async function checkGithubRepositoryExists(repoOrg: string, repoName: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/repositories/exists?repository_name=' +
    repoName +
    '&git_organization=' +
    repoOrg
  const uri = API_URL + subPath;

  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  return result;
}

async function getConfiguration(teamName: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/teams/get/configuration';
  const uri = API_URL + subPath;

  await axios
    .get(uri, {
      headers: {},
      params: {
        team_name: teamName,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function isJqlQueryValid(jqlQuery: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/jira/jql-query/valid';
  const uri = API_URL + subPath;

  await axios
    .get(uri, {
      headers: {},
      params: {
        jql_query: jqlQuery,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getJiraKeysByTeam(teamName: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/teams/get/jira_keys';
  const uri = API_URL + subPath;

  await axios
    .get(uri, {
      headers: {},
      params: {
        team_name: teamName,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getTeam(teamName: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/teams/get/team';
  const uri = API_URL + subPath;

  await axios
    .get(uri, {
      headers: {},
      params: {
        team_name: teamName,
      },
    })
    .then((res: AxiosResponse) => {
      if (res != undefined) {
        result.code = res.status;
        result.data = res.data;
      }
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  return result;
}

export {
  getVersion,
  getRepositories,
  createRepository,
  getWorkflowByRepositoryName,
  getAllRepositoriesWithOrgs,
  listTeamRepos,
  getLatestProwJob,
  getProwJobStatistics,
  getProwJobs,
  getTeams,
  createTeam,
  getJiras,
  getJobTypes,
  getJobNamesAndTypes,
  deleteInApi,
  updateTeam,
  checkDbConnection,
  getJirasResolutionTime,
  getJirasOpen,
  listJiraProjects,
  getPullRequests,
  listBugsAffectingCI,
  getFailures,
  createFailure,
  bugExists,
  getBugSLIs,
  getRepositoriesWithJobs,
  getFlakyData,
  getGlobalImpactData,
  getProwJobMetricsDaily,
  createUser,
  listUsers,
  getUser,
  checkGithubRepositoryUrl,
  checkGithubRepositoryExists,
  getConfiguration,
  isJqlQueryValid,
  getJiraKeysByTeam,
  getTeam
};
