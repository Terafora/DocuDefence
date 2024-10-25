import React from 'react';

function CreateUserForm({ newUser, handleInputChange, handleSubmit }) {
  return (
    <div>
      <h2>Create a New User</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          name="first_name"
          placeholder="First Name"
          value={newUser.first_name}
          onChange={handleInputChange}
        />
        <input
          type="text"
          name="surname"
          placeholder="Surname"
          value={newUser.surname}
          onChange={handleInputChange}
        />
        <input
          type="email"
          name="email"
          placeholder="Email"
          value={newUser.email}
          onChange={handleInputChange}
        />
        <input
          type="date"
          name="birthdate"
          value={newUser.birthdate}
          onChange={handleInputChange}
        />
        <input
          type="password"
          name="password"
          placeholder="Password"
          value={newUser.password}
          onChange={handleInputChange}
        />
        <button type="submit">Create User</button>
      </form>
    </div>
  );
}

export default CreateUserForm;
