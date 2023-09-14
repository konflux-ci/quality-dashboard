import * as React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import { BrowserRouter as Router } from 'react-router-dom';
import { AppLayout } from '@app/AppLayout/AppLayout';
import { AppRoutes } from '@app/routes';
import Store from './store/store';
import '@app/app.css';
import { Login } from './Login/Login'
import { Route, Switch } from 'react-router-dom';

const App: React.FunctionComponent = () => (
  <Store>
    <Router>
      <Switch>
        <Route path="/login" component={Login}/>
        <AppLayout>
          <AppRoutes />
        </AppLayout>
      </Switch>
    </Router>
  </Store>
);

export default App;
