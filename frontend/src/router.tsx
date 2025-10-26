import { createBrowserRouter, Outlet } from "react-router-dom";
import Auth from "@/features/auth/auth";
import Dashboard from "@/features/dashboard/Dashboard";
import { ProtectedRoute } from "@/components/ProtectedRoute";
import { ROUTES } from "@/types/routes";
import { DashboardLayout } from "./components/DashboardLayout";
import { GenerateHooks } from "./features/hooks/generate-hooks/GenerateHooks";
import { YourVideos } from "./features/videos/user-generated-videos/YourVideos";
import { GenerateAiAvatarVideo } from "./features/videos/generate-ai-avatar-video/GenerateAiAvatarVideo";
import { UserGeneratedVideos } from "./features/videos/user-generated-videos/components/UserGeneratedVideos";

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
          <Outlet />
        </DashboardLayout>
      </ProtectedRoute>
    ),
    children: [
      {
        path: ROUTES.userGeneratedVideos,
        element: <UserGeneratedVideos />,
      },
      {
        path: ROUTES.generateAiAvatarVideo,
        element: <GenerateAiAvatarVideo />,
      },
      {
        path: ROUTES.generateHooks,
        element: <GenerateHooks />,
      },
    ],
  },
]);
