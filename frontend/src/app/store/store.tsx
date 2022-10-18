import React, { createContext, Dispatch, useContext, useReducer } from "react";
import rootReducer, { StateContext } from './reducer'
import { getTeams } from '@app/utils/APIService';
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import { initialState } from '@app/store/initState';
import { loadStateItem, stateItemExists } from '@app/utils/utils'

interface IContextProps {
    state: StateContext;
    dispatch: ({ type }: { type: string, data: any }) => void;
}


export const Context = React.createContext({} as IContextProps);


const Store = ({ children }) => {
    
    let store = configureStore({ reducer: rootReducer, preloadedState: initialState });
    
    React.useEffect(() => {
        const state = store.getState()
        const dispatch = store.dispatch
        
        getTeams().then(data => {
            let selectedTeam = data.data[0].team_name;
            if (stateItemExists('TEAM')){
                selectedTeam = loadStateItem('TEAM');
            }

            if (data.data.length > 0) {
                dispatch({ type: "SET_TEAMS_AVAILABLE", data: data.data });
                dispatch({ type: "SET_TEAM", data: selectedTeam });
            }
            
        }
        )

    }, []);

    return (
        <Provider store={store}>
            {children}
        </Provider>
    )

};

export default Store;
