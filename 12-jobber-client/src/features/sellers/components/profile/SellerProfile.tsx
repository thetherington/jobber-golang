import { type FC, type ReactElement, useState } from "react";
import { useParams } from "react-router-dom";
import GigViewReviews from "src/features/gigs/components/view/components/GigViewLeft/GigViewReviews";
import type { ISellerGig } from "src/features/gigs/interfaces/gig.interface";
import { useGetGigsBySellerIdQuery } from "src/features/gigs/services/gigs.service";
import { IReviewDocument } from "src/features/order/interfaces/review.interface";
import { useGetReviewsBySellerIdQuery } from "src/features/order/services/review.service";
import Breadcrumb from "src/shared/breadcrumb/Breadcrumb";
import GigCardDisplayItem from "src/shared/gigs/GigCardDisplayItem";
import CircularPageLoader from "src/shared/page-loader/CircularPageLoader";
import { v4 as uuidv4 } from "uuid";

import { useGetSellerByIdQuery } from "../../services/seller.service";
import ProfileHeader from "./components/ProfileHeader";
import ProfileTabs from "./components/ProfileTabs";
import SellerOverview from "./components/SellerOverview";

const SellerProfile: FC = (): ReactElement => {
  const { sellerId } = useParams();

  const [type, setType] = useState<string>("Overview");

  const { data: sellerData, isLoading } = useGetSellerByIdQuery(`${sellerId}`);
  const { data, isSuccess: isSellerGigSuccess, isLoading: isSellerGigLoading } = useGetGigsBySellerIdQuery(`${sellerId}`);

  const {
    data: sellerReviewData,
    isSuccess: isGigReviewSuccess,
    isLoading: isGigReviewLoading,
  } = useGetReviewsBySellerIdQuery(`${sellerId}`);

  let reviews: IReviewDocument[] = [];

  if (isGigReviewSuccess) {
    reviews = sellerReviewData.reviews as IReviewDocument[];
  }

  const isDataLoading: boolean = isSellerGigLoading && !isSellerGigSuccess && isGigReviewLoading && !isGigReviewSuccess;

  return (
    <div className="relative w-full pb-6">
      <Breadcrumb breadCrumbItems={["Seller", `${sellerData && sellerData.seller ? sellerData.seller.username : ""}`]} />
      {isLoading || isDataLoading ? (
        <CircularPageLoader />
      ) : (
        <div className="container mx-auto px-2 md:px-0">
          <ProfileHeader sellerProfile={sellerData?.seller} showHeaderInfo={true} showEditIcons={false} />
          <div className="my-4 cursor-pointer">
            <ProfileTabs type={type} setType={setType} />
          </div>

          <div className="flex flex-wrap bg-white">
            {type === "Overview" && <SellerOverview sellerProfile={sellerData?.seller} showEditIcons={false} />}
            {type === "Active Gigs" && (
              <div className="grid gap-x-6 pt-6 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
                {data?.gigs &&
                  data.gigs.map((gig: ISellerGig) => (
                    <GigCardDisplayItem key={uuidv4()} gig={gig} linkTarget={false} showEditIcon={false} />
                  ))}
              </div>
            )}
            {type === "Ratings & Reviews" && <GigViewReviews showRatings={false} reviews={reviews} hasFetchedReviews={true} />}
          </div>
        </div>
      )}
    </div>
  );
};

export default SellerProfile;
