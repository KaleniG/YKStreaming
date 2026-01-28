import * as React from "react";
import AuthRouter from "./AuthRouter";
import { PageLoadingProvider } from "./LoadingContext";
import { AuthProvider } from "./AuthContext";
import LoadingGate from "./LoadingGate";
import Bar from "../components/Nav/Bar";

const App = () => {
  return (
    <React.StrictMode>
      <PageLoadingProvider>
        <AuthProvider>
          {/* Root container: full height flex column */}
          <div className="flex flex-col h-screen bg-white overflow-hidden">
            {/* Navbar stays at the top */}
            <Bar />

            {/* Scrollable content only */}
            <div className="flex-1 overflow-auto">
              <LoadingGate>
                <AuthRouter />
              </LoadingGate>
            </div>
          </div>
        </AuthProvider>
      </PageLoadingProvider>
    </React.StrictMode>
  );
};

export default App;
