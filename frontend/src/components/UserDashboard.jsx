// src/components/UserDashboard.js

import React, { useState, useEffect, useCallback } from 'react';
import { uploadFile, getUserFiles, fetchUserIDByEmail, downloadFile, deleteFile } from '../services/userService';
import { getUserEmail } from '../services/authService';

function UserDashboard() {
    const [selectedFile, setSelectedFile] = useState(null);
    const [files, setFiles] = useState([]);
    const [loading, setLoading] = useState(false);
    const [userId, setUserId] = useState(null);

    const initializeUserID = useCallback(async () => {
        const email = getUserEmail();
        if (!email) {
            console.error("No email found for user.");
            return;
        }
        try {
            const userData = await fetchUserIDByEmail(email);
            setUserId(userData.id);
            console.log("Fetched user ID:", userData.id);
        } catch (error) {
            console.error('Error fetching user ID:', error);
        }
    }, []);

    useEffect(() => {
        initializeUserID();
    }, [initializeUserID]);

    const fetchUserFiles = useCallback(async () => {
        if (!userId) {
            console.error("User ID is undefined. Unable to fetch files.");
            return;
        }
        try {
            setLoading(true);
            const userData = await getUserFiles(userId);
            setFiles(userData.file_names || []);
            console.log("Fetched files for user ID:", userId);
        } catch (error) {
            console.error('Error fetching user files:', error);
        } finally {
            setLoading(false);
        }
    }, [userId]);

    useEffect(() => {
        if (userId) {
            fetchUserFiles();
        }
    }, [userId, fetchUserFiles]);

    const handleFileChange = (e) => setSelectedFile(e.target.files[0]);

    const handleUpload = async () => {
        if (!selectedFile) {
            alert('Please select a file first');
            return;
        }
        if (!userId) {
            console.error("User ID is undefined. Unable to upload file.");
            return;
        }
        try {
            setLoading(true);
            console.log("Uploading file for user ID:", userId);
            await uploadFile(userId, selectedFile);
            alert('File uploaded successfully');
            setSelectedFile(null);
            fetchUserFiles();
        } catch (error) {
            console.error('Error uploading file:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleDownload = async (filename) => {
        if (!userId) {
            console.error("User ID is undefined. Unable to download file.");
            return;
        }
        try {
            console.log("Downloading file:", filename);
            await downloadFile(userId, filename);
        } catch (error) {
            console.error('Error downloading file:', error);
        }
    };

    const handleDelete = async (filename) => {
        if (!userId) {
            console.error("User ID is undefined. Unable to delete file.");
            return;
        }
        try {
            console.log("Deleting file:", filename);
            await deleteFile(userId, filename);
            alert('File deleted successfully');
            fetchUserFiles();
        } catch (error) {
            console.error('Error deleting file:', error);
        }
    };

    return (
        <div>
            <h2>User Dashboard</h2>
            <div>
                <input type="file" accept="application/pdf" onChange={handleFileChange} />
                <button onClick={handleUpload} disabled={loading}>Upload PDF</button>
            </div>

            {loading && <p>Loading...</p>}

            <h3>My Files</h3>
            <ul>
                {files.length > 0 ? (
                    files.map((filename, index) => (
                        <li key={index}>
                            {filename}
                            <button onClick={() => handleDownload(filename)}>Download</button>
                            <button onClick={() => handleDelete(filename)}>Delete</button>
                        </li>
                    ))
                ) : (
                    <p>No files uploaded.</p>
                )}
            </ul>
        </div>
    );
}

export default UserDashboard;
