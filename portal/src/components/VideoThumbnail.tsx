import * as React from "react";
import { useThumbnail as useThumbnail } from "@/hooks/useThumbnail";
import { Link } from "react-router-dom";

import DefaultThumbnail from "../assets/default_thumbnail.jpg";

interface Stream {
  streamer_name: string;
  key: string;
  name: string;
  has_custom_thumbnail: boolean;
  is_live: boolean;
  is_vod: boolean;
  live_viewers: number;
}

interface VideoThumbnailProps {
  stream: Stream;
}

const VideoThumbnail: React.FC<VideoThumbnailProps> = ({ stream }) => {
  const liveThumbnail = useThumbnail(`http://localhost/thumbnails/live/${stream.key}.jpg`);
  const [hovered, setHovered] = React.useState(false);

  const thumbnailSrc =
    hovered && liveThumbnail.exists
      ? `http://localhost/thumbnails/live/${stream.key}.jpg`
      : stream.has_custom_thumbnail
        ? `http://localhost/thumbnails/custom/${stream.key}.jpg`
        : DefaultThumbnail;

  return (
    <Link
      key={stream.key}
      to={`/stream/${stream.key}`}
      className="group block rounded-md overflow-hidden bg-gradient-to-b from-zinc-100 to-zinc-200 shadow-[0_2px_4px_rgba(0,0,0,0.1)] hover:shadow-[0_4px_6px_rgba(0,0,0,0.15)] transition-all duration-200 select-none mb-6"
    >
      <img
        src={thumbnailSrc}
        onMouseEnter={() => setHovered(true)}
        onMouseLeave={() => setHovered(false)}
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
        <h4
          className={`flex items-center text-sm ${stream.is_live ? "visible" : "invisible"
            }`}
        >
          <span className="w-3 h-3 rounded-full mr-2 bg-red-400 inline-block" />
          {`${stream.live_viewers ?? 0} viewers`}
        </h4>
      </div>
    </Link>
  );
};

export default VideoThumbnail;
