import { type FC, lazy, type LazyExoticComponent, type ReactElement, Suspense, useCallback, useEffect, useState } from "react";
import { type NavigateFunction, useNavigate } from "react-router-dom";
import HomeHeader from "src/shared/header/components/HomeHeader";
import CircularPageLoader from "src/shared/page-loader/CircularPageLoader";
import { applicationLogout, getDataFromLocalStorage, saveToSessionStorage } from "src/shared/utils/utils.service";
// import { socket } from "src/sockets/socket.service";
import { WebSocketProvider } from "src/sockets/socketContext";
import { useAppDispatch, useAppSelector } from "src/store/store";
import { type IReduxState } from "src/store/store.interface";

import { addAuthUser } from "./auth/reducers/auth.reducer";
import { useCheckCurrentUserQuery } from "./auth/services/auth.service";
import { addBuyer } from "./buyer/reducers/buyer.reducer";
import { useGetCurrentBuyerByUsernameQuery } from "./buyer/services/buyer.service";
import { addSeller } from "./sellers/reducers/seller.reducer";
import { useGetSellerByUsernameQuery } from "./sellers/services/seller.service";

const Index: LazyExoticComponent<FC> = lazy(() => import("./index/Index"));
const Home: LazyExoticComponent<FC> = lazy(() => import("./home/components/Home"));

const AppPage: FC = (): ReactElement => {
  const navigate: NavigateFunction = useNavigate();

  const dispatch = useAppDispatch();

  const authUser = useAppSelector((state: IReduxState) => state.authUser);
  const appLogout = useAppSelector((state: IReduxState) => state.logout);
  const showCategoryContainer = useAppSelector((state: IReduxState) => state.showCategoryContainer);

  const [tokenIsValid, setTokenIsValid] = useState<boolean>(false);

  const { data: currentUserData, isError } = useCheckCurrentUserQuery(undefined, { skip: authUser.id === null });
  const { data: buyerData, isLoading: isBuyerLoading } = useGetCurrentBuyerByUsernameQuery(undefined, { skip: authUser.id === null });
  const { data: sellerData, isLoading: isSellerLoading } = useGetSellerByUsernameQuery(`${authUser.username}`, {
    skip: authUser.id === null,
  });

  const checkUser = useCallback(async () => {
    try {
      if (currentUserData && currentUserData.user && !appLogout) {
        setTokenIsValid(true);

        // dispatch user info into state
        dispatch(addAuthUser({ authInfo: currentUserData.user }));

        // dispatch buyer info into state
        dispatch(addBuyer(buyerData?.buyer));

        // dispatch seller info into state
        dispatch(addSeller(sellerData?.seller));

        saveToSessionStorage(JSON.stringify(true), JSON.stringify(authUser.username));

        const becomeASeller = getDataFromLocalStorage("becomeASeller");
        if (becomeASeller) {
          navigate("/seller_onboarding");
        }

        if (authUser.username !== null) {
          // socket.emit("loggedInUsers", authUser.username);
        }
      }
    } catch (error) {
      console.log(error);
    }
  }, [currentUserData, dispatch, appLogout, authUser.username, buyerData, sellerData, navigate]);

  const logoutUser = useCallback(async () => {
    if ((!currentUserData && appLogout) || isError) {
      setTokenIsValid(false);

      applicationLogout(dispatch, navigate);
    }
  }, [currentUserData, dispatch, appLogout, isError, navigate]);

  useEffect(() => {
    checkUser();
    logoutUser();
  }, [checkUser, logoutUser]);

  if (authUser) {
    return !tokenIsValid && !authUser.id ? (
      <Suspense>
        <Index />
      </Suspense>
    ) : (
      <>
        {isBuyerLoading && isSellerLoading ? (
          <CircularPageLoader />
        ) : (
          <Suspense>
            <WebSocketProvider>
              <HomeHeader showCategoryContainer={showCategoryContainer} />
              <Home />
            </WebSocketProvider>
          </Suspense>
        )}
      </>
    );
  } else {
    return (
      <Suspense>
        <Index />
      </Suspense>
    );
  }
};

export default AppPage;
