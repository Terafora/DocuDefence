// src/services/userService.js

import { getToken, setToken, clearToken } from './authService';

const BASE_URL = 'http://localhost:8000';

// Helper function to add token to headers
async function fetchWithAuth(url, options = {}) {
  const token = getToken();
  const headers = {
      ...options.headers,
      'Content-Type': 'application/json',
      ...(token && { Authorization: token }), // Use the full "Bearer <token>" string
  };

  return fetch(url, { ...options, headers });
}


// Fetch all users with token authentication
export async function getUsers() {
  const response = await fetchWithAuth(`${BASE_URL}/users`);
  if (!response.ok) throw new Error('Failed to fetch users');
  return response.json();
}

// Create a new user
export async function createUser(user) {
  const response = await fetchWithAuth(`${BASE_URL}/users`, {
    method: 'POST',
    body: JSON.stringify(user),
  });
  if (!response.ok) throw new Error('Failed to create user');
  return response.json();
}

// Update an existing user
export async function updateUser(userId, updatedUserData) {
  const response = await fetchWithAuth(`${BASE_URL}/users/${userId}`, {
    method: 'PUT',
    body: JSON.stringify(updatedUserData),
  });
  if (!response.ok) throw new Error('Failed to update user');
  return response.json();
}

// Delete a user
export async function deleteUser(userId) {
  const response = await fetchWithAuth(`${BASE_URL}/users/${userId}`, {
    method: 'DELETE',
  });
  if (!response.ok) throw new Error('Failed to delete user');
  return response.json();
}

// Login user and store JWT token in localStorage
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

// Log out the user and clear token
export function logoutUser() {
  clearToken();
}
