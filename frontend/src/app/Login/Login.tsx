import * as React from 'react';
import { Button } from '@patternfly/react-core';
import * as oauth2 from 'oauth4webapi'
import {
  LoginMainFooterBandItem,
  LoginPage,
  ListVariant
} from '@patternfly/react-core';
import { Spinner } from '@patternfly/react-core';
import { useHistory } from 'react-router-dom';
import { ReactReduxContext, useSelector } from 'react-redux';
import { initOauthFlow, completeOauthFlow, OauthData, refreshTokenFlow } from '@app/utils/oauth'

let authorizationUrl: URL
async function callLogin(){
  document.location.href = authorizationUrl.toString()
}

// eslint-disable-next-line prefer-const
let Login = () => {
  const { store } = React.useContext(ReactReduxContext);
  const redux_state = store.getState();
  const redux_dispatch = store.dispatch;
  const history = useHistory();
  const [isCallback, setIsCallback] = React.useState(false);
  const [callbackError, setCallbackError] = React.useState<string>("");

  const goBackToLogin = () => {
    setCallbackError("")
    setIsCallback(false)
    history.push("/login");
  }

  const signUpForAccountMessage = (
    <LoginMainFooterBandItem>
      Need an account? Ask to be let in.
    </LoginMainFooterBandItem>
  );

  const images = {
    lg: '/assets/images/pfbg_1200.jpg',
    sm: '/assets/images/pfbg_768.jpg',
    sm2x: '/assets/images/pfbg_768@2x.jpg',
    xs: '/assets/images/pfbg_576.jpg',
    xs2x: '/assets/images/pfbg_576@2x.jpg'
  };

  React.useEffect(() => {
    (async () => {
      const authURL = await initOauthFlow()
      authorizationUrl = authURL

      if(document.location.href.includes('code=')){
        setIsCallback(true)

        const data:OauthData = await completeOauthFlow()

        redux_dispatch({ type: "SET_REFRESH_TOKEN", data: data.RT });
        redux_dispatch({ type: "SET_ACCESS_TOKEN", data: data.AT });
        redux_dispatch({ type: "SET_ID_TOKEN", data: data.IDT });
        redux_dispatch({ type: "SET_AT_EXPIRATION", data: data.AT_EXPIRATION });
        redux_dispatch({ type: "SET_USERNAME", data: data.USERNAME });

        setTimeout(function() {
          setCallbackError("")
          const API_URL = process.env.REACT_APP_API_SERVER_URL || 'http://localhost:9898'
          document.location.href = "/home/overview"
        }, 3000);
      }

      if(document.location.href.includes('session_expired=true')){
        setCallbackError("Your session has expired. Please login again.")
      }

    })();
  }, []);

  return (
    <React.Fragment>
      <div>
        <LoginPage
            style={{textAlign: "center"}}
            footerListVariants={ListVariant.inline}
            backgroundImgSrc={images}
            loginTitle="Log in to your account"
            loginSubtitle="Choose your provider to login"
            signUpForAccountMessage={signUpForAccountMessage}
          >
            { !isCallback && <Button isBlock onClick={callLogin}>Continue to Login...</Button>}
            { isCallback  &&
              <div>
                <div> <Spinner isSVG aria-label="Contents of the basic example" /> </div>
                <div> You will be redirected in a few seconds.. </div>
              </div>
            }
            { callbackError != ""  &&
              <div>
                <div style={{color: "red", margin: "5px"}}> {callbackError} </div>
                <Button onClick={goBackToLogin}>Go Back</Button>
              </div>
            }
          </LoginPage>
      </div>
    </React.Fragment>
  )
};

export { Login };