import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useHealth } from "./queries/useHealth";
import { useUser } from "./queries/useUser";

export default function Dashboard() {
  const { user, session, signOut } = useAuth();
  const {
    data: health,
    isLoading: healthLoading,
    error: healthError,
  } = useHealth();
  const {
    data: userAccount,
    isLoading: userLoading,
    error: userError,
  } = useUser();

  const handleSignOut = async () => {
    try {
      await signOut();
    } catch (error) {
      console.error("Error signing out:", error);
    }
  };

  const accessToken = session?.access_token;

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <Button onClick={handleSignOut} variant="outline">
            Sign Out
          </Button>
        </div>

        <div className="space-y-6">
          {/* User Info Card */}
          <Card>
            <CardHeader>
              <CardTitle>Welcome to Reel Farm!</CardTitle>
              <CardDescription>
                You're successfully authenticated
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <p>
                  <strong>Email:</strong> {user?.email}
                </p>
                <p>
                  <strong>User ID:</strong> {user?.id}
                </p>
                <p>
                  <strong>Last Sign In:</strong>{" "}
                  {user?.last_sign_in_at
                    ? new Date(user.last_sign_in_at).toLocaleString()
                    : "N/A"}
                </p>
                <p>
                  <strong>Access Token:</strong>
                </p>
                <div className="bg-gray-100 p-2 rounded border font-mono text-sm break-all">
                  {accessToken || "N/A"}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* API Health Card */}
          <Card>
            <CardHeader>
              <CardTitle>API Health</CardTitle>
            </CardHeader>
            <CardContent>
              {healthLoading && (
                <p className="text-blue-600">Checking API health...</p>
              )}
              {healthError && (
                <p className="text-red-600">Error: {healthError.message}</p>
              )}
              {health && (
                <div className="bg-green-50 p-3 rounded border">
                  <p>
                    <strong>Status:</strong> {health.status}
                  </p>
                  <p>
                    <strong>Message:</strong> {health.message}
                  </p>
                  <p>
                    <strong>Port:</strong> {health.port}
                  </p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* User Account Card */}
          <Card>
            <CardHeader>
              <CardTitle>User Account (from API)</CardTitle>
            </CardHeader>
            <CardContent>
              {userLoading && (
                <p className="text-blue-600">Loading user account...</p>
              )}
              {userError && (
                <p className="text-red-600">Error: {userError.message}</p>
              )}
              {userAccount && (
                <div className="bg-green-50 p-3 rounded border">
                  <p>
                    <strong>Plan:</strong> {userAccount.plan}
                  </p>
                  <p>
                    <strong>Plan Started:</strong>{" "}
                    {new Date(userAccount.plan_started_at).toLocaleDateString()}
                  </p>
                  <p>
                    <strong>Created:</strong>{" "}
                    {new Date(userAccount.created_at).toLocaleDateString()}
                  </p>
                  {userAccount.plan_ends_at && (
                    <p>
                      <strong>Plan Ends:</strong>{" "}
                      {new Date(userAccount.plan_ends_at).toLocaleDateString()}
                    </p>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
