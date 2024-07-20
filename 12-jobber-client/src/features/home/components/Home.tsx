import { type FC, type ReactElement } from "react";
import type { ISellerGig } from "src/features/gigs/interfaces/gig.interface";
import { useGetGigsByCategoryQuery, useGetTopRatedGigsByCategoryQuery } from "src/features/gigs/services/gigs.service";
import type { ISellerDocument } from "src/features/sellers/interfaces/seller.interface";
import { useGetRandomSellersQuery } from "src/features/sellers/services/seller.service";
import TopGigsView from "src/shared/gigs/TopGigsView";
import { lowerCase } from "src/shared/utils/utils.service";
import { useAppSelector } from "src/store/store";
import type { IReduxState } from "src/store/store.interface";

import FeaturedExperts from "./FeaturedExperts";
import HomeGigsView from "./HomeGigsView";
import HomeSlider from "./HomeSlider";

const Home: FC = (): ReactElement => {
  const authUser = useAppSelector((state: IReduxState) => state.authUser);

  const { data, isSuccess } = useGetRandomSellersQuery("10");
  const { data: categoryData, isSuccess: isCategorySuccess } = useGetGigsByCategoryQuery(`${authUser.username}`);
  const { data: topGigsData, isSuccess: isTopGigsSuccess } = useGetTopRatedGigsByCategoryQuery(`${authUser.username}`);
  // const { data: topGigsData, isSuccess: isTopGigsSuccess } = useGetGigsBySellerIdQuery(`1`);

  let sellers: ISellerDocument[] = [];
  let categoryGigs: ISellerGig[] = [];
  let topGigs: ISellerGig[] = [];

  if (isSuccess) {
    sellers = data.sellers as ISellerDocument[];
  }

  if (isCategorySuccess) {
    categoryGigs = categoryData.gigs as ISellerGig[];
  }

  if (isTopGigsSuccess) {
    topGigs = topGigsData.gigs as ISellerGig[];
  }

  // useEffect(() => {
  //   // socketService.setupSocketConnection();
  // }, []);

  return (
    <div className="m-auto px-6 w-screen relative min-h-screen xl:container md:px-12 lg:px-6">
      <HomeSlider />
      {topGigs.length > 0 && (
        <TopGigsView
          gigs={topGigs}
          title="Top Rated Services in"
          subTitle={`Highest rated telents for all your ${lowerCase(topGigs[0].categories)} needs`}
          category={topGigs[0].categories}
          width="w-72"
          type="home"
        />
      )}
      {categoryGigs.length > 0 && (
        <HomeGigsView gigs={categoryGigs} title="Because you viewed a gig on" subTitle="" category={categoryGigs[0].categories} />
      )}
      <FeaturedExperts sellers={sellers} />
    </div>
  );
};

export default Home;
