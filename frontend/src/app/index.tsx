import * as React from 'react';
import '@patternfly/react-core/dist/styles/base.css';
import { BrowserRouter as Router } from 'react-router-dom';
import { AppLayout } from '@app/AppLayout/AppLayout';
import { AppRoutes } from '@app/routes';
import Store from './store/store';
import '@app/app.css';

const App: React.FunctionComponent = () => (
  <Store>
    <Router>
      <AppLayout>
        <AppRoutes />
      </AppLayout>
    </Router>
  </Store>
);

export default App;
