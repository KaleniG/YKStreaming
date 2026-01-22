import * as React from "react";
import axios from "axios";

interface GetStreamsResponse {
  streams: null | {
    streamer_name: string;
    key: string;
    name: string;
    has_custom_thumbnail: boolean;
    is_live: boolean;
    is_vod: boolean;
    live_viewers: number;
  }[];
}

export function usePublicStreams(intervalMs: number = 45000) {
  const [streams, setStreams] = React.useState<Array<any>>([]);

  React.useEffect(() => {
    let isMounted = true;

    const fetchStreams = async () => {
      try {
        const res = await axios.post<GetStreamsResponse>(
          "http://localhost/api/get-streams",
          {},
          { withCredentials: true }
        );

        if (isMounted) {
          if (res.data) {
            if (!res.data.streams) {
              setStreams([]);
            } else {
              setStreams(res.data.streams);
            }
          }
        }
      } catch (err: any) {
        if (err?.response?.data) {
          if (isMounted) {
            setStreams([]);
            console.warn(err?.response?.data.error);
          }
        }
      }
    }

    // Initial fetch
    fetchStreams();

    // Set up interval
    const interval = setInterval(fetchStreams, intervalMs);

    // Cleanup
    return () => {
      isMounted = false;
      clearInterval(interval);
    };
  }, [intervalMs]);

  return { streams };
}
