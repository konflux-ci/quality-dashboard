import { ITeam } from '@app/Teams/TeamsSelect';
import { initial, reduceRight } from 'lodash';
import { combineReducers } from 'redux';
//import { initialState } from './store';
import { Reducer } from 'react';
import { configureStore } from '@reduxjs/toolkit';


export interface StateContext {
    APIData: [];
    E2E_KNOWN_ISSUES: [];
    error: string;
    alerts: [];
    version: string;
    repositories: [];
    workflows: [];
    Allrepositories: [];
    Team: string;
    TeamsAvailable: ITeam[]
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

        default: return action;
    };

};

const jirasReducer = (state, action) => {
    switch (action.type) {
        case 'SET_JIRAS':
            return {
                ...state,
                E2E_KNOWN_ISSUES: action.data
            };
        default: return action;
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
        default: return action;
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
        default: return action;
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
        default: return action;
    }
};

export const rootReducer = combineReducers({
    general : generalReducer,
    jira: jirasReducer,
    repo: repositoriesReducer,
    alert: alertsReducer,
    team: teamsReducer
});

export default rootReducer;



    