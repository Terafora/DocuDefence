// src/App.js
import React, { useEffect, useState } from 'react';
import { getUsers, createUser } from './services/userService';
import { isLoggedIn, setToken, clearToken, getUserEmail } from './services/authService';
import UserList from './components/UserList';
import CreateUserForm from './components/CreateUserForm';
import Login from './components/Login';
import Logout from './components/Logout';

function App() {
    const [users, setUsers] = useState([]);
    const [newUser, setNewUser] = useState({ first_name: '', surname: '', email: '', birthdate: '', password: '' });
    const [loggedIn, setLoggedIn] = useState(isLoggedIn());
    const userEmail = getUserEmail();

    useEffect(() => {
        if (loggedIn) {
            const fetchUsers = async () => {
                try {
                    const users = await getUsers();
                    setUsers(users);
                } catch (error) {
                    console.error('Error fetching users:', error);
                }
            };
            fetchUsers();
        }
    }, [loggedIn]);

    const handleInputChange = (e) => setNewUser({ ...newUser, [e.target.name]: e.target.value });

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await createUser(newUser);
            setNewUser({ first_name: '', surname: '', email: '', birthdate: '', password: '' });
            const updatedUsers = await getUsers();
            setUsers(updatedUsers);
        } catch (error) {
            console.error('Error creating user:', error);
        }
    };

    const handleLogin = (token) => {
        setToken(token);
        setLoggedIn(true);
    };

    const handleLogout = () => {
        clearToken();
        setLoggedIn(false);
    };

    return (
        <div className="App">
            <h1>DocuDefense Frontend</h1>
            {loggedIn ? (
                <>
                    <Logout onLogout={handleLogout} />
                    <UserList users={users} userEmail={userEmail} /> {/* Pass userEmail */}
                    <CreateUserForm 
                        newUser={newUser} 
                        handleInputChange={handleInputChange} 
                        handleSubmit={handleSubmit} 
                    />
                </>
            ) : (
                <Login onLogin={handleLogin} />
            )}
        </div>
    );
}

export default App;
