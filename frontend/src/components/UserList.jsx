import React, { useState, useEffect } from 'react';

function UserList({ users = [], userEmail, onUpdate, onDelete, searchUsers, fetchUsers }) {
    const [searchTerm, setSearchTerm] = useState('');

    // Update search and fetch users when the search term changes
    useEffect(() => {
        if (searchTerm.trim() === '') {
            fetchUsers(); // Load all users when search bar is cleared
        } else {
            searchUsers({ term: searchTerm });
        }
    }, [searchTerm, searchUsers, fetchUsers]);

    return (
        <div>
            <h2>Users</h2>

            {/* Unified Search Form */}
            <input
                type="text"
                placeholder="Search by First Name or Surname"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
            />

            {/* Filtered User List */}
            <ul>
                {users.map((user) => (
                    <li key={user.id}>
                        {user.first_name} {user.surname}
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default UserList;
