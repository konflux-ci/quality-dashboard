import * as React from 'react';
import { Redirect, Route, RouteComponentProps, Switch } from 'react-router-dom';
import { accessibleRouteChangeHandler } from '@app/utils/utils';
import { Dashboard } from '@app/Dashboard/Dashboard';
import { Support } from '@app/Support/Support';
import { GeneralSettings } from '@app/Settings/General/GeneralSettings';
import { ProfileSettings } from '@app/Settings/Profile/ProfileSettings';
import { useDocumentTitle } from '@app/utils/useDocumentTitle';
import { LastLocationProvider, useLastLocation } from 'react-router-last-location';
import { JobsComponent } from './Jobs/Jobs';
import { Teams } from '@app/Teams/Teams';
import { Context } from '@app/store/store';

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
    label: 'Home',
    routes: [
      {
        component: Dashboard,
        exact: true,
        label: 'Overview',
        path: '/home/overview',
        title: 'Quality Studio | Overview',
      },
    ],
  },

  {
    component: JobsComponent,
    exact: true,
    isAsync: true,
    label: 'CI Jobs',
    path: '/ci/jobs',
    title: 'PatternFly Seed | Support Page',
  },
  {
    component: Support,
    exact: true,
    isAsync: true,
    label: 'Tests Reports',
    path: '/reports/test',
    title: 'Tests Reports | Page',
  },
  {
    label: 'Settings',
    routes: [
      {
        component: GeneralSettings,
        exact: true,
        label: 'General',
        path: '/settings/general',
        title: 'PatternFly Seed | General Settings',
      },
      {
        component: ProfileSettings,
        exact: true,
        label: 'Profile',
        path: '/settings/profile',
        title: 'PatternFly Seed | Profile Settings',
      },
      {
        component: Teams,
        exact: true,
        isAsync: true,
        label: 'Create New Team',
        path: '/settings/teams/onboarding',
        title: 'Teams | Onboarding',
      }
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
  return <Route render={routeWithTitle} {...rest}/>;
};

const flattenedRoutes: IAppRoute[] = routes.reduce(
  (flattened, route) => [...flattened, ...(route.routes ? route.routes : [route])],
  [] as IAppRoute[]
);

const AppRoutes = (): React.ReactElement => {
  const { state } = React.useContext(Context)
  const [TeamsNotSet, setTeamsNotSet] = React.useState(false)

  React.useEffect(() => {
    if(state.Team == undefined || state.Team == "Select Team" || state.Team == ""){
      setTeamsNotSet(true)
    } else {
      setTeamsNotSet(false)
    }
  }, [location.pathname, state.Team]);


  return (
  <LastLocationProvider>
    <Switch>
      {flattenedRoutes.map(({ path, exact, component, title, isAsync }, idx) => (
        <RouteWithTitleUpdates
          path={path}
          exact={exact}
          component={TeamsNotSet ? Teams : component}
          key={idx}
          title={title}
          isAsync={isAsync}
        />
      ))}
      <Redirect to='/home/overview'  />
    </Switch>
  </LastLocationProvider>
)};

export { AppRoutes, routes };
