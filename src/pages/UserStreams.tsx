import * as React from "react";
import axios from "axios";
import { useStreamerStatus } from "../hooks/useStreamerStatus";
import { AiOutlineClose } from "react-icons/ai";
import { RxStop } from "react-icons/rx";
import { MdOutlineOpenInNew } from "react-icons/md";
import { FaPlus } from "react-icons/fa";
import {
  IoCopyOutline,
  IoCopy,
  IoDocumentAttachOutline,
} from "react-icons/io5";
import useAuth from "@/core/AuthContext";
import { useNavigate } from "react-router-dom";

const UserStreams: React.FC = () => {
  const statusAuth = useAuth();
  const navigate = useNavigate();

  const { streamsState, streaming } = useStreamerStatus();
  const [streams, setStreams] = streamsState;

  const [copiedKeys, setCopiedKeys] = React.useState<Record<string, boolean>>(
    {}
  );

  const [showAddModal, setShowAddModal] = React.useState(false);
  const [newStreamName, setNewStreamName] = React.useState("");
  const [newThumbnail, setNewThumbnail] = React.useState<File | null>(null);
  const [isVod, setIsVod] = React.useState(false);

  const copyToClipboard = (text: string, key: string) => {
    navigator.clipboard.writeText(text);
    setCopiedKeys((prev) => ({ ...prev, [key]: true }));
    setTimeout(
      () => setCopiedKeys((prev) => ({ ...prev, [key]: false })),
      2000
    );
  };

  const removeStreamKey = (keyToRemove: string) => {
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
          setCopiedKeys((prev) => {
            const copy = { ...prev };
            delete copy[keyToRemove];
            return copy;
          });
        }
      } catch (err: any) {
        if (err?.response?.data) {
          console.warn(err?.response?.data.error);
        }
        if (err.response?.status == 401) {
          statusAuth.setAuthenticated(false)
          navigate("/login")
        }
      }
    };
    fetchRemoveStreamKey();
  };

  const openStream = (streamKey: string) => {
    const baseUrl = window.location.origin;
    window.open(`${baseUrl}/stream/${streamKey}`, "_blank");
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
          statusAuth.setAuthenticated(false)
          navigate("/login")
        }
      }
    };
    fetchStopStream();
  };

  const inputBaseStyle =
    "w-[550px] caret-zinc-500 selection:bg-zinc-300 selection:text-black h-9 rounded-md bg-gradient-to-b from-white to-zinc-200 border px-3 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.2)] focus:outline-none focus:ring-1 border-zinc-400";

  const buttonBaseStyle =
    "mb-2 full-h px-4 text-sm bg-gradient-to-b from-zinc-100 to-zinc-300 border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)] hover:from-zinc-200 hover:to-zinc-400 transition select-none";

  return (
    <>
      {showAddModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 shadow-md">
          <div className="bg-white rounded-lg shadow-lg p-6 w-[400px] relative border border-zinc-300">
            <h3 className="text-lg font-semibold mb-4">Add New Stream</h3>

            <label
              className="block mb-2 font-medium text-zinc-700"
              htmlFor="stream_name"
            >
              Stream Name
            </label>
            <input
              id="stream_name"
              type="text"
              value={newStreamName}
              onChange={(e) => setNewStreamName(e.target.value)}
              className={inputBaseStyle + " mb-4 w-full focus:ring-zinc-500"}
              autoCorrect="off"
              autoCapitalize="off"
              spellCheck={false}
            />

            <label className="block mb-2 font-medium text-zinc-700 select-none pt-2">
              Thumbnail
            </label>
            <div className="relative w-full mb-4 ">
              <input
                type="text"
                value={newThumbnail ? newThumbnail.name : ""}
                placeholder="Stream screenshots by default"
                disabled
                className={`${inputBaseStyle} w-full pr-10 bg-white`}
                onClick={() =>
                  document.getElementById("hiddenFileInput")?.click()
                }
              />
              <button
                type="button"
                className="absolute right-2 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
                onClick={() =>
                  document.getElementById("hiddenFileInput")?.click()
                }
              >
                <IoDocumentAttachOutline size={20} />
              </button>
              <input
                id="hiddenFileInput"
                type="file"
                className="hidden"
                accept=".jpg,.jpeg,.png,.webp,.gif,.svg"
                onChange={(e) =>
                  setNewThumbnail(e.target.files ? e.target.files[0] : null)
                }
              />
            </div>

            <label className="flex items-center mb-100 cursor-pointer select-none mt-6">
              <input
                type="checkbox"
                checked={isVod}
                onChange={(e) => setIsVod(e.target.checked)}
                className="mr-2 mt-1 accent-zinc-700 cursor-pointer"
              />
              <span className="text-zinc-700 text-sm">
                At stream end convert into a VOD
              </span>
            </label>

            <div className="flex justify-end space-x-2 pt-6">
              <button
                onClick={() => setShowAddModal(false)}
                className="h-8 px-4 text-sm rounded bg-gradient-to-b from-zinc-100 to-zinc-300 border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)] hover:from-zinc-200 hover:to-zinc-400 transition select-none"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={async () => {
                  if (!newStreamName) return;

                  const formData = new FormData();
                  formData.append("name", newStreamName);
                  formData.append("is_vod", isVod ? "1" : "0");
                  if (newThumbnail) formData.append("thumbnail", newThumbnail);

                  try {
                    const res = await axios.post(
                      "http://localhost/api/user/streams/add",
                      formData,
                      {
                        withCredentials: true,
                        headers: { "Content-Type": "multipart/form-data" },
                      }
                    );

                    if (res.data) {
                      setStreams((prev) => {
                        if (prev) {
                          return [
                            ...prev,
                            {
                              key: res.data.key,
                              active: false,
                              name: newStreamName,
                              is_vod: isVod,
                              total_views: 0,
                              live_viewers: 0
                            },
                          ]
                        } else {
                          return [
                            {
                              key: res.data.key,
                              active: false,
                              name: newStreamName,
                              is_vod: isVod,
                              total_views: 0,
                              live_viewers: 0
                            },
                          ]
                        }
                      });
                      setShowAddModal(false);
                      setNewStreamName("");
                      setNewThumbnail(null);
                      setIsVod(false);
                    }
                  } catch (err: any) {
                    if (err?.response?.data) {
                      console.warn(err?.response?.data.error);
                    }
                    if (err.response?.status == 401) {
                      statusAuth.setAuthenticated(false)
                      navigate("/login")
                    }
                  }
                }}
                className="h-8 px-4 text-sm rounded bg-gradient-to-b from-zinc-100 to-zinc-300 border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)] hover:from-zinc-200 hover:to-zinc-400 transition select-none"
              >
                Add
              </button>
            </div>
          </div>
        </div>
      )}
      <div className="flex flex-col h-full flex-1 bg-gradient-to-b from-zinc-100 via-zinc-50 to-zinc-100">
        {/* Main Container without rounded borders */}
        <div className="h-full bg-gradient-to-b from-zinc-100 to-zinc-200 border border-zinc-300 shadow-md p-6 m-6">
          <h2 className="text-2xl font-semibold text-zinc-700 mb-6 text-left select-none">
            Live Stream Info
          </h2>

          {/* Streaming Indicator */}
          <div className="flex items-center mb-4">
            <span
              className={`w-3 h-3 rounded-full mr-2 ${streaming ? "bg-red-500" : "bg-gray-400"
                }`}
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
                <div className="flex" key={stream.key}>
                  <div className="inline-flex items-center relative mb-2 mr-2">
                    <span
                      className={`mr-3 w-3 h-3 rounded-full mr-2 ${stream.is_active
                        ? "bg-green-500"
                        : stream.ended_at
                          ? "bg-red-400"
                          : "bg-gray-400"
                        }`}
                    />
                    <input
                      type="text"
                      value={stream.name}
                      disabled
                      className={`${inputBaseStyle} rounded-l-md w-[125px] mr-2`}
                    />
                    <input
                      type="text"
                      value={`Views: ${stream.total_views}`}
                      disabled
                      className={`${inputBaseStyle} rounded-l-md w-[125px] mr-2`}
                    />
                    <input
                      type="text"
                      value={`Live viewers: ${stream.live_viewers ? stream.live_viewers : 0
                        }`}
                      disabled
                      className={`${inputBaseStyle} rounded-l-md w-[125px] mr-2`}
                    />
                    <input
                      type="text"
                      value={stream.key}
                      disabled
                      className={`${inputBaseStyle} rounded-l-md`}
                    />
                    <button
                      onClick={() => copyToClipboard(stream.key, stream.key)}
                      className="absolute right-2 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
                    >
                      {copiedKeys[stream.key] ? (
                        <IoCopy size={18} />
                      ) : (
                        <IoCopyOutline size={18} />
                      )}
                    </button>
                  </div>
                  <div className="inline-flex items-center relative mb-2">
                    <input
                      type="text"
                      value="rtmp://localhost/live"
                      disabled
                      className={`${inputBaseStyle} rounded-l-md w-[200px]`}
                    />
                    <button
                      onClick={() =>
                        copyToClipboard(
                          stream.is_vod
                            ? "rtmp://localhost/vodlive"
                            : "rtmp://localhost/live",
                          "link"
                        )
                      }
                      className="absolute right-2 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
                    >
                      {copiedKeys["link"] ? (
                        <IoCopy size={18} />
                      ) : (
                        <IoCopyOutline size={18} />
                      )}
                    </button>
                  </div>
                  <button
                    onClick={() => removeStreamKey(stream.key)}
                    className={`${buttonBaseStyle} ml-2 bg-gray-300 text-gray-800 hover:bg-gray-400`}
                  >
                    <AiOutlineClose size={18} />
                  </button>
                  {stream.is_active ? (
                    <>
                      <button
                        onClick={() => stopStream(stream.key)}
                        className={`${buttonBaseStyle} ml-2 bg-gray-300 text-gray-800 hover:bg-gray-400`}
                      >
                        <RxStop size={18} />
                      </button>
                      <button
                        onClick={() => openStream(stream.key)}
                        className={`${buttonBaseStyle} ml-2 bg-gray-300 text-gray-800 hover:bg-gray-400`}
                      >
                        <MdOutlineOpenInNew size={18} />
                      </button>
                    </>
                  ) : null}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </>
  );
};

export default UserStreams;
