import React from 'react';

const UserProfileDelete = ({ onDelete }) => {
  const handleDelete = async () => {
    const confirmation = window.confirm('Are you sure you want to delete your account?');
    if (confirmation) {
      try {
        await onDelete(); // Call the onDelete prop function
        alert('Account deleted successfully');
        // Optionally implement logout or redirect logic after deletion
      } catch (error) {
        console.error('Error deleting account:', error);
      }
    }
  };

  return (
    <div>
      <h3>Account Management</h3>
      <button onClick={handleDelete} style={{ color: 'red' }}>
        Delete Account
      </button>
    </div>
  );
};

export default UserProfileDelete;
