import React, { useState } from 'react';

const UserProfileForm = ({ user, onUpdate }) => {
  const [updatedData, setUpdatedData] = useState(user);
  const [newPassword, setNewPassword] = useState('');

  const handleChange = (e) => {
    const { name, value } = e.target;
    setUpdatedData((prevData) => ({
      ...prevData,
      [name]: value,
    }));
  };

  const handlePasswordChange = (e) => {
    setNewPassword(e.target.value);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const dataToUpdate = { ...updatedData };
    if (newPassword) {
      dataToUpdate.password = newPassword; // Include the new password only if changed
    }
    try {
      await onUpdate(user.id, dataToUpdate);
      alert('Profile updated successfully');
    } catch (error) {
      console.error('Error updating profile:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="d-flex flex-column gap-3 px-5 py-3">
      <input
        type="text"
        name="first_name"
        value={updatedData.first_name}
        onChange={handleChange}
        placeholder="First Name"
        className="form-control"
      />
      <input
        type="text"
        name="surname"
        value={updatedData.surname}
        onChange={handleChange}
        placeholder="Surname"
        className="form-control"
      />
      <input
        type="email"
        name="email"
        value={updatedData.email}
        onChange={handleChange}
        placeholder="Email"
        className="form-control"
      />
      <input
        type="password"
        name="password"
        value={newPassword}
        onChange={handlePasswordChange}
        placeholder="New Password (leave blank to keep current)"
        className="form-control"
      />
      <button type="submit" className="btn btn-primary mt-3">Update Profile</button>
    </form>
  );
};

export default UserProfileForm;
