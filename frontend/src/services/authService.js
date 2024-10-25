export function setToken(token) {
    localStorage.setItem('jwtToken', token);
  }
  
  export function getToken() {
    return localStorage.getItem('jwtToken');
  }
  
  export function clearToken() {
    localStorage.removeItem('jwtToken');
  }
  
  export function isLoggedIn() {
    return !!getToken();
  }
  