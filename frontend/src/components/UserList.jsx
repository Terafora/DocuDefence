import React from 'react';

function UserList({ users }) {
  return (
    <div>
      <h2>Users</h2>
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
