import { createBrowserRouter } from "react-router-dom";
import Auth from "@/features/auth/auth";
import Dashboard from "@/features/dashboard/Dashboard";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import { ROUTES } from "@/types/routes";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Auth />,
  },
  {
    path: ROUTES.dashboard,
    element: (
      <ProtectedRoute>
        <Dashboard />
      </ProtectedRoute>
    ),
  },
]);
