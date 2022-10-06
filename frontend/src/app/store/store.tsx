import React, {createContext, Dispatch, useReducer} from "react";
import rootReducer, {StateContext} from './reducer'
import { getTeams } from '@app/utils/APIService';

interface IContextProps {
    state: StateContext;
    dispatch: ({type}:{type:string, data: any}) => void;
  }

const initialState = {
    APIData: [],
    E2E_KNOWN_ISSUES: [],
    error: 'error',
    alerts: [],
    version: '',
    repositories: [],
    Allrepositories: [],
    Team: "",
    TeamsAvailable: []
};

export const Context = React.createContext({} as IContextProps);

const Store = ({children}) => {
    React.useEffect(() => {
        getTeams().then(data => {
            if( data.data.length > 0){ 
                dispatch({ type: "SET_TEAM", data:  data.data[0].team_name });
                dispatch({ type: "SET_TEAMS_AVAILABLE", data:  data.data });
            }
          })
    }, []);

    const [state, dispatch] = useReducer(Reducer, initialState);
    const value = { state, dispatch };
    return (
        <Context.Provider value={value}>
            {children}
        </Context.Provider>
    )
};

export default Store;