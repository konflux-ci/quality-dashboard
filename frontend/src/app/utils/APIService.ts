/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import axios, { AxiosResponse, AxiosError } from "axios";
import _ from 'lodash';

type ApiResponse = {
    code: number;
    data: any;
};

type RepositoriesApiResponse = {
    code: number;
    data: any;
    all: any;
};

const API_URL = (process.env.REACT_APP_API_SERVER_URL || 'http://localhost:9898')

async function getVersion(){
    const result: ApiResponse = { code: 0, data: {} };
    const subPath ='/api/quality/server/info';
    const uri = API_URL + subPath;
    await axios.get(uri).then((res: AxiosResponse) => {
        result.code = res.status;
        result.data = res.data;
    }).catch((err) => {
        result.code = err.response.status;
        result.data = err.response.data;
    });
    return result;
}

async function getRepositories(perPage = 5){
    const REPOS_IN_PAGE = perPage 
    const result: RepositoriesApiResponse = { code: 0, data: [], all: [] };
    const subPath ='/api/quality/repositories/list';
    const uri = API_URL + subPath;
    await axios.get(uri, {
      }).then((res: AxiosResponse) => {
        result.code = res.status;
        result.all = res.data;
        if(res.data.length >= REPOS_IN_PAGE){
            result.data = _.chunk(res.data, REPOS_IN_PAGE);
        }
        else{
            result.data = _.chunk(res.data, res.data.length)
        }
    }).catch((err) => {
        result.code = err.response.status;
        result.data = err.response.data;
    });
    return result;
}

async function getAllRepositoriesWithOrgs(){
    const subPath ='/api/quality/repositories/list';
    const uri = API_URL + subPath;
    const response = await fetch(uri)
    let repoAndOrgs = []
    if (!response.ok) {
        throw "Error fetching data from server"
    }
    else {
        const data = await response.json()
        repoAndOrgs = data.map((row, index) => {return {"repoName": row.repository_name, "organization": row.git_organization }})
    }
    return repoAndOrgs
}

async function getWorkflowByRepositoryName(repositoryName:string){
    const result: ApiResponse = { code: 0, data: {} };
    const subPath ='/api/quality/workflows/get';
    const uri = API_URL + subPath;
    await axios.get(uri, {
        headers: {},
        params: {
            "repository_name": repositoryName
        }
      }).then((res: AxiosResponse) => {
        result.code = res.status;
        result.data = res.data;
    }).catch((err) => {
        result.code = err.response.status;
        result.data = err.response.data;
    });

    return result;
}

// eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types
async function deleteRepositoryAPI(data = {}){
    const result: ApiResponse = { code: 0, data: {} };
    const subPath ='/api/quality/repositories/delete';
    const uri = API_URL + subPath;
    await axios.delete(uri, {
        headers: {},
        data: data
      }).then((res: AxiosResponse) => {
        console.log(res)
        result.code = res.status;
        result.data = res.data;
    }).catch((err) => {
        result.code = err.response.status;
        result.data = err.response.data;
    });
    return result;
}

async function createRepository(data = {}) {
    const result: ApiResponse = { code: 0, data: {} };
    const subPath ='/api/quality/repositories/create';
    const uri = API_URL + subPath;
    axios.request({
        method: 'POST',
        url: uri,
        data: {...data},
      }).then((res: AxiosResponse) => {
        result.code = res.status;
        result.data = res.data;
    }).catch((err) => {
        result.code = err.response.status;
        result.data = err.response.data;
    });
    return result;
}

async function getLatestProwJob(repoName: string, repoOrg:string, jobType:string){
    const response = await fetch("http://127.0.0.1:9898/api/quality/prow/results/latest/get?repository_name="+repoName+"&git_organization="+repoOrg+"&job_type="+jobType)
    if(!response.ok){
        throw "Error fetching data from server. "
    }
    const data = await response.json()
    return data
}

async function getProwJobStatisticsMOCK(repoName: string, repoOrg:string, jobType:string){
    const statistics = {
        "repository_name": "infra-deployments",
        "type": "periodic",
        "git_org": "redhat-appstudio",
        "jobs": [
          {
            "name": "periodic-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-deployment-periodic",
            "metrics": [
              {
                "success_rate": "99.00",
                "failure_rate": "1.00",
                "ci_failed_rate": "2.00",
                "date": "2022-08-13 07:14:29.000 +0200"
              },
              {
                "success_rate": "96.00",
                "failure_rate": "4.00",
                "ci_failed_rate": "5.00",
                "date": "2022-08-12 09:54:44.000 +0200"
              },
              {
                "success_rate": "98.00",
                "failure_rate": "2.00",
                "ci_failed_rate": "5.00",
                "date": "2022-08-11 01:54:44.000 +0200"
              },
              {
                "success_rate": "100.00",
                "failure_rate": "0.00",
                "ci_failed_rate": "1.00",
                "date": "2022-08-10 07:14:29.000 +0200"
              },
              {
                "success_rate": "97.00",
                "failure_rate": "3.00",
                "ci_failed_rate": "2.00",
                "date": "2022-08-09 09:54:44.000 +0200"
              },
              {
                "success_rate": "80.00",
                "failure_rate": "20.00",
                "ci_failed_rate": "5.00",
                "date": "2022-08-08 01:54:44.000 +0200"
              },
              {
                "success_rate": "99.00",
                "failure_rate": "1.00",
                "ci_failed_rate": "2.00",
                "date": "2022-08-07 07:14:29.000 +0200"
              },
              {
                "success_rate": "96.00",
                "failure_rate": "4.00",
                "ci_failed_rate": "3.00",
                "date": "2022-08-06 09:54:44.000 +0200"
              },
              {
                "success_rate": "89.00",
                "failure_rate": "11.00",
                "ci_failed_rate": "5.00",
                "date": "2022-08-05 01:54:44.000 +0200"
              },
              {
                "success_rate": "100.00",
                "failure_rate": "0.00",
                "ci_failed_rate": "0.00",
                "date": "2022-08-04 07:14:29.000 +0200"
              },
              {
                "success_rate": "98.00",
                "failure_rate": "2.00",
                "ci_failed_rate": "1.00",
                "date": "2022-08-03 09:54:44.000 +0200"
              },
              {
                "success_rate": "99.00",
                "failure_rate": "1.00",
                "ci_failed_rate": "2.00",
                "date": "2022-08-02 01:54:44.000 +0200"
              },
              {
                "success_rate": "99.99",
                "failure_rate": "0.01",
                "ci_failed_rate": "2.00",
                "date": "2022-08-01 07:14:29.000 +0200"
              }
            ],
            "summary": {
              "success_rate_avg": "85.15",
              "failure_rate_avg": "20.07",
              "ci_failed_rate_avg": "3.5",
              "date_from": "2022-08-01 07:14:29.000 +0200",
              "date_to": "2022-08-13 07:14:29.000 +0200"
            }
          },
          {
            "name": "periodic-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-deployment-periodic-2",
            "metrics": [
              {
                "success_rate": "66.00",
                "failure_rate": "2.00",
                "ci_failed_rate": "5.00",
                "date": "2022-08-03 09:54:44.000 +0200"
              },
              {
                "success_rate": "89.00",
                "failure_rate": "78.00",
                "ci_failed_rate": "5.00",
                "date": "2022-08-02 01:54:44.000 +0200"
              },
              {
                "success_rate": "1.00",
                "failure_rate": "100.00",
                "ci_failed_rate": "60.00",
                "date": "2022-08-01 07:14:29.000 +0200"
              }
            ],
            "summary": {
              "success_rate_avg": "52",
              "failure_rate_avg": "0",
              "ci_failed_rate_avg": "23.33",
              "date_from": "2022-08-02 07:14:29.000 +0200",
              "date_to": "2022-08-13 07:14:29.000 +0200"
            }
          }
        ]
    }
    
    statistics.jobs.forEach((job, j_idx) => {
      let j = job.metrics.sort(function(a,b){
        return +new Date(a.date) - +new Date(b.date);
      })
      console.log(j)
      statistics.jobs[j_idx].metrics = j
    })
    
    return statistics
      
}

export { getVersion, getRepositories, createRepository, deleteRepositoryAPI, getWorkflowByRepositoryName, getAllRepositoriesWithOrgs, getLatestProwJob, getProwJobStatisticsMOCK}
