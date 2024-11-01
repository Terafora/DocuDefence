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
        <div className="d-flex justify-content-center align-items-center vh-100">
            <div className="container my-5 text-center">
                <h2 className="display-6 mb-4">Users</h2>

                {/* Search Form */}
                <form className="d-flex justify-content-center mb-3" onSubmit={handleSearch}>
                    <input
                        type="text"
                        className="form-control me-2"
                        placeholder="Search by First Name or Surname"
                        value={localSearchTerm}
                        onChange={(e) => setLocalSearchTerm(e.target.value)}
                        style={{ maxWidth: "400px" }}
                    />
                    <button type="submit" className="btn btn-primary">Search</button>
                    {localSearchTerm && <button type="button" className="btn btn-secondary ms-2" onClick={handleClearSearch}>Clear</button>}
                </form>

                {/* User List */}
                <ul className="list-unstyled">
                    {users.map((user) => (
                        <li key={user.id} className="mb-2">{user.first_name} {user.surname}</li>
                    ))}
                </ul>

                {/* Pagination Controls */}
                <div className="d-flex justify-content-center mt-4">
                    <button className="btn btn-secondary me-2" onClick={handlePreviousPage} disabled={page === 1}>Previous</button>
                    <button className="btn btn-secondary" onClick={handleNextPage}>Next</button>
                </div>
            </div>
        </div>
    );
}

export default UserList;
