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
      const issuer = new URL('http://127.0.0.1:5556/dex')
      const state = "Login to DEX server"
      const as = await oauth2
        .discoveryRequest(issuer)
        .then((response) => oauth2.processDiscoveryResponse(issuer, response))
  
      const client: oauth2.Client = {
        client_id: 'example-app',
        token_endpoint_auth_method: 'none',
      }
  
      const redirect_uri = 'http://localhost:9000/login'
  
      if (as.code_challenge_methods_supported?.includes('S256') !== true) {
        setCallbackError("An error occurred")
        throw new Error()
      }
  
      const d = new Date().toDateString()
      const code_verifier = btoa(d);
      const code_challenge = await oauth2.calculatePKCECodeChallenge(code_verifier)
      const code_challenge_method = 'S256'
    
      authorizationUrl = new URL(as.authorization_endpoint!)
      authorizationUrl.searchParams.set('client_id', client.client_id)
      authorizationUrl.searchParams.set('code_challenge', code_challenge)
      authorizationUrl.searchParams.set('code_challenge_method', code_challenge_method)
      authorizationUrl.searchParams.set('redirect_uri', redirect_uri)
      authorizationUrl.searchParams.set('response_type', 'code')
      authorizationUrl.searchParams.set('scope', 'openid email')
      authorizationUrl.searchParams.set('state', state)
      
      if(document.location.href.includes('code=')){
        setIsCallback(true)
        let sub: string
        let access_token: string
        let id_token: string
        let refresh_token: string
        let expiration: Number

        {
          const currentUrl: URL = new URL(document.location.href)
          const params = oauth2.validateAuthResponse(as, client, currentUrl, state)
          if (oauth2.isOAuth2Error(params)) {
            console.log('error', params)
            let err: string = params.error_description || ""
            setCallbackError(err)
            throw new Error()
          }
  
          const response = await oauth2.authorizationCodeGrantRequest(
            as,
            client,
            params,
            redirect_uri,
            code_verifier,
          )
  
          let challenges: oauth2.WWWAuthenticateChallenge[] | undefined
          if ((challenges = oauth2.parseWwwAuthenticateChallenges(response))) {
            throw new Error() // Handle www-authenticate challenges as needed
          }
  
          const result = await oauth2.processAuthorizationCodeOpenIDResponse(as, client, response)
          if (oauth2.isOAuth2Error(result)) {
            let err: string = result.error_description || ""
            setCallbackError(err)
            throw new Error() // Handle OAuth 2.0 response body error
          }
  
          const claims = oauth2.getValidatedIdTokenClaims(result)
          access_token = result.access_token
          id_token = result.id_token

          if(result.refresh_token){
            refresh_token = result.refresh_token
            redux_dispatch({ type: "SET_REFRESH_TOKEN", data: id_token });
          }
          redux_dispatch({ type: "SET_ACCESS_TOKEN", data: access_token });
          redux_dispatch({ type: "SET_ID_TOKEN", data: id_token });

          if(result.expires_in){
            expiration = 1683204020
            redux_dispatch({ type: "SET_AT_EXPIRATION", data: expiration });
          }

          sub = claims.sub
          // fetch userinfo response
          {
            const response = await oauth2.userInfoRequest(as, client, access_token)
            let challenges: oauth2.WWWAuthenticateChallenge[] | undefined
            if ((challenges = oauth2.parseWwwAuthenticateChallenges(response))) {
              throw new Error() // Handle www-authenticate challenges as needed
            }
            const result = await oauth2.processUserInfoResponse(as, client, sub, response)
            if(result.email){
              let email = result.email
              redux_dispatch({ type: "SET_USERNAME", data: email });
            }
          }
                    
          setTimeout(function() { 
            setCallbackError("")
            setIsCallback(false)
            history.push("/");
          }, 3000);
          
        }
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
            { !isCallback && <Button isBlock onClick={callLogin}>Login with Github</Button>}
            { isCallback  && 
              <div>
                <div> <Spinner isSVG aria-label="Contents of the basic example" /> </div>
                <div> You will be redirected in a few seconds.. </div>
              </div>
            }
            { callbackError != ""  && 
              <div>
                <div style={{color: "red"}}> {callbackError} </div>
                <Button onClick={goBackToLogin}>Go Back</Button>
              </div>
            }
          </LoginPage>
      </div>
    </React.Fragment>
  )
};

export { Login };