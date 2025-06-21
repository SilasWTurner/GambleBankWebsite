import React, { createContext, useState } from 'react';

export const InviteContext = createContext();

export const InviteProvider = ({ children }) => {
  const [newInvite, setNewInvite] = useState(null);

  return (
    <InviteContext.Provider value={{ newInvite, setNewInvite }}>
      {children}
    </InviteContext.Provider>
  );
};