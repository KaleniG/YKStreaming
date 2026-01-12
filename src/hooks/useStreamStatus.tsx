import * as React from "react";
import axios from "axios";

export function useStreamStatus(streamKey: string, intervalMs: number = 5000) {
  const [isLive, setIsLive] = React.useState(false);
  const [isVOD, setIsVOD] = React.useState(false);
  const [exists, setExists] = React.useState(false);

  React.useEffect(() => {
    let isMounted = true;

    const fetchStatus = async () => {
      try {
        const res = await axios.post<{
          success: boolean;
          exists: boolean;
          is_live: boolean;
          is_vod: boolean;
        }>(
          "http://localhost/api/get_stream_status.php",
          { stream_key: streamKey },
          { withCredentials: true }
        );

        if (isMounted && res.data) {
          setExists(res.data.exists); // Always update
          setIsLive(res.data.exists && res.data.is_live); // Update based on exists & active
          setIsVOD(res.data.exists && res.data.is_vod); // Update based on exists & active
        }
      } catch (err) {
        console.error("Failed to fetch stream status:", err);
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
