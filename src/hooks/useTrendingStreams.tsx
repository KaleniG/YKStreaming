import * as React from "react";
import axios from "axios";

export function useTrendingStreams(intervalMs: number = 45000) {
  const [streams, setStreams] = React.useState<Array<any>>([]);

  React.useEffect(() => {
    let isMounted = true;

    const fetchStreams = async () => {
      const res = await axios.post<{
        success: boolean;
        streams: Array<any>;
      }>(
        "http://localhost/api/get_trending_streams.php",
        {},
        { withCredentials: true }
      );

      if (isMounted) {
        setStreams((prev) => (res.data ? res.data.streams : prev));
      }
    };

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
