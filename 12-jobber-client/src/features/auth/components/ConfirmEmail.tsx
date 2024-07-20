import { type FC, type ReactElement, useCallback, useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import Alert from "src/shared/alert/Alert";
import { type IResponse } from "src/shared/shared.interface";
import { useAppDispatch } from "src/store/store";

import { AUTH_FETCH_STATUS } from "../interfaces/auth.interface";
import { addAuthUser } from "../reducers/auth.reducer";
import { useVerifyEmailMutation } from "../services/auth.service";

const ConfirmEmail: FC = (): ReactElement => {
  const [searchParams] = useSearchParams({});
  const dispatch = useAppDispatch();

  const [alertMessage, setAlertMessage] = useState<string>("");
  const [status, setStatus] = useState<string>(AUTH_FETCH_STATUS.IDLE);

  const [verifyEmail] = useVerifyEmailMutation();

  const onVerifyEmail = useCallback(async (): Promise<void> => {
    try {
      const result: IResponse = await verifyEmail(`${searchParams.get("v_token")}`).unwrap();

      setAlertMessage(`Email verified successfully`);
      setStatus(AUTH_FETCH_STATUS.SUCCESS);
      dispatch(addAuthUser({ authInfo: result.user }));
    } catch (error) {
      setStatus(AUTH_FETCH_STATUS.ERROR);
      setAlertMessage(error?.data.message);
    }
  }, [dispatch, searchParams, verifyEmail]);

  useEffect(() => {
    onVerifyEmail();
  }, [onVerifyEmail]);

  return (
    <div className="relative mt-24 mx-auto w-11/12 max-w-md rounded-lg bg-white md:w-2/3">
      <div className="relative px-5 py-5">
        <Alert type={status} message={alertMessage} />
        <Link
          to="/"
          className="rounded bg-sky-500 px-6 py-3 mt-5 text-center text-sm font-bold text-white hover:bg-sky-400 focus:outline-none md:px-4 md:py-2 md:text-base"
        >
          Continue to Home
        </Link>
      </div>
    </div>
  );
};

export default ConfirmEmail;
