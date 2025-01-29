import { Alpine } from './alpine.js';
import { getAccessToken } from './jinya-alpine-tools.js';

export class HttpError extends Error {
  constructor(status, error) {
    super();
    this.status = status;
    this.message = error.message;
    this.type = error.type;
    this.plainError = error;
  }
}

export class BadRequestError extends HttpError {
  constructor(error) {
    super(400, error);
  }
}

export class ConflictError extends HttpError {
  constructor(error) {
    super(409, error);
  }
}

export class NotAllowedError extends HttpError {
  constructor(error) {
    super(403, error);
  }
}

export class NotFoundError extends HttpError {
  constructor(error) {
    super(404, error);
  }
}

export class UnauthorizedError extends HttpError {
  constructor(error) {
    super(401, error);
  }
}

export async function send(
  verb,
  url,
  data = undefined,
  contentType = 'application/json',
  additionalHeaders = {},
  plain = false,
) {
  const headers = { 'Content-Type': contentType, 'Authorization': `Bearer ${getAccessToken()}`, ...additionalHeaders };

  const request = {
    headers,
    credentials: 'same-origin',
    method: verb,
  };

  if (data) {
    if (data instanceof Blob) {
      request.body = data;
    } else {
      request.body = JSON.stringify(data);
    }
  }

  const response = await fetch(url, request);
  if (response.ok) {
    if (response.status !== 204) {
      if (plain) {
        return await response.text();
      }

      return await response.json();
    }

    return null;
  }

  const httpError = (await response.json()).error;
  switch (response.status) {
    case 400:
      throw new BadRequestError(httpError);
    case 401:
      if (httpError.type === 'invalid-api-key') {
        if (window.document) {
          Alpine.store('authentication').logout();
        }

        return null;
      }

      throw new UnauthorizedError(httpError);
    case 403:
      throw new NotAllowedError(httpError);
    case 404:
      throw new NotFoundError(httpError);
    case 409:
      throw new ConflictError(httpError);
    default:
      throw new HttpError(response.status, httpError);
  }
}

export function get(url) {
  return send('get', url);
}

export function getPlain(url) {
  return send('get', url, null, null, null, true);
}

export function head(url) {
  return send('head', url);
}

export function put(url, data) {
  if (data) {
    return send('put', url, data);
  }

  return send('put', url, data, '');
}

export function post(url, data) {
  if (data) {
    return send('post', url, data);
  }

  return send('post', url, data, '');
}

export function httpDelete(url) {
  return send('delete', url);
}

export function upload(url, file) {
  return send('put', url, file, file.type);
}

export function uploadPost(url, file) {
  return send('post', url, file, file.type);
}
