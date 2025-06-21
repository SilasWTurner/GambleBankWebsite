import React, { useState, useEffect, useContext } from 'react';
import axios from 'axios';
import { AuthContext } from './AuthContext';
import { FaUserPlus, FaTimes } from 'react-icons/fa';
import './LobbyPage.css';
import WebSocketInstance from './WebSocketService';

function LobbyPage() {
  const { token } = useContext(AuthContext);
  const [invites, setInvites] = useState([]);
  const [username, setUsername] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    const fetchInvites = async () => {
      try {
        const response = await axios.get('http://localhost:8080/list_invites', {
          headers: { Authorization: token },
        });
        setInvites(response.data || []);
      } catch (error) {
        console.error('Error fetching invites:', error);
      }
    };

    fetchInvites();
  }, [token]);

  const handleSendInvite = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post(
        'http://localhost:8080/send_invite',
        { receiver_username: username },
        { headers: { Authorization: token } }
      );

      if (response.status === 200) {
        setSuccess('Invite sent successfully!');
        setError('');
        setUsername('');
      } else {
        setError(response.data.error || 'Failed to send invite');
        setSuccess('');
      }
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to send invite');
      setSuccess('');
    }
  };

  const handleRejectInvite = async (inviteId) => {
    try {
      const response = await axios.post(
        'http://localhost:8080/reject_invite',
        { invite_id: inviteId },
        { headers: { Authorization: token } }
      );

      if (response.status === 200) {
        setInvites(invites.filter((invite) => invite.id !== inviteId));
      } else {
        setError(response.data.error || 'Failed to reject invite');
      }
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to reject invite');
    }
  };

  const handleAcceptInvite = (inviteId) => {
    WebSocketInstance.sendMessage({
      action: 'accept_invite',
      args: [inviteId],
    });
  };

  return (
    <div className="lobby-page">
      <h2>Lobby</h2>
      <div className="lobby-container">
        <div className="invites-section">
          <h3>Invites</h3>
          {invites.length === 0 ? (
            <p className="no-invites">No invites</p>
          ) : (
            <ul>
              {invites.map((invite) => (
                <li key={invite.id} className="invite-item">
                  <span>{invite.sender_name}</span>
                  <button onClick={() => handleAcceptInvite(invite.id)}>Accept</button>
                  <FaTimes
                    className="icon"
                    title="Dismiss"
                    onClick={() => handleRejectInvite(invite.id)}
                  />
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="send-invite-section">
          <h3>Send Invite</h3>
          <form onSubmit={handleSendInvite} className="invite-form">
            <div className="form-group">
              <label htmlFor="username">Username</label>
              <input
                type="text"
                id="username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
            </div>
            {error && <p className="error">{error}</p>}
            {success && <p className="success">{success}</p>}
            <button type="submit" className="invite-button">
              <FaUserPlus className="icon" /> Send Invite
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

export default LobbyPage;
