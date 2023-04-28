import * as React from 'react';
import { Route, RouteComponentProps, Switch } from 'react-router-dom';
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
import { Login } from './Login/Login';

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
  routes?: undefined;
}

export interface IAppRouteGroup {
  label: string;
  routes: IAppRoute[];
}

export type AppRouteConfig = IAppRoute | IAppRouteGroup;

const routes: AppRouteConfig[] = [
  {
    label: 'Login',
    routes: [
      {
        component: Login,
        exact: true,
        label: 'Login',
        path: '/login',
        title: 'Login | Quality Studio',
      },
    ],
  },
  {
    label: 'Home',
    routes: [
      {
        component: Overview,
        exact: true,
        label: 'Overview',
        path: '/home/overview',
        title: 'Overview | Quality Studio',
      },
      {
        component: Teams,
        exact: true,
        isAsync: true,
        label: 'Teams',
        path: '/home/teams',
        title: 'Teams | Quality Studio',
      },
      {
        component: Config,
        exact: true,
        isAsync: true,
        label: 'Config',
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
        label: 'Github',
        path: '/home/github',
        title: 'Github | Quality Studio',
      },
      {
        component: Jira,
        exact: true,
        isAsync: true,
        label: 'Jira',
        path: '/home/jira',
        title: 'Jira  | Quality Studio',
      },
      {
        component: Reports,
        exact: true,
        isAsync: true,
        label: 'Openshift CI',
        path: '/reports/test',
        title: 'Openshift CI | Quality Studio',
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

const RouteWithTitleUpdates = ({ component: Component, isAsync = false, title, ...rest }: IAppRoute) => {
  useA11yRouteChange(isAsync);
  useDocumentTitle(title);

  function routeWithTitle(routeProps: RouteComponentProps) {
    return <Component {...rest} {...routeProps} />;
  }
  return <Route render={routeWithTitle} {...rest} />;
};

const flattenedRoutes: IAppRoute[] = routes.reduce(
  (flattened, route) => [...flattened, ...(route.routes ? route.routes : [route])],
  [] as IAppRoute[]
);

const AppRoutes = (): React.ReactElement => {
  const { store } = React.useContext(ReactReduxContext);
  const state = store.getState();

  const [TeamsNotSet, setTeamsNotSet] = React.useState(false);

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
        {flattenedRoutes.map(({ path, exact, component, title, isAsync }, idx) => (
          <RouteWithTitleUpdates
            path={path}
            exact={exact}
            component={component}
            key={idx}
            title={title}
            isAsync={isAsync}
          />
        ))}
      </Switch>
    </LastLocationProvider>
  );
};

export { AppRoutes, routes };
