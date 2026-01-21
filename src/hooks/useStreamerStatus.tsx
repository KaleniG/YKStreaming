import * as React from "react";
import axios from "axios";
import useAuth from "@/core/AuthContext";
import { useNavigate } from "react-router-dom";

interface GetUserStreamsResponse {
  streams: null | {
    name: string;
    key: string;
    is_active: boolean;
    ended_at: Date;
    started_at: Date;
    total_views: number;
    is_vod: boolean;
    live_viewers: number;
  }[];
}

export function useStreamerStatus(intervalMs: number = 5000) {
  const statusAuth = useAuth();
  const navigate = useNavigate();

  const streamsState = React.useState<Array<any>>([]);
  const [streams, setStreams] = streamsState;
  const [streaming, setStreaming] = React.useState<boolean>(false);

  React.useEffect(() => {
    let isMounted = true;

    const fetchStatus = async () => {
      try {
        const res = await axios.post<GetUserStreamsResponse>(
          "http://localhost/api/user/streams/",
          {},
          { withCredentials: true }
        );

        if (isMounted) {
          setStreams((prev) => (res.data ? res.data.streams : prev));
          const streaming = streams ? streams.some((stream) => stream.is_active) : false;
          setStreaming((prev) => (res.data ? streaming : prev));
        }
      } catch (err: any) {
        if (isMounted) {
          if (err?.response?.data) {
            setStreams([]);
            setStreaming(false);
            console.warn(err?.response?.data.error);
          }
        }
        if (err.response?.status == 401) {
          try {
            const res = await axios.post(
              "http://localhost/api/user/logout",
              {},
              { withCredentials: true }
            );
          } catch (err: any) {
            if (err?.response?.data) {
              console.warn(err?.response?.data.error)
            }
          }
          statusAuth.setAuthenticated(false);
          navigate("/login")
        }
      }

    };

    // Initial fetch
    fetchStatus();

    // Set up interval
    const interval = setInterval(fetchStatus, intervalMs);

    // Cleanup
    return () => {
      isMounted = false;
      clearInterval(interval);
    };
  }, [intervalMs]);

  React.useEffect(() => {
    const streaming = streams ? streams.some((stream) => stream.is_active) : false;
    setStreaming(streaming);
  }, [streams]);

  return { streamsState, streaming };
}
