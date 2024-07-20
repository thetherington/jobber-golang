import { createSlice, type Slice } from "@reduxjs/toolkit";

import type { IReduxHeader } from "../interfaces/header.interface";

const initialState: string = "index";

const headerSlice: Slice = createSlice({
  name: "header",
  initialState,
  reducers: {
    updateHeader: (state: string, action: IReduxHeader): string => {
      state = action.payload;
      return state;
    },
  },
});

export const { updateHeader } = headerSlice.actions;

export default headerSlice.reducer;
