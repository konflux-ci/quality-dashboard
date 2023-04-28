import * as React from 'react';
import * as oauth from 'oauth4webapi'
import { Button } from '@patternfly/react-core';

async function callLogin(){
  const API_URL = process.env.REACT_APP_API_SERVER_URL || 'http://localhost:9898';
  fetch(API_URL + "/api/quality/login")
}

// eslint-disable-next-line prefer-const
let Login = () => {
  return (
    <React.Fragment>
      <div style={{padding: '150px'}}>
        <h1>Login</h1>
        <Button variant="primary" onClick={callLogin}>Login</Button>
      </div>
    </React.Fragment>
  )
};

export { Login };