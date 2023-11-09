import React, { useEffect } from 'react';
import { NavLink, useLocation, useHistory } from 'react-router-dom';
import {
  Nav,
  NavList,
  NavItem,
  NavExpandable,
  Page,
  PageHeader,
  PageSidebar,
  SkipToContent,
  Alert,
  AlertGroup,
  MastheadContent,
} from '@patternfly/react-core';
import { routes, IAppRoute, IAppRouteGroup, isPluginRouteAuthorized } from '@app/routes';
import logo from '@app/bgimages/Logo-RedHat-A-Reverse-RGB.svg';
import { BasicMasthead } from '@app/Teams/TeamsSelect';
import { checkDbConnection, getTeams, getVersion, listInstalledPlugins } from '@app/utils/APIService';
import { ReactReduxContext, useSelector } from 'react-redux';

interface IAppLayout {
  children: React.ReactNode;
}

const AppLayout: React.FunctionComponent<IAppLayout> = ({ children }) => {
  const locationHistory = useLocation();
  const history = useHistory();
  const [isNavOpen, setIsNavOpen] = React.useState(true);
  const [isMobileView, setIsMobileView] = React.useState(true);
  const [isNavOpenMobile, setIsNavOpenMobile] = React.useState(false);
  const { store } = React.useContext(ReactReduxContext);
  const state = store.getState();
  const currentTeam = useSelector((state: any) => state.teams.Team);
  const pluginList = useSelector((state: any) => state.teams.FlattenedPlugins);
  const dispatch = store.dispatch
  
  const onNavToggleMobile = () => {
    setIsNavOpenMobile(!isNavOpenMobile);
  };
  const onNavToggle = () => {
    setIsNavOpen(!isNavOpen);
  };
  const onPageResize = (props: { mobileView: boolean; windowSize: number }) => {
    setIsMobileView(props.mobileView);
  };

  const params = new URLSearchParams(window.location.search);
  const team = params.get('team');

  const listAllPlugins = () =>  {
    if(currentTeam == ''){
      return
    }
    listInstalledPlugins(currentTeam).then(res => {
      if(res.code == 200 && res.data){
        dispatch({ type: "SET_INSTALLED_PLUGINS", data: res.data }) 
      } else {
        throw("Error getting plugins list for team")
      }
    })
  };

  function LogoImg() {
    const history = useHistory();
    function handleClick() {
      history.push('/home/overview');
    }
    return <img onClick={handleClick} style={{ height: '32px' }} src={logo} alt="Red Hat logo" />;
  }

  const [areTeamsEmpty, setAreTeamsEmpty] = React.useState(false);
  const [serverUnavailable, setServerUnavailable] = React.useState(false);
  const [dbUnavailable, setDbUnavailable] = React.useState(false);
  const [alerts, setAlerts] = React.useState<React.ReactNode[]>([]);
  const [Navigation, setNavigation] = React.useState<JSX.Element>();

  const Header = (
    <PageHeader
      logo={<LogoImg />}
      headerTools={
        <React.Fragment>
          <BasicMasthead></BasicMasthead>
          <MastheadContent>
            <AlertGroup isToast>
              {serverUnavailable && <Alert
                variant="danger"
                timeout={5000}
                title="Quality Studio unable to connect to backend server"
                key={0}
              />}
              {dbUnavailable && <Alert
                variant="danger"
                timeout={5000}
                title="Quality Studio unable to connect to database"
                key={0}
              />}
            </AlertGroup>
          </MastheadContent>
        </React.Fragment>
      }
      showNavToggle
      isNavOpen={isNavOpen}
      onNavToggle={isMobileView ? onNavToggleMobile : onNavToggle}
    />
  );

  const location = useLocation();

  // Handles the route in other to keep the query parameters in nav item multiples clicks
  const handleRoute = (route) => {
    if (history.location.pathname == route && team != null) {
      return history.location.pathname + history.location.search;
    }
    return route;
  };

  const renderNavItem = (route: IAppRoute, index: number, allowedPlugins: string[]) => {
    const navRoute = (<NavItem key={`${route.label}-${index}`} id={`${route.label}-${index}`}>
      <NavLink exact={route.exact} to={handleRoute(route.path)} activeClassName="pf-m-current">
        {route.label}
      </NavLink>
    </NavItem>)

    if(route.isPlugin){
      if(route.label && isPluginRouteAuthorized(route.label, allowedPlugins)){
        return navRoute
      } else {
        return false
      }
    } else {
      return navRoute
    }
  };

  const renderNavGroup = (group: IAppRouteGroup, groupIndex: number) => {
    return (
      <NavExpandable
        key={`${group.label}-${groupIndex}`}
        id={`${group.label}-${groupIndex}`}
        title={group.label}
        isActive={group.routes.some((route) => route.path === location.pathname)}
      >
        {group.routes.map((route, idx) => route.label && renderNavItem(route, idx, state.teams.FlattenedPlugins))}
      </NavExpandable>
  )};


  const toRender = (label) => {
    if (label == "Plugins") {
      getVersion().then(res => {
        if (!(res.code == 200)) {
          setServerUnavailable(true)
        }
      })
      if (serverUnavailable) {
        return false
      }

      checkDbConnection().then(res => {
        if (!(res.code == 200)) {
          setDbUnavailable(true)
        }
      })
      if (dbUnavailable) {
        return false
      }

      getTeams().then(res => {
        if (res.data.length == 0) {
          setAreTeamsEmpty(true)
        }
      })
      if (areTeamsEmpty) {
        return false
      }
    }
    return true
  }

  React.useEffect(() => {
    const N = (
      <Nav id="nav-primary-simple" theme="dark">
        <NavList id="nav-list-simple">
          {routes.map(
            (route, idx) => route.label && toRender(route.label) && (!route.routes ? renderNavItem(route, idx, state.teams.FlattenedPlugins) : renderNavGroup(route, idx))
          )}
        </NavList>
      </Nav>
    );
    setNavigation(N)
  }, [currentTeam, pluginList]);

  useEffect(() => {
    if (locationHistory.pathname === '/') {
      history.push('/home/overview');
    }
    listAllPlugins();
  }, [currentTeam]);

  const Sidebar = <PageSidebar theme="dark" nav={Navigation} isNavOpen={isMobileView ? isNavOpenMobile : isNavOpen} />;

  const pageId = 'primary-app-container';

  const PageSkipToContent = (
    <SkipToContent
      onClick={(event) => {
        event.preventDefault();
        const primaryContentContainer = document.getElementById(pageId);
        primaryContentContainer && primaryContentContainer.focus();
      }}
      href={`#${pageId}`}
    >
      Skip to Content
    </SkipToContent>
  );

  return (
    <Page
      mainContainerId={pageId}
      header={Header}
      sidebar={Sidebar}
      onPageResize={onPageResize}
      skipToContent={PageSkipToContent}
    >
      {children}
    </Page>
  );
};

export { AppLayout };
