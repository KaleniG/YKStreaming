import useAuth from "@/core/AuthContext";
import axios from "axios";
import * as React from "react";
import { IoDocumentAttachOutline } from "react-icons/io5";
import { useNavigate } from "react-router-dom";

type Props = {
  setShowModal: (boolean) => void;
  addStream: (key: string, name: string, isVOD: boolean) => void;
};

export const AddStreamModal: React.FC<Props> = ({ setShowModal, addStream }) => {
  const statusAuth = useAuth();
  const navigate = useNavigate();

  const [newStreamName, setNewStreamName] = React.useState("");
  const [newThumbnail, setNewThumbnail] = React.useState<File | null>(null);
  const [isVod, setIsVod] = React.useState(false);

  const hiddenFileInputRef = React.useRef(null);

  const handleAddStream = async () => {
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
        addStream(res.data.key, newStreamName, isVod)
        setShowModal(false);
        setNewStreamName("");
        setNewThumbnail(null);
        setIsVod(false);
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
      }
    }
  };

  const inputBaseStyle =
    "caret-zinc-500 selection:bg-zinc-300 selection:text-black h-9 rounded-md bg-gradient-to-b from-white to-zinc-200 border px-3 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.2)] focus:outline-none focus:ring-1 border-zinc-400";

  return (
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
            onClick={() => hiddenFileInputRef?.current?.click()}
          />
          <button
            type="button"
            className="absolute right-2 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
            onClick={() => hiddenFileInputRef?.current?.click()}
          >
            <IoDocumentAttachOutline size={20} />
          </button>
          <input
            id="hiddenFileInput"
            type="file"
            ref={hiddenFileInputRef}
            className="hidden"
            accept=".jpg,.jpeg,.png,.webp,.gif,.svg"
            onChange={(e) => setNewThumbnail(e.target.files ? e.target.files[0] : null)}
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
            onClick={() => setShowModal(false)}
            className="h-8 px-4 text-sm rounded bg-gradient-to-b from-zinc-100 to-zinc-300 border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)] hover:from-zinc-200 hover:to-zinc-400 transition select-none"
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={handleAddStream}
            className="h-8 px-4 text-sm rounded bg-gradient-to-b from-zinc-100 to-zinc-300 border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)] hover:from-zinc-200 hover:to-zinc-400 transition select-none"
          >
            Add
          </button>
        </div>
      </div>
    </div>
  );
};
