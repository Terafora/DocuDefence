import React, { useState } from 'react';

function UserList({ users = [], page, searchUsers, fetchUsers, handleNextPage, handlePreviousPage }) {
    const [searchTerm, setSearchTerm] = useState('');

    const handleSearch = (e) => {
        e.preventDefault();
        searchUsers({ term: searchTerm });
    };

    return (
        <div>
            <h2>Users</h2>

            {/* Search Form */}
            <form onSubmit={handleSearch}>
                <input
                    type="text"
                    placeholder="Search by First Name or Surname"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                />
                <button type="submit">Search</button>
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
