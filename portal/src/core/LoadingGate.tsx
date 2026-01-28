import * as React from "react";

import Splash from "../pages/Splash";
import usePageLoading from "./LoadingContext";

const LoadingGate: React.FC<React.PropsWithChildren> = ({ children }) => {
  const page = usePageLoading();
  return page.isLoading ? <Splash /> : <>{children}</>;
};

export default LoadingGate;
