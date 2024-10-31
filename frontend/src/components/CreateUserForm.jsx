import React from 'react';

function CreateUserForm({ newUser, handleInputChange, handleSubmit }) {
  return (
    <div className="custom-modal show">
      <div className="custom-modal-dialog">
        <div className="custom-modal-content">
          <div className="custom-modal-header">
            <h2 className="modal-title">Create a New User</h2>
          </div>
          <div className="modal-body custom-modal-body">
            <form onSubmit={handleSubmit} className="custom-modal-form">
              <input
                type="text"
                name="first_name"
                placeholder="First Name"
                value={newUser.first_name}
                onChange={handleInputChange}
                className="custom-input"
              />
              <input
                type="text"
                name="surname"
                placeholder="Surname"
                value={newUser.surname}
                onChange={handleInputChange}
                className="custom-input"
              />
              <input
                type="email"
                name="email"
                placeholder="Email"
                value={newUser.email}
                onChange={handleInputChange}
                className="custom-input"
              />
              <input
                type="date"
                name="birthdate"
                value={newUser.birthdate}
                onChange={handleInputChange}
                className="custom-input"
              />
              <input
                type="password"
                name="password"
                placeholder="Password"
                value={newUser.password}
                onChange={handleInputChange}
                className="custom-input"
              />
              <button type="submit" className="custom-btn primary-btn">Create User</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}

export default CreateUserForm;
