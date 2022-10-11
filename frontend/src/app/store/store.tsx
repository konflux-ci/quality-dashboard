import React, {createContext, Dispatch, useContext, useReducer} from "react";
import rootReducer, {StateContext} from './reducer'
import { getTeams } from '@app/utils/APIService';
import { $CombinedState, configureStore } from '@reduxjs/toolkit';
import { Provider, useDispatch } from 'react-redux';
import { ReactReduxContext } from 'react-redux'
import { initialState } from './initState'



interface IContextProps {
    state: StateContext;
    dispatch: ({type}:{type:string, data: any}) => void;
  }


export const Context = React.createContext({} as IContextProps);


const Store = ({children}) => {
    
    React.useEffect(() => {
        
        getTeams().then(data => {
            if( data.data.length > 0){ 
                store.dispatch({ type: "SET_TEAM", data:  data.data[0].team_name });
                store.dispatch({ type: "SET_TEAMS_AVAILABLE", data:  data.data });
            }
          }
        )
        //const { store } = useContext(ReactReduxContext); 
        const state = store.getState(); 

    }, []);

    const store = configureStore({reducer :rootReducer, preloadedState : initialState});    
    return (
        <Provider store={store}>
            {children}
        </Provider>
    )

};
export default Store;
