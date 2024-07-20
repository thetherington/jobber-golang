import { type FC, type ReactElement, type ReactNode, useCallback, useEffect, useState } from "react";
import { Navigate, NavigateFunction, useNavigate } from "react-router-dom";
import HomeHeader from "src/shared/header/components/HomeHeader";
import { applicationLogout, saveToSessionStorage } from "src/shared/utils/utils.service";
import { WebSocketProvider } from "src/sockets/socketContext";
import { useAppDispatch, useAppSelector } from "src/store/store";
import type { IReduxState } from "src/store/store.interface";

import { addAuthUser } from "./auth/reducers/auth.reducer";
import { useCheckCurrentUserQuery } from "./auth/services/auth.service";

export interface IProtectedRouteProps {
  children: ReactNode;
}

const ProtectedRoute: FC<IProtectedRouteProps> = ({ children }): ReactElement => {
  const navigate: NavigateFunction = useNavigate();
  const [tokenIsValid, setTokenIsValid] = useState<boolean>(false);

  const dispatch = useAppDispatch();

  const authUser = useAppSelector((state: IReduxState) => state.authUser);
  const showCategoryContainer = useAppSelector((state: IReduxState) => state.showCategoryContainer);
  const header = useAppSelector((state: IReduxState) => state.header);

  const { data, isError } = useCheckCurrentUserQuery();

  const checkUser = useCallback(async () => {
    if (data && data.user) {
      setTokenIsValid(true);
      dispatch(addAuthUser({ authInfo: data.user }));

      saveToSessionStorage(JSON.stringify(true), JSON.stringify(authUser.username));
    }

    if (isError) {
      setTokenIsValid(false);
      applicationLogout(dispatch, navigate);
    }
  }, [data, dispatch, authUser.username, isError, navigate]);

  useEffect(() => {
    checkUser();
  }, [checkUser]);

  if ((data && data.user) || authUser) {
    if (tokenIsValid) {
      return (
        <>
          <WebSocketProvider>
            {header && header === "home" && <HomeHeader showCategoryContainer={showCategoryContainer} />}
            {children}
          </WebSocketProvider>
        </>
      );
    } else {
      return <></>;
    }
  } else {
    return <>{<Navigate to="/" />}</>;
  }
};

export default ProtectedRoute;
