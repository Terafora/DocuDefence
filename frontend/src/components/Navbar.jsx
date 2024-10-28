import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/navbar.scss';

function Navbar({ loggedIn, currentUser, onLogout, onShowLogin }) {
    const navigate = useNavigate();

    const handleLogout = () => {
        onLogout();
        // Redirect to home or another appropriate page after logout
        navigate('/');
    };

    return (
        <nav className="navbar-wrapper">
            <div className="navbar">
                <div className="navbar-list">
                    <Link to="/" className="navbar-item">Home</Link>
                    <Link to="/about" className="navbar-item">About</Link>
                    {loggedIn && <Link to="/dashboard" className="navbar-item">Dashboard</Link>}
                    {loggedIn ? (
                        <button onClick={handleLogout} className="navbar-item">Logout</button>
                    ) : (
                        <button onClick={onShowLogin} className="navbar-item">Login</button>
                    )}
                    {loggedIn && currentUser && <span className="navbar-item">Welcome back, {currentUser.first_name}</span>}
                </div>
                <div className="wave-layer wave1"></div>
                <div className="wave-layer wave2"></div>
                <div className="wave-layer wave3"></div>
            </div>
        </nav>
    );
}

export default Navbar;
