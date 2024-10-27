import { getToken, setToken, clearToken } from './authService';

const BASE_URL = 'http://localhost:8000';

async function fetchWithAuth(url, options = {}) {
    const token = getToken();
    const headers = {
        ...options.headers,
        'Content-Type': 'application/json',
        ...(token && { Authorization: token }),
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
    setToken(data.token); // Store JWT token
    return data;
}

export function logoutUser() {
    clearToken();
}
