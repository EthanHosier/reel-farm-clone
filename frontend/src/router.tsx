import { createBrowserRouter } from "react-router-dom";
import Auth from "@/features/auth/auth";
import Dashboard from "@/features/dashboard/Dashboard";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import { ROUTES } from "@/types/routes";
import { DashboardLayout } from "./components/DashboardLayout";

export const router = createBrowserRouter([
  {
    path: ROUTES.auth,
    element: <Auth />,
  },
  {
    path: ROUTES.dashboard,
    element: (
      <ProtectedRoute>
        <DashboardLayout>
          <Dashboard />
        </DashboardLayout>
      </ProtectedRoute>
    ),
  },
]);
