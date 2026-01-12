import * as React from "react";
import { useStreamThumbnail } from "@/hooks/useStreamThumbnail";
import { Link } from "react-router-dom";

import DefaultThumbnail from "../assets/default_thumbnail.jpg";

interface Stream {
  key: string;
  name: string;
  uses_thumbnail: boolean;
  thumbnail_format: string;
  is_live: boolean;
  is_vod: boolean;
  streamer_name: string;
}

interface VideoThumbnailProps {
  stream: Stream;
}

const VideoThumbnail: React.FC<VideoThumbnailProps> = ({ stream }) => {
  const { exists } = useStreamThumbnail(stream.key);

  return (
    <Link
      key={stream.key}
      to={`/stream/${stream.key}`}
      className="group block rounded-md overflow-hidden bg-gradient-to-b from-zinc-100 to-zinc-200 shadow-[0_2px_4px_rgba(0,0,0,0.1)] hover:shadow-[0_4px_6px_rgba(0,0,0,0.15)] transition-all duration-200 select-none"
    >
      <img
        src={
          stream.uses_thumbnail
            ? `http://localhost/thumbnails/${stream.key}.${stream.thumbnail_format}`
            : exists
            ? `http://localhost/"stream_screenshots/${stream.key}.jpg`
            : DefaultThumbnail
        }
        alt={`${stream.streamer_name} screenshot`}
        className="w-full h-60 object-cover group-hover:brightness-105 transition duration-200"
      />
      <div className="p-3 bg-zinc-100">
        <h2 className="text-sm font-medium text-zinc-900 truncate">
          {stream.name}
        </h2>
        <h3 className="text-sm font-medium text-zinc-900 truncate">
          {stream.streamer_name}
        </h3>
        <h4>{stream.is_live ? "Live" : stream.is_vod ? "VOD" : null}</h4>
      </div>
    </Link>
  );
};

export default VideoThumbnail;
