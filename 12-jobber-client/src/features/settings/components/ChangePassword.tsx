import { type ChangeEvent, type FC, type ReactElement, useState } from "react";
import { FaEye, FaEyeSlash } from "react-icons/fa";
import { type NavigateFunction, useNavigate } from "react-router-dom";
import Alert from "src/shared/alert/Alert";
import Button from "src/shared/button/Button";
import TextInput from "src/shared/input/TextInput";
import { PASSWORD_TYPE } from "src/shared/utils/static-data";
import { applicationLogout, isFetchBaseQueryError, showErrorToast, showSuccessToast } from "src/shared/utils/utils.service";
import { useAppDispatch } from "src/store/store";

import { useChangePasswordMutation } from "../services/settings.service";

interface IPasswordItem {
  currentPassword: string;
  newPassword: string;
  passwordType: string;
}

const ChangePassword: FC = (): ReactElement => {
  const navigate: NavigateFunction = useNavigate();
  const dispatch = useAppDispatch();

  const [passwordItem, setPasswordItem] = useState<IPasswordItem>({
    currentPassword: "",
    newPassword: "",
    passwordType: PASSWORD_TYPE.PASSWORD,
  });

  const [alertMessage, setAlertMessage] = useState<string>("");

  const [changePassword] = useChangePasswordMutation();

  const updatePassword = async (): Promise<void> => {
    try {
      setAlertMessage("");

      await changePassword({ currentPassword: passwordItem.currentPassword, newPassword: passwordItem.newPassword }).unwrap();

      showSuccessToast("Password Updated Successfully");

      setTimeout(() => {
        applicationLogout(dispatch, navigate);
      }, 3000);
    } catch (error) {
      if (isFetchBaseQueryError(error)) {
        setAlertMessage(error?.data?.message);
      }

      showErrorToast("Error updating password");
    }
  };

  return (
    <div>
      {alertMessage && <Alert type="error" message={alertMessage} />}
      <>
        <label htmlFor="currentPassword" className="text-sm font-bold leading-tight tracking-normal text-gray-800">
          Current Password
        </label>
        <TextInput
          id="currentPassword"
          name="currentPassword"
          type="password"
          value={passwordItem.currentPassword}
          onChange={(event: ChangeEvent) => {
            setPasswordItem({ ...passwordItem, currentPassword: (event.target as HTMLInputElement).value });
          }}
          className="mb-5 mt-2 flex h-10 w-full items-center rounded border border-gray-300 pl-3 text-sm font-normal text-gray-600 focus:border focus:border-sky-500/50 focus:outline-none"
          placeholder="Enter current password"
        />
      </>
      <>
        <label htmlFor="newPassword" className="text-sm font-bold leading-tight tracking-normal text-gray-800">
          New Password
        </label>
        <div className="relative flex gap-4">
          <TextInput
            id="newPassword"
            name="newPassword"
            type={passwordItem.passwordType}
            value={passwordItem.newPassword}
            onChange={(event: ChangeEvent) => {
              setPasswordItem({ ...passwordItem, newPassword: (event.target as HTMLInputElement).value });
            }}
            className="mb-5 mt-2 flex h-10 w-full items-center rounded border border-gray-300 pl-3 text-sm font-normal text-gray-600 focus:border focus:border-sky-500/50 focus:outline-none"
            placeholder="Enter new password"
          />
          <div className="absolute right-0  flex h-full cursor-pointer items-center pr-3 text-gray-600">
            {passwordItem.passwordType === PASSWORD_TYPE.PASSWORD ? (
              <FaEyeSlash className="mb-2" onClick={() => setPasswordItem({ ...passwordItem, passwordType: PASSWORD_TYPE.TEXT })} />
            ) : (
              <FaEye className="mb-2" onClick={() => setPasswordItem({ ...passwordItem, passwordType: PASSWORD_TYPE.PASSWORD })} />
            )}
          </div>
        </div>
        <div className="flex w-full items-center justify-center">
          <Button
            className={`text-md block w-full cursor-pointer rounded  px-8 py-2 text-center font-bold text-white focus:outline-none ${
              !passwordItem.currentPassword || !passwordItem.newPassword
                ? "cursor-not-allowed bg-sky-200"
                : "bg-sky-500 hover:bg-sky-400 cursor-pointer"
            }`}
            label="Save Changes"
            onClick={updatePassword}
          />
        </div>
      </>
    </div>
  );
};

export default ChangePassword;
