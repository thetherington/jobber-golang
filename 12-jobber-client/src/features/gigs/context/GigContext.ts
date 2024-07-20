import { type Context, createContext, useContext } from "react";
import { emptyGigData, emptySellerData } from "src/shared/utils/static-data";

import type { IGigContext } from "../interfaces/gig.interface";

export const GigContext: Context<IGigContext> = createContext({
  gig: emptyGigData,
  seller: emptySellerData,
}) as Context<IGigContext>;

export const useGig = () => {
  const context = useContext(GigContext);
  if (context === undefined || context === null) {
    throw new Error("Cannot use GigContext without using GigContext Provider");
  }

  return context;
};
