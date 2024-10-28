import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AuthPanel from './components/AuthPanel';
import UserDashboard from './components/UserDashboard';
import About from './components/About';
import Navbar from './components/Navbar';
import Home from './components/Homepage';
import Footer from './components/Footer';
import { isLoggedIn, clearToken, getUserEmail } from './services/authService';
import './App.scss';

function App() {
    const [loggedIn, setLoggedIn] = useState(isLoggedIn());
    const [currentUser, setCurrentUser] = useState(null);
    const [showModal, setShowModal] = useState(false); // New state to control modal visibility

    useEffect(() => {
        if (loggedIn) {
            fetchCurrentUser();
        }
    }, [loggedIn]);

    const fetchCurrentUser = async () => {
        try {
            const userEmail = getUserEmail();
            setCurrentUser({ email: userEmail });
        } catch (error) {
            console.error('Error fetching current user:', error);
        }
    };

    const handleLogout = () => {
        clearToken();
        setLoggedIn(false);
        setCurrentUser(null);
    };

    const handleShowLogin = () => {
        setShowModal(true); // Show the login modal
    };

    const handleCloseLogin = () => {
        setShowModal(false); // Close the login modal
    };

    return (
        <Router>
            <div className="App d-lg-flex">
                <Navbar loggedIn={loggedIn} currentUser={currentUser} onLogout={handleLogout} onShowLogin={handleShowLogin} />
                <div className="main-content flex-grow-1 p-4">
                    <Routes>
                        <Route path="/" element={<Home />} />
                        <Route path="/about" element={<About />} />
                        <Route path="/dashboard" element={loggedIn ? <UserDashboard /> : <Navigate to="/" />} />
                        {/* Removed the /auth route completely */}
                        <Route path="*" element={<Navigate to="/" />} />
                    </Routes>
                </div>
                {showModal && <AuthPanel onLogin={() => setLoggedIn(true)} onClose={handleCloseLogin} />}
                <Footer />
            </div>
        </Router>
    );
}

export default App;
