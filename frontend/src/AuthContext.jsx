import React, { createContext, useState, useEffect, useContext } from 'react';
import WebSocketInstance from './WebSocketService';
import { InviteContext } from './InviteContext';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [websocketUrl, setWebsocketUrl] = useState('');
  const { setNewInvite } = useContext(InviteContext);

  useEffect(() => {
    if (token) {
      const decodedToken = JSON.parse(atob(token.split('.')[1]));
      const isTokenExpired = decodedToken.exp * 1000 < Date.now();

      if (isTokenExpired) {
        localStorage.removeItem('token');
        setToken(null);
        return;
      }

      const username = decodedToken.username;
      const wsUrl = `ws://localhost:8080/ws?username=${username}&token=${token}`;
      setWebsocketUrl(wsUrl);
      WebSocketInstance.connect(wsUrl);
      WebSocketInstance.addCallbacks({
        new_invite: handleNewInvite,
        accept_invite: handleAcceptInvite,
        game_action: handleGameAction,
      });
    }
  }, [token]);

  const handleNewInvite = (args) => {
    console.log('New invite received:', args);
    const [inviteId, senderName] = args;
    setNewInvite({ id: inviteId, senderName });
  };

  const handleAcceptInvite = (args) => {
    console.log('Invite accepted:', args);
    // Handle invite acceptance
  };

  const handleGameAction = (args) => {
    console.log('Game action received:', args);
    // Handle game-related actions
  };

  const login = (newToken) => {
    localStorage.setItem('token', newToken);
    setToken(newToken);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    WebSocketInstance.socketRef.close();
  };

  return (
    <AuthContext.Provider value={{ token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};