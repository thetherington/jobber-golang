import { createSlice, type Slice } from "@reduxjs/toolkit";

import type { IReduxShowCategory } from "../interfaces/header.interface";

const initialState: boolean = true;

const categoryContainerSlice: Slice = createSlice({
  name: "showCategoryContainer",
  initialState,
  reducers: {
    updateCategoryContainer: (state: boolean, action: IReduxShowCategory): boolean => {
      state = action.payload;
      return state;
    },
  },
});

export const { updateCategoryContainer } = categoryContainerSlice.actions;

export default categoryContainerSlice.reducer;
