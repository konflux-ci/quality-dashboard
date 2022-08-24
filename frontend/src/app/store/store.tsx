import React, {createContext, Dispatch, useReducer} from "react";
import Reducer, {StateContext} from './reducer'

interface IContextProps {
    state: StateContext;
    dispatch: ({type}:{type:string, data: any}) => void;
  }

const initialState = {
    APIData: [],
    error: 'error',
    alerts: [],
    version: '',
    repositories: [],
    Allrepositories: [],
    Team: ""
};

export const Context = React.createContext({} as IContextProps);

const Store = ({children}) => {
    let default_team = localStorage.getItem("default_team");
    if(default_team) initialState.Team = default_team;

    const [state, dispatch] = useReducer(Reducer, initialState);
    const value = { state, dispatch };
    return (
        <Context.Provider value={value}>
            {children}
        </Context.Provider>
    )
};

export default Store;