import React from 'react';

const UserProfileDelete = ({ userId, onDelete }) => {
  const handleDelete = async () => {
    const confirmation = window.confirm('Are you sure you want to delete your account?');
    if (confirmation) {
      try {
        await onDelete(userId); // Call the onDelete function from props
        alert('Account deleted successfully');
      } catch (error) {
        console.error('Error deleting account:', error);
      }
    }
  };

  return (
    <div className="d-flex flex-column px-5 py-3">
      <button onClick={handleDelete} className="btn btn-danger mt-3">
        Delete Account
      </button>
    </div>
  );
};

export default UserProfileDelete;
