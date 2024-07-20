import { filter } from "lodash";
import { type FC, FormEvent, type ReactElement, useEffect, useState } from "react";
import { FaReact } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import type { IBuyerDocument } from "src/features/buyer/interfaces/buyer.interface";
import { addBuyer } from "src/features/buyer/reducers/buyer.reducer";
import type {
  ICertificate,
  IEducation,
  IExperience,
  ILanguage,
  IPersonalInfoData,
  ISellerDocument,
} from "src/features/sellers/interfaces/seller.interface";
import Breadcrumb from "src/shared/breadcrumb/Breadcrumb";
import Button from "src/shared/button/Button";
import CircularPageLoader from "src/shared/page-loader/CircularPageLoader";
import type { IResponse } from "src/shared/shared.interface";
import { deleteFromLocalStorage, lowerCase } from "src/shared/utils/utils.service";
import { useAppDispatch, useAppSelector } from "src/store/store";
import type { IReduxState } from "src/store/store.interface";

import { useSellerSchema } from "../../hooks/useSellerSchema";
import { addSeller } from "../../reducers/seller.reducer";
import { useCreateSellerMutation } from "../../services/seller.service";
import PersonalInfo from "./components/PersonalInfo";
import SellerCertificateFields from "./components/SellerCertificateFields";
import SellerEducationFields from "./components/SellerEductionFields";
import SellerExperienceFields from "./components/SellerExperienceFields";
import SellerLanguageFields from "./components/SellerLanguageFields";
import SellerSkillField from "./components/SellerSkillField";
import SellerSocialLinksFields from "./components/SellerSocialLinksfields";

const AddSeller: FC = (): ReactElement => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const authUser = useAppSelector((state: IReduxState) => state.authUser);
  const buyer = useAppSelector((state: IReduxState) => state.buyer);
  const [createSeller, { isLoading }] = useCreateSellerMutation();

  const [personalInfo, setPersonalInfo] = useState<IPersonalInfoData>({
    fullName: "",
    profilePicture: `${authUser.profilePicture}`,
    description: "",
    responseTime: "",
    oneliner: "",
  });

  const [experienceFields, setExperienceFields] = useState<IExperience[]>([
    {
      title: "",
      company: "",
      startDate: "Start Year",
      endDate: "End Year",
      currentlyWorkingHere: false,
      description: "",
    },
  ]);

  const [educationFields, setEducationFields] = useState<IEducation[]>([
    {
      country: "Country",
      university: "",
      title: "Title",
      major: "",
      year: "Year",
    },
  ]);

  const [skillsFields, setSkillsFields] = useState<string[]>([""]);

  const [languageFields, setLanguageFields] = useState<ILanguage[]>([
    {
      language: "",
      level: "Level",
    },
  ]);

  const [certificateFields, setCertificateFields] = useState<ICertificate[]>([
    {
      name: "",
      from: "",
      year: "Year",
    },
  ]);

  const [socialFields, setSocialFields] = useState<string[]>([""]);

  const [schemaValidation, personalInfoErrors, experienceErrors, eductionErrors, skillsErrors, languagesErrors] = useSellerSchema({
    personalInfo,
    experienceFields,
    educationFields,
    skillsFields,
    languageFields,
  });

  const onCreateSeller = async (event: FormEvent): Promise<void> => {
    event.preventDefault();

    try {
      const isValid: boolean = await schemaValidation();

      if (isValid) {
        const skills: string[] = filter(skillsFields, (skill: string) => skill !== "") as string[];
        const socialLinks: string[] = filter(socialFields, (item: string) => item !== "") as string[];

        const certificates: ICertificate[] = filter(
          certificateFields,
          (item: ICertificate) => item.name !== "" && item.from !== "" && item.year !== "",
        ) as ICertificate[];

        const sellerData: ISellerDocument = {
          email: `${authUser.email}`,
          profilePublicId: `${authUser.profilePublicId}`,
          profilePicture: `${authUser.profilePicture}`,
          fullName: personalInfo.fullName,
          description: personalInfo.description,
          country: `${authUser.country}`,
          skills,
          oneliner: personalInfo.oneliner,
          languages: languageFields,
          responseTime: parseInt(personalInfo.responseTime, 10),
          experience: experienceFields,
          education: educationFields,
          socialLinks,
          certificates,
        };

        const updateBuyer: IBuyerDocument = { ...buyer, isSeller: true };

        const response: IResponse = await createSeller(sellerData).unwrap();

        dispatch(addSeller(response.seller));
        dispatch(addBuyer(updateBuyer));
        navigate(`/seller_profile/${lowerCase(`${authUser.username}`)}/${response.seller?._id}/edit`);
      }
    } catch (error) {
      console.log(error);
    }
  };

  const populateState = (): void => {
    setPersonalInfo({
      fullName: "Thomas Hetherington",
      profilePicture: `${authUser.profilePicture}`,
      description:
        "Lorem ipsum dolor sit amet consectetur adipisicing elit. Ipsam blanditiis rerum eum sequi fuga veniam pariatur accusantium, sunt quas ex.",
      responseTime: "5",
      oneliner: "Engineering Specialist",
    });

    setExperienceFields([
      {
        title: "Customer Service",
        company: "Evertz",
        startDate: "2004",
        endDate: "2006",
        currentlyWorkingHere: false,
        description:
          "Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat, laboriosam libero laudantium reiciendis dolorem eos.",
      },
      {
        title: "Engineering Specialist",
        company: "Evertz",
        startDate: "2006",
        endDate: "2024",
        currentlyWorkingHere: true,
        description:
          "Lorem ipsum dolor sit amet consectetur adipisicing elit. Placeat, laboriosam libero laudantium reiciendis dolorem eos.",
      },
    ]);

    setEducationFields([
      {
        country: "Canada",
        university: "RCC",
        title: "Computer Networks Engineering Technologist",
        major: "-",
        year: "2004",
      },
    ]);

    setSkillsFields(["Devops", "NMS", "Analytics"]);

    setLanguageFields([
      {
        language: "English",
        level: "Expert",
      },
    ]);

    setCertificateFields([
      {
        name: "CCNA",
        from: "Cisco",
        year: "2004",
      },
    ]);

    setSocialFields(["facebook.com/thetherington", "http://www.linkedin.com/in/tom-hetherington-0913801b2"]);
  };

  // remove the local storage of "becomeASeller" key when navigate away from this page.
  useEffect(() => {
    return () => {
      deleteFromLocalStorage("becomeASeller");
    };
  }, []);

  return (
    <div className="relative w-full">
      <Breadcrumb breadCrumbItems={["Seller", "Create Profile"]} />

      <div className="container mx-auto my-5 overflow-hidden px-2 pb-12 md:px-0">
        {isLoading && <CircularPageLoader />}
        {authUser && !authUser.emailVerified && (
          <div className="absolute left-0 top-0 z-50 flex h-full w-full justify-center bg-white/[0.8] text-sm font-bold md:text-base lg:text-xl">
            <span className="mt-20">Please verify your email.</span>
          </div>
        )}

        <div className="left-0 top-0 z-10 mt-4 block h-full bg-white">
          <PersonalInfo personalInfo={personalInfo} setPersonalInfo={setPersonalInfo} personalInfoErrors={personalInfoErrors} />
          <SellerExperienceFields
            experienceFields={experienceFields}
            setExperienceFields={setExperienceFields}
            experienceErrors={experienceErrors}
          />
          <SellerEducationFields
            educationFields={educationFields}
            setEducationFields={setEducationFields}
            educationErrors={eductionErrors}
          />
          <SellerSkillField skillsFields={skillsFields} setSkillsFields={setSkillsFields} skillsErrors={skillsErrors} />
          <SellerLanguageFields languageFields={languageFields} setLanguageFields={setLanguageFields} languagesErrors={languagesErrors} />
          <SellerCertificateFields certificatesFields={certificateFields} setCertificatesFields={setCertificateFields} />
          <SellerSocialLinksFields socialFields={socialFields} setSocialFields={setSocialFields} />
          <div className="flex justify-end p-6">
            <Button
              className="rounded bg-sky-500 px-8 text-center text-sm font-bold text-white hover:bg-sky-400 focus:outline-none md:py-3 md:text-base"
              onClick={onCreateSeller}
              label="Create Profile"
            />
            <span
              onClick={populateState}
              className="inline-flex rounded mx-2 items-center text-sm font-bold bg-sky-500 uppercase hover:text-blue-600 dark:text-white dark:hover:text-white hover:bg-sky-400 focus:outline-none cursor-pointer"
            >
              <FaReact className="mx-2 h-5 w-5" />
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AddSeller;
