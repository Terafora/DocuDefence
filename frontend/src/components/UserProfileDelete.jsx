import React from 'react';
import { deleteUser } from '../services/userService';

const UserProfileDelete = ({ userId }) => {

  const handleDelete = async () => {
    const confirmation = window.confirm('Are you sure you want to delete your account?');
    if (confirmation) {
      try {
        await deleteUser(userId);
        alert('Account deleted successfully');
        // Implement logout or redirect logic after deletion
      } catch (error) {
        console.error('Error deleting account:', error);
      }
    }
  };

  return (
    <button onClick={handleDelete} style={{ color: 'red' }}>
      Delete Account
    </button>
  );
};

export default UserProfileDelete;
