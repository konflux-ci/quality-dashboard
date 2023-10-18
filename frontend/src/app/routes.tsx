import * as React from 'react';
import { Route, RouteComponentProps, Switch, Redirect, useLocation, useHistory } from 'react-router-dom';
import { accessibleRouteChangeHandler } from '@app/utils/utils';
import { Overview } from '@app/Overview/Overview';
import { Reports } from '@app/Reports/Reports';
import { useDocumentTitle } from '@app/utils/useDocumentTitle';
import { LastLocationProvider, useLastLocation } from 'react-router-last-location';
import { Teams } from '@app/Teams/Teams';
import { ReactReduxContext } from 'react-redux';
import { Jira } from './Jira/Jira';
import { GitHub } from './Github/Github';
import { Config } from './Config/Config';
import { initOauthFlow, completeOauthFlow, OauthData, refreshTokenFlow } from '@app/utils/oauth'
import { CiFailures } from './CiFailures/CiFailures';
import { BugSLIs } from './BugSLIs/MainPage';

let routeFocusTimer: number;
export interface IAppRoute {
  label?: string; // Excluding the label will exclude the route from the nav sidebar in AppLayout
  /* eslint-disable @typescript-eslint/no-explicit-any */
  component: React.ComponentType<RouteComponentProps<any>> | React.ComponentType<any>;
  /* eslint-enable @typescript-eslint/no-explicit-any */
  exact?: boolean;
  path: string;
  title: string;
  isAsync?: boolean;
  isProtected?: boolean;
  isAuthenticated?: boolean;
  routes?: undefined;
}

export interface IAppRouteGroup {
  label: string;
  routes: IAppRoute[];
}

export type AppRouteConfig = IAppRoute | IAppRouteGroup;

const routes: AppRouteConfig[] = [
  {
    label: 'Home',
    routes: [
      {
        component: Overview,
        exact: true,
        label: 'Overview',
        isProtected: true,
        path: '/home/overview',
        title: 'Overview | Quality Studio',
      },
      {
        component: Teams,
        exact: true,
        isAsync: true,
        label: 'Teams',
        path: '/home/teams',
        isProtected: true,
        title: 'Teams | Quality Studio',
      },
      {
        component: Config,
        exact: true,
        isAsync: true,
        label: 'Config',
        isProtected: true,
        path: '/home/config',
        title: 'Config | Quality Studio',
      },
    ],
  },
  {
    label: 'Plugins',
    routes: [
      {
        component: GitHub,
        exact: true,
        isAsync: true,
        isProtected: true,
        label: 'Github',
        path: '/home/github',
        title: 'Github | Quality Studio',
      },
      {
        component: Jira,
        exact: true,
        isAsync: true,
        isProtected: true,
        label: 'Jira',
        path: '/home/jira',
        title: 'Jira  | Quality Studio',
      },
      {
        component: BugSLIs,
        exact: true,
        isAsync: true,
        isProtected: true,
        label: 'RHTAP Bug SLIs',
        path: '/home/bug-slis',
        title: 'RHTAP Bug SLIs | Quality Studio',
      },
      {
        component: Reports,
        exact: true,
        isAsync: true,
        isProtected: true,
        label: 'Openshift CI',
        path: '/reports/test',
        title: 'Openshift CI | Quality Studio',
      },
      {
        component: CiFailures,
        exact: true,
        isAsync: true,
        isProtected: true,
        label: 'RHTAPBUGS Impact on CI',
        path: '/home/rhtapbugs-impact',
        title: 'RHTAPBUGS Impact on CI | Quality Studio',
      },
    ],
  },
];

// a custom hook for sending focus to the primary content container
// after a view has loaded so that subsequent press of tab key
// sends focus directly to relevant content
const useA11yRouteChange = (isAsync: boolean) => {
  const lastNavigation = useLastLocation();
  React.useEffect(() => {
    if (!isAsync && lastNavigation !== null) {
      routeFocusTimer = accessibleRouteChangeHandler();
    }
    return () => {
      window.clearTimeout(routeFocusTimer);
    };
  }, [isAsync, lastNavigation]);
};

const RouteWithTitleUpdates = ({ component: Component, isAsync = false, title, isProtected, isAuthenticated, ...rest }: IAppRoute) => {
  useA11yRouteChange(isAsync);
  useDocumentTitle(title);

  function routeWithTitle(routeProps: RouteComponentProps) {
    return <Component {...rest} {...routeProps} />;
  }

  return <Route {...rest} render={(props) => (
    isProtected == true && isAuthenticated == true ? routeWithTitle(props) : <Redirect to={{ pathname: '/login', state: { from: props.location } }} />
  )} />
};

const flattenedRoutes: IAppRoute[] = routes.reduce(
  (flattened, route) => [...flattened, ...(route.routes ? route.routes : [route])],
  [] as IAppRoute[]
);

const AppRoutes = (): React.ReactElement => {
  const { store } = React.useContext(ReactReduxContext);
  const state = store.getState();
  const dispatch = store.dispatch

  const [TeamsNotSet, setTeamsNotSet] = React.useState(false);
  let isAuthenticated = false

  if (state.auth.IDT) {
    isAuthenticated = true
  }

  const location = useLocation()
  const history = useHistory();

  React.useEffect(() => {
    // when the route changes, check for token expired and refresh it
    if (state.auth.AT && state.auth.AT_expiration) {
      const now = new Date()
      if (new Date(state.auth.AT_expiration * 1000) < now) {
        (async () => {
          try {
            let data = await refreshTokenFlow(state.auth.RT)
            dispatch({ type: "SET_REFRESH_TOKEN", data: data.RT });
            dispatch({ type: "SET_ACCESS_TOKEN", data: data.AT });
            dispatch({ type: "SET_ID_TOKEN", data: data.IDT });
            dispatch({ type: "SET_AT_EXPIRATION", data: data.AT_EXPIRATION });
          } catch (err) {
            dispatch({ type: "SET_REFRESH_TOKEN", data: "" });
            dispatch({ type: "SET_ACCESS_TOKEN", data: "data.AT" });
            dispatch({ type: "SET_ID_TOKEN", data: "" });
            dispatch({ type: "SET_AT_EXPIRATION", data: "" });
            localStorage.clear();
            history.push("/login")
            window.location.reload();
          }
        })();
      }
    }
  }, [location])

  React.useEffect(() => {
    if (state.teams.Team == undefined || state.teams.Team == 'Select Team' || state.teams.Team == '') {
      setTeamsNotSet(true);
    } else {
      setTeamsNotSet(false);
    }
  }, [location.pathname, state.teams.Team]);

  return (
    <LastLocationProvider>
      <Switch>
        {flattenedRoutes.map(({ path, exact, component, title, isAsync, isProtected }, idx) => (
          <RouteWithTitleUpdates
            path={path}
            exact={exact}
            component={component}
            key={idx}
            title={title}
            isAsync={isAsync}
            isProtected={isProtected}
            isAuthenticated={isAuthenticated}
          />
        ))}
      </Switch>
    </LastLocationProvider>
  );
};

export { AppRoutes, routes };
