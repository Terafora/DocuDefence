import React, { useState, useEffect, useCallback } from 'react';
import { uploadFile, getUserFiles, fetchUserIDByEmail, deleteFile } from '../services/userService';
import { getUserEmail, getToken } from '../services/authService';
import PDFPreview from './pdfPreview';
import '../styles/userDashboard.scss';

function UserDashboard() {
    const [selectedFile, setSelectedFile] = useState(null);
    const [files, setFiles] = useState([]);
    const [loading, setLoading] = useState(false);
    const [userId, setUserId] = useState(null);
    const [expandedFiles, setExpandedFiles] = useState({});
    const [message, setMessage] = useState(null);
    const [showModal, setShowModal] = useState(false);
    const [fileToDelete, setFileToDelete] = useState(null);
    const [previewFile, setPreviewFile] = useState(null);

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

            const groupedFiles = {};
            userFiles.forEach((file) => {
                if (!groupedFiles[file.filename]) {
                    groupedFiles[file.filename] = [];
                }
                groupedFiles[file.filename].push(file);
            });

            Object.keys(groupedFiles).forEach(filename => {
                groupedFiles[filename].sort((a, b) => b.version - a.version);
            });

            setFiles(groupedFiles);
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

    const displayMessage = (msg, isError = false) => {
        setMessage({ text: msg, isError });
        setTimeout(() => setMessage(null), 3000);
    };

    const handleUpload = async () => {
        if (!selectedFile || !userId) return;
        try {
            setLoading(true);
            await uploadFile(userId, selectedFile);
            fetchUserFiles();
            displayMessage('File uploaded successfully!');
        } catch (error) {
            console.error('Error uploading file:', error);
            displayMessage('Error uploading file. Please try again.', true);
        } finally {
            setLoading(false);
        }
    };

    const handleDeleteConfirmed = async () => {
        if (!fileToDelete) return;

        setFiles(prevFiles => {
            const updatedFiles = { ...prevFiles };
            delete updatedFiles[fileToDelete];
            return updatedFiles;
        });

        try {
            await deleteFile(userId, fileToDelete);
            displayMessage('File deleted successfully!');
        } catch (error) {
            console.warn('Minor error deleting file:', error);
        } finally {
            setShowModal(false);
            setFileToDelete(null);
        }
    };

    const toggleExpand = (filename) => {
        setExpandedFiles(prevState => ({
            ...prevState,
            [filename]: !prevState[filename]
        }));
    };

    const confirmDelete = (filename) => {
        setFileToDelete(filename);
        setShowModal(true);
    };

    const handlePreview = async (filename, version) => {
        try {
            const token = getToken().replace("Bearer ", ""); // Remove the duplicate "Bearer" part if already present
            const response = await fetch(`http://localhost:8000/users/${userId}/files/${filename}/download`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/pdf'
                },
            });
    
            if (!response.ok) {
                throw new Error(`Failed to fetch PDF file: ${response.statusText}`);
            }
    
            const arrayBuffer = await response.arrayBuffer();
            setPreviewFile(new Uint8Array(arrayBuffer));
        } catch (error) {
            console.error('Error fetching PDF for preview:', error);
        }
    };
    

    return (
        <div>
            <h2>User Dashboard</h2>
            <input type="file" accept="application/pdf" onChange={handleFileChange} />
            <button onClick={handleUpload} disabled={loading}>Upload PDF</button>

            {loading && <p>Loading...</p>}
            {message && (
                <p className={`message ${message.isError ? 'error' : 'success'}`}>
                    {message.text}
                </p>
            )}

            <h3>My Files</h3>
            <ul style={{ listStyleType: 'none' }}>
                {Object.keys(files).length > 0 ? (
                    Object.keys(files).map((filename, index) => (
                        <li key={index} style={{ marginBottom: '20px', padding: '10px', border: '1px solid #ccc', borderRadius: '5px' }}>
                            <p>Filename: {filename}</p>
                            <p>Version: {files[filename][0].version}</p>
                            <p>Upload Date: {new Date(files[filename][0].upload_date).toLocaleDateString()}</p>
                            <button onClick={() => confirmDelete(filename)}>Delete</button>
                            <button onClick={() => handlePreview(filename, files[filename][0].version)}>Preview</button>
                            <button onClick={() => toggleExpand(filename)}>
                                {expandedFiles[filename] ? 'Hide Previous Versions' : 'Show Previous Versions'}
                            </button>

                            {expandedFiles[filename] && (
                                <ul style={{ paddingLeft: '20px', marginTop: '10px', listStyleType: 'none' }}>
                                    {files[filename].slice(1).map((versionedFile, versionIndex) => (
                                        <li key={versionIndex}>
                                            <p>Version: {versionedFile.version}</p>
                                            <p>Upload Date: {new Date(versionedFile.upload_date).toLocaleDateString()}</p>
                                            <button onClick={() => handlePreview(filename, versionedFile.version)}>Preview</button>
                                        </li>
                                    ))}
                                </ul>
                            )}
                        </li>
                    ))
                ) : (
                    <p>No files uploaded.</p>
                )}
            </ul>

            {/* PDF Preview */}
            {previewFile && (
                <div className="pdf-preview-modal">
                    <PDFPreview fileBlob={previewFile} />
                    <button onClick={() => setPreviewFile(null)}>Close Preview</button>
                </div>
            )}

            {/* Confirmation Modal */}
            {showModal && (
                <div className="modal">
                    <div className="modal-content">
                        <p>Are you sure you want to delete this file?</p>
                        <button onClick={handleDeleteConfirmed}>Yes, Delete</button>
                        <button onClick={() => setShowModal(false)}>Cancel</button>
                    </div>
                </div>
            )}
        </div>
    );
}

export default UserDashboard;
