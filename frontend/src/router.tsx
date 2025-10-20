import { createBrowserRouter } from "react-router-dom";
import Auth from "@/features/auth/auth";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Auth />,
  },
]);
