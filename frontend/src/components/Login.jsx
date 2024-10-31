import React, { useState } from 'react';
import { setToken } from '../services/authService';
import { loginUser } from '../services/userService';

function Login({ onLogin }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const { token } = await loginUser({ email, password });
      setToken(token);
      onLogin();
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  return (
    <form className="custom-modal-form" onSubmit={handleSubmit}>
    <input 
        type="email" 
        placeholder="Email" 
        value={email} 
        onChange={(e) => setEmail(e.target.value)} 
        className="custom-input"
    />
    <input 
        type="password" 
        placeholder="Password" 
        value={password} 
        onChange={(e) => setPassword(e.target.value)} 
        className="custom-input"
    />
    <button type="submit" className="custom-btn primary-btn">Login</button>
</form>

  );
}

export default Login;
