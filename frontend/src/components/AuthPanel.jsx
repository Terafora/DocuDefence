import React, { useState, useEffect } from 'react';
import { loginUser, createUser } from '../services/userService';
import { setToken } from '../services/authService';

function AuthPanel({ onLogin, onClose, isRegistering: initialRegisteringState }) {
  const [isRegistering, setIsRegistering] = useState(initialRegisteringState);
  const [formData, setFormData] = useState({
    first_name: '',
    surname: '',
    email: '',
    birthdate: '',
    password: '',
  });

  useEffect(() => {
    setIsRegistering(initialRegisteringState);
  }, [initialRegisteringState]);

  const toggleForm = () => setIsRegistering(!isRegistering);
    
  const handleInputChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (isRegistering) {
        await createUser(formData);
        alert('Account created successfully. You may now log in.');
        setIsRegistering(false);
      } else {
        const { token } = await loginUser({ email: formData.email, password: formData.password });
        setToken(token);
        onLogin();
      }
      onClose();
    } catch (error) {
      console.error(isRegistering ? 'Registration failed:' : 'Login failed:', error);
    }
  };

  return (
    <div className="custom-modal show">
      <div className="custom-modal-dialog">
        <div className="custom-modal-content">
          <div className="custom-modal-header">
            <h5 className="modal-title">{isRegistering ? 'Register' : 'Login'}</h5>
            <button type="button" className="close-btn" onClick={onClose}>&times;</button>
          </div>
          <div className="custom-modal-body">
            <form className="custom-modal-form" onSubmit={handleSubmit}>
              {isRegistering && (
                <>
                  <input type="text" name="first_name" placeholder="First Name" value={formData.first_name} onChange={handleInputChange} className="custom-input" />
                  <input type="text" name="surname" placeholder="Surname" value={formData.surname} onChange={handleInputChange} className="custom-input" />
                  <input type="date" name="birthdate" value={formData.birthdate} onChange={handleInputChange} className="custom-input" />
                </>
              )}
              <input type="email" name="email" placeholder="Email" value={formData.email} onChange={handleInputChange} className="custom-input" />
              <input type="password" name="password" placeholder="Password" value={formData.password} onChange={handleInputChange} className="custom-input" />
              <button type="submit" className="custom-btn primary-btn w-100">{isRegistering ? 'Register' : 'Login'}</button>
            </form>
          </div>
          <div className="custom-modal-footer">
            <button className="custom-btn secondary-btn" onClick={toggleForm}>{isRegistering ? 'Already have an account? Login' : 'Create an account'}</button>
            <button className="custom-btn secondary-btn" onClick={onClose}>Close</button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AuthPanel;
