import React, { useEffect, useState } from 'react';
import { isLoggedIn, clearToken, getUserEmail } from './services/authService';
import AuthPanel from './components/AuthPanel';
import Logout from './components/Logout';
import UserDashboard from './components/UserDashboard';

function App() {
    const [loggedIn, setLoggedIn] = useState(isLoggedIn());
    const [currentUser, setCurrentUser] = useState(null);

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

    return (
        <div className="App">
            <h1>DocuDefense Frontend</h1>
            {loggedIn ? (
                <>
                    <Logout onLogout={handleLogout} />
                    <UserDashboard currentUser={currentUser} />
                </>
            ) : (
                <AuthPanel onLogin={() => setLoggedIn(true)} />
            )}
        </div>
    );
}

export default App;
