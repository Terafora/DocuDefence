import React, { useState, useEffect, useCallback } from 'react';
import { uploadFile, getUserFiles, fetchUserIDByEmail, deleteFile, updateUser, deleteUser } from '../services/userService';
import { getUserEmail, getToken } from '../services/authService';
import PDFPreview from './pdfPreview';
import UserProfileForm from './UserProfileForm';
import UserProfileDelete from './UserProfileDelete';
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
    const [userProfile, setUserProfile] = useState(null);

    const [showProfileModal, setShowProfileModal] = useState(false);
    const [showDeleteModal, setShowDeleteModal] = useState(false);

    const initializeUserID = useCallback(async () => {
        const email = getUserEmail();
        if (!email) {
            console.error("No email found for user.");
            return;
        }
        try {
            const userData = await fetchUserIDByEmail(email);
            setUserId(userData.id);
            setUserProfile(userData); // Store user data for the profile form
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
            const token = getToken().replace("Bearer ", "");
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

    const handleDownload = async (filename, version) => {
        try {
            const token = getToken().replace("Bearer ", "");
            const response = await fetch(`http://localhost:8000/users/${userId}/files/${filename}/download`, {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/pdf'
                },
            });

            if (!response.ok) {
                throw new Error(`Failed to download file: ${response.statusText}`);
            }

            const blob = await response.blob();
            const downloadUrl = window.URL.createObjectURL(blob);
            const link = document.createElement('a');
            link.href = downloadUrl;
            link.download = `${filename}_v${version}.pdf`;
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
        } catch (error) {
            console.error('Error downloading file:', error);
        }
    };

    const handleProfileUpdate = async (id, updatedData) => {
        try {
            await updateUser(id, updatedData);
            setUserProfile(updatedData);
            displayMessage('Profile updated successfully!');
            setShowProfileModal(false);
        } catch (error) {
            console.error('Error updating profile:', error);
            displayMessage('Error updating profile. Please try again.', true);
        }
    };

    const handleProfileDelete = async (id) => {
        try {
            await deleteUser(id);
            displayMessage('Account deleted successfully!');
            setShowDeleteModal(false);
        } catch (error) {
            console.error('Error deleting account:', error);
            displayMessage('Error deleting account. Please try again.', true);
        }
    };

    return (
        <div className="user-dashboard">
            <h2 className="dashboard-title">User Dashboard</h2>

            <div className="file-upload-section custom-card">
                <input
                    type="file"
                    className="custom-input file-input"
                    accept="application/pdf"
                    onChange={handleFileChange}
                    style={{ maxWidth: "400px", margin: "0 auto" }}
                />
                <button className="custom-btn primary-btn mt-3" onClick={handleUpload} disabled={loading}>
                    Upload PDF
                </button>
            </div>

            {loading && <p className="loading-text">Loading...</p>}
            {message && (
                <p className={`message ${message.isError ? 'error-message' : 'success-message'}`}>
                    {message.text}
                </p>
            )}

            <h3 className="my-files-title">My Files</h3>
            <ul className="file-list list-unstyled mt-4">
                {Object.keys(files).length > 0 ? (
                    Object.keys(files).map((filename, index) => (
                        <li key={index} className="file-item custom-card">
                            <p><strong>Filename:</strong> {filename}</p>
                            <p><strong>Version:</strong> {files[filename][0].version}</p>
                            <p><strong>Upload Date:</strong> {new Date(files[filename][0].upload_date).toLocaleDateString()}</p>
                            <div className="file-actions">
                                <button className="custom-btn danger-btn" onClick={() => confirmDelete(filename)}>Delete</button>
                                <button className="custom-btn secondary-btn" onClick={() => handlePreview(filename, files[filename][0].version)}>Preview</button>
                                <button className="custom-btn info-btn" onClick={() => handleDownload(filename, files[filename][0].version)}>Download</button>
                                <button className="custom-btn info-btn" onClick={() => toggleExpand(filename)}>
                                    {expandedFiles[filename] ? 'Hide Previous Versions' : 'Show Previous Versions'}
                                </button>
                            </div>

                            {expandedFiles[filename] && (
                                <ul className="version-list list-unstyled">
                                    {files[filename].slice(1).map((versionedFile, versionIndex) => (
                                        <li key={versionIndex} className="version-item">
                                            <p><strong>Version:</strong> {versionedFile.version}</p>
                                            <p><strong>Upload Date:</strong> {new Date(versionedFile.upload_date).toLocaleDateString()}</p>
                                            <button className="custom-btn secondary-btn small-btn" onClick={() => handlePreview(filename, versionedFile.version)}>Preview</button>
                                            <button className="custom-btn info-btn small-btn" onClick={() => handleDownload(filename, versionedFile.version)}>Download</button>
                                        </li>
                                    ))}
                                </ul>
                            )}
                        </li>
                    ))
                ) : (
                    <p className="no-files-message">No files uploaded.</p>
                )}
            </ul>

            {previewFile && (
                <div className="custom-modal show" tabIndex="-1">
                    <div className="custom-modal-dialog">
                        <div className="custom-modal-content">
                            <div className="custom-modal-header">
                                <h5 className="modal-title">PDF Preview</h5>
                                <button type="button" className="close-btn" onClick={() => setPreviewFile(null)}>&times;</button>
                            </div>
                            <div className="modal-body">
                                <PDFPreview fileBlob={previewFile} />
                            </div>
                            <div className="modal-footer custom-modal-footer">
                                <button className="custom-btn secondary-btn" onClick={() => setPreviewFile(null)}>Close Preview</button>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            <div className="user-profile-section mt-5">
                <h3 className="profile-section-title">User Profile</h3>
                <button className="custom-btn info-btn" onClick={() => setShowProfileModal(true)}>Edit Profile</button>
                <button className="custom-btn danger-btn ms-2" onClick={() => setShowDeleteModal(true)}>Delete Account</button>
            </div>

            {showProfileModal && (
                <div className="custom-modal show" tabIndex="-1">
                    <div className="custom-modal-dialog">
                        <div className="custom-modal-content">
                            <div className="custom-modal-header">
                                <h5 className="modal-title">Update Profile</h5>
                                <button type="button" className="close-btn" onClick={() => setShowProfileModal(false)}>&times;</button>
                            </div>
                            <div className="modal-body custom-modal-body">
                                <UserProfileForm user={userProfile} onUpdate={handleProfileUpdate} />
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {showDeleteModal && (
                <div className="custom-modal show" tabIndex="-1">
                    <div className="custom-modal-dialog">
                        <div className="custom-modal-content">
                            <div className="custom-modal-header">
                                <h2 style={{ color: 'black', fontSize: '1.5rem' }}>Are you sure you want to delete your account and files?</h2>
                                <button type="button" className="close-btn" onClick={() => setShowDeleteModal(false)}>&times;</button>
                            </div>
                            <div className="modal-body custom-modal-body">
                                <UserProfileDelete userId={userId} onDelete={handleProfileDelete} />
                            </div>
                            <div className="modal-footer custom-modal-footer">
                                <button className="custom-btn secondary-btn" onClick={() => setShowDeleteModal(false)}>Cancel</button>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {showModal && (
                <div className="custom-modal show" tabIndex="-1">
                    <div className="custom-modal-dialog">
                        <div className="custom-modal-content">
                            <div className="modal-body custom-modal-body">
                                <p>Are you sure you want to delete this file?</p>
                                <button className="custom-btn danger-btn" onClick={handleDeleteConfirmed}>Yes, Delete</button>
                                <button className="custom-btn secondary-btn ms-2" onClick={() => setShowModal(false)}>Cancel</button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

export default UserDashboard;
