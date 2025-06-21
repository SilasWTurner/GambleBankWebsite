import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import NavBar from './NavBar';
import LandingPage from './LandingPage';
import SignUpPage from './SignUpPage';
import LoginPage from './LoginPage';
import LobbyPage from './LobbyPage';
import { AuthProvider } from './AuthContext';
import { InviteProvider } from './InviteContext';
import Notification from './Notification';
import './App.css';

function App() {
  return (
    <InviteProvider>
      <AuthProvider>
        <Router>
          <div className="app">
            <NavBar />
            <Notification />
            <Routes>
              <Route path="/" element={<LandingPage />} />
              <Route path="/signup" element={<SignUpPage />} />
              <Route path="/login" element={<LoginPage />} />
              <Route path="/lobby" element={<LobbyPage />} />
            </Routes>
          </div>
        </Router>
      </AuthProvider>
    </InviteProvider>
  );
}

export default App;
