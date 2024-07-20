import { cloneDeep, filter, findIndex, lowerCase, orderBy, remove } from "lodash";
import { type FC, type ReactElement, useEffect, useRef, useState } from "react";
import { FaCheck, FaCheckDouble, FaCircle } from "react-icons/fa";
import { LazyLoadImage } from "react-lazy-load-image-component";
import { type Location, type NavigateFunction, useLocation, useNavigate, useParams } from "react-router-dom";
import { updateNotification } from "src/shared/header/reducers/notification.reducer";
import { TimeAgo } from "src/shared/utils/timeago.utils";
import { isFetchBaseQueryError, showErrorToast } from "src/shared/utils/utils.service";
// import { socket } from "src/sockets/socket.service";
import { useSocket } from "src/sockets/socketHook";
import { useAppDispatch, useAppSelector } from "src/store/store";
import { type IReduxState } from "src/store/store.interface";
import { v4 as uuidv4 } from "uuid";

import type { IMessageDocument } from "../../interfaces/chat.interface";
import { useGetConversationListQuery, useMarkMultipleMessagesAsReadMutation } from "../../services/chat.service";

const ChatList: FC = (): ReactElement => {
  const navigate: NavigateFunction = useNavigate();
  const location: Location = useLocation();

  const dispatch = useAppDispatch();

  const { message: socketMessage, getLoggedInUsers, msgUpdated: socketUpdatedMessage } = useSocket();

  const authUser = useAppSelector((state: IReduxState) => state.authUser);
  const { username, conversationId } = useParams<string>();

  const [selectedUser, setSelectedUser] = useState<IMessageDocument>();
  const [chatList, setChatList] = useState<IMessageDocument[]>([]);
  const conversationListRef = useRef<IMessageDocument[]>([]);

  const { data, isSuccess } = useGetConversationListQuery(`${authUser.username}`);
  const [markMultipleMessagesAsRead] = useMarkMultipleMessagesAsReadMutation();

  const selectUserFromList = async (user: IMessageDocument): Promise<void> => {
    try {
      setSelectedUser(user);

      const pathList: string[] = location.pathname.split("/");
      pathList.splice(-2, 2);

      const locationPathname: string = !pathList.join("/") ? location.pathname : pathList.join("/");
      const chatUsername: string = (user.receiverUsername !== authUser?.username ? user.receiverUsername : user.senderUsername) as string;

      navigate(`${locationPathname}/${lowerCase(chatUsername)}/${user.conversationId}`);

      getLoggedInUsers();

      if (user.receiverUsername === authUser?.username && lowerCase(`${user.senderUsername}`) === username && !user.isRead) {
        const list: IMessageDocument[] = filter(
          chatList,
          (item: IMessageDocument) => !item.isRead && item.receiverUsername === authUser.username,
        );

        if (list.length > 0) {
          await markMultipleMessagesAsRead({
            receiverUsername: `${user.receiverUsername}`,
            senderUsername: `${user.senderUsername}`,
            messageId: `${user._id}`,
          });
        }
      }
    } catch (error) {
      if (isFetchBaseQueryError(error)) {
        showErrorToast(error?.data?.message);
      } else {
        showErrorToast("Error selecting chat user.");
      }
    }
  };

  useEffect(() => {
    if (isSuccess) {
      const sortedConverstations: IMessageDocument[] = orderBy(data.conversations, ["createdAt"], ["desc"]) as IMessageDocument[];
      setChatList(sortedConverstations);

      if (!sortedConverstations.length) {
        dispatch(updateNotification({ hasUncreadMessage: false }));
      }
    }
  }, [isSuccess, data?.conversations, dispatch]);

  // useEffect(() => {
  // chatListMessageReceived(`${authUser.username}`, chatList, conversationListRef.current, dispatch, setChatList);
  // chatListMessageUpdated(`${authUser.username}`, chatList, conversationListRef.current, dispatch, setChatList);
  // }, [authUser.username, conversationId, chatList, dispatch]);

  // Effect for when a chat message is updated
  useEffect(() => {
    if (socketUpdatedMessage === undefined) {
      return;
    }

    conversationListRef.current = cloneDeep(chatList);

    if (
      lowerCase(`${socketUpdatedMessage!.receiverUsername}`) === lowerCase(`${authUser.username}`) ||
      lowerCase(`${socketUpdatedMessage!.senderUsername}`) === lowerCase(`${authUser.username}}`)
    ) {
      const messageIndex = findIndex(chatList, ["conversationId", socketUpdatedMessage.conversationId]);

      if (messageIndex > -1) {
        conversationListRef.current.splice(messageIndex, 1, socketUpdatedMessage);
      }

      if (lowerCase(`${socketUpdatedMessage.receiverUsername}`) === lowerCase(`${username}`)) {
        const list: IMessageDocument[] = filter(
          conversationListRef.current,
          (item: IMessageDocument) => !item.isRead && item.receiverUsername === username,
        );

        // console.log(list);
        dispatch(updateNotification({ hasUnreadMessage: list.length > 0 }));
      }

      setChatList(conversationListRef.current);
    }
  }, [socketUpdatedMessage]); // eslint-disable-line

  // Effect for when a chat message is received
  useEffect(() => {
    if (socketMessage === undefined) {
      return;
    }

    conversationListRef.current = cloneDeep(chatList);

    if (
      lowerCase(`${socketMessage!.receiverUsername}`) === lowerCase(`${authUser.username}`) ||
      lowerCase(`${socketMessage!.senderUsername}`) === lowerCase(`${authUser.username}}`)
    ) {
      const messageIndex = findIndex(chatList, ["conversationId", socketMessage!.conversationId]);

      if (messageIndex > -1) {
        remove(conversationListRef.current, (chat: IMessageDocument) => chat.conversationId === socketMessage.conversationId);
      } else {
        remove(conversationListRef.current, (chat: IMessageDocument) => chat.receiverUsername === socketMessage.receiverUsername);
      }

      conversationListRef.current = [socketMessage, ...conversationListRef.current];

      if (lowerCase(`${socketMessage!.receiverUsername}`) === lowerCase(`${authUser.username}`)) {
        const list: IMessageDocument[] = filter(
          conversationListRef.current,
          (item: IMessageDocument) => !item.isRead && item.receiverUsername === username,
        );

        // console.log(list);
        dispatch(updateNotification({ hasUnreadMessage: list.length > 0 }));
      }

      setChatList(conversationListRef.current);
    }
  }, [socketMessage]); // eslint-disable-line

  return (
    <>
      <div className="border-grey truncate border-b px-5 py-3 text-base font-medium">
        <h2 className="w-6/12 truncate text-sm md:text-base lg:text-lg">All Conversations</h2>
      </div>
      <div className="absolute h-full w-full overflow-scroll pb-14">
        {chatList.map((data: IMessageDocument, index: number) => (
          <div
            key={uuidv4()}
            onClick={() => selectUserFromList(data)}
            className={`flex w-full cursor-pointer items-center space-x-4 px-5 py-4 hover:bg-gray-50 ${index !== chatList.length - 1 ? "border-grey border-b" : ""} ${!data.isRead} ? 'bg-[#f5fbff]':''} ${data.conversationId === conversationId ? "bg-[#f5fbff]" : ""}`}
          >
            <LazyLoadImage
              src={data.receiverUsername !== authUser?.username ? data.receiverPicture : data.senderPicture}
              alt="profile image"
              className="h-10 w-10 object-cover rounded-full"
              placeholderSrc="https://placehold.co/330x220?text=Profile+Image"
              // effect="blur"
            />
            <div className="w-full text-sm dark:text-white">
              <div className="flex justify-between pb-1 font-bold text-[#777d74]">
                <span className={`${selectedUser && !data.body ? "flex items-center" : ""}`}>
                  {data.receiverUsername !== authUser?.username ? data.receiverUsername : data.senderUsername}
                </span>
                {data.createdAt && <span className="font-normal">{TimeAgo.transform(`${data.createdAt}`)}</span>}
              </div>
              <div className="flex justify-between text-[#777d74]">
                <span>
                  {data.receiverUsername === authUser.username ? "" : "Me: "}
                  {data.body}
                </span>
                {!data.isRead ? (
                  <>
                    {data.receiverUsername === authUser.username ? (
                      <FaCircle className="mt-2 text-sky-500" size={8} />
                    ) : (
                      <FaCheck className="mt-2" size={8} />
                    )}
                  </>
                ) : (
                  <FaCheckDouble className="mt-2 text-sky-500" size={8} />
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </>
  );
};

export default ChatList;
