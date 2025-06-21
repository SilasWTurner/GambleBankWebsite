import React, { useContext } from 'react';
import { Link } from 'react-router-dom';
import { AuthContext } from './AuthContext';
import './NavBar.css';

function NavBar() {
  const { token, logout } = useContext(AuthContext);
  const username = token ? JSON.parse(atob(token.split('.')[1])).username : '';

  return (
    <div className="navbar">
      <Link to="/">Home</Link>
      {!token ? (
        <>
          <Link to="/login">Login</Link>
          <Link to="/signup">Sign Up</Link>
        </>
      ) : (
        <>
          <Link to="/lobby">Lobby</Link>
          <div className="navbar-bottom">
            <span className="username">{username}</span>
            <button onClick={logout}>Logout</button>
          </div>
        </>
      )}
    </div>
  );
}

export default NavBar;