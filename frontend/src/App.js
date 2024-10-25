import React, { useEffect, useState } from 'react';
import { getUsers, createUser, updateUser, deleteUser } from './services/userService'; // Include updateUser and deleteUser
import { isLoggedIn, setToken, clearToken, getUserEmail } from './services/authService';
import UserList from './components/UserList';
import CreateUserForm from './components/CreateUserForm';
import Login from './components/Login';
import Logout from './components/Logout';
import UserProfileForm from './components/UserProfileForm';  // For updating user info
import UserProfileDelete from './components/UserProfileDelete'; // For deleting user account

function App() {
    const [users, setUsers] = useState([]);
    const [newUser, setNewUser] = useState({ first_name: '', surname: '', email: '', birthdate: '', password: '' });
    const [loggedIn, setLoggedIn] = useState(isLoggedIn());
    const userEmail = getUserEmail(); // Get the logged-in user's email
    const [currentUser, setCurrentUser] = useState(null); // Store the current user

    useEffect(() => {
        if (loggedIn) {
            const fetchUsers = async () => {
                try {
                    const users = await getUsers();
                    setUsers(users);

                    // Find the logged-in user from the list of users
                    const currentUser = users.find(user => user.email === userEmail);
                    setCurrentUser(currentUser);
                } catch (error) {
                    console.error('Error fetching users:', error);
                }
            };
            fetchUsers();
        }
    }, [loggedIn, userEmail]);

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
        setCurrentUser(null); // Reset current user on logout
    };

    return (
        <div className="App">
            <h1>DocuDefense Frontend</h1>
            {loggedIn ? (
                <>
                    <Logout onLogout={handleLogout} />
                    <h2>Logged in as: {userEmail || 'No current user found'}</h2> {/* Display the current user */}
                    <UserList users={users} userEmail={userEmail} />
                    
                    {currentUser && (
                        <>
                            <h3>Edit Your Profile</h3>
                            {/* Render update and delete buttons */}
                            <UserProfileForm user={currentUser} />
                            <UserProfileDelete userId={currentUser.id} />
                        </>
                    )}
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
