import React, { useState } from "react";
import rootReducer, { StateContext } from './reducer'
import { getTeams } from '@app/utils/APIService';
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import { initialState } from '@app/store/initState';
import { loadStateContext, stateContextExists } from '@app/utils/utils'
import axios, { AxiosResponse } from 'axios';
import { Route, RouteComponentProps, Switch, Redirect, useLocation, useHistory } from 'react-router-dom';

interface IContextProps {
    state: StateContext;
    dispatch: ({ type }: { type: string, data: any }) => void;
}

export const Context = React.createContext({} as IContextProps);

const loadTeamSelection = (data) => {
    if (stateContextExists('TEAM')) {
        const teamPersisted = loadStateContext('TEAM')
        if (data.map((team) => team.team_name).includes(teamPersisted)) {
            return teamPersisted;
        }
    }
    return null;
}

const Store = ({ children }) => {
    const store = configureStore({ reducer: rootReducer, preloadedState: initialState });
    const state = store.getState()
    const dispatch = store.dispatch
    const history = useHistory();

    // Request interceptor for API calls
    const interceptor = axios.interceptors.request.use(
        async config => {
            config.headers = config.headers ?? {};
            config.headers.Authorization = state.auth.AT;
            return config;
        },
        error => {
            Promise.reject(error)
    });

    axios.interceptors.response.use(
        response => response,
        error => {
          if (error.response.status === 500) {
            window.location.href = '/login?session_expired=true';
          }
    });

    React.useEffect(() => {
        if(state.auth.IDT){
            getTeams().then(data => {
                if (data.data.length > 0) {
                    data.data.sort((a, b) => (a.team_name < b.team_name ? -1 : 1));
                    const loadedTeam = loadTeamSelection(data.data);
                    if (loadedTeam == null) { dispatch({ type: "SET_TEAM", data: data.data[0].team_name }) }
                    else { dispatch({ type: "SET_TEAM", data: loadedTeam }) }

                    const params = new URLSearchParams(window.location.search)
                    const team = params.get('team')

                    // team query param exists in the teams available
                    if (team != null && data.data.find(t => t.team_name === team)) {
                        dispatch({ type: "SET_TEAM", data: team })
                    }

                    dispatch({ type: "SET_TEAMS_AVAILABLE", data: data.data })
                }
            }
            )
        }

    }, []);

    return (
        <Provider store={store}>
            {children}
        </Provider>
    )
};

export default Store;
