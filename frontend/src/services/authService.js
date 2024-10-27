import {jwtDecode} from 'jwt-decode';

const BASE_URL = 'http://localhost:8000';

export function setToken(token) {
    if (!token) {
        console.error("Attempted to set an undefined token");
        return;
    }
    const bearerToken = `Bearer ${token}`;
    console.log("Setting token in localStorage:", bearerToken); // Added detailed logging
    localStorage.setItem('jwtToken', bearerToken);
}

export function getToken() {
    const token = localStorage.getItem('jwtToken');
    console.log("Retrieved token from localStorage:", token); // Check if we're getting the token correctly
    return token;
}

export function clearToken() {
    console.log("Clearing token from localStorage");
    localStorage.removeItem('jwtToken');
}

export function isLoggedIn() {
    return !!getToken();
}

export function getUserEmail() {
    const token = getToken();
    if (!token) {
        console.error("No token found in localStorage.");
        return null;
    }

    try {
        const decoded = jwtDecode(token.split(" ")[1]); // Only decode the actual token, without "Bearer"
        console.log("Decoded token, extracted email:", decoded.email); // Logging the decoded email
        return decoded.email;
    } catch (error) {
        console.error("Invalid token:", error.message);
        clearToken();  // Clear the token if it's invalid
        return null;
    }
}

async function fetchWithAuth(url, options = {}) {
    const token = getToken();
    const headers = {
        ...options.headers,
        'Content-Type': 'application/json',
        ...(token && { Authorization: token }), // Use the full "Bearer <token>" string
    };
    return fetch(url, { ...options, headers });
}

export async function getUsers() {
    const response = await fetchWithAuth(`${BASE_URL}/users`);
    if (!response.ok) throw new Error('Failed to fetch users');
    return response.json();
}

export async function createUser(user) {
    const response = await fetchWithAuth(`${BASE_URL}/users`, {
        method: 'POST',
        body: JSON.stringify(user),
    });
    if (!response.ok) throw new Error('Failed to create user');
    return response.json();
}

export async function updateUser(userId, updatedUserData) {
    const response = await fetchWithAuth(`${BASE_URL}/users/${userId}`, {
        method: 'PUT',
        body: JSON.stringify(updatedUserData),
    });
    if (!response.ok) throw new Error('Failed to update user');
    return response.json();
}

export async function deleteUser(userId) {
    const response = await fetchWithAuth(`${BASE_URL}/users/${userId}`, {
        method: 'DELETE',
    });
    if (!response.ok) throw new Error('Failed to delete user');
    return response.json();
}

export async function loginUser(credentials) {
    const response = await fetch(`${BASE_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(credentials),
    });
    if (!response.ok) throw new Error('Login failed');
    
    const data = await response.json();
    console.log("Login response data:", data); // Log the received response data
    if (data.token) {
        setToken(data.token); // Store JWT token if present
    } else {
        console.error("Login response did not include a token.");
    }
    return data;
}

export function logoutUser() {
    clearToken();
}
