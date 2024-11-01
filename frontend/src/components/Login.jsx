import React, { useState } from 'react';
import { setToken } from '../services/authService';
import { loginUser } from '../services/userService';

function Login({ onLogin }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const { token } = await loginUser({ email, password });
      setToken(token);
      setError(false);
      onLogin();
    } catch (error) {
      console.error('Login failed:', error);
      setError(true);
    }
  };

  return (
    <form className="custom-modal-form custom-modal add-skew" onSubmit={handleSubmit}>
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
      {error && (
        <p className="text-danger mt-2">Invalid credentials</p>
      )}
      <button type="submit" className="custom-btn primary-btn">Login</button>
    </form>
  );
}

export default Login;
