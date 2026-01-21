import * as React from "react";
import { useParams } from "react-router-dom";
import axios from "axios";
import VideoPlayer from "../components/VideoPlayer";
import { useStreamStatus } from "@/hooks/useStreamStatus";

const Stream: React.FC = () => {
  const { streamKey } = useParams<{ streamKey: string }>();
  const { exists, isLive, isVOD } = useStreamStatus(streamKey);

  const viewerRegisteredRef = React.useRef(false);

  if (!streamKey) {
    return <div className="pt-16 text-zinc-700">Invalid stream</div>;
  }

  const sources = React.useMemo(() => {
    if (isLive) {
      return [
        {
          src: `http://localhost/streams/${streamKey}.m3u8`,
          type: "application/x-mpegURL",
        },
      ];
    } else if (isVOD) {
      return [
        {
          src: `http://localhost/vods/${streamKey}.mp4`,
          type: "video/mp4",
        },
      ];
    }

    return [];
  }, [streamKey, isLive, isVOD]);

  React.useEffect(() => {
    if (!streamKey || !isLive) return;

    const registerViewer = async () => {
      try {
        const res = await axios.post(
          `http://localhost/api/stream/view/${streamKey}`,
          {},
          { withCredentials: true }
        );

        viewerRegisteredRef.current = (res.status == 200);
      } catch (err: any) {
        if (err?.response?.data) {
          viewerRegisteredRef.current = false;
          console.warn(err?.response?.data.error);
        }
      }
    };

    registerViewer();

    return () => {
      if (viewerRegisteredRef.current) {
        try {
          void axios.post(
            `http://localhost/api/stream/unview/${streamKey}`,
            {},
            { withCredentials: true }
          );
        } catch (err: any) {
          if (err?.response?.data) {
            viewerRegisteredRef.current = false;
            console.warn(err?.response?.data.error);
          }
        }
      }
    };
  }, [streamKey, isLive]);

  const videoJsOptions = React.useMemo(
    () => ({
      autoplay: true,
      controls: true,
      responsive: true,
      liveui: isLive,
      fluid: true,
      sources,
    }),
    [sources, isLive]
  );

  return (
    <div className="flex flex-col h-full w-full bg-gradient-to-b from-zinc-100 via-zinc-50 to-zinc-100 items-center justify-center p-6">
      {!exists ? (
        <h1 className="text-zinc-700 text-lg">This stream does not exist</h1>
      ) : !isLive && !isVOD ? (
        <h1 className="text-zinc-700 text-lg">
          This has ended or has not started yet
        </h1>
      ) : (
        <div className="w-full max-w-4xl aspect-video bg-black overflow-hidden border border-zinc-300">
          <VideoPlayer
            options={videoJsOptions}
            onReady={(player) => {
              player.on("error", () => {
                setTimeout(() => {
                  player.src(videoJsOptions.sources);
                  player.play();
                }, 3000);
              });
            }}
          />
        </div>
      )}
    </div>
  );
};

export default Stream;
