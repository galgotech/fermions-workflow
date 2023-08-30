import { configureStore } from '@reduxjs/toolkit';

import { createRootReducer } from '../core/reducers/root';

export const store = configureStore({
  reducer: createRootReducer(),
  middleware: (getDefaultMiddleware) =>
  getDefaultMiddleware({
    serializableCheck: false,
  }),
  devTools: process.env.NODE_ENV !== 'production',
});
