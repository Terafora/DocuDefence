import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AuthPanel from './components/AuthPanel';
import UserDashboard from './components/UserDashboard';
import About from './components/About';
import Navbar from './components/Navbar';
import Home from './components/Homepage';
import Footer from './components/Footer';
import UserList from './components/UserList';
import { isLoggedIn, clearToken, getUserEmail } from './services/authService';
import { getUsers } from './services/userService';
import './App.scss';

function App() {
    const [loggedIn, setLoggedIn] = useState(isLoggedIn());
    const [currentUser, setCurrentUser] = useState(null);
    const [users, setUsers] = useState([]);
    const [showModal, setShowModal] = useState(false);

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

    // Fetch users function for initial load
    const fetchUsers = async () => {
        try {
            const fetchedUsers = await getUsers();
            setUsers(fetchedUsers);
        } catch (error) {
            console.error('Error fetching users:', error);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, []);

    const handleLogout = () => {
        clearToken();
        setLoggedIn(false);
        setCurrentUser(null);
    };

    const handleShowLogin = () => {
        setShowModal(true);
    };

    const handleCloseLogin = () => {
        setShowModal(false);
    };

    // Search users function
    const searchUsers = async (searchCriteria) => {
        const query = new URLSearchParams();
        if (searchCriteria.term) query.append("term", searchCriteria.term);

        try {
            // Adjust the path to the correct endpoint
            const response = await fetch(`http://localhost:8000/api/users/search?${query.toString()}`);
            if (!response.ok) throw new Error("Error fetching search results");

            const data = await response.json();
            setUsers(data); // Set the filtered data to state
        } catch (error) {
            console.error("Error fetching search results:", error);
        }
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
                        <Route path="/allusers" element={<UserList users={users} userEmail={currentUser?.email} searchUsers={searchUsers} fetchUsers={fetchUsers} />}/>
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
