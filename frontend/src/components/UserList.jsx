// src/components/UserList.jsx
import React from 'react';

function UserList({ users, userEmail }) {
    return (
        <div>
            <h2>Users</h2>
            <ul>
                {users.map((user) => (
                    <li key={user.id}>
                        {user.first_name} {user.surname} - {user.email}
                        {user.email === userEmail && ( // Show buttons if user is the owner
                            <>
                                <button onClick={() => console.log("Update clicked for", user.id)}>Update</button>
                                <button onClick={() => console.log("Delete clicked for", user.id)}>Delete</button>
                            </>
                        )}
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default UserList;
