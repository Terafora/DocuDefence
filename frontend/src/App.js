import React, { useEffect, useState, useCallback } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import AuthPanel from './components/AuthPanel';
import UserDashboard from './components/UserDashboard';
import About from './components/About';
import Navbar from './components/Navbar';
import Home from './components/Homepage';
import Footer from './components/Footer';
import UserList from './components/UserList';
import { isLoggedIn, clearToken, getUserEmail } from './services/authService';
import bgElement from './assets/bg-element.svg';
import './App.scss';

function MainContentWrapper({ children }) {
    const [animateContent, setAnimateContent] = useState(false);
    const location = useLocation();

    useEffect(() => {
        setAnimateContent(true);
        const timeout = setTimeout(() => setAnimateContent(false), 1000);
        return () => clearTimeout(timeout);
    }, [location]);

    useEffect(() => {
        const handleScroll = () => {
            const scrollTop = window.scrollY;
            const mainContent = document.querySelector('.main-content');

            if (window.innerWidth > 992) {
                if (mainContent) {
                    mainContent.style.transform = `translateX(${scrollTop * 0.09}px) skewX(-5deg)`;
                }
            } else {
                if (mainContent) {
                    mainContent.style.transform = 'none';
                }
            }
        };

        const handleResize = () => {
            const mainContent = document.querySelector('.main-content');
            if (window.innerWidth <= 992 && mainContent) {
                mainContent.style.transform = 'none';
            }
        };

        window.addEventListener('scroll', handleScroll);
        window.addEventListener('resize', handleResize);

        return () => {
            window.removeEventListener('scroll', handleScroll);
            window.removeEventListener('resize', handleResize);
        };
    }, []);

    return (
        <div className={`main-content flex-grow-1 p-4 ${animateContent ? 'animate' : ''}`}>
            {children}
        </div>
    );
}

function App() {
    const [loggedIn, setLoggedIn] = useState(isLoggedIn());
    const [currentUser, setCurrentUser] = useState(null);
    const [users, setUsers] = useState([]);
    const [showModal, setShowModal] = useState(false);
    const [isRegistering, setIsRegistering] = useState(false);
    const [page, setPage] = useState(1);
    const limit = 10;

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

    const fetchUsers = useCallback(async () => {
        try {
            const query = new URLSearchParams();
            query.append("page", page);
            query.append("limit", limit);
    
            const response = await fetch(`http://localhost:8000/api/users?${query.toString()}`);
            if (!response.ok) throw new Error("Error fetching users");
    
            const data = await response.json();
            setUsers(data);
        } catch (error) {
            console.error('Error fetching users:', error);
        }
    }, [page, limit]);
    
    const searchUsers = useCallback(async (searchCriteria) => {
        try {
            const query = new URLSearchParams();
            if (searchCriteria.term) query.append("term", searchCriteria.term);
            query.append("page", page);
            query.append("limit", limit);
    
            const response = await fetch(`http://localhost:8000/api/users?${query.toString()}`);
            if (!response.ok) throw new Error("Error fetching search results");
    
            const data = await response.json();
            setUsers(data);
        } catch (error) {
            console.error("Error fetching search results:", error);
        }
    }, [page, limit]);
    
    const handleLogout = () => {
        clearToken();
        setLoggedIn(false);
        setCurrentUser(null);
    };

    const handleShowLogin = (registerMode = false) => {
        setIsRegistering(registerMode);
        setShowModal(true);
    };    

    const handleCloseLogin = () => {
        setShowModal(false);
    };

    const handleNextPage = () => {
        setPage((prevPage) => prevPage + 1);
    };

    const handlePreviousPage = () => {
        setPage((prevPage) => Math.max(prevPage - 1, 1));
    };

    return (
        <Router>
            <div className="App d-lg-flex">
            <div className="background-animation" style={{backgroundImage: `url(${bgElement})`,}}></div>
                {/* Main App Content */}
                <Navbar 
                    loggedIn={loggedIn} 
                    currentUser={currentUser} 
                    onLogout={handleLogout} 
                    onShowLogin={handleShowLogin} 
                />
                <MainContentWrapper>
                    <Routes>
                        <Route path="/" element={<Home />} />
                        <Route path="/about" element={<About />} />
                        <Route path="/dashboard" element={loggedIn ? <UserDashboard /> : <Navigate to="/" />} />
                        <Route 
                            path="/allusers" 
                            element={<UserList 
                                        users={users} 
                                        page={page} 
                                        fetchUsers={fetchUsers} 
                                        searchUsers={searchUsers} 
                                        handleNextPage={handleNextPage} 
                                        handlePreviousPage={handlePreviousPage} 
                                    />} 
                        />
                        <Route path="*" element={<Navigate to="/" />} />
                    </Routes>
                </MainContentWrapper>
                {showModal && (
                    <AuthPanel 
                        onLogin={() => setLoggedIn(true)} 
                        onClose={handleCloseLogin} 
                        isRegistering={isRegistering}
                    />
                )}
                <Footer />
            </div>
        </Router>

    );
}

export default App;
