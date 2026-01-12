import * as React from "react";
import axios from "axios";

export function useStreamerStatus(intervalMs: number = 5000) {
  const streamsState = React.useState<Array<any>>([]);
  const [streams, setStreams] = streamsState;
  const [streaming, setStreaming] = React.useState<boolean>(false);

  React.useEffect(() => {
    let isMounted = true;

    const fetchStatus = async () => {
      const res = await axios.post<{
        success: boolean;
        streams: Array<any>;
      }>(
        "http://localhost/api/get_streamer_status.php",
        {},
        { withCredentials: true }
      );

      if (isMounted) {
        setStreams((prev) => (res.data ? res.data.streams : prev));
        const streaming = res.data.streams.some((stream) => stream.active);
        setStreaming((prev) => (res.data ? streaming : prev));
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
    const streaming = streams.some((stream) => stream.active);
    setStreaming(streaming);
  }, [streams]);

  return { streamsState, streaming };
}
