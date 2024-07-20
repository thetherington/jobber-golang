import { useEffect } from "react";
// import { useSocket } from "src/sockets/socketHook";

// import { socket } from "src/sockets/socket.service";
// import { getDataFromSessionStorage } from "../utils/utils.service";

const useBeforeWindowUnload = (): void => {
  // const { removeUser } = useSocket();

  useEffect(() => {
    // If the user closes the browser or tab, we emit the socketio event
    window.addEventListener("beforeunload", () => {
      // const loggedInUsername: string = getDataFromSessionStorage("loggedInuser");
      // removeUser(loggedInUsername);
      // socket.emit("removeLoggedInUser", loggedInUsername);
    });
  }, []); // eslint-disable-line
};

export default useBeforeWindowUnload;
