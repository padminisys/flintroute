import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';
import bgpReducer from './bgpSlice';
import alertsReducer from './alertsSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    bgp: bgpReducer,
    alerts: alertsReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;