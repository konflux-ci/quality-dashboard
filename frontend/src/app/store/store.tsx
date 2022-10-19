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


const Store = ({ children }) => {
    let store = configureStore({ reducer: rootReducer, preloadedState: initialState });
    //let initState = initialState;

    React.useEffect(() => {
        const state = store.getState()
        const dispatch = store.dispatch

        getTeams().then(data => {
            if (data.data.length > 0) {

                

                if (stateContextExists('TEAM')) {
                    const teamPersist = loadStateContext('TEAM')
                    /* 
                        team selected has a persisted state And 
                        team exists in teamsAvailable: persist this team 
                    */
                    if (data.data.map(( team ) => team.team_name).includes( teamPersist )){ dispatch({ type: "SET_TEAM", data: teamPersist }) }
                    /* 
                        team selected is not persisted Or
                        team is not valid within teamsAvailable: set to the first available team
                    */
                   
                    else { dispatch({ type: "SET_TEAM", data: data.data[0].team_name }) }
                }

                dispatch({ type: "SET_TEAMS_AVAILABLE", data: data.data })
            }


        }
        )
        store = configureStore({ reducer: rootReducer, preloadedState: store.getState() });
    }, []);

    

    return (
        <Provider store={store}>
            {children}
        </Provider>
    )

};

export default Store;
