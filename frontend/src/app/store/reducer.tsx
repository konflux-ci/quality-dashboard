
export interface StateContext {
    APIData: [];
    E2E_KNOWN_ISSUES: [];
    error: string;
    alerts: [];
    version: string;
    repositories: []
    workflows: []
    Allrepositories: []
}

const Reducer = (state, action) => {
    switch (action.type) {
        case 'APIData':
            return {
                ...state,
                APIData: action.data
            };
        case 'SET_JIRAS':
            return {
                ...state,
                E2E_KNOWN_ISSUES: action.data
            }
        case 'SET_Version':
            return {
                ...state,
                version: action.data
            };
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
        default:
            return state;
    }
};

export default Reducer;
