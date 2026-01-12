import * as ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./core/App";
import "./index.css";
import reportWebVitals from "./reportWebVitals";

const rootElement = document.getElementById("root");

if (!rootElement) {
  throw new Error("Root container missing in index.html");
}

ReactDOM.createRoot(rootElement).render(
  <BrowserRouter>
    <App />
  </BrowserRouter>
);

reportWebVitals();
