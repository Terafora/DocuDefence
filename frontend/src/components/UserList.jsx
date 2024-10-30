import React, { useState, useEffect } from 'react';

function UserList({ users = [], page, searchUsers, fetchUsers, handleNextPage, handlePreviousPage, searchTerm }) {
    const [localSearchTerm, setLocalSearchTerm] = useState('');

    useEffect(() => {
        if (localSearchTerm === '') {
            fetchUsers();
        }
    }, [localSearchTerm, fetchUsers]);

    const handleSearch = (e) => {
        e.preventDefault();
        searchUsers({ term: localSearchTerm });
    };

    const handleClearSearch = () => {
        setLocalSearchTerm('');
        fetchUsers();
    };

    return (
        <div>
            <h2>Users</h2>

            {/* Search Form */}
            <form onSubmit={handleSearch}>
                <input
                    type="text"
                    placeholder="Search by First Name or Surname"
                    value={localSearchTerm}
                    onChange={(e) => setLocalSearchTerm(e.target.value)}
                />
                <button type="submit">Search</button>
                {localSearchTerm && <button type="button" onClick={handleClearSearch}>Clear</button>}
            </form>

            {/* User List */}
            <ul>
                {users.map((user) => (
                    <li key={user.id}>
                        {user.first_name} {user.surname}
                    </li>
                ))}
            </ul>

            {/* Pagination Controls */}
            <div>
                <button onClick={handlePreviousPage} disabled={page === 1}>Previous</button>
                <button onClick={handleNextPage}>Next</button>
            </div>
        </div>
    );
}

export default UserList;
