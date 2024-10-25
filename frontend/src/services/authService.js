import { jwtDecode } from 'jwt-decode';

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

// New function to get the logged-in user's email with error handling
export function getUserEmail() {
    const token = getToken();
    if (!token) {
        console.error("No token found in localStorage.");
        return null;
    }

    try {
        const decoded = jwtDecode(token); 
        return decoded.email;
    } catch (error) {
        console.error("Invalid token:", error.message);
        clearToken();  // Clear the token if it's invalid
        return null;
    }
}
