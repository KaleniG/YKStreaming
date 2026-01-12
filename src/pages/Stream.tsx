import * as React from "react";
import { useParams } from "react-router-dom";
import VideoPlayer from "../components/VideoPlayer";
import { useStreamStatus } from "@/hooks/useStreamStatus";

const Stream: React.FC = () => {
  const { streamKey } = useParams<{ streamKey: string }>();
  const { exists, isLive, isVOD } = useStreamStatus(streamKey);

  if (!streamKey) {
    return <div className="pt-16 text-zinc-700">Invalid stream</div>;
  }

  const sources = React.useMemo(() => {
    if (isLive) {
      return [
        {
          src: `http://localhost/hls/${streamKey}.m3u8`,
          type: "application/x-mpegURL",
        },
      ];
    }

    if (isVOD) {
      return [
        {
          src: `http://localhost/recordings/${streamKey}.mp4`,
          type: "video/mp4",
        },
      ];
    }

    return [];
  }, [streamKey, isLive, isVOD]);

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
              player.on("playing", () => console.log("🔴 LIVE"));
              player.on("waiting", () => console.log("⏳ BUFFERING"));
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
