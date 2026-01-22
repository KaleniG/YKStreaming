import * as React from "react";
import axios from "axios";
import { useUserStreams } from "../hooks/useUserStreams";
import { FaPlus } from "react-icons/fa";
import useAuth from "@/core/AuthContext";
import { useNavigate } from "react-router-dom";
import { StatusIndicator } from "@/components/StatusIndicator";
import { StreamInfoRow } from "@/components/StreamInfoRow";
import { AddStreamModal } from "@/components/AddStreamModal";

const UserStreams: React.FC = () => {
  const statusAuth = useAuth();
  const navigate = useNavigate();

  const { streamsState, streaming } = useUserStreams();
  const [streams, setStreams] = streamsState;
  const [showAddModal, setShowAddModal] = React.useState(false);

  const deleteStream = (keyToRemove: string) => {
    const fetchRemoveStreamKey = async () => {
      try {
        const res = await axios.post(
          `http://localhost/api/user/streams/remove/${keyToRemove}`,
          {},
          { withCredentials: true }
        );
        if (res.status == 200) {
          setStreams((prev) =>
            prev.filter((stream) => stream.key !== keyToRemove)
          );
        }
      } catch (err: any) {
        if (err?.response?.data) {
          console.warn(err?.response?.data.error);
        }
        if (err.response?.status == 401) {
          try {
            const res = await axios.post(
              "http://localhost/api/user/logout",
              {},
              { withCredentials: true }
            );
          } catch (err: any) {
            if (err?.response?.data) {
              console.warn(err?.response?.data.error)
            }
          }
          statusAuth.setAuthenticated(false);
          navigate("/login")
          statusAuth.setAuthenticated(false)
          navigate("/login")
        }
      }
    };
    fetchRemoveStreamKey();
  };

  const stopStream = (keyToStop: string) => {
    const fetchStopStream = async () => {
      try {
        const res = await axios.post(
          `http://localhost/api/user/streams/stop/${keyToStop}`,
          {},
          { withCredentials: true }
        );
        if (res.status == 200) {
          setStreams((prev) =>
            prev.map((stream) => {
              if (stream.key == keyToStop) {
                stream.active = false;
              }
              return stream;
            })
          );
        }
      } catch (err: any) {
        if (err?.response?.data) {
          console.warn(err?.response?.data.error);
        }
        if (err.response?.status == 401) {
          try {
            await axios.post(
              "http://localhost/api/user/logout",
              {},
              { withCredentials: true }
            );
          } catch (err: any) {
            if (err?.response?.data) {
              console.warn(err?.response?.data.error)
            }
          }
          statusAuth.setAuthenticated(false);
          navigate("/login");
        }
      }
    };
    fetchStopStream();
  };

  const addStreamRecord = (key: string, name: string, isVOD: boolean) => {
    setStreams((prev) => {
      if (prev) {
        return [
          ...prev,
          {
            key: key,
            active: false,
            name: name,
            is_vod: isVOD,
            total_views: 0,
            live_viewers: 0
          },
        ]
      } else {
        return [
          {
            key: key,
            active: false,
            name: name,
            is_vod: isVOD,
            total_views: 0,
            live_viewers: 0
          },
        ]
      }
    });
  }

  return (
    <>
      {showAddModal && (
        <AddStreamModal setShowModal={setShowAddModal} addStream={addStreamRecord} />
      )}
      <div className="flex flex-col h-full flex-1 bg-gradient-to-b from-zinc-100 via-zinc-50 to-zinc-100">
        {/* Main Container without rounded borders */}
        <div className="h-full bg-gradient-to-b from-zinc-100 to-zinc-200 border border-zinc-300 shadow-md p-6 m-6">
          <h2 className="text-2xl font-semibold text-zinc-700 mb-6 text-left select-none">
            Live Stream Info
          </h2>

          {/* Streaming Indicator */}
          <div className="flex items-center mb-4 gap-2">
            <StatusIndicator
              className={streaming ? "bg-red-500" : "bg-gray-400"}
            />
            <span className="text-zinc-700 font-medium">
              {streaming ? "Live Streaming" : "Offline"}
            </span>
          </div>

          {/* Stream Keys Header with "+" button */}
          <div className="flex items-center justify-between mb-4">
            <label className="text-zinc-700 font-medium select-none flex items-center">
              Stream Keys
              <button
                onClick={() => setShowAddModal(true)}
                className="h-6 px-2 text-sm rounded bg-gradient-to-b from-zinc-100 to-zinc-300 border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)] hover:from-zinc-200 hover:to-zinc-400 transition select-none mx-3 mt-1"
              >
                <FaPlus size={10} />
              </button>
            </label>
          </div>

          {/* Stream Keys List */}
          {streams && streams.length > 0 && (
            <div className="max-h-72 w-full overflow-y-auto border border-zinc-400 pt-3 pr-3 pl-3 pb-1">
              {streams.map((stream) => (
                <StreamInfoRow stream={stream} handleDeleteStream={deleteStream} handleStopStream={stopStream} />
              ))}
            </div>
          )}
        </div>
      </div>
    </>
  );
};

export default UserStreams;
