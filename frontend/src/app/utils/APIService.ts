/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import axios, { AxiosResponse } from 'axios';
import _ from 'lodash';
import { JobsStatistics } from '@app/utils/sharedComponents';
import { sortGlobalSLI, teamIsNotEmpty } from '@app/utils/utils';
import { formatDate } from '@app/Reports/utils';
import { PrsStatistics } from '@app/Github/PullRequests';
import { Job } from '@app/Reports/FailedE2ETests';

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
        result.code = res.status;
        result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
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
  const start_date = formatDate(rangeDateTime[0]);
  const end_date = formatDate(rangeDateTime[1]);

  await axios
    .post(uri, {
      priority: priority,
      team_name: team,
      start: start_date,
      end: end_date,
    })
    .then((res: AxiosResponse) => {
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getJirasOpen(priority: string, team: string, rangeDateTime: Date[]) {
  const start_date = formatDate(rangeDateTime[0]);
  const end_date = formatDate(rangeDateTime[1]);

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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.all = res.data;
      if (res.data.length >= REPOS_IN_PAGE) {
        result.data = _.chunk(res.data, REPOS_IN_PAGE);
      } else {
        result.data = _.chunk(res.data, res.data.length);
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
  const start_date = formatDate(rangeDateTime[0]);
  const end_date = formatDate(rangeDateTime[1]);

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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
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

async function getProwJobStatistics(repoName: string, repoOrg: string, jobType: string, rangeDateTime: Date[]) {
  const start_date = formatDate(rangeDateTime[0]);
  const end_date = formatDate(rangeDateTime[1]);
  const result: ApiResponse = { code: 0, data: {} };
  const uri = API_URL + '/api/quality/prow/metrics/get?repository_name=' +
    repoName +
    '&git_organization=' +
    repoOrg +
    '&job_type=' +
    jobType +
    '&start_date=' +
    start_date +
    '&end_date=' +
    end_date;

  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      result.code = res.status;
      result.data = res.data;
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  if (result.code != 200) {
    throw 'Error fetching data from server. ';
  }

  const statistics: JobsStatistics = result.data;
  if (statistics.jobs == null) {
    throw 'No jobs detected in OpenShift CI';
  }

  statistics.jobs.forEach((job, j_idx) => {
    let j = job.metrics.sort(function (a, b) {
      return +new Date(a.date) - +new Date(b.date);
    });
    statistics.jobs[j_idx].metrics = j;
  });

  return statistics;
}

async function getTeams() {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/teams/list/all';
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });

  return result;
}

async function getPullRequests(repoName: string, repoOrg: string, rangeDateTime: Date[]) {
  const start_date = formatDate(rangeDateTime[0]);
  const end_date = formatDate(rangeDateTime[1]);

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
      result.code = res.status;
      result.data = res.data;
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

async function getProwJobs(repoName: string, repoOrg: string) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/prow/results/get?repository_name=' + repoName + '&git_organization=' + repoOrg
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      result.code = res.status;
      result.data = res.data;
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

async function listE2EBugsKnown() {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/jira/bugs/e2e';
  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      result.code = res.status;
      result.data = res.data;
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getFailures(team: string, rangeDateTime: Date[]) {
  const start_date = formatDate(rangeDateTime[0]);
  const end_date = formatDate(rangeDateTime[1]);
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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
    })
    .catch((err) => {
      result.code = err.status;
      result.data = err.data;
    });

  return result;
}

async function getBugSLIs(team: string, rangeDateTime: Date[]) {
  const result: ApiResponse = { code: 0, data: {} };
  const start_date = formatDate(rangeDateTime[0]);
  const end_date = formatDate(rangeDateTime[1]);

  const subPath = '/api/quality/jira/slis/list?team_name=' + team +
    '&start_date=' +
    start_date +
    '&end_date=' +
    end_date

  const uri = API_URL + subPath;
  await axios
    .get(uri)
    .then((res: AxiosResponse) => {
      result.code = res.status;
      result.data = res.data;
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
      result.code = res.status;
      result.data = res.data;
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getFlakyData(team:string, job: string, repo: string, startDate: string, endDate: string, gitOrg: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/suites/ocurrencies';
  const uri = API_URL + subPath;

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
      result.code = res.status;
      result.data = res.data;
    })
    .catch((err) => {
      result.code = err.response.status;
      result.data = err.response.data;
    });
  return result;
}

async function getGlobalImpactData(team:string, job: string, repo: string, startDate: string, endDate: string, gitOrg: string) {
  const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
  const subPath = '/api/quality/suites/flaky/trends';
  const uri = API_URL + subPath;

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
      result.code = res.status;
      result.data = res.data;
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
  getLatestProwJob,
  getProwJobStatistics,
  getProwJobs,
  getTeams,
  createTeam,
  getJiras,
  getJobTypes,
  deleteInApi,
  updateTeam,
  checkDbConnection,
  getJirasResolutionTime,
  getJirasOpen,
  listJiraProjects,
  getPullRequests,
  listE2EBugsKnown,
  getFailures,
  createFailure,
  bugExists,
  getBugSLIs,
  getRepositoriesWithJobs,
  getFlakyData,
  getGlobalImpactData
};
