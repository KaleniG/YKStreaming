import * as React from "react";
import { usePublicStreams } from "../hooks/usePublicStreams";
import VideoThumbnail from "@/components/VideoThumbnail";

const Home: React.FC = () => {
  const { streams } = usePublicStreams();

  return (
    <div className="flex flex-col h-full pt-6">
      <div className="w-full max-w-8xl px-4">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {streams.map((stream, i) => (
            <VideoThumbnail stream={stream} key={i} />
          ))}
        </div>
      </div>
    </div>
  );
};

export default Home;
