/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import axios, { AxiosResponse, AxiosError } from "axios";
import _ from 'lodash';
import {JobsStatistics} from '@app/utils/sharedComponents'

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

async function getProwJobStatistics(repoName: string, repoOrg:string, jobType:string){
    
    const response = await fetch("http://127.0.0.1:9898/api/quality/prow/metrics/get?repository_name="+repoName+"&git_organization="+repoOrg+"&job_type="+jobType)
    if(!response.ok){
        throw "Error fetching data from server. "
    }
    const statistics:JobsStatistics = await response.json()
    
    statistics.jobs.forEach((job, j_idx) => {
      let j = job.metrics.sort(function(a,b){
        return +new Date(a.date) - +new Date(b.date);
      })
      statistics.jobs[j_idx].metrics = j
    })
    
    return statistics
      
}

export { getVersion, getRepositories, createRepository, deleteRepositoryAPI, getWorkflowByRepositoryName, getAllRepositoriesWithOrgs, getLatestProwJob, getProwJobStatistics}
