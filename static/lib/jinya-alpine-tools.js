import Alpine from './alpine.js';
import PineconeRouter from './pinecone-router.js';
import * as client from './openid-client/index.js';

let authenticationConfiguration = {
  openIdUrl: '',
  openIdClientId: '',
  openIdCallbackUrl: ''
};
let scriptBasePath = '/static/js/';
let languages = {};

export function setRedirect(redirect) {
  sessionStorage.setItem('/jinya/login/redirect', redirect);
}

export function getRedirect() {
  return sessionStorage.getItem('/jinya/login/redirect');
}

export function deleteRedirect() {
  sessionStorage.removeItem('/jinya/login/redirect');
}

export function hasAccessToken() {
  return !!localStorage.getItem('/jinya/api/access-token');
}

export function getAccessToken() {
  return localStorage.getItem('/jinya/api/access-token');
}

export function setAccessToken(code) {
  localStorage.setItem('/jinya/api/access-token', code);
}

export function deleteAccessToken() {
  localStorage.removeItem('/jinya/api/access-token');
}

function setCodeVerifier(code) {
  localStorage.setItem('/jinya/login/code-verifier', code);
}

function getCodeVerifier() {
  return localStorage.getItem('/jinya/login/code-verifier');
}

function deleteCodeVerifier() {
  localStorage.removeItem('/jinya/login/code-verifier');
}

export async function needsLogin(context) {
  if (await checkLogin()) {
    return null;
  }

  const redirect = context.path.substring('/admin'.length);
  setRedirect(redirect);

  return context.redirect('/login');
}

export async function needsLogout(context) {
  if (await checkLogin()) {
    return context.redirect('/');
  }

  return null;
}

export async function performLogin(context) {
  const config = await client.discovery(
    new URL(window.jinyaConfig.openIdUrl),
    window.jinyaConfig.openIdClientId
  );

  const tokenResponse = await client.authorizationCodeGrant(
    config,
    new URL(location.href),
    {
      pkceCodeVerifier: getCodeVerifier()
    }
  );
  setAccessToken(tokenResponse.access_token);
  Alpine.store('authentication').login();
  context.redirect(getRedirect() ?? '/');
}

export async function checkLogin() {
  if (!hasAccessToken()) {
    return false;
  }

  const config = await client.discovery(
    new URL(authenticationConfiguration.openIdUrl),
    authenticationConfiguration.openIdClientId
  );

  try {
    const response = await fetch(config.serverMetadata().userinfo_endpoint, {
      method: 'GET',
      mode: 'cors',
      cache: 'no-cache',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAccessToken()}`
      }
    });

    return response.status === 200;
  } catch (error) {
    console.error(error);
    return false;
  }
}

export async function fetchScript({ route }) {
  const [, , area, page] = route.split('/');
  if (area) {
    await import(`${scriptBasePath}${area}/${page?.replaceAll(':', '') ?? 'index'}.js`);
    Alpine.store('navigation').navigate({
      area,
      page: page ?? 'index'
    });
  }
}

export function getLanguage() {
  if (navigator.language.startsWith('de')) {
    return 'de';
  }

  return 'en';
}

/**
 * Localizes the given key and returns the matching string
 * @param key {string}
 * @param values {Object}
 * @return string
 */
export function localize({ key, values = {} }) {
  let transformed = languages[getLanguage()][key];
  for (const valueKey of Object.keys(values)) {
    transformed = transformed.replaceAll(`{${valueKey}}`, values[valueKey]);
  }

  return transformed;
}

export async function openIdLogin() {
  const config = await client.discovery(
    new URL(authenticationConfiguration.openIdUrl),
    authenticationConfiguration.openIdClientId
  );
  const redirectUrl = authenticationConfiguration.openIdCallbackUrl;
  const codeVerifier = client.randomPKCECodeVerifier();
  const codeChallenge = await client.calculatePKCECodeChallenge(codeVerifier);
  const parameters = {
    redirect_uri: redirectUrl,
    scope: 'openid profile',
    code_challenge: codeChallenge,
    code_challenge_method: 'S256'
  };
  const redirectTo = client.buildAuthorizationUrl(config, parameters);
  setCodeVerifier(codeVerifier);
  window.location.href = redirectTo;
}

export function setupLocalization(Alpine, langs) {
  languages = langs;

  Alpine.directive('localize', (el, { value, expression, modifiers }, { evaluateLater, effect }) => {
    const getValues = expression ? evaluateLater(expression) : (load) => load();
    effect(() => {
      getValues((values) => {
        const localized = localize({
          key: value,
          values
        });

        if (modifiers.includes('html')) {
          el.innerHTML = localized;
        } else if (modifiers.includes('title')) {
          el.setAttribute('title', localized);
        } else {
          el.textContent = localized;
        }
      });
    });
  });
}

function setupAuthentication(openIdUrl, openIdClientId, openIdCallbackUrl) {
  authenticationConfiguration.openIdClientId = openIdClientId;
  authenticationConfiguration.openIdUrl = openIdUrl;
  authenticationConfiguration.openIdCallbackUrl = openIdCallbackUrl;
}

function setupRouting(baseScriptPath, routerBasePath = '') {
  scriptBasePath = baseScriptPath;

  document.addEventListener('alpine:init', () => {
    window.PineconeRouter.settings.basePath = routerBasePath;
    window.PineconeRouter.settings.templateTargetId = 'app';
  });
}

async function setupAlpine(alpine, defaultArea, defaultPage) {
  Alpine.directive('active-route', (el, { expression, modifiers }, { Alpine, effect }) => {
    effect(() => {
      const { page, area } = Alpine.store('navigation');
      if ((modifiers.includes('area') && area === expression) || (!modifiers.includes('area') && page === expression)) {
        el.classList.add('is--active');
      } else {
        el.classList.remove('is--active');
      }
    });
  });

  Alpine.store('loaded', false);
  Alpine.store('authentication', {
    needsLogin,
    needsLogout,
    performLogin,
    loggedIn: await checkLogin(),
    login() {
      this.loggedIn = true;
      history.replaceState(null, null, location.href.split('?')[0]);
    },
    logout() {
      deleteAccessToken();
      setRedirect(location.pathname.substring(0, 6));
      window.PineconeRouter.context.navigate('/login');
      this.loggedIn = false;
      this.roles = [];
    }
  });
  Alpine.store('navigation', {
    fetchScript,
    area: defaultArea,
    page: defaultPage,
    navigate({ area, page }) {
      this.area = area;
      this.page = page;
    }
  });
}

export async function setup({
                              defaultArea,
                              defaultPage,
                              baseScriptPath,
                              routerBasePath = '',
                              openIdUrl = undefined,
                              openIdClientId = undefined,
                              openIdCallbackUrl = undefined,
                              languages = [],
                              afterSetup = () => {
                              }
                            }) {
  window.Alpine = Alpine;

  Alpine.plugin(PineconeRouter);

  if (openIdUrl && openIdClientId && openIdCallbackUrl) {
    setupAuthentication(openIdUrl, openIdClientId, openIdCallbackUrl);
  }
  if (Object.keys(languages ?? {}).length > 0) {
    setupLocalization(Alpine, languages);
  }
  await setupAlpine(Alpine, defaultArea, defaultPage);

  setupRouting(baseScriptPath, routerBasePath);

  Alpine.start();

  await afterSetup();

  Alpine.store('loaded', true);
}
