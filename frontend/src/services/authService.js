import jwt_decode from "jwt-decode";

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

// New function to get the logged-in user's email
export function getUserEmail() {
    const token = getToken();
    if (!token) return null;
    const decoded = jwt_decode(token);
    return decoded.email;
}
