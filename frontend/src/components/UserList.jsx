import React from 'react';

function UserList({ users, userEmail, onUpdate, onDelete }) {
    return (
        <div>
            <h2>Users</h2>
            <ul>
                {users.map((user) => (
                    <li key={user.id}>
                        {user.first_name} {user.surname} - {user.email}
                        {user.email === userEmail && ( 
                            <>
                                <button onClick={() => onUpdate(user.id, { first_name: user.first_name, surname: user.surname })}>
                                    Update
                                </button>
                                <button onClick={() => onDelete(user.id)}>
                                    Delete
                                </button>
                            </>
                        )}
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default UserList;
