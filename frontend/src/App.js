import React, { useEffect, useState } from 'react';
import { getUsers, createUser } from './services/userService';
import UserList from './components/UserList';
import CreateUserForm from './components/CreateUserForm';

function App() {
  const [users, setUsers] = useState([]);
  const [newUser, setNewUser] = useState({ first_name: '', surname: '', email: '', birthdate: '', password: '' });

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const users = await getUsers();
        setUsers(users);
      } catch (error) {
        console.error('Error fetching users:', error);
      }
    };

    fetchUsers();
  }, []);

  const handleInputChange = (e) => {
    setNewUser({ ...newUser, [e.target.name]: e.target.value });
  };

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

  return (
    <div className="App">
      <h1>DocuDefense Frontend</h1>
      <UserList users={users} />
      <CreateUserForm 
        newUser={newUser} 
        handleInputChange={handleInputChange} 
        handleSubmit={handleSubmit} 
      />
    </div>
  );
}

export default App;
