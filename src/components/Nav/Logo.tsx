import * as React from "react";
import { Link } from "react-router-dom";

import YKStreamingLogo from "../../assets/logo.png";

const Logo: React.FC = () => {
  return (
    <Link
      to="/"
      className="group h-full px-4 flex items-center gap-2 bg-transparent"
    >
      <img
        src={YKStreamingLogo}
        alt="YKStreamingLogo"
        className="h-10 filter drop-shadow-[0_1px_1px_rgba(0,0,0,0.35)] group-hover:drop-shadow-[0_2px_3px_rgba(0,0,0,0.45)] transition select-none"
      />
    </Link>
  );
};

export default React.memo(Logo);
