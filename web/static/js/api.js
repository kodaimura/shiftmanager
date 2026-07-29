const BASE_URL = '/api';
let accessToken = null;
let refreshPromise = null;


class HttpError extends Error {
  status;

  constructor(status, message) {
    super(message);
    this.status = status;
  }
}

class Api {
  #url;

  constructor(url) {
    this.#url = url;
  }

  apiFetch = async (endpoint, method, body, retry = true) => {
    if (endpoint.startsWith('/')) {
      endpoint = endpoint.slice(1);
    }
    try {
      if (this.#needsAccessToken(endpoint)) {
        await this.refreshAccessTokenIfNeeded();
      }

      let header = {
        method: method,
        headers: {
          'Content-Type': 'application/json',
        },
      };

      if (accessToken && this.#needsAccessToken(endpoint)) {
        header.headers['X-Access-Token'] = accessToken;
      }

      if (body) {
        header.body = JSON.stringify(body);
      }
      const response = await fetch(`${this.#url}/${endpoint}`, header);

      if (!response.ok) {
        const errorData = await response.json();
        if (response.status === 401 && retry && this.#needsAccessToken(endpoint)) {
          accessToken = null;
          await this.refreshAccessTokenIfNeeded();
          return this.apiFetch(endpoint, method, body, false);
        }
        throw new HttpError(response.status, errorData.error);
      }

      let data;
      try {
        data = await response.json();
      } catch (error) {
        if (response.status !== 204 && response.status !== 200) {
          throw new HttpError(response.status, 'Error parsing JSON');
        }
      }
      return data;
    } catch (error) {
      if (error instanceof HttpError) {
        this.handleHttpError(error);
      } else {
        console.error(error);
        throw error;
      }
    }
  };

  #needsAccessToken = (endpoint) => {
    return !['login', 'signup', 'refresh'].includes(endpoint);
  };

  refreshAccessTokenIfNeeded = async () => {
    if (accessToken) {
      return;
    }
    if (!refreshPromise) {
      refreshPromise = this.post('refresh', null)
        .then((data) => {
          accessToken = data.access_token;
        })
        .finally(() => {
          refreshPromise = null;
        });
    }
    return refreshPromise;
  };

  setAccessToken = (token) => {
    accessToken = token;
  };

  clearAccessToken = () => {
    accessToken = null;
  };

  get = async (endpoint) => {
    return this.apiFetch(endpoint, 'GET', null);
  };

  post = async (endpoint, body) => {
    return this.apiFetch(endpoint, 'POST', body);
  };

  put = async (endpoint, body) => {
    return this.apiFetch(endpoint, 'PUT', body);
  };

  delete = async (endpoint) => {
    return this.apiFetch(endpoint, 'DELETE', null);
  };

  handleHttpError = (error) => {
    console.error(error);
    throw error;
  };
}

const api = new Api(BASE_URL);
api.handleHttpError = (error) => {
  const status = error.status;
  if (status === 401) {
    if (window.location.pathname !== '/login') {
      window.location.replace('/login');
    }
  } else if (status === 403) {
    alert("権限がありません。");
  } else if (status === 500) {
    alert("予期せぬエラーが発生しました。");
  }
  throw error;
}

export { HttpError, Api, BASE_URL, api };
