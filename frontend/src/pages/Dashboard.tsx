import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Container,
  Box,
  Paper,
  Typography,
  Card,
  CardContent,
  Stack,
} from '@mui/material';
import { bgpPeersAPI, bgpSessionsAPI } from '../services/api';
import { setPeers, setSessions } from '../store/bgpSlice';
import type { RootState } from '../store';

export default function Dashboard() {
  const dispatch = useDispatch();
  const peers = useSelector((state: RootState) => (state.bgp as any).peers || []);
  const sessions = useSelector((state: RootState) => (state.bgp as any).sessions || []);
  const unacknowledgedCount = useSelector((state: RootState) => (state.alerts as any).unacknowledgedCount || 0);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [peersData, sessionsData] = await Promise.all([
        bgpPeersAPI.list(),
        bgpSessionsAPI.list(),
      ]);
      dispatch(setPeers(peersData || []));
      dispatch(setSessions(sessionsData || []));
    } catch (error) {
      console.error('Failed to load data:', error);
    }
  };

  const establishedSessions = sessions.filter((s: any) => s.state === 'Established').length;
  const activePeers = peers.filter((p: any) => p.enabled).length;

  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>

      <Stack spacing={3}>
        {/* Stats Cards */}
        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
          <Card sx={{ flex: 1 }}>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Total Peers
              </Typography>
              <Typography variant="h3">{peers.length}</Typography>
            </CardContent>
          </Card>

          <Card sx={{ flex: 1 }}>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Active Peers
              </Typography>
              <Typography variant="h3">{activePeers}</Typography>
            </CardContent>
          </Card>

          <Card sx={{ flex: 1 }}>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Established Sessions
              </Typography>
              <Typography variant="h3">{establishedSessions}</Typography>
            </CardContent>
          </Card>

          <Card sx={{ flex: 1, bgcolor: unacknowledgedCount > 0 ? 'warning.light' : 'inherit' }}>
            <CardContent>
              <Typography color="text.secondary" gutterBottom>
                Unacknowledged Alerts
              </Typography>
              <Typography variant="h3">{unacknowledgedCount}</Typography>
            </CardContent>
          </Card>
        </Stack>

        {/* Recent Sessions */}
        <Paper sx={{ p: 2 }}>
          <Typography variant="h6" gutterBottom>
            BGP Sessions
          </Typography>
          <Box sx={{ mt: 2 }}>
            {sessions.length === 0 ? (
              <Typography color="text.secondary">No sessions found</Typography>
            ) : (
              <Stack spacing={1}>
                {sessions.slice(0, 5).map((session: any) => (
                  <Box
                    key={session.id}
                    sx={{
                      p: 2,
                      border: 1,
                      borderColor: 'divider',
                      borderRadius: 1,
                    }}
                  >
                    <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2}>
                      <Box sx={{ flex: 1 }}>
                        <Typography variant="subtitle2">
                          {session.peer?.name || 'Unknown'}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {session.peer?.ip_address}
                        </Typography>
                      </Box>
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          State
                        </Typography>
                        <Typography
                          variant="body1"
                          color={session.state === 'Established' ? 'success.main' : 'error.main'}
                        >
                          {session.state}
                        </Typography>
                      </Box>
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          Prefixes
                        </Typography>
                        <Typography variant="body1">
                          RX: {session.prefixes_received} / TX: {session.prefixes_sent}
                        </Typography>
                      </Box>
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          Uptime
                        </Typography>
                        <Typography variant="body1">
                          {Math.floor(session.uptime / 3600)}h {Math.floor((session.uptime % 3600) / 60)}m
                        </Typography>
                      </Box>
                    </Stack>
                  </Box>
                ))}
              </Stack>
            )}
          </Box>
        </Paper>
      </Stack>
    </Container>
  );
}