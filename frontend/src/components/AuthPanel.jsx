import React, { useState } from 'react';
import Login from './Login';
import CreateUserForm from './CreateUserForm';
import { createUser } from '../services/userService'; // Import createUser from userService

function AuthPanel({ onLogin }) {
    const [newUser, setNewUser] = useState({ first_name: '', surname: '', email: '', birthdate: '', password: '' });

    const handleInputChange = (e) => setNewUser({ ...newUser, [e.target.name]: e.target.value });

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await createUser(newUser);
            setNewUser({ first_name: '', surname: '', email: '', birthdate: '', password: '' });
            alert('Account created successfully. Please log in.');
        } catch (error) {
            console.error('Error creating user:', error);
        }
    };

    return (
        <div style={{ display: 'flex', justifyContent: 'space-between', padding: '20px' }}>
            {/* Login form */}
            <div style={{ width: '45%' }}>
                <h2>Login</h2>
                <Login onLogin={onLogin} />
            </div>

            {/* Create New User form */}
            <div style={{ width: '45%' }}>
                <h2>Create a New User</h2>
                <CreateUserForm 
                    newUser={newUser} 
                    handleInputChange={handleInputChange} 
                    handleSubmit={handleSubmit} 
                />
            </div>
        </div>
    );
}

export default AuthPanel;
