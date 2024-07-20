import { cloneDeep } from "lodash";
import { type FC, type ReactElement, useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { useSocket } from "src/sockets/socketHook";

import type { IMessageDocument } from "../interfaces/chat.interface";
import { useGetUserMessagesQuery } from "../services/chat.service";
// import { chatMessageReceived } from "../services/chat.utils";
import ChatList from "./ChatList/ChatList";
import ChatWindow from "./ChatWindow/ChatWindow";

const Chat: FC = (): ReactElement => {
  const { conversationId } = useParams<string>();
  const { message: socketMessage } = useSocket();

  const [skip, setSkip] = useState<boolean>(false);
  const [chatMessagesData, setChatMessagesData] = useState<IMessageDocument[]>([]);

  const chatMessages = useRef<IMessageDocument[]>([]);

  const { data, isSuccess, isLoading, isError } = useGetUserMessagesQuery(`${conversationId}`, { skip, refetchOnMountOrArgChange: true });

  useEffect(() => {
    console.log(conversationId);
  }, [conversationId]);

  useEffect(() => {
    if (isSuccess) {
      setChatMessagesData(data?.messages as IMessageDocument[]);
    }
  }, [isSuccess, data?.messages]);

  // useEffect(() => {
  // chatMessageReceived(`${conversationId}`, chatMessagesData, chatMessages.current, setChatMessagesData);
  // }, [chatMessagesData, conversationId]);

  useEffect(() => {
    if (socketMessage === undefined) {
      return;
    }

    if (socketMessage?.conversationId === conversationId) {
      chatMessages.current = cloneDeep(chatMessagesData);

      chatMessages.current.push(socketMessage!);

      // remove duplicates from chat messages
      const uniq = chatMessages.current.filter((item: IMessageDocument, index: number, list: IMessageDocument[]) => {
        const itemIndex = list.findIndex((listItem: IMessageDocument) => listItem._id === item._id);
        return itemIndex === index;
      });

      setChatMessagesData(uniq);
    }
  }, [socketMessage]); // eslint-disable-line

  return (
    <div className="border-grey mx-2 my-5 flex max-h-[90%] flex-wrap border lg:container lg:mx-auto">
      <div className="lg:border-grey relative w-full overflow-hidden lg:w-1/3 lg:border-r">
        <ChatList />
      </div>

      <div className="relative hidden w-full overflow-hidden md:w-2/3 lg:flex">
        {conversationId && chatMessagesData.length > 0 ? (
          <ChatWindow setSkip={setSkip} chatMessages={chatMessagesData} isLoading={isLoading} isError={isError} />
        ) : (
          <div className="flex w-full items-center justify-center">Select a user to chat with.</div>
        )}
      </div>
    </div>
  );
};

export default Chat;
