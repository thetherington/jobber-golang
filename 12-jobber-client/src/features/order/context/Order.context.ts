import { type Context, createContext, useContext } from "react";
import type { IAuthUser } from "src/features/auth/interfaces/auth.interface";

import type { IOrderContext, IOrderDocument, IOrderInvoice } from "../interfaces/order.interface";

export const OrderContext: Context<IOrderContext> = createContext({
  order: {} as IOrderDocument,
  authUser: {} as IAuthUser,
  orderInvoice: {} as IOrderInvoice,
}) as Context<IOrderContext>;

export const useOrder = () => {
  const context = useContext(OrderContext);
  if (context === undefined || context === null) {
    throw new Error("Cannot use OrderContext without using OrderContext Provider");
  }

  return context;
};
