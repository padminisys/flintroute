import { createSlice, type PayloadAction } from '@reduxjs/toolkit';

interface BGPPeer {
  id: number;
  name: string;
  ip_address: string;
  asn: number;
  remote_asn: number;
  description: string;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

interface BGPSession {
  id: number;
  peer_id: number;
  peer?: BGPPeer;
  state: string;
  uptime: number;
  prefixes_received: number;
  prefixes_sent: number;
  messages_received: number;
  messages_sent: number;
  last_error: string;
  created_at: string;
  updated_at: string;
}

interface BGPState {
  peers: BGPPeer[];
  sessions: BGPSession[];
  selectedPeer: BGPPeer | null;
  loading: boolean;
  error: string | null;
}

const initialState: BGPState = {
  peers: [],
  sessions: [],
  selectedPeer: null,
  loading: false,
  error: null,
};

const bgpSlice = createSlice({
  name: 'bgp',
  initialState,
  reducers: {
    setPeers: (state, action: PayloadAction<BGPPeer[]>) => {
      state.peers = action.payload;
    },
    addPeer: (state, action: PayloadAction<BGPPeer>) => {
      state.peers.push(action.payload);
    },
    updatePeer: (state, action: PayloadAction<BGPPeer>) => {
      const index = state.peers.findIndex((p) => p.id === action.payload.id);
      if (index !== -1) {
        state.peers[index] = action.payload;
      }
    },
    removePeer: (state, action: PayloadAction<number>) => {
      state.peers = state.peers.filter((p) => p.id !== action.payload);
    },
    setSessions: (state, action: PayloadAction<BGPSession[]>) => {
      state.sessions = action.payload;
    },
    updateSession: (state, action: PayloadAction<BGPSession>) => {
      const index = state.sessions.findIndex((s) => s.id === action.payload.id);
      if (index !== -1) {
        state.sessions[index] = action.payload;
      } else {
        state.sessions.push(action.payload);
      }
    },
    setSelectedPeer: (state, action: PayloadAction<BGPPeer | null>) => {
      state.selectedPeer = action.payload;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
  },
});

export const {
  setPeers,
  addPeer,
  updatePeer,
  removePeer,
  setSessions,
  updateSession,
  setSelectedPeer,
  setLoading,
  setError,
} = bgpSlice.actions;

export default bgpSlice.reducer;