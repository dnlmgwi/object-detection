import type { PossessionStats, PerformanceMetrics } from '~/hooks/useWebSocket';

interface StatsDashboardProps {
  possessionStats: PossessionStats | null;
  metrics: PerformanceMetrics | null;
  teamAName: string;
  teamBName: string;
  teamAColor: string;
  teamBColor: string;
}

export const StatsDashboard = ({
  possessionStats,
  metrics,
  teamAName,
  teamBName,
  teamAColor,
  teamBColor,
}: StatsDashboardProps) => {
  const teamAPercent = possessionStats?.teamA || 0;
  const teamBPercent = possessionStats?.teamB || 0;

  return (
    <div className="space-y-6">
      {/* Possession Gauge */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold mb-4">Ball Possession</h2>

        {/* Progress Bar */}
        <div className="relative h-12 bg-gray-200 rounded-full overflow-hidden mb-4">
          <div
            className="absolute top-0 left-0 h-full transition-all duration-500"
            style={{
              width: `${teamAPercent}%`,
              backgroundColor: teamAColor,
            }}
          />
          <div
            className="absolute top-0 right-0 h-full transition-all duration-500"
            style={{
              width: `${teamBPercent}%`,
              backgroundColor: teamBColor,
            }}
          />

          {/* Team Labels */}
          <div className="absolute inset-0 flex items-center justify-between px-4 text-white font-bold">
            <span>{teamAName} {teamAPercent.toFixed(1)}%</span>
            <span>{teamBName} {teamBPercent.toFixed(1)}%</span>
          </div>
        </div>

        {/* Current Holder */}
        <div className="text-center">
          <span className="text-gray-600">Current possession: </span>
          <span className="font-bold">
            {possessionStats?.currentHolder === 'teamA' && teamAName}
            {possessionStats?.currentHolder === 'teamB' && teamBName}
            {possessionStats?.currentHolder === 'none' && 'None'}
          </span>
        </div>

        {/* Total Duration */}
        {possessionStats && (
          <div className="text-center text-sm text-gray-500 mt-2">
            Total tracked: {possessionStats.totalDuration.toFixed(1)}s
          </div>
        )}
      </div>

      {/* Performance Metrics */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold mb-4">Performance</h2>

        <div className="grid grid-cols-2 gap-4">
          <div className="text-center p-4 bg-gray-50 rounded">
            <div className="text-3xl font-bold text-blue-600">
              {metrics?.fps.toFixed(1) || '0.0'}
            </div>
            <div className="text-sm text-gray-600 mt-1">FPS</div>
          </div>

          <div className="text-center p-4 bg-gray-50 rounded">
            <div className="text-3xl font-bold text-green-600">
              {metrics?.latency || '0'}
            </div>
            <div className="text-sm text-gray-600 mt-1">Latency (ms)</div>
          </div>
        </div>

        {/* Status Indicators */}
        <div className="mt-4 flex items-center justify-center space-x-4 text-sm">
          <div className="flex items-center">
            <div
              className={`w-3 h-3 rounded-full mr-2 ${
                metrics ? 'bg-green-500' : 'bg-red-500'
              }`}
            />
            <span>{metrics ? 'Processing' : 'Stopped'}</span>
          </div>

          {metrics && metrics.fps < 10 && (
            <div className="text-orange-500">⚠️ Low FPS</div>
          )}

          {metrics && metrics.latency > 200 && (
            <div className="text-orange-500">⚠️ High Latency</div>
          )}
        </div>
      </div>

      {/* Player Counts */}
      {possessionStats && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold mb-4">Statistics</h2>

          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-gray-600">{teamAName} Possession Time:</span>
              <span className="font-bold">
                {((possessionStats.teamA / 100) * possessionStats.totalDuration).toFixed(1)}s
              </span>
            </div>

            <div className="flex justify-between">
              <span className="text-gray-600">{teamBName} Possession Time:</span>
              <span className="font-bold">
                {((possessionStats.teamB / 100) * possessionStats.totalDuration).toFixed(1)}s
              </span>
            </div>

            {possessionStats.lastChange && (
              <div className="flex justify-between text-sm text-gray-500">
                <span>Last Change:</span>
                <span>{new Date(possessionStats.lastChange).toLocaleTimeString()}</span>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
