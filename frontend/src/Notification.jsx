import React, { useContext } from 'react';
import { InviteContext } from './InviteContext';
import { AuthContext } from './AuthContext';
import axios from 'axios';
import { FaTimes } from 'react-icons/fa';
import './Notification.css';
import WebSocketInstance from './WebSocketService';

const Notification = () => {
  const { newInvite, setNewInvite } = useContext(InviteContext);
  const { token } = useContext(AuthContext);

  const handleAcceptInvite = (inviteId) => {
    WebSocketInstance.sendMessage({
      action: 'accept_invite',
      args: [inviteId],
    });
    setNewInvite(null);
  };

  const handleRejectInvite = async (inviteId) => {
    try {
      const response = await axios.post(
        'http://localhost:8080/reject_invite',
        { invite_id: inviteId },
        { headers: { Authorization: token } }
      );

      if (response.status === 200) {
        setNewInvite(null);
      } else {
        console.error('Failed to reject invite');
      }
    } catch (error) {
      console.error('Failed to reject invite', error);
    }
  };

  if (!newInvite) return null;

  return (
    <div className="new-invite-box">
      <p>New invite from {newInvite.senderName}</p>
      <button onClick={() => handleAcceptInvite(newInvite.id)}>Accept</button>
      <button onClick={() => handleRejectInvite(newInvite.id)}>Decline</button>
    </div>
  );
};

export default Notification;