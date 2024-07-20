import { type FC, lazy, type LazyExoticComponent, type PropsWithChildren, Suspense } from "react";
import { type RouteObject, useRoutes } from "react-router-dom";

import type { IGigsProps } from "./features/gigs/interfaces/gig.interface";
import ProtectedRoute from "./features/ProtectedRoute";

const AppPage: LazyExoticComponent<FC> = lazy(() => import("./features/AppPage"));
const ConfirmEmail: LazyExoticComponent<FC> = lazy(() => import("./features/auth/components/ConfirmEmail"));
const ResetPassword: LazyExoticComponent<FC> = lazy(() => import("./features/auth/components/ResetPassword"));
const Home: LazyExoticComponent<FC> = lazy(() => import("./features/home/components/Home"));
const Error: LazyExoticComponent<FC> = lazy(() => import("./features/error/Error"));
const BuyerDashboard: LazyExoticComponent<FC> = lazy(() => import("./features/buyer/components/Dashboard"));
const AddSeller: LazyExoticComponent<FC> = lazy(() => import("./features/sellers/components/add/AddSeller"));
const CurrentSellerProfile: LazyExoticComponent<FC> = lazy(() => import("./features/sellers/components/profile/CurrentSellerProfile"));
const SellerProfile: LazyExoticComponent<FC> = lazy(() => import("./features/sellers/components/profile/SellerProfile"));
const Seller: LazyExoticComponent<FC> = lazy(() => import("./features/sellers/components/dashboard/Seller"));
const SellerDashboard: LazyExoticComponent<FC> = lazy(() => import("./features/sellers/components/dashboard/SellerDashboard"));
const ManageOrders: LazyExoticComponent<FC> = lazy(() => import("./features/sellers/components/dashboard/ManageOrders"));
const ManageEarnings: LazyExoticComponent<FC> = lazy(() => import("./features/sellers/components/dashboard/ManageEarnings"));
const AddGig: LazyExoticComponent<FC> = lazy(() => import("./features/gigs/components/gig/AddGig"));
const EditGig: LazyExoticComponent<FC> = lazy(() => import("./features/gigs/components/gig/EditGig"));
const GigView: LazyExoticComponent<FC> = lazy(() => import("./features/gigs/components/view/GigView"));
const Gigs: LazyExoticComponent<FC<IGigsProps>> = lazy(() => import("./features/gigs/components/gigs/Gigs"));
const Chat: LazyExoticComponent<FC> = lazy(() => import("./features/chat/components/Chat"));
const Checkout: LazyExoticComponent<FC> = lazy(() => import("./features/order/components/Checkout"));
const Requirement: LazyExoticComponent<FC> = lazy(() => import("./features/order/components/Requirement"));
const Order: LazyExoticComponent<FC> = lazy(() => import("./features/order/components/Order"));
const Settings: LazyExoticComponent<FC> = lazy(() => import("./features/settings/components/Settings"));
const GigsIndexDisplay: LazyExoticComponent<FC<IGigsProps>> = lazy(() => import("./features/index/gig-tabs/GigsIndexDisplay"));
const GigInfoDisplay: LazyExoticComponent<FC> = lazy(() => import("./features/index/gig-tabs/GigInfoDisplay"));

type LayoutProps = PropsWithChildren<{
  backgroundColor?: string;
}>;

const Layout = ({ backgroundColor = "#ffffff", children }: LayoutProps): JSX.Element => (
  <div style={{ backgroundColor }} className="flex flex-grow">
    {children}
  </div>
);

const AppRouter: FC = () => {
  const routes: RouteObject[] = [
    {
      path: "/",
      element: (
        <Suspense>
          <AppPage />
        </Suspense>
      ),
    },
    {
      path: "reset_password",
      element: (
        <Suspense>
          <ResetPassword />
        </Suspense>
      ),
    },
    {
      path: "confirm_email",
      element: (
        <Suspense>
          <ConfirmEmail />
        </Suspense>
      ),
    },
    {
      path: "/gig/:gigId/:title",
      element: (
        <Suspense>
          <Layout backgroundColor="#ffffff">
            <GigInfoDisplay />
          </Layout>
        </Suspense>
      ),
    },
    {
      path: "/search/categories/:category",
      element: (
        <Suspense>
          <Layout>
            <GigsIndexDisplay type="categories" />
          </Layout>
        </Suspense>
      ),
    },
    {
      path: "/gigs/search",
      element: (
        <Suspense>
          <Layout>
            <GigsIndexDisplay type="search" />
          </Layout>
        </Suspense>
      ),
    },
    {
      path: "/",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <Home />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/users/:username/:buyerId/orders",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <BuyerDashboard />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/seller_onboarding",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <AddSeller />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/seller_profile/:username/:sellerId/edit",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <CurrentSellerProfile />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/seller_profile/:username/:sellerId/view",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <SellerProfile />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/:username/:sellerId",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <Seller />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
      children: [
        {
          path: "seller_dashboard",
          element: <SellerDashboard />,
        },
        {
          path: "manage_orders",
          element: <ManageOrders />,
        },
        {
          path: "manage_earnings",
          element: <ManageEarnings />,
        },
      ],
    },
    {
      path: "/manage_gigs/new/:sellerId",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <AddGig />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/manage_gigs/edit/:gigId",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <EditGig />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/gig/:username/:title/:sellerId/:gigId/view",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <GigView />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/categories/:category",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <Gigs type="categories" />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/search/gigs",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <Gigs type="search" />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/inbox",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <Chat />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/inbox/:username/:conversationId",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <Chat />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/gig/checkout/:gigId",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <Checkout />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/gig/order/requirement/:gigId",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout>
              <Requirement />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/orders/:orderId/activities",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout backgroundColor="#f5f5f5">
              <Order />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "/:username/edit",
      element: (
        <Suspense>
          <ProtectedRoute>
            <Layout backgroundColor="#f5f5f5">
              <Settings />
            </Layout>
          </ProtectedRoute>
        </Suspense>
      ),
    },
    {
      path: "*",
      element: (
        <Suspense>
          <Error />
        </Suspense>
      ),
    },
  ];

  return useRoutes(routes);
};

export default AppRouter;
