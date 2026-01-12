import * as React from "react";
import { useTrendingStreams } from "../hooks/useTrendingStreams";
import VideoThumbnail from "@/components/VideoThumbnail";

const Home: React.FC = () => {
  const { streams } = useTrendingStreams();

  return (
    <div className="flex flex-col h-full bg-gradient-to-b from-zinc-100 via-zinc-50 to-zinc-100 pt-6">
      {/* Title */}
      <div className="w-full max-w-6xl px-4 mb-6">
        <h1 className="text-1xl font-semibold text-zinc-600 select-none">
          Live on YKStreaming
        </h1>
      </div>

      {/* Grid of streams */}
      <div className="w-full max-w-6xl px-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {streams.map((stream) => (
            <VideoThumbnail stream={stream} />
          ))}
        </div>
      </div>
    </div>
  );
};

export default Home;
