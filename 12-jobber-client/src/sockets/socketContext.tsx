import { createContext, type FC, type ReactElement, type ReactNode, useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { IMessageDocument } from "src/features/chat/interfaces/chat.interface";
import type { IOrderDocument, IOrderNotifcation } from "src/features/order/interfaces/order.interface";

type ContextProviderProps = {
  children: ReactNode;
};

export interface SocketMessageInterface {
  type: string;
  payload: any; // eslint-disable-line
}

enum PayloadType {
  EventMessageReceived = "message_received",
  EventMessageUpdated = "message_updated",
  GetLoggedInUsers = "getLoggedInUsers",
  LoggedInUsers = "loggedInUsers",
  RemoveLoggedInUser = "removeLoggedInUser",
  Category = "category",
  OrderNotification = "order_notification",
  OrderUpdate = "order_update",
  UpdateUsername = "update_username",
  Online = "online",
}

interface SocketInterface {
  isReady: boolean;
  online: string[];
  message: IMessageDocument | undefined;
  notification: IOrderNotifcation | undefined;
  orderDocument: IOrderDocument | undefined;
  msgUpdated: IMessageDocument | undefined;
  getLoggedInUsers: () => void;
  setLoggedInUser: (username: string) => void;
  setUsername: (username: string) => void;
  setCategory: (username: string, category: string) => void;
  removeUser: (username: string) => void;
}

export const WebSocketContext = createContext<SocketInterface | null>(null);

export const WebSocketProvider: FC<ContextProviderProps> = ({ children }): ReactElement => {
  const [isReady, setIsReady] = useState(false);
  const [online, setOnline] = useState<string[]>([]);
  const [message, setMessage] = useState<IMessageDocument>();
  const [notification, setNotification] = useState<IOrderNotifcation>();
  const [msgUpdated, setMsgUpdated] = useState<IMessageDocument>();
  const [orderDocument, setOrderDocument] = useState<IOrderDocument>();

  const { sendMessage, readyState, lastMessage } = useWebSocket("ws://localhost:4000/ws", {
    onOpen: () => console.log("opened"),
    shouldReconnect: () => true,
    share: true,
  });

  useEffect(() => {
    const connectionStatus = {
      [ReadyState.CONNECTING]: "Connecting",
      [ReadyState.OPEN]: "Open",
      [ReadyState.CLOSING]: "Closing",
      [ReadyState.CLOSED]: "Closed",
      [ReadyState.UNINSTANTIATED]: "Uninstantiated",
    }[readyState];

    if (connectionStatus == "Open") {
      setIsReady(true);
    } else {
      setIsReady(false);
    }
  }, [readyState]);

  useEffect(() => {
    if (lastMessage !== null) {
      RouteEvent(lastMessage?.data);
    }
  }, [lastMessage]);

  const RouteEvent = (event: string): void => {
    try {
      const e = JSON.parse(event) as SocketMessageInterface;

      switch (e.type) {
        case PayloadType.Online:
          setOnline(e.payload as string[]);
          break;

        case PayloadType.EventMessageReceived:
          setMessage(e.payload as IMessageDocument);
          break;

        case PayloadType.EventMessageUpdated:
          setMsgUpdated(e.payload as IMessageDocument);
          break;

        case PayloadType.OrderNotification:
          setNotification(e.payload as IOrderNotifcation);
          break;

        case PayloadType.OrderUpdate:
          setOrderDocument(e.payload as IOrderDocument);
          break;

        default:
          break;
      }
    } catch (error) {
      console.log(error);
    }
  };

  const getLoggedInUsers = () => {
    sendMessage(JSON.stringify({ type: PayloadType.GetLoggedInUsers }));
  };

  const setLoggedInUser = (username: string): void => {
    sendMessage(JSON.stringify({ type: PayloadType.LoggedInUsers, payload: username }));
  };

  const setUsername = (username: string): void => {
    sendMessage(JSON.stringify({ type: PayloadType.UpdateUsername, payload: username }));
  };

  const setCategory = (username: string, category: string): void => {
    sendMessage(JSON.stringify({ type: PayloadType.Category, payload: { username: username, category: category } }));
  };

  const removeUser = (username: string): void => {
    sendMessage(JSON.stringify({ type: PayloadType.RemoveLoggedInUser, payload: username }));
  };

  return (
    <WebSocketContext.Provider
      value={{
        isReady,
        online,
        message,
        notification,
        orderDocument,
        msgUpdated,
        getLoggedInUsers,
        setLoggedInUser,
        setUsername,
        setCategory,
        removeUser,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};
