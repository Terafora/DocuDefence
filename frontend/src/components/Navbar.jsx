import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import '../styles/navbar.scss';

function Navbar({ loggedIn, currentUser, onLogout, onShowLogin }) {
    const navigate = useNavigate();
    const [isOpen, setIsOpen] = useState(false);

    useEffect(() => {
        if (currentUser) console.log("Current User:", currentUser);
    }, [currentUser]);

    const handleLogout = () => {
        onLogout();
        navigate('/');
        setIsOpen(false);
    };

    return (
        <>
            {/* Shadow for Sidebar version for large screens (lg and up) */}
            <div className="navbar-shadow d-none d-lg-block"></div>

            {/* Sidebar version for large screens (lg and up) */}
            <nav className="custom-navbar d-none d-lg-flex flex-column">
                <Link to="/" className="navbar-brand text-center mb-3 navbar-money">DocuDefense</Link>
                <ul className="navbar-menu navbar-nav flex-column w-100 justify-content-between">
                    {loggedIn && currentUser && (
                        <li className="nav-item">
                            <span className="nav-link">Welcome back!</span>
                        </li>
                    )}
                    <li className="nav-item">
                        <Link to="/" className="nav-link">Home</Link>
                    </li>
                    <li className="nav-item">
                        <Link to="/about" className="nav-link">About</Link>
                    </li>
                    <li className="nav-item">
                        <Link to="/allusers" className="nav-link">All Users</Link>
                    </li>
                    {loggedIn && (
                        <li className="nav-item">
                            <Link to="/dashboard" className="nav-link">Dashboard</Link>
                        </li>
                    )}
                    <li className="nav-item">
                        {loggedIn ? (
                            <button onClick={handleLogout} className="nav-link btn-link">Logout</button>
                        ) : (
                            <button onClick={onShowLogin} className="nav-link btn-link">Login</button>
                        )}
                    </li>
                </ul>
            </nav>

            {/* Standard Navbar for medium and smaller screens */}
            <nav className="navbar navbar-expand-lg navbar-sc navbar-dark d-lg-none">
                <div className="container-fluid">
                    <Link to="/" className="navbar-brand">DocuDefense</Link>
                    
                    {/* Toggler button for dropdown */}
                    <button
                        className="navbar-toggler"
                        type="button"
                        data-bs-toggle="collapse"
                        data-bs-target="#navbarNav"
                        aria-controls="navbarNav"
                        aria-expanded={isOpen ? "true" : "false"}
                        aria-label="Toggle navigation"
                        onClick={() => setIsOpen(!isOpen)}
                    >
                        <span className="navbar-toggler-icon"></span>
                    </button>

                    {/* Collapsible navbar content */}
                    <div className={`collapse navbar-collapse ${isOpen ? "show" : ""}`} id="navbarNav">
                        <ul className="navbar-nav ms-auto mb-2 mb-lg-0">
                            {loggedIn && currentUser && (
                                <li className="nav-item">
                                    <span className="nav-link">Welcome back!</span>
                                </li>
                            )}
                            <li className="nav-item">
                                <Link to="/" className="nav-link" onClick={() => setIsOpen(false)}>Home</Link>
                            </li>
                            <li className="nav-item">
                                <Link to="/about" className="nav-link" onClick={() => setIsOpen(false)}>About</Link>
                            </li>
                            {loggedIn && (
                                <li className="nav-item">
                                    <Link to="/dashboard" className="nav-link" onClick={() => setIsOpen(false)}>Dashboard</Link>
                                </li>
                            )}
                            <li className="nav-item">
                                <Link to="/allusers" className="nav-link" onClick={() => setIsOpen(false)}>All Users</Link>
                            </li>
                            <li className="nav-item">
                                {loggedIn ? (
                                    <button onClick={handleLogout} className="nav-link btn-link">Logout</button>
                                ) : (
                                    <button onClick={onShowLogin} className="nav-link btn-link">Login</button>
                                )}
                            </li>
                        </ul>
                    </div>
                </div>
            </nav>
        </>
    );
}

export default Navbar;
