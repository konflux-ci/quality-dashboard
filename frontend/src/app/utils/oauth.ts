import * as React from 'react';
import * as oauth2 from 'oauth4webapi'
import { useHistory } from 'react-router-dom';
import { ReactReduxContext, useSelector } from 'react-redux';

export interface OauthData {
  AT: string;
  RT: string;
  IDT: string;
  AT_EXPIRATION: Number;
  USERNAME: string;
}

const issuer = new URL('http://127.0.0.1:5556/dex')
const state = "Login to DEX server"

const client: oauth2.Client = {
  client_id: 'example-app',
  token_endpoint_auth_method: 'none',
}

const redirect_uri = 'http://localhost:9000/login'

export async function initOauthFlow():Promise<URL> {
  let authorizationUrl: URL
  let as = await oauth2.discoveryRequest(issuer).then((response) => oauth2.processDiscoveryResponse(issuer, response))

  if (as.code_challenge_methods_supported?.includes('S256') !== true) {
    throw new Error("An error occurred: S256 not supported")
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
  authorizationUrl.searchParams.set('scope', 'openid email offline_access', )
  authorizationUrl.searchParams.set('state', state)

  return authorizationUrl
}

export async function completeOauthFlow():Promise<OauthData> {
  const d = new Date().toDateString()
  const code_verifier = btoa(d);
  let as = await oauth2.discoveryRequest(issuer).then((response) => oauth2.processDiscoveryResponse(issuer, response))

  let sub: string
  let access_token: string
  let id_token: string
  let refresh_token: string
  let expiration: Number = -1
  let username: string = ""

  const currentUrl: URL = new URL(document.location.href)
  const params = oauth2.validateAuthResponse(as, client, currentUrl, state)

  if (oauth2.isOAuth2Error(params)) {
    console.log('error', params)
    let err: string = params.error_description || ""
    throw new Error(err)
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
    throw new Error("An error occurred processing www-authenticate challenges") // Handle www-authenticate challenges as needed
  }

  const result = await oauth2.processAuthorizationCodeOpenIDResponse(as, client, response)
  if (oauth2.isOAuth2Error(result)) {
    let err: string = result.error_description || ""
    throw new Error(err)
  }

  const claims = oauth2.getValidatedIdTokenClaims(result)
  access_token = result.access_token
  id_token = result.id_token
  expiration = claims.exp
  sub = claims.sub
  refresh_token = result.refresh_token ? result.refresh_token : ""

  const resp = await oauth2.userInfoRequest(as, client, access_token)
  let ch: oauth2.WWWAuthenticateChallenge[] | undefined
  if ((ch = oauth2.parseWwwAuthenticateChallenges(resp))) {
    throw new Error() // Handle www-authenticate challenges as needed
  }
  const res = await oauth2.processUserInfoResponse(as, client, sub, resp)
  if(result.email){
    username = result.email as string
  }

  let Data:OauthData = {
    AT: access_token,
    IDT: id_token,
    AT_EXPIRATION: expiration,
    USERNAME: username,
    RT: refresh_token
  }

  return Data
}

export async function refreshTokenFlow(rt:string):Promise<OauthData> {
  const d = new Date().toDateString()
  const code_verifier = btoa(d);
  let as = await oauth2.discoveryRequest(issuer).then((response) => oauth2.processDiscoveryResponse(issuer, response))

  let sub: string
  let access_token: string
  let id_token: string
  let refresh_token: string
  let expiration: Number = -1
  let username: string = ""

  const response = await oauth2.refreshTokenGrantRequest(
    as,
    client,
    rt
  )

  let challenges: oauth2.WWWAuthenticateChallenge[] | undefined
  if ((challenges = oauth2.parseWwwAuthenticateChallenges(response))) {
    throw new Error("An error occurred processing www-authenticate challenges") // Handle www-authenticate challenges as needed
  }

  const result = await oauth2.processRefreshTokenResponse(as, client, response)
  if (oauth2.isOAuth2Error(result)) {
    let err: string = result.error_description || ""
    throw new Error(err)
  }

  // better handle null tokens and values
  const claims = oauth2.getValidatedIdTokenClaims(result)
  access_token = result.access_token
  id_token = result.id_token ? result.id_token : ""
  refresh_token = result.refresh_token ? result.refresh_token : ""
  expiration = claims?.exp ? claims.exp : -1
  username = ""
  
  let Data:OauthData = {
    AT: access_token,
    IDT: id_token,
    AT_EXPIRATION: expiration,
    USERNAME: username,
    RT: refresh_token
  }
  return Data
}

