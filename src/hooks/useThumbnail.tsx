import * as React from "react";

export function useThumbnail(
  thumbnailUrl: string,
  intervalMs: number = 5000
) {
  const [exists, setExists] = React.useState<boolean>(false);

  React.useEffect(() => {
    let isMounted = true;

    const checkThumbnail = async () => {
      const res = await new Promise<boolean>((resolve) => {
        const img = new Image();
        img.onload = () => resolve(true);
        img.onerror = () => resolve(false);

        img.src = thumbnailUrl;
      });

      if (isMounted) {
        setExists(res);
      }
    };

    checkThumbnail();
    const interval = setInterval(checkThumbnail, intervalMs);

    return () => {
      isMounted = false;
      clearInterval(interval);
    };
  }, [thumbnailUrl, intervalMs]);

  return { exists };
}
