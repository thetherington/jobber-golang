import { type Dispatch, type MutableRefObject, SetStateAction, useCallback, useEffect, useState } from "react";

const useDetectOutsideClick = (
  ref: MutableRefObject<HTMLDivElement | null>,
  initialState: boolean,
): [boolean, Dispatch<SetStateAction<boolean>>] => {
  const [isActive, setIsActive] = useState<boolean>(initialState);

  const handleClick = useCallback(
    (event: MouseEvent): void => {
      // if the drop down is visiable and the click is not inside the dropdown
      if (ref.current !== null && !ref.current.contains(event.target as HTMLDivElement)) {
        setIsActive(!isActive);
      }
    },
    [isActive, ref],
  );

  useEffect(() => {
    if (isActive) {
      window.addEventListener("click", handleClick);
    }

    return () => {
      window.removeEventListener("click", handleClick);
    };
  }, [isActive, handleClick]);

  return [isActive, setIsActive];
};

export default useDetectOutsideClick;
