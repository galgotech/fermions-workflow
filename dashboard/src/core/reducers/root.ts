import { AnyAction, combineReducers } from '@reduxjs/toolkit';

import cloudEventSliceReducer from '../../components/CloudEvents/states/slice';

const rootReducers = {
  cloudEvents: cloudEventSliceReducer
};

export const createRootReducer = () => {
  const appReducer = combineReducers({
    ...rootReducers,
  });

  return (state: any, action: AnyAction) => {
    return appReducer(state, action);
  };
};
