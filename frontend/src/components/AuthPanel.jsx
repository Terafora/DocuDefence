
import React, { useState } from 'react';
import { loginUser, createUser } from '../services/userService';
import { setToken } from '../services/authService';

function AuthPanel({ onLogin }) {
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
        } catch (error) {
            console.error(isRegistering ? 'Registration failed:' : 'Login failed:', error);
        }
    };

    return (
        <div>
            <h2>{isRegistering ? 'Register' : 'Login'}</h2>
            <form onSubmit={handleSubmit}>
                {isRegistering && (
                    <>
                        <input
                            type="text"
                            name="first_name"
                            placeholder="First Name"
                            value={formData.first_name}
                            onChange={handleInputChange}
                        />
                        <input
                            type="text"
                            name="surname"
                            placeholder="Surname"
                            value={formData.surname}
                            onChange={handleInputChange}
                        />
                        <input
                            type="date"
                            name="birthdate"
                            value={formData.birthdate}
                            onChange={handleInputChange}
                        />
                    </>
                )}
                <input
                    type="email"
                    name="email"
                    placeholder="Email"
                    value={formData.email}
                    onChange={handleInputChange}
                />
                <input
                    type="password"
                    name="password"
                    placeholder="Password"
                    value={formData.password}
                    onChange={handleInputChange}
                />
                <button type="submit">{isRegistering ? 'Register' : 'Login'}</button>
            </form>
            <button onClick={toggleForm}>
                {isRegistering ? 'Already have an account? Login' : 'Create an account'}
            </button>
        </div>
    );
}

export default AuthPanel;
