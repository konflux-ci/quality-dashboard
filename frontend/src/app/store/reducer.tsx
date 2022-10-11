import { ITeam } from '@app/Teams/TeamsSelect';
import { initial, reduceRight } from 'lodash';
import { combineReducers } from 'redux';
//import { initialState } from './store';
import { Reducer } from 'react';
import { configureStore } from '@reduxjs/toolkit';


export interface StateContext {
    general : {
        APIData: [];
        error: string;
        workflows: []
    }, 
    alerts : {
        alerts: [];
        version: string;
    }, jira : {
        E2E_KNOWN_ISSUES: [];
    },
    repos : {
        repositories: [];
        Allrepositories: [];
    }, 
    teams : {
        Team: string;
        TeamsAvailable: ITeam[];
    }, 
};

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
    };

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
            return {
                ...state,
                Team: action.data
            };
        case 'SET_TEAMS_AVAILABLE':
            return {
                ...state,
                TeamsAvailable: action.data
            };
        case 'SET_TEAMS_AVAIBALE':
            return {
                ...state,
                Team: action.data
            };
        default: return state || null;
    }
};

export const rootReducer = combineReducers({
    general : generalReducer,
    jiras : jirasReducer,
    repos : repositoriesReducer,
    alerts : alertsReducer,
    teams : teamsReducer
});

export default rootReducer;



    