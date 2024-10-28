// src/components/AuthPanel.js
import React, { useState } from 'react';
import { loginUser, createUser } from '../services/userService';
import { setToken } from '../services/authService';

function AuthPanel({ onLogin, onClose }) {
    const [isRegistering, setIsRegistering] = useState(false);
    const [formData, setFormData] = useState({
        first_name: '',
        surname: '',
        email: '',
        birthdate: '',
        password: '',
    });

    const toggleForm = () => setIsRegistering(!isRegistering);
    
    const handleInputChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (isRegistering) {
                await createUser(formData);
                alert('Account created successfully. You may now log in.');
                setIsRegistering(false);
            } else {
                const { token } = await loginUser({ email: formData.email, password: formData.password });
                setToken(token);
                onLogin();
            }
            onClose(); // Close modal after submission
        } catch (error) {
            console.error(isRegistering ? 'Registration failed:' : 'Login failed:', error);
        }
    };

    return (
        <div className="modal show d-block" role="dialog" tabIndex="-1">
            <div className="modal-dialog" role="document">
                <div className="modal-content">
                    <div className="modal-header">
                        <h5 className="modal-title">{isRegistering ? 'Register' : 'Login'}</h5>
                        <button type="button" className="close" onClick={onClose}>&times;</button>
                    </div>
                    <div className="modal-body">
                        <form onSubmit={handleSubmit}>
                            {isRegistering && (
                                <>
                                    <input type="text" name="first_name" placeholder="First Name" value={formData.first_name} onChange={handleInputChange} className="form-control mb-2" />
                                    <input type="text" name="surname" placeholder="Surname" value={formData.surname} onChange={handleInputChange} className="form-control mb-2" />
                                    <input type="date" name="birthdate" value={formData.birthdate} onChange={handleInputChange} className="form-control mb-2" />
                                </>
                            )}
                            <input type="email" name="email" placeholder="Email" value={formData.email} onChange={handleInputChange} className="form-control mb-2" />
                            <input type="password" name="password" placeholder="Password" value={formData.password} onChange={handleInputChange} className="form-control mb-2" />
                            <button type="submit" className="btn btn-primary w-100">{isRegistering ? 'Register' : 'Login'}</button>
                        </form>
                    </div>
                    <div className="modal-footer">
                        <button className="btn btn-secondary" onClick={toggleForm}>{isRegistering ? 'Already have an account? Login' : 'Create an account'}</button>
                        <button className="btn btn-secondary" onClick={onClose}>Close</button>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default AuthPanel;
