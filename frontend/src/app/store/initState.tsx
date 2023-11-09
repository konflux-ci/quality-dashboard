import { loadStateContext, stateContextExists } from '@app/utils/utils'

export const initialState = { 
    general : {
        APIData: [],
        error: 'error',
        workflows: []
    }, 
    alerts : {
        alerts: [],
        version: ''
    }, jiras : {
        E2E_KNOWN_ISSUES: []
    },
    repos : {
        repositories: [],
        Allrepositories: []
    }, 
    teams : {
        Team: "",
        TeamsAvailable: [],
        InstalledPlugins: [],
        FlattenedPlugins: []
    },
    auth : {
        AT: stateContextExists("AT") ? loadStateContext("AT") : "",
        RT: stateContextExists("RT") ? loadStateContext("RT") : "",
        IDT: stateContextExists("IDT") ? loadStateContext("IDT") : "",
        AT_expiration: stateContextExists("AT_EXPIRATION") ? loadStateContext("AT_EXPIRATION") : 0,
        Username: stateContextExists("USERNAME") ? loadStateContext("USERNAME") : "",
    }, 
};