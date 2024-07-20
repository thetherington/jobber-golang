import { type ChangeEvent, type FC, type KeyboardEvent, type ReactElement, useState } from "react";
import type { IPersonalInfoProps } from "src/features/sellers/interfaces/seller.interface";
import TextAreaInput from "src/shared/input/TextAreaInput";
import TextInput from "src/shared/input/TextInput";
import type { validationErrorsType } from "src/shared/shared.interface";

interface IAllowedLength {
  description: string;
  oneliner: string;
}

const PersonalInfo: FC<IPersonalInfoProps> = ({ personalInfo, setPersonalInfo, personalInfoErrors }): ReactElement => {
  const [allowedInfoLength, setAllowedInfoLength] = useState<IAllowedLength>({
    description: "600/600",
    oneliner: "70/70",
  });

  const maxDescriptionCharacters = 600;
  const maxOneLinerCharacters = 70;

  const isError = (key: string) => personalInfoErrors.some((obj: validationErrorsType) => Object.keys(obj).includes(key));

  return (
    <div className="border-b border-grey p-6">
      <div className="mb-6 grid md:grid-cols-5">
        <div className="pb-2 text-base font-medium">
          Fullname<sup className="top-[-0.3em] text-base text-red-500">*</sup>
        </div>
        <div className="col-span-4 w-full">
          <TextInput
            className={`border-grey mb-1 w-full rounded border p-2.5 text-sm font-normal text-gray-600 focus:outline-none ${isError("fullName") ? "border-red-500" : ""}`}
            type="text"
            name="fullname"
            value={personalInfo.fullName}
            onChange={(event: ChangeEvent) => {
              setPersonalInfo((prev) => ({ ...prev, fullName: (event.target as HTMLInputElement).value }));
            }}
          />
        </div>
      </div>
      <div className="grid md:grid-cols-5 mb-6">
        <div className="text-base font-medium pb-2 mt-6 md:mt-0">
          Oneliner<sup className="text-red-500 text-base top-[-0.3em]">*</sup>
        </div>
        <div className="w-full col-span-4">
          <TextInput
            className={`w-full rounded border border-grey p-2.5 mb-1 text-sm font-normal text-gray-600 focus:outline-none ${isError("oneliner") ? "border-red-500" : ""}`}
            type="text"
            name="oneliner"
            value={personalInfo.oneliner}
            onChange={(event: ChangeEvent) => {
              const onelinerValue: string = (event.target as HTMLInputElement).value;
              setPersonalInfo((prev) => ({ ...prev, oneliner: onelinerValue }));
              const counter: number = maxOneLinerCharacters - onelinerValue.length;
              setAllowedInfoLength((prev) => ({ ...prev, oneliner: `${counter}/70` }));
            }}
            onKeyDown={(event: KeyboardEvent) => {
              const currentTextLength = (event.target as HTMLInputElement).value.length;
              if (currentTextLength === maxOneLinerCharacters && event.key !== "Backspace") {
                event.preventDefault();
              }
            }}
            placeholder="E.g. Expert Mobile and Web Developer"
          />
          <span className="flex justify-end text-[#95979d] text-xs">{allowedInfoLength.oneliner} Characters</span>
        </div>
      </div>
      <div className="grid md:grid-cols-5 mb-6">
        <div className="text-base font-medium pb-2">
          Description<sup className="text-red-500 text-base top-[-0.3em]">*</sup>
        </div>
        <div className="w-full col-span-4">
          <TextAreaInput
            className={`w-full rounded border border-grey p-2.5 mb-1 text-sm font-normal text-gray-600 focus:outline-none ${isError("description") ? "border-red-500" : ""}`}
            name="description"
            value={personalInfo.description}
            onChange={(event: ChangeEvent) => {
              const descriptionValue: string = (event.target as HTMLInputElement).value;
              setPersonalInfo((prev) => ({ ...prev, description: descriptionValue }));
              const counter: number = maxDescriptionCharacters - descriptionValue.length;
              setAllowedInfoLength((prev) => ({ ...prev, description: `${counter}/600` }));
            }}
            onKeyDown={(event: KeyboardEvent) => {
              const currentTextLength = (event.target as HTMLInputElement).value.length;
              if (currentTextLength === maxDescriptionCharacters && event.key !== "Backspace") {
                event.preventDefault();
              }
            }}
            rows={5}
          />
          <span className="flex justify-end text-[#95979d] text-xs">{allowedInfoLength.description} Characters</span>
        </div>
      </div>
      <div className="grid md:grid-cols-5 mb-6">
        <div className="text-base font-medium pb-2">
          Response Time<sup className="text-red-500 text-base top-[-0.3em]">*</sup>
        </div>
        <div className="w-full col-span-4">
          <TextInput
            className={`w-full rounded border border-grey p-2.5 mb-1 text-sm font-normal text-gray-600 focus:outline-none ${isError("responseTime") ? "border-red-500" : ""}  `}
            type="number"
            name="responseTime"
            value={personalInfo.responseTime}
            onChange={(event: ChangeEvent) => {
              const value = (event.target as HTMLInputElement).value;
              setPersonalInfo((prev) => ({ ...prev, responseTime: parseInt(value) > 0 ? value : "" }));
            }}
            placeholder="E.g. 1"
          />
        </div>
      </div>
    </div>
  );
};

export default PersonalInfo;
