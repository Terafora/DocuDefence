import React, { useState } from 'react';

const UserProfileForm = ({ user, onUpdate }) => {
  const [updatedData, setUpdatedData] = useState(user);

  const handleChange = (e) => {
    setUpdatedData({
      ...updatedData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await onUpdate(user.id, updatedData); // Use the onUpdate prop function
      alert('Profile updated successfully');
    } catch (error) {
      console.error('Error updating profile:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        name="first_name"
        value={updatedData.first_name}
        onChange={handleChange}
      />
      <input
        type="text"
        name="surname"
        value={updatedData.surname}
        onChange={handleChange}
      />
      <button type="submit">Update Profile</button>
    </form>
  );
};

export default UserProfileForm;
