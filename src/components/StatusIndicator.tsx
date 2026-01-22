import * as React from "react";

type Props = {
  className: string;
};

export const StatusIndicator: React.FC<Props> = ({ className }) => {
  return (
    <span
      className={`w-3 h-3 rounded-full ${className}`}
    />
  );
};
