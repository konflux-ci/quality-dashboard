/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import axios, { AxiosResponse, AxiosError } from 'axios';
import _ from 'lodash';
import { JobsStatistics } from '@app/utils/sharedComponents';
import { teamIsNotEmpty } from '@app/utils/utils';

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

async function getJiras() {
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

async function getRepositories(perPage = 5, team: string) {
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

async function getAllRepositoriesWithOrgs(team: string, openshift: boolean) {
  let repoAndOrgs = [];

  if (!teamIsNotEmpty(team)) return repoAndOrgs;

  const response = await fetch(
    API_URL + '/api/quality/repositories/list?team_name=' + team + '&openshift_ci=' + (openshift ? 'true' : 'false')
  );

  if (!response.ok) {
    throw 'Error fetching data from server';
  } else {
    const data = await response.json();
    data.sort((a, b) => (a.repository_name < b.repository_name ? -1 : 1));
    repoAndOrgs = data.map((row, index) => {
      return { repoName: row.repository_name, organization: row.git_organization };
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

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
async function deleteRepositoryAPI(data = {}) {
  const result: ApiResponse = { code: 0, data: {} };
  const subPath = '/api/quality/repositories/delete';
  const uri = API_URL + subPath;
  await axios
    .delete(uri, {
      headers: {},
      data: data,
    })
    .then((res: AxiosResponse) => {
      console.log(res);
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
  const response = await fetch(
    API_URL +
      '/api/quality/prow/results/latest/get?repository_name=' +
      repoName +
      '&git_organization=' +
      repoOrg +
      '&job_type=' +
      jobType
  );
  if (!response.ok) {
    throw 'Error fetching data from server. ';
  }
  const data = await response.json();
  return data;
}

async function getProwJobStatistics(repoName: string, repoOrg: string, jobType: string) {
  const response = await fetch(
    API_URL +
      '/api/quality/prow/metrics/get?repository_name=' +
      repoName +
      '&git_organization=' +
      repoOrg +
      '&job_type=' +
      jobType
  );
  if (!response.ok) {
    throw 'Error fetching data from server. ';
  }
  const statistics: JobsStatistics = await response.json();
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
  const response = await fetch(
    API_URL +
      '/api/quality/repositories/getJobTypesFromRepo?repository_name=' +
      repoName +
      '&git_organization=' +
      repoOrg
  );
  if (!response.ok) {
    throw 'Error fetching data from server. ';
  }
  const data = await response.json();
  return data.sort((a, b) => (a < b ? -1 : 1));
}

export {
  getVersion,
  getRepositories,
  createRepository,
  deleteRepositoryAPI,
  getWorkflowByRepositoryName,
  getAllRepositoriesWithOrgs,
  getLatestProwJob,
  getProwJobStatistics,
  getTeams,
  createTeam,
  getJiras,
  getJobTypes,
};
