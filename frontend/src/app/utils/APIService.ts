/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
//import { Repositories } from "@app/Repositories/Repositories";
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
    const subPath ='/api/version';
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

<<<<<<< Updated upstream
export { getVersion, getRepositories, createRepository, deleteRepositoryAPI, getWorkflowByRepositoryName }
=======
export { getVersion, getRepositories, createRepository, deleteRepositoryAPI, getWorkflowByRepositoryName, getJiras }
>>>>>>> Stashed changes
