import { useContext } from "react";

import { WebSocketContext } from "./socketContext";

export const useSocket = () => {
  const context = useContext(WebSocketContext);
  if (context === undefined || context === null) {
    throw new Error("Cannot use WebSocketContext without using WebSocket Provider");
  }

  return context;
};
