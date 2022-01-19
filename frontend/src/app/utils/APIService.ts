/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import axios, { AxiosResponse, AxiosError } from "axios";

type ApiResponse = {
    code: number;
    data: any;
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

async function getRepositories(){
    const result: ApiResponse = { code: 0, data: {} };
    const subPath ='/api/quality/repositories/list';
    const uri = API_URL + subPath;
    await axios.get(uri, {
      }).then((res: AxiosResponse) => {
        result.code = res.status;
        result.data = res.data;
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

export { getVersion, getRepositories, createRepository, deleteRepositoryAPI, getWorkflowByRepositoryName }
