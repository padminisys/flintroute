import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

interface Alert {
  id: number;
  type: string;
  severity: string;
  message: string;
  details: string;
  peer_id?: number;
  acknowledged: boolean;
  acknowledged_at?: string;
  created_at: string;
  updated_at: string;
}

interface AlertsState {
  alerts: Alert[];
  unacknowledgedCount: number;
  loading: boolean;
  error: string | null;
}

const initialState: AlertsState = {
  alerts: [],
  unacknowledgedCount: 0,
  loading: false,
  error: null,
};

const alertsSlice = createSlice({
  name: 'alerts',
  initialState,
  reducers: {
    setAlerts: (state, action: PayloadAction<Alert[]>) => {
      state.alerts = action.payload;
      state.unacknowledgedCount = action.payload.filter((a) => !a.acknowledged).length;
    },
    addAlert: (state, action: PayloadAction<Alert>) => {
      state.alerts.unshift(action.payload);
      if (!action.payload.acknowledged) {
        state.unacknowledgedCount++;
      }
    },
    acknowledgeAlert: (state, action: PayloadAction<number>) => {
      const alert = state.alerts.find((a) => a.id === action.payload);
      if (alert && !alert.acknowledged) {
        alert.acknowledged = true;
        alert.acknowledged_at = new Date().toISOString();
        state.unacknowledgedCount--;
      }
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
  },
});

export const { setAlerts, addAlert, acknowledgeAlert, setLoading, setError } = alertsSlice.actions;

export default alertsSlice.reducer;