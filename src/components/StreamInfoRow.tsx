import { UserStreamData } from "@/hooks/useUserStreams";
import * as React from "react";
import { StatusIndicator } from "./StatusIndicator";
import { TextField } from "./TextField";
import { CopyableField } from "./CopyableField";
import { StreamButton } from "./StreamButton";
import { AiOutlineClose } from "react-icons/ai";
import { RxStop } from "react-icons/rx";
import { MdOutlineOpenInNew } from "react-icons/md";

type Props = {
  stream: UserStreamData;
  handleDeleteStream: (string) => void;
  handleStopStream: (string) => void;
};

export const StreamInfoRow: React.FC<Props> = ({ stream, handleDeleteStream, handleStopStream }) => {
  const openStream = React.useCallback(() => {
    const baseUrl = window.location.origin;
    window.open(`${baseUrl}/stream/${stream.key}`, "_blank");
  }, [stream.key]);

  return (
    <div className="flex" key={stream.key}>
      <div className="inline-flex items-center mb-2 gap-2">
        <StatusIndicator
          className={stream.is_active
            ? "bg-green-500"
            : stream.ended_at
              ? "bg-red-400"
              : "bg-gray-400"
          }
        />
        <TextField value={stream.name} className={"w-[200px]"} />
        {stream.started_at ?
          <>
            <TextField value={`Views: ${stream.total_views}`} className={"w-[150px]"} />
            {!stream.ended_at ?
              <TextField
                value={`Live viewers: ${stream.live_viewers ? stream.live_viewers : 0}`}
                className={"w-[150px]"}
              /> : null
            }
          </> : null
        }
        <CopyableField value={stream.key} className="w-[415px]" />
        <CopyableField value={"rtmp://localhost/live"} className="w-[200px]" />
      </div>
      <StreamButton onClick={() => handleDeleteStream(stream.key)}>
        <AiOutlineClose size={18} />
      </StreamButton>
      {stream.is_active ? (
        <>
          <StreamButton onClick={() => handleStopStream(stream.key)}>
            <RxStop size={18} />
          </StreamButton>
          <StreamButton onClick={openStream}>
            <MdOutlineOpenInNew size={18} />
          </StreamButton>
        </>
      ) : null}
    </div>
  );
};
