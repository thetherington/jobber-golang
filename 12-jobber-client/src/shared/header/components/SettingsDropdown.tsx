import type { FC, ReactElement } from "react";
import { Link, type NavigateFunction, useNavigate } from "react-router-dom";
import { applicationLogout, lowerCase } from "src/shared/utils/utils.service";
// import { socket } from "src/sockets/socket.service";
import { useAppDispatch } from "src/store/store";

import type { IHomeHeaderProps } from "../interfaces/header.interface";
import { updateCategoryContainer } from "../reducers/category.reducer";
import { updateHeader } from "../reducers/header.reducer";

const SettingsDropdown: FC<IHomeHeaderProps> = ({ seller, authUser, buyer, type, setIsDropdownOpen }): ReactElement => {
  const navigate: NavigateFunction = useNavigate();
  const dispatch = useAppDispatch();

  const onLogout = (): void => {
    if (setIsDropdownOpen) {
      setIsDropdownOpen(false);
    }

    applicationLogout(dispatch, navigate);

    // force components to check if user is online
    // socket.emit("getLoggedInUsers");
  };

  return (
    <div className="border-grey w-44 divide-y divide-gray-100 rounded border bg-white shadow-md">
      <ul className="text-gray-700s py-2 text-sm" aria-labelledby="avatarButton">
        {buyer && buyer.isSeller && (
          <li className="mx-3 mb-1">
            <Link
              to={`${type === "buyer" ? `/${lowerCase(`${authUser?.username}`)}/${seller?._id}/seller_dashboard` : "/"}`}
              onClick={() => {
                if (setIsDropdownOpen) {
                  setIsDropdownOpen(false);
                }
                dispatch(updateHeader(`${type === "buyer" ? "sellerDashboard" : "home"}`));
                dispatch(updateCategoryContainer(true));
              }}
              className="block w-full cursor-pointer rounded bg-sky-500 px-4s py-2 text-center font-bold text-white hover:bg-sky-400 focus:outline-none"
            >
              {type === "buyer" ? "Switch to Selling" : "Switch to Buying"}
            </Link>
          </li>
        )}
        {buyer && buyer.isSeller && type === "buyer" && (
          <li>
            <Link
              to={`/manage_gigs/new/${seller?._id}`}
              onClick={() => {
                if (setIsDropdownOpen) {
                  setIsDropdownOpen(false);

                  dispatch(updateHeader("home"));
                  dispatch(updateCategoryContainer(true));
                }
              }}
              className="block px-4 py-2 hover:text-sky-400"
            >
              Add a new gig
            </Link>
          </li>
        )}
        {type === "buyer" && (
          <li>
            <Link
              to={`/users/${buyer?.username}/${buyer?._id}/orders`}
              className="block px-4 py-2 hover:text-sky-400"
              onClick={() => {
                if (setIsDropdownOpen) {
                  setIsDropdownOpen(false);
                }
                dispatch(updateHeader("home"));
                dispatch(updateCategoryContainer(true));
              }}
            >
              Dashboard
            </Link>
          </li>
        )}
        {buyer && buyer.isSeller && type === "buyer" && (
          <li>
            <Link
              to={`/seller_profile/${lowerCase(`${seller?.username}`)}/${seller?._id}/edit`}
              className="block px-4 py-2 hover:text-sky-400"
              onClick={() => {
                if (setIsDropdownOpen) {
                  setIsDropdownOpen(false);
                }
                dispatch(updateHeader("home"));
                dispatch(updateCategoryContainer(true));
              }}
            >
              Profile
            </Link>
          </li>
        )}
        <li>
          <Link
            to={`/${lowerCase(`${buyer?.username}/edit`)}`}
            className="block px-4 py-2 hover:text-sky-400"
            onClick={() => {
              if (setIsDropdownOpen) {
                setIsDropdownOpen(false);
              }
              dispatch(updateHeader("home"));
              dispatch(updateCategoryContainer(false));
            }}
          >
            Settings
          </Link>
        </li>
      </ul>
      <div className="py-1">
        <div onClick={() => onLogout()} className="block px-4 py-2 text-sm hover:text-sky-400">
          Sign out
        </div>
      </div>
    </div>
  );
};

export default SettingsDropdown;
