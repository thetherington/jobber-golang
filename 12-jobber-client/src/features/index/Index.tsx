import { type FC, lazy, type LazyExoticComponent, type ReactElement, Suspense, useEffect } from "react";
import { IHeader } from "src/shared/header/interfaces/header.interface";
import CircularPageLoader from "src/shared/page-loader/CircularPageLoader";
import { saveToSessionStorage } from "src/shared/utils/utils.service";

const IndexHeader: LazyExoticComponent<FC<IHeader>> = lazy(() => import("src/shared/header/components/Header"));
const Hero: LazyExoticComponent<FC> = lazy(() => import("./Hero"));
const GigTabs: LazyExoticComponent<FC> = lazy(() => import("./gig-tabs/GigTabs"));
const HowItWorks: LazyExoticComponent<FC> = lazy(() => import("./HowItWorks"));
const Categories: LazyExoticComponent<FC> = lazy(() => import("./Categories"));

const Index: FC = (): ReactElement => {
  useEffect(() => {
    saveToSessionStorage(JSON.stringify(false), JSON.stringify(""));
  }, []);

  return (
    <div className="flex flex-col">
      <Suspense fallback={<CircularPageLoader />}>
        <IndexHeader navClass="navbar peer-checked:navbar-active fixed z-20 w-full border-b border-gray-100 bg-white/90 shadow-2xl shadow-gray-600/5 backdrop-blur dark:border-gray-800 dark:bg-gray-900/80 dark:shadow-none" />
        <Hero />
        <GigTabs />
        <HowItWorks />
        <Categories />
      </Suspense>
    </div>
  );
};

export default Index;
