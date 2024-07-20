import { lowerCase } from "lodash";
import Quill from "quill";
import { type ChangeEvent, FC, ReactElement, useRef, useState } from "react";
import equal from "react-fast-compare";
import { FaCamera, FaReact } from "react-icons/fa";
import ReactQuill, { type UnprivilegedEditor } from "react-quill";
import { useNavigate, useParams } from "react-router-dom";
import type { ISellerDocument } from "src/features/sellers/interfaces/seller.interface";
import { addSeller } from "src/features/sellers/reducers/seller.reducer";
import Breadcrumb from "src/shared/breadcrumb/Breadcrumb";
import Button from "src/shared/button/Button";
import Dropdown from "src/shared/dropdown/Dropdown";
import TextAreaInput from "src/shared/input/TextAreaInput";
import TextInput from "src/shared/input/TextInput";
import ApprovalModal from "src/shared/modals/ApprovalModal";
import type { IApprovalModalContent } from "src/shared/modals/interfaces/modal.interface";
import CircularPageLoader from "src/shared/page-loader/CircularPageLoader";
import type { IResponse, validationErrorsType } from "src/shared/shared.interface";
import { checkImage, readAsBase64 } from "src/shared/utils/image-utils.service";
import { categories, expectedGigDelivery, reactQuillUtils, replaceSpacesWithDash, showErrorToast } from "src/shared/utils/utils.service";
import { useAppDispatch, useAppSelector } from "src/store/store";
import type { IReduxState } from "src/store/store.interface";

import { useGigSchema } from "../../hooks/ useGigSchema";
import { GIG_MAX_LENGTH, IAllowedGigItem, type ICreateGig, IShowGigModal } from "../../interfaces/gig.interface";
import { gigInfoSchema } from "../../schemes/gig.schema";
import { useCreateGigMutation } from "../../services/gigs.service";
import TagsInput from "./components/TagsInput";

const defaultGigInfo: ICreateGig = {
  title: "",
  categories: "",
  description: "",
  subCategories: [],
  tags: [],
  price: 0,
  coverImage: "https://placehold.co/330x220?text=Cover+Image",
  expectedDelivery: "Expected delivery",
  basicTitle: "",
  basicDescription: "",
};

const AddGig: FC = (): ReactElement => {
  const navigate = useNavigate();
  const { sellerId } = useParams();

  const dispatch = useAppDispatch();

  const [createGig, { isLoading }] = useCreateGigMutation();

  const authUser = useAppSelector((state: IReduxState) => state.authUser);
  const seller = useAppSelector((state: IReduxState) => state.seller);

  const [gigInfo, setGigInfo] = useState<ICreateGig>(defaultGigInfo);

  const [subCategory, setSubCategory] = useState<string[]>([]);
  const [subCategoryInput, setSubCategoryInput] = useState<string>("");

  const [tags, setTags] = useState<string[]>([]);
  const [tagsInput, setTagsInput] = useState<string>("");

  const [showGigModal, setShowGigModal] = useState<IShowGigModal>({
    image: false,
    cancel: false,
  });

  const reactQuillRef = useRef<ReactQuill | null>(null);
  const fileRef = useRef<HTMLInputElement>(null);

  const [allowedGigItemLength, setAllowedGigItemLength] = useState<IAllowedGigItem>({
    gigTitle: "80/80",
    basicTitle: "40/40",
    basicDescription: "100/100",
    descriptionCharacters: "1200/1200",
  });

  const gigInfoRef = useRef<ICreateGig>(defaultGigInfo);
  const [approvalModalContent, setApprovalModalContent] = useState<IApprovalModalContent>();

  const [schemaValidation, gigInfoErrors] = useGigSchema({ schema: gigInfoSchema, gigInfo });
  const isError = (key: string) => gigInfoErrors.some((obj: validationErrorsType) => Object.keys(obj).includes(key));

  const handleFileChange = async (event: ChangeEvent): Promise<void> => {
    const target: HTMLInputElement = event.target as HTMLInputElement;

    if (target.files) {
      const file: File = target.files[0];
      const isValid = checkImage(file, "image");

      if (isValid) {
        const dataImage: string | ArrayBuffer | null = await readAsBase64(file);
        setGigInfo({ ...gigInfo, coverImage: `${dataImage}` });
      }

      setShowGigModal({ ...showGigModal, image: false });
    }
  };

  const onCreateGig = async (): Promise<void> => {
    try {
      const editor: Quill | undefined = reactQuillRef.current?.editor;
      // In React, it is not recommended to mutate objects directly. It is better to update with useState method.
      // The reason it is not recommended is because if the object is mutated directly,
      // 1) React is not able to keep track of the change
      // 2) There will be no re-renderng of the component.
      // In our case, we don't care about the above reasons because we update a property, validate and send to the backend.
      // The updated properly is not reflected in the component and we don't need to keep track of the object.
      // We are not using the useState method inside useEffect because it causes too many rerender errors.
      // Also, we are not updating the property inside the onChange method because editor?.getText() causes too many rerender errors.
      // The only option we have right now is to directly mutate the gigInfo useState object.
      gigInfo.description = editor?.getText().trim() as string;

      console.log(gigInfo);

      const isValid: boolean = await schemaValidation();
      console.log("valid", isValid);
      if (isValid) {
        const gig: ICreateGig = {
          profilePicture: `${authUser.profilePicture}`,
          sellerId,
          title: gigInfo.title,
          categories: gigInfo.categories,
          description: gigInfo.description,
          subCategories: subCategory,
          tags,
          price: gigInfo.price,
          coverImage: gigInfo.coverImage,
          expectedDelivery: gigInfo.expectedDelivery,
          basicTitle: gigInfo.basicTitle,
          basicDescription: gigInfo.basicDescription,
        };

        const response: IResponse = await createGig(gig).unwrap();

        // increment seller state with total gigs + 1
        const updatedSeller: ISellerDocument = { ...seller, totalGigs: (seller.totalGigs as number) + 1 };
        dispatch(addSeller(updatedSeller));

        const title: string = replaceSpacesWithDash(gig.title);
        navigate(`/gig/${lowerCase(`${authUser.username}`)}/${title}/${response?.gig?.sellerId}/${response?.gig?.id}/view`);
      }
    } catch (error) {
      console.log(error);
      showErrorToast("Error creating gig");
    }
  };

  const onCancelCreate = (): void => {
    navigate(`/seller_profile/${lowerCase(`${authUser.username}`)}/${sellerId}/edit`);
  };

  const populateState = (): void => {
    setGigInfo({
      title: "I will build a complete web app using React and NodeJS",
      basicTitle: "I will build a web application for you",
      basicDescription: "I am a web developer with lots of experience in React",
      description: "",
      categories: "Programming & Tech",
      subCategories: ["web dev", "programming"],
      tags: ["react", "nodejs", "web"],
      price: 69,
      coverImage: "https://placehold.co/330x220?text=Cover+Image",
      expectedDelivery: "2 Days Delivery",
    });
    setSubCategory(["web dev", "programming"]);
    setTags(["react", "nodejs", "web"]);
  };

  const populateDescription = (): void => {
    const editor: Quill | undefined = reactQuillRef.current?.editor;
    editor?.clipboard.dangerouslyPasteHTML(
      "<p><strong>Skills</strong></p><ol><li>Web Dev</li><li>Backend</li><li>API Development</li></ol><p><br></p><p>Willing to build a full stack application with these languages</p><ul><li>GO</li><li>Typescript</li><li>React</li><li>GraphQL</li></ul>",
    );
  };

  return (
    <>
      {showGigModal.cancel && (
        <ApprovalModal
          approvalModalContent={approvalModalContent}
          onClose={() => setShowGigModal({ ...showGigModal, cancel: false })}
          onClick={onCancelCreate}
        />
      )}
      <div className="relative w-screen">
        <Breadcrumb breadCrumbItems={["Seller", "Create new gig"]} />
        <div className="container relative mx-auto my-5 px-2 pb-12 md:px-0">
          {isLoading && <CircularPageLoader />}
          {authUser && !authUser.emailVerified && (
            <div className="absolute left-0 top-0 z-[80] flex h-full w-full justify-center bg-white/[0.8] text-sm font-bold md:text-base lg:text-xl">
              <span className="mt-40">Please verify your email.</span>
            </div>
          )}

          <div className="border-grey left-0 top-0 z-10 mt-4 block rounded border bg-white p-6">
            <div className="mb-6 grid md:grid-cols-5">
              <div className="pb-2 text-base font-medium">
                Gig title<sup className="top-[-0.3em] text-base text-red-500">*</sup>
              </div>
              <div className="col-span-4 md:w-11/12 lg:w-8/12">
                <TextInput
                  className={`border-grey mb-1 w-full rounded border p-2.5 text-sm font-normal text-gray-600 focus:outline-none ${isError("title") ? "border-red-500" : ""}`}
                  type="text"
                  name="gigTitle"
                  value={gigInfo.title}
                  placeholder="I will build something I'm good at."
                  maxLength={80}
                  onChange={(event: ChangeEvent) => {
                    const gigTitleValue: string = (event.target as HTMLInputElement).value;
                    setGigInfo({ ...gigInfo, title: gigTitleValue });

                    const counter: number = GIG_MAX_LENGTH.gigTitle - gigTitleValue.length;
                    setAllowedGigItemLength({ ...allowedGigItemLength, gigTitle: `${counter}/80` });
                  }}
                />
                <span className="flex justify-end text-xs text-[#95979d]">{allowedGigItemLength.gigTitle} Characters</span>
              </div>
            </div>
            <div className="mb-6 grid md:grid-cols-5">
              <div className="pb-2 text-base font-medium">
                Basic title<sup className="top-[-0.3em] text-base text-red-500">*</sup>
              </div>
              <div className="col-span-4 md:w-11/12 lg:w-8/12">
                <TextInput
                  className={`border-grey mb-1 w-full rounded border p-2.5 text-sm font-normal text-gray-600 focus:outline-none ${isError("basicTitle") ? "border-red-500" : ""}`}
                  placeholder="Write what exactly you'll do in short."
                  type="text"
                  name="basicTitle"
                  value={gigInfo.basicTitle}
                  maxLength={40}
                  onChange={(event: ChangeEvent) => {
                    const basicTitleValue: string = (event.target as HTMLInputElement).value;
                    setGigInfo({ ...gigInfo, basicTitle: basicTitleValue });

                    const counter: number = GIG_MAX_LENGTH.basicTitle - basicTitleValue.length;
                    setAllowedGigItemLength({ ...allowedGigItemLength, basicTitle: `${counter}/40` });
                  }}
                />
                <span className="flex justify-end text-xs text-[#95979d]">{allowedGigItemLength.basicTitle} Characters</span>
              </div>
            </div>
            <div className="mb-6 grid md:grid-cols-5">
              <div className="pb-2 text-base font-medium">
                Brief description<sup className="top-[-0.3em] text-base text-red-500">*</sup>
              </div>
              <div className="col-span-4 md:w-11/12 lg:w-8/12">
                <TextAreaInput
                  className={`border-grey mb-1 w-full rounded border p-2.5 text-sm font-normal text-gray-600 focus:outline-none ${isError("basicDescription") ? "border-red-500" : ""}`}
                  placeholder="Write a brief description..."
                  name="basicDescription"
                  value={gigInfo.basicDescription}
                  rows={5}
                  maxLength={100}
                  onChange={(event: ChangeEvent) => {
                    const basicDescription: string = (event.target as HTMLInputElement).value;
                    setGigInfo({ ...gigInfo, basicDescription: basicDescription });

                    const counter: number = GIG_MAX_LENGTH.basicDescription - basicDescription.length;
                    setAllowedGigItemLength({ ...allowedGigItemLength, basicDescription: `${counter}/100` });
                  }}
                />
                <span className="flex justify-end text-xs text-[#95979d]">{allowedGigItemLength.basicDescription} Characters</span>
              </div>
            </div>
            <div className="mb-6 grid md:grid-cols-5">
              <div className="pb-2 text-base font-medium">
                Full description<sup className="top-[-0.3em] text-base text-red-500">*</sup>
              </div>
              <div className="col-span-4 md:w-11/12 lg:w-8/12">
                <ReactQuill
                  theme="snow"
                  value={gigInfo.description}
                  className={`border-grey border rounded ${isError("description") ? "border-red-500" : ""}`}
                  modules={reactQuillUtils().modules}
                  formats={reactQuillUtils().formats}
                  ref={(element: ReactQuill | null) => {
                    reactQuillRef.current = element;
                    const reactQuillEditor = reactQuillRef.current?.getEditor();
                    reactQuillEditor?.on("text-change", () => {
                      if (reactQuillEditor.getLength() > GIG_MAX_LENGTH.fullDescription) {
                        reactQuillEditor.deleteText(GIG_MAX_LENGTH.fullDescription, reactQuillEditor.getLength());
                      }
                    });
                  }}
                  onChange={(event: string, _, __, editor: UnprivilegedEditor) => {
                    setGigInfo({ ...gigInfo, description: event });
                    const counter: number = GIG_MAX_LENGTH.fullDescription - editor.getText().length;
                    setAllowedGigItemLength({ ...allowedGigItemLength, descriptionCharacters: `${counter}/1200` });
                  }}
                />
                <span className="flex justify-end text-xs text-[#95979d]">{allowedGigItemLength.descriptionCharacters} Characters</span>
              </div>
            </div>
            <div className="mb-12 grid md:grid-cols-5">
              <div className="pb-2 text-base font-medium">
                Category<sup className="top-[-0.3em] text-base text-red-500">*</sup>
              </div>
              <div className="relative col-span-4 md:w-11/12 lg:w-8/12">
                <Dropdown
                  text={gigInfo.categories}
                  maxHeight="300"
                  mainClassNames={`absolute bg-white ${isError("categories") ? "border-red-500" : ""}`}
                  values={categories()}
                  onClick={(item: string) => setGigInfo({ ...gigInfo, categories: item })}
                />
              </div>
            </div>

            <TagsInput
              title="SubCategory"
              placeholder="E.g. Website development, Mobile Apps"
              gigInfo={gigInfo}
              setGigInfo={setGigInfo}
              tags={subCategory}
              itemInput={subCategoryInput}
              itemName="subCategories"
              counterText="Subcategories"
              inputErrorMessage={false}
              setItem={setSubCategory}
              setItemInput={setSubCategoryInput}
            />

            <TagsInput
              title="Tags"
              placeholder="enter search terms for your gig"
              gigInfo={gigInfo}
              setGigInfo={setGigInfo}
              tags={tags}
              itemInput={tagsInput}
              itemName="tags"
              counterText="Tags"
              inputErrorMessage={false}
              setItem={setTags}
              setItemInput={setTagsInput}
            />

            <div className="mb-6 grid md:grid-cols-5">
              <div className="pb-2 text-base font-medium">
                Price<sup className="top-[-0.3em] text-base text-red-500">*</sup>
              </div>
              <div className="col-span-4 md:w-11/12 lg:w-8/12">
                <TextInput
                  type="number"
                  className={`border-grey mb-1 w-full rounded border p-2.5 text-sm font-normal text-gray-600 focus:outline-none ${isError("price") ? "border-red-500" : ""}`}
                  placeholder="Enter minimum price"
                  name="price"
                  value={`${gigInfo.price}`}
                  onChange={(event: ChangeEvent) => {
                    const price: string = (event.target as HTMLInputElement).value;
                    setGigInfo({ ...gigInfo, price: parseInt(price) > 0 ? +price : 0 });
                  }}
                />
              </div>
            </div>
            <div className="mb-12 grid md:grid-cols-5">
              <div className="pb-2 text-base font-medium">
                Expected delivery<sup className="top-[-0.3em] text-base text-red-500">*</sup>
              </div>
              <div className="relative col-span-4 md:w-11/12 lg:w-8/12">
                <Dropdown
                  text={gigInfo.expectedDelivery}
                  maxHeight="300"
                  mainClassNames={`absolute bg-white z-40 ${isError("expectedDelivery") ? "border-red-500" : ""}`}
                  values={expectedGigDelivery()}
                  onClick={(item: string) => setGigInfo({ ...gigInfo, expectedDelivery: item })}
                />
              </div>
            </div>
            <div className="mb-6 grid md:grid-cols-5">
              <div className="mt-6 pb-2 text-base font-medium lg:mt-0">
                Cover image<sup className="top-[-0.3em] text-base text-red-500">*</sup>
              </div>
              <div
                className="relative col-span-4 cursor-pointer md:w-11/12 lg:w-8/12"
                onMouseEnter={() => {
                  setShowGigModal((item) => ({ ...item, image: !item.image }));
                }}
                onMouseLeave={() => {
                  setShowGigModal((item) => ({ ...item, image: false }));
                }}
              >
                {gigInfo.coverImage && (
                  <img src={gigInfo.coverImage} alt="Cover Image" className="left-0 top-0 h-[220px] w-[320px] bg-white object-cover" />
                )}
                {!gigInfo.coverImage && (
                  <div className="left-0 top-0 flex h-[220px] w-[320px] cursor-pointer justify-center bg-[#dee1e7]"></div>
                )}
                {showGigModal.image && (
                  <div
                    onClick={() => fileRef.current?.click()}
                    className="absolute left-0 top-0 flex h-[220px] w-[320px] cursor-pointer justify-center bg-[#dee1e7]"
                  >
                    <FaCamera className="flex self-center" />
                  </div>
                )}
                <TextInput
                  ref={fileRef}
                  name="image"
                  type="file"
                  style={{ display: "none" }}
                  onClick={() => {
                    if (fileRef.current) {
                      fileRef.current.value = "";
                    }
                  }}
                  onChange={handleFileChange}
                />
              </div>
            </div>
            <div className="grid xs:grid-cols-1 md:grid-cols-5">
              <div className="pb-2 text-base font-medium lg:mt-0"></div>
              <div className="col-span-4 flex gap-x-4 md:w-11/12 lg:w-8/12">
                <Button
                  disabled={isLoading}
                  className="rounded bg-sky-500 px-8 py-3 text-center text-sm font-bold text-white hover:bg-sky-400 focus:outline-none md:py-3 md:text-base"
                  label="Create Gig"
                  onClick={onCreateGig}
                />
                <Button
                  disabled={isLoading}
                  className="rounded bg-red-500 px-8 py-3 text-center text-sm font-bold text-white hover:bg-red-400 focus:outline-none md:py-3 md:text-base"
                  label="Cancel"
                  onClick={() => {
                    const isEqual: boolean = equal(gigInfo, gigInfoRef.current);
                    if (!isEqual) {
                      setApprovalModalContent({
                        header: "Cancel Gig Creation",
                        body: "Are you sure you want to cancel?",
                        btnText: "Yes, Cancel",
                        btnColor: "bg-red-500 hover:bg-red-400",
                      });
                      setShowGigModal({ ...showGigModal, cancel: true });
                    } else {
                      onCancelCreate();
                    }
                  }}
                />
                <span
                  onClick={populateState}
                  onMouseLeave={populateDescription}
                  className="inline-flex rounded mx-2 items-center text-sm font-bold focus:outline-none cursor-pointer"
                >
                  <FaReact className="mx-2 h-5 w-5" />
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default AddGig;
