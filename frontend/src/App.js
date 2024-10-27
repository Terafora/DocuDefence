// src/App.js
import React, { useState } from 'react';
import { setToken, clearToken, isLoggedIn } from './services/authService';
import AuthPanel from './components/AuthPanel';
import UserDashboard from './components/UserDashboard';

function App() {
    const [loggedIn, setLoggedIn] = useState(isLoggedIn());

    const handleLogin = (token) => {
        setToken(token);
        setLoggedIn(true);
    };

    const handleLogout = () => {
        clearToken();
        setLoggedIn(false);
    };

    return (
        <div className="App">
            <h1>DocuDefense Frontend</h1>
            {loggedIn ? (
                <UserDashboard onLogout={handleLogout} />
            ) : (
                <AuthPanel onLogin={handleLogin} />
            )}
        </div>
    );
}

export default App;
