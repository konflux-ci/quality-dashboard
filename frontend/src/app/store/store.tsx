import React, { createContext, Dispatch, useContext, useReducer } from "react";
import rootReducer, { StateContext } from './reducer'
import { getTeams } from '@app/utils/APIService';
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import { initialState } from '@app/store/initState';
import { loadStateContext, stateContextExists } from '@app/utils/utils'

interface IContextProps {
    state: StateContext;
    dispatch: ({ type }: { type: string, data: any }) => void;
}


export const Context = React.createContext({} as IContextProps);

const loadTeamSelection = (data) => {
    if (stateContextExists('TEAM')){
        const teamPersisted = loadStateContext('TEAM')
        if (data.data.map(( team ) => team.team_name).includes( teamPersisted )){
            return teamPersisted;
        }
    }
    return null;
}

const Store = ({ children }) => {
    const store = configureStore({ reducer: rootReducer, preloadedState: initialState });

    React.useEffect(() => {
        const state = store.getState()
        const dispatch = store.dispatch

        getTeams().then(data => {
            if (data.data.length > 0) {
                const loadedTeam = loadTeamSelection(data);
                switch (loadedTeam){
                    case null: 
                        dispatch({ type: "SET_TEAM", data: data.data[0].team_name })
                    default:
                        dispatch({ type: "SET_TEAM", data: loadedTeam })
                }  
                dispatch({ type: "SET_TEAMS_AVAILABLE", data: data.data })       
            }
        }
    )}, []);

    

    return (
        <Provider store={store}>
            {children}
        </Provider>
    )

};

export default Store;
