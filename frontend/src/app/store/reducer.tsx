import { ITeam } from '@app/Teams/TeamsSelect';
import { saveStateContext } from '@app/utils/utils';
import { combineReducers } from 'redux';

export interface StateContext {
    general: {
        APIData: []
        error: string
        workflows: []
    },
    alerts: {
        alerts: []
        version: string
    }, jira: {
        E2E_KNOWN_ISSUES: []
    },
    repos: {
        repositories: []
        Allrepositories: []
    },
    teams: {
        Team: string
        TeamsAvailable: ITeam[]
    },
    auth : {
        AT: string,
        RT: string,
        IDT: string, 
        AT_expiration: Date,
        Username: string,
    }, 
}

const generalReducer = (state, action) => {
    switch (action.type) {
        case 'APIData':
            return {
                ...state,
                APIData: action.data
            };
        case 'SET_Version':
            return {
                ...state,
                version: action.data
            };
        case 'SET_WORKFLOWS':
            return {
                ...state,
                workflows: action.data
            };
        case 'SET_ERROR':
            return {
                ...state,
                error: action.data
            };

        default: return state || null;
    }
};

const jirasReducer = (state, action) => {
    switch (action.type) {
        case 'SET_JIRAS':
            return {
                ...state,
                E2E_KNOWN_ISSUES: action.data
            };
        default: return state || null;
    }
};

const repositoriesReducer = (state, action) => {
    switch (action.type) {
        case 'SET_REPOSITORIES':
            return {
                ...state,
                repositories: action.data
            };
        case 'SET_REPOSITORIES_ALL':
            return {
                ...state,
                Allrepositories: action.data
            };
        default: return state || null;
    }
};

const alertsReducer = (state, action) => {
    switch (action.type) {
        case 'ADD_Alert':
            return {
                ...state,
                alerts: state.alerts.concat(action.data)
            }
        case 'REMOVE_Alert':
            return {
                ...state,
                alerts: state.alerts.filter(el => el.key !== action.data)
            };
        default: return state || null;
    }
};

const teamsReducer = (state, action) => {
    switch (action.type) {
        case 'SET_TEAM':
            // Change the persisted 'saved' team when its state has been changed
            saveStateContext('TEAM', action.data)
            return {
                ...state,
                Team: action.data
            };
        case 'SET_TEAMS_AVAILABLE':
            return {
                ...state,
                TeamsAvailable: action.data
            };
        default: return state || null;
    }
};

const authReducer = (state, action) => {
    switch (action.type) {
        case 'SET_ACCESS_TOKEN':
            // Change the persisted 'saved' team when its state has been changed
            saveStateContext('AT', action.data)
            return {
                ...state,
                AT: action.data
            };
        case 'SET_REFRESH_TOKEN':
            // Change the persisted 'saved' team when its state has been changed
            saveStateContext('RT', action.data)
            return {
                ...state,
                RT: action.data
            };
        case 'SET_ID_TOKEN':
            saveStateContext('IDT', action.data)
            return {
                ...state,
                IDT: action.data
            };
        case 'SET_AT_EXPIRATION':
            saveStateContext('AT_EXPIRATION', action.data)
            return {
                ...state,
                AT_expiration: action.data
            };
        case 'SET_USERNAME':
            saveStateContext('USERNAME', action.data)
            return {
                ...state,
                Username: action.data
            };
        default: return state || null;
    }
};

export const rootReducer = combineReducers({
    general: generalReducer,
    jiras: jirasReducer,
    repos: repositoriesReducer,
    alerts: alertsReducer,
    teams: teamsReducer,
    auth: authReducer,
});

export default rootReducer;



