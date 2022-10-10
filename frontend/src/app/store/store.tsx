import React, {createContext, Dispatch, useReducer} from "react";
import rootReducer, {StateContext} from './reducer'
import { getTeams } from '@app/utils/APIService';
import { configureStore } from '@reduxjs/toolkit';
import { Provider, useDispatch } from 'react-redux';

interface IContextProps {
    state: StateContext;
    dispatch: ({type}:{type:string, data: any}) => void;
  }

const initialState = { general : {
        APIData: [],
        E2E_KNOWN_ISSUES: [],
        error: 'error',
        alerts: [],
        version: '',
        repositories: [],
        Allrepositories: [],
        Team: "",
        TeamsAvailable: []
    }
};

export const Context = React.createContext({} as IContextProps);

const Store = ({children}) => {

    const dispatch = useDispatch();
    const store = configureStore({reducer:rootReducer, preloadedState : initialState});
    
    React.useEffect(() => {
        
        getTeams().then(data => {
            if( data.data.length > 0){ 
                dispatch({ type: "SET_TEAM", data:  data.data[0].team_name });
                dispatch({ type: "SET_TEAMS_AVAILABLE", data:  data.data });
            }
          }
        )
        
        const state = store.getState();
        const value = { state , dispatch };
        
    }, []);
    
    return (
        <Provider store={store}>
            {children}
        </Provider>
    )

};

export default Store;