import { Card, CardContent } from "@/components/ui/card";
import type { HealthResponse } from "@/api/models/HealthResponse";

interface HealthStatusProps {
  health: HealthResponse | undefined;
}

export function HealthStatus({ health }: HealthStatusProps) {
  if (!health) return null;

  return (
    <div className="mb-6">
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-green-500 rounded-full"></div>
            <span className="text-sm text-gray-600">
              {health.message} - Port: {health.port}
            </span>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
