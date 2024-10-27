// src/components/UserDashboard.js
import React, { useEffect, useState } from 'react';
import { getUsers, updateUser, deleteUser } from '../services/userService';
import UserList from './UserList';
import UserProfileForm from './UserProfileForm';
import UserProfileDelete from './UserProfileDelete';
import Logout from './Logout';
import { getUserEmail } from '../services/authService';

function UserDashboard({ onLogout }) {
    const [users, setUsers] = useState([]);
    const [currentUser, setCurrentUser] = useState(null);
    const userEmail = getUserEmail();

    useEffect(() => {
        const fetchUsers = async () => {
            try {
                const users = await getUsers();
                setUsers(users);
                
                // Find the logged-in user from the list of users
                const loggedInUser = users.find(user => user.email === userEmail);
                setCurrentUser(loggedInUser);
            } catch (error) {
                console.error('Error fetching users:', error);
            }
        };
        fetchUsers();
    }, [userEmail]);

    const handleDelete = async (userId) => {
        try {
            await deleteUser(userId);
            const updatedUsers = await getUsers();
            setUsers(updatedUsers);
            setCurrentUser(null); // Clear current user if they delete their own account
        } catch (error) {
            console.error('Error deleting account:', error);
        }
    };

    const handleUpdate = async (userId, updatedData) => {
        try {
            await updateUser(userId, updatedData);
            const updatedUsers = await getUsers();
            setUsers(updatedUsers);
        } catch (error) {
            console.error('Error updating user:', error);
        }
    };

    return (
        <div>
            <Logout onLogout={onLogout} />
            <h2>Logged in as: {userEmail || 'No current user found'}</h2>
            <UserList users={users} userEmail={userEmail} onUpdate={handleUpdate} onDelete={handleDelete} />

            {currentUser && (
                <>
                    <h3>Edit Your Profile</h3>
                    <UserProfileForm user={currentUser} onUpdate={handleUpdate} />
                    <UserProfileDelete userId={currentUser.id} onDelete={handleDelete} />
                </>
            )}
        </div>
    );
}

export default UserDashboard;
