import * as React from "react";
import axios from "axios";

interface GetStreamStatusResponse {
  stream: { is_vod: boolean; is_live: boolean; };
}

export function useStreamStatus(streamKey: string, intervalMs: number = 5000) {
  const [isLive, setIsLive] = React.useState(false);
  const [isVOD, setIsVOD] = React.useState(false);
  const [exists, setExists] = React.useState(false);

  React.useEffect(() => {
    let isMounted = true;

    const fetchStatus = async () => {
      try {
        const res = await axios.post<GetStreamStatusResponse>(
          `http://localhost/api/stream/${streamKey}`,
          {},
          { withCredentials: true }
        );

        if (isMounted && res.data) {
          setExists(res.data.stream.is_live || res.data.stream.is_vod);
          setIsLive(res.data.stream.is_live);
          setIsVOD(res.data.stream.is_vod);
        }
      } catch (err: any) {
        if (err?.response?.data) {
          if (isMounted) {
            setExists(false);
            setIsLive(false);
            setIsVOD(false);
            console.warn(err?.response?.data.error);
          }
        }
      }
    };

    // Initial fetch
    fetchStatus();

    // Poll interval
    const interval = setInterval(fetchStatus, intervalMs);

    // Cleanup
    return () => {
      isMounted = false;
      clearInterval(interval);
    };
  }, [streamKey, intervalMs]); // ✅ Add streamKey to dependencies

  return { exists, isLive, isVOD };
}
