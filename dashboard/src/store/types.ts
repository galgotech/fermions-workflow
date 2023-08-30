/* eslint-disable no-restricted-imports */
import { PayloadAction, ThunkAction } from '@reduxjs/toolkit';
import {
    TypedUseSelectorHook,
    useDispatch as useDispatchUntyped,
    useSelector as useSelectorUntyped,
} from 'react-redux';

import { createRootReducer } from '../core/reducers/root';

import { store } from './configure';

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch

export const useDispatch: () => AppDispatch = useDispatchUntyped
export const useSelector: TypedUseSelectorHook<RootState> = useSelectorUntyped

export type StoreState = ReturnType<ReturnType<typeof createRootReducer>>;
export type ThunkResult<R> = ThunkAction<R, StoreState, undefined, PayloadAction<any>>;
