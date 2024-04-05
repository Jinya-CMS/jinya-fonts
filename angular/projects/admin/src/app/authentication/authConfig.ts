import { AuthConfig } from 'angular-oauth2-oidc';
import { environment } from '../../environments/environment';

export const authConfig: AuthConfig = {
  scope: 'openid profile offline_access',
  responseType: 'code',
  oidc: true,
  clientId: jinyaConfig.openIdClientId,
  issuer: jinyaConfig.openIdUrl,
  redirectUri: jinyaConfig.openIdCallbackUrl,
  postLogoutRedirectUri: jinyaConfig.openIdLogoutRedirectUrl,
  requireHttps: environment.production
};
