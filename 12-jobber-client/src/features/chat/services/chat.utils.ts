// import { cloneDeep, filter, findIndex, lowerCase, remove } from "lodash";
// import { Dispatch, SetStateAction } from "react";
// import { updateNotification } from "src/shared/header/reducers/notification.reducer";
// import { socket } from "src/sockets/socket.service";
// import { type AppDispatch } from "src/store/store";

// import type { IMessageDocument } from "../interfaces/chat.interface";

// export const chatMessageReceived = (
//   conversationId: string,
//   chatMessagesData: IMessageDocument[],
//   chatMessages: IMessageDocument[],
//   setChatMessagesData: Dispatch<SetStateAction<IMessageDocument[]>>,
// ): void => {
// socket.on("message received", (data: IMessageDocument) => {
//   chatMessages = cloneDeep(chatMessagesData);
//   if (data.conversationId === conversationId) {
//     chatMessages.push(data);
//     // remove duplicates from chat messages
//     const uniq = chatMessages.filter((item: IMessageDocument, index: number, list: IMessageDocument[]) => {
//       const itemIndex = list.findIndex((listItem: IMessageDocument) => listItem._id === item._id);
//       return itemIndex === index;
//     });
//     setChatMessagesData(uniq);
//   }
// });
// };

// export const chatListMessageReceived = (
//   username: string,
//   chatList: IMessageDocument[],
//   conversationListRef: IMessageDocument[],
//   dispatch: AppDispatch,
//   setChatList: Dispatch<SetStateAction<IMessageDocument[]>>,
// ): void => {
// socket.on("message received", (data: IMessageDocument) => {
//   conversationListRef = cloneDeep(chatList);
//   if (
//     lowerCase(`${data.receiverUsername}`) === lowerCase(`${username}`) ||
//     lowerCase(`${data.senderUsername}`) === lowerCase(`${username}`)
//   ) {
//     const messageIndex = findIndex(chatList, ["conversationId", data.conversationId]);
//     if (messageIndex > -1) {
//       remove(conversationListRef, (chat: IMessageDocument) => chat.conversationId === data.conversationId);
//     } else {
//       remove(conversationListRef, (chat: IMessageDocument) => chat.receiverUsername === data.receiverUsername);
//     }
//     conversationListRef = [data, ...conversationListRef];
//     if (lowerCase(`${data.receiverUsername}`) === lowerCase(`${username}`)) {
//       const list: IMessageDocument[] = filter(
//         conversationListRef,
//         (item: IMessageDocument) => !item.isRead && item.receiverUsername === username,
//       );
//       // console.log(list);
//       dispatch(updateNotification({ hasUnreadMessage: list.length > 0 }));
//     }
//     setChatList(conversationListRef);
//   }
// });
// };

// export const chatListMessageUpdated = (
//   username: string,
//   chatList: IMessageDocument[],
//   conversationListRef: IMessageDocument[],
//   dispatch: AppDispatch,
//   setChatList: Dispatch<SetStateAction<IMessageDocument[]>>,
// ): void => {
// socket.on("message updated", (data: IMessageDocument) => {
//   conversationListRef = cloneDeep(chatList);
//   if (
//     lowerCase(`${data.receiverUsername}`) === lowerCase(`${username}`) ||
//     lowerCase(`${data.senderUsername}`) === lowerCase(`${username}`)
//   ) {
//     const messageIndex = findIndex(chatList, ["conversationId", data.conversationId]);
//     if (messageIndex > -1) {
//       conversationListRef.splice(messageIndex, 1, data);
//     }
//     if (lowerCase(`${data.receiverUsername}`) === lowerCase(`${username}`)) {
//       const list: IMessageDocument[] = filter(
//         conversationListRef,
//         (item: IMessageDocument) => !item.isRead && item.receiverUsername === username,
//       );
//       // console.log(list);
//       dispatch(updateNotification({ hasUnreadMessage: list.length > 0 }));
//     }
//     setChatList(conversationListRef);
//   }
// });
// };
