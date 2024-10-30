import { getToken } from './authService';

const BASE_URL = 'http://localhost:8000';

// Fetch user ID by email
export async function fetchUserIDByEmail(email) {
    const token = getToken();
    if (!token) {
        throw new Error("Authorization token missing.");
    }

    const response = await fetch(`${BASE_URL}/users/email?email=${encodeURIComponent(email)}`, {
        method: 'GET',
        headers: {
            Authorization: token,
        },
    });

    if (!response.ok) {
        console.error('Fetch user ID response:', await response.text());
        throw new Error('Failed to fetch user ID');
    }
    return response.json();
}

// Upload file and create a new document version
export async function uploadFile(userId, file) {
    const formData = new FormData();
    formData.append('contract', file);

    const token = getToken();
    const response = await fetch(`${BASE_URL}/users/${userId}/upload`, {
        method: 'POST',
        headers: {
            Authorization: token,
        },
        body: formData,
    });

    if (!response.ok) {
        console.error('Upload response:', await response.text());
        throw new Error('Failed to upload file');
    }
    return response.json();
}

// Fetch all versions of a user's files
export async function getUserFiles(userId) {
    const token = getToken();
    const response = await fetch(`${BASE_URL}/users/${userId}/files`, {
        headers: {
            Authorization: token,
        },
    });

    if (!response.ok) {
        console.error('Fetch files response:', await response.text());
        throw new Error('Failed to fetch user files');
    }
    return response.json();
}

// Download a specific version of a file by filename
export async function downloadFile(userId, filename) {
    const token = getToken();
    const encodedFilename = encodeURIComponent(filename);
    const response = await fetch(`${BASE_URL}/users/${userId}/files/${encodedFilename}/download`, {
        method: 'GET',
        headers: {
            Authorization: token,
        },
    });

    if (!response.ok) {
        console.error('Download response:', await response.text());
        throw new Error('Failed to download file');
    }

    // Create a blob URL for the downloaded file
    const blob = await response.blob();
    const downloadUrl = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = downloadUrl;
    a.download = filename;  // Set the filename for download
    document.body.appendChild(a);
    a.click();
    a.remove();
}

// Delete a specific version of a file by filename
export async function deleteFile(userId, filename) {
    const token = getToken();
    const encodedFilename = encodeURIComponent(filename);
    const response = await fetch(`${BASE_URL}/users/${userId}/files/${encodedFilename}/delete`, {
        method: 'DELETE',
        headers: {
            Authorization: token,
        },
    });

    if (!response.ok) {
        console.error('Delete response:', await response.text());
        throw new Error('Failed to delete file');
    }
    return response.json();
}

// Other user service functions for handling users
export async function getUsers() {
    const token = getToken();
    const response = await fetch(`${BASE_URL}/users`, {
        headers: {
            Authorization: token,
        },
    });
    if (!response.ok) throw new Error('Failed to fetch users');
    return response.json();
}

export async function createUser(user) {
    const response = await fetch(`${BASE_URL}/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(user),
    });
    if (!response.ok) throw new Error('Failed to create user');
    return response.json();
}

export async function updateUser(userId, updatedUserData) {
    const response = await fetch(`${BASE_URL}/users/${userId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            Authorization: getToken(),
        },
        body: JSON.stringify(updatedUserData),
    });
    if (!response.ok) throw new Error('Failed to update user');
    return response.json();
}

export async function deleteUser(userId) {
    const response = await fetch(`${BASE_URL}/users/${userId}`, {
        method: 'DELETE',
        headers: {
            Authorization: getToken(),
        },
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
    return data;
}
