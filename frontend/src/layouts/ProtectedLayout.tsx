import { Outlet } from "react-router-dom";

export default function ProtectedLayout() {
  // TODO: Add authentication check
  return <Outlet />;
}
