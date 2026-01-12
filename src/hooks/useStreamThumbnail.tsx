import * as React from "react";

export function useStreamThumbnail(
  thmbnailUrl: string,
  intervalMs: number = 45000
) {
  const [exists, setExists] = React.useState<boolean>(false);

  React.useEffect(() => {
    let isMounted = true;

    const fetchStreams = async () => {
      const res = new Promise((resolve) => {
        const img = new Image();

        img.onload = () => resolve(true);
        img.onerror = () => resolve(false);

        img.src = thmbnailUrl;
      });

      res.then((res) => {
        if (isMounted) {
          setExists((prev) => (res ? true : prev));
        }
      });
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

  return { exists };
}
