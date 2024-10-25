import React from 'react';
import { clearToken } from '../services/authService';

function Logout({ onLogout }) {
  const handleLogout = () => {
    clearToken();
    onLogout();
  };

  return <button onClick={handleLogout}>Logout</button>;
}

export default Logout;
