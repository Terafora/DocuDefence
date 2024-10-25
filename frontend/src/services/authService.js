import {jwtDecode} from 'jwt-decode';

export function setToken(token) {
    console.log("Setting token in localStorage:", token);
    const bearerToken = `Bearer ${token}`; // Adding Bearer prefix
    localStorage.setItem('jwtToken', bearerToken); // Store with Bearer prefix
}

export function getToken() {
    const token = localStorage.getItem('jwtToken');
    console.log("Retrieved token from localStorage:", token);
    return token; // No change, retrieves the entire Bearer <token> string
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
        const decoded = jwtDecode(token.split(" ")[1]); // Only decode the actual token, without "Bearer"
        return decoded.email;
    } catch (error) {
        console.error("Invalid token:", error.message);
        clearToken();  // Clear the token if it's invalid
        return null;
    }
}


// Function to handle user login
export async function loginUser(credentials) {
    const response = await fetch(`http://localhost:8000/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(credentials),
    });
    if (!response.ok) throw new Error('Login failed');
    
    const data = await response.json();
    console.log("Received token:", data.token); // Log the received token
    setToken(data.token); // Store JWT token in localStorage
    return data;
}

// Log out the user and clear token
export function logoutUser() {
    clearToken();
}
