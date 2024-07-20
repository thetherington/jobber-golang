import { Context, createContext, useContext } from "react";
import { emptySellerData } from "src/shared/utils/static-data";

import { ISellerContext } from "../interfaces/seller.interface";

export const SellerContext: Context<ISellerContext> = createContext({
  showEditIcons: false,
  sellerProfile: emptySellerData,
}) as Context<ISellerContext>;

export const useSeller = () => {
  const context = useContext(SellerContext);
  if (context === undefined || context === null) {
    throw new Error("Cannot use SellerContext without using SellerOverview Provider");
  }

  return context;
};
