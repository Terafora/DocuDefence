import React, { useState, useEffect, useCallback } from 'react';
import { uploadFile, getUserFiles, fetchUserIDByEmail, downloadFile, deleteFile } from '../services/userService';
import { getUserEmail } from '../services/authService';

function UserDashboard() {
    const [selectedFile, setSelectedFile] = useState(null);
    const [files, setFiles] = useState([]); // Default to an empty array
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
        } catch (error) {
            console.error('Error fetching user ID:', error);
        }
    }, []);

    useEffect(() => {
        initializeUserID();
    }, [initializeUserID]);

    const fetchUserFiles = useCallback(async () => {
        if (!userId) return;
        try {
            setLoading(true);
            const userFiles = await getUserFiles(userId);
            setFiles(userFiles || []); // Ensure files is an array even if null is returned
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
        if (!selectedFile || !userId) return;
        try {
            setLoading(true);
            await uploadFile(userId, selectedFile);
            fetchUserFiles();
        } catch (error) {
            console.error('Error uploading file:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleDownload = async (filename) => {
        if (!filename) {
            console.error("Filename is undefined. Cannot download.");
            return;
        }
        try {
            await downloadFile(userId, filename);
        } catch (error) {
            console.error('Error downloading file:', error);
        }
    };

    const handleDelete = async (filename) => {
        if (!filename) {
            console.error("Filename is undefined. Cannot delete.");
            return;
        }
        try {
            const encodedFilename = encodeURIComponent(filename);
            await deleteFile(userId, encodedFilename);
            fetchUserFiles();
        } catch (error) {
            console.error('Error deleting file:', error);
        }
    };

    return (
        <div>
            <h2>User Dashboard</h2>
            <input type="file" accept="application/pdf" onChange={handleFileChange} />
            <button onClick={handleUpload} disabled={loading}>Upload PDF</button>

            {loading && <p>Loading...</p>}

            <h3>My Files</h3>
            <ul>
                {files.length > 0 ? (
                    files.map((file, index) => (
                        <li key={file.id || index}>
                            <p>Filename: {file.filename || "N/A"}</p>
                            <p>Version: {file.version || "N/A"}</p>
                            <p>Upload Date: {file.upload_date ? new Date(file.upload_date).toLocaleDateString() : "Unknown"}</p>
                            <button onClick={() => handleDownload(file.filename)}>Download</button>
                            <button onClick={() => handleDelete(file.filename)}>Delete</button>
                        </li>
                    ))
                ) : (
                    <p>No files uploaded.</p> // Fallback message when files array is empty
                )}
            </ul>
        </div>
    );
}

export default UserDashboard;
