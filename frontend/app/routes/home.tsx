import { useState, useCallback } from 'react';
import type { Route } from './+types/home';
import { VideoPlayer } from '~/components/VideoPlayer';
import { StatsDashboard } from '~/components/StatsDashboard';
import { ConfigPanel } from '~/components/ConfigPanel';
import {
  useWebSocket,
  type DetectionResult,
  type PossessionStats,
  type PerformanceMetrics,
} from '~/hooks/useWebSocket';

export function meta({}: Route.MetaArgs) {
  return [
    { title: 'Football Object Detection' },
    { name: 'description', content: 'Real-time football game object detection and possession tracking' },
  ];
}

export default function Home() {
  const [detection, setDetection] = useState<DetectionResult | null>(null);
  const [possessionStats, setPossessionStats] = useState<PossessionStats | null>(null);
  const [metrics, setMetrics] = useState<PerformanceMetrics | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [teamAName, setTeamAName] = useState('Team A');
  const [teamBName, setTeamBName] = useState('Team B');
  const [teamAColor, setTeamAColor] = useState('#FF0000');
  const [teamBColor, setTeamBColor] = useState('#0000FF');

  // WebSocket configuration
  const wsUrl = typeof window !== 'undefined'
    ? `ws://${window.location.hostname}:8080/ws`
    : 'ws://localhost:8080/ws';

  const { isConnected } = useWebSocket({
    url: wsUrl,
    onDetection: useCallback((data: DetectionResult) => {
      setDetection(data);
    }, []),
    onPossession: useCallback((data: PossessionStats) => {
      setPossessionStats(data);
    }, []),
    onMetrics: useCallback((data: PerformanceMetrics) => {
      setMetrics(data);
    }, []),
    onError: useCallback((error: string) => {
      setError(error);
      console.error('WebSocket error:', error);
    }, []),
  });

  const apiCall = async (endpoint: string, method: string = 'POST', body?: any) => {
    const apiUrl = typeof window !== 'undefined'
      ? `http://${window.location.hostname}:8080/api${endpoint}`
      : `http://localhost:8080/api${endpoint}`;

    const response = await fetch(apiUrl, {
      method,
      headers: body ? { 'Content-Type': 'application/json' } : undefined,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'API request failed');
    }

    return response.json();
  };

  const handleUpdateTeams = async (teamA: any, teamB: any) => {
    try {
      await apiCall('/config/teams', 'POST', { teamA, teamB });
      setTeamAName(teamA.name);
      setTeamBName(teamB.name);
      setTeamAColor(
        `#${teamA.primaryColor.r.toString(16).padStart(2, '0')}${teamA.primaryColor.g.toString(16).padStart(2, '0')}${teamA.primaryColor.b.toString(16).padStart(2, '0')}`
      );
      setTeamBColor(
        `#${teamB.primaryColor.r.toString(16).padStart(2, '0')}${teamB.primaryColor.g.toString(16).padStart(2, '0')}${teamB.primaryColor.b.toString(16).padStart(2, '0')}`
      );
      alert('Teams configured successfully!');
    } catch (error: any) {
      alert(`Failed to update teams: ${error.message}`);
    }
  };

  const handleUpdateVideo = async (video: any) => {
    try {
      await apiCall('/config/video', 'POST', video);
      alert('Video source configured successfully!');
    } catch (error: any) {
      alert(`Failed to update video config: ${error.message}`);
    }
  };

  const handleUploadVideo = async (file: File): Promise<string> => {
    const formData = new FormData();
    formData.append('video', file);

    const apiUrl = typeof window !== 'undefined'
      ? `http://${window.location.hostname}:8080/api/video/upload`
      : 'http://localhost:8080/api/video/upload';

    const response = await fetch(apiUrl, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      throw new Error('Upload failed');
    }

    const result = await response.json();
    return result.filePath;
  };

  const handleStart = async () => {
    try {
      await apiCall('/control/start', 'POST');
      setIsProcessing(true);
      setError(null);
    } catch (error: any) {
      alert(`Failed to start processing: ${error.message}`);
    }
  };

  const handleStop = async () => {
    try {
      await apiCall('/control/stop', 'POST');
      setIsProcessing(false);
      setDetection(null);
      setMetrics(null);
    } catch (error: any) {
      alert(`Failed to stop processing: ${error.message}`);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 py-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-900">
              ⚽ Football Object Detection
            </h1>
            <div className="flex items-center space-x-4">
              <div
                className={`flex items-center text-sm ${
                  isConnected ? 'text-green-600' : 'text-red-600'
                }`}
              >
                <div
                  className={`w-2 h-2 rounded-full mr-2 ${
                    isConnected ? 'bg-green-600' : 'bg-red-600'
                  }`}
                />
                {isConnected ? 'Connected' : 'Disconnected'}
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 py-6 sm:px-6 lg:px-8">
        {error && (
          <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column - Video Player */}
          <div className="lg:col-span-2 space-y-6">
            <VideoPlayer
              detection={detection}
              teamAColor={teamAColor}
              teamBColor={teamBColor}
            />

            {/* Info */}
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <h3 className="font-semibold text-blue-900 mb-2">How to use:</h3>
              <ol className="text-sm text-blue-800 space-y-1 list-decimal list-inside">
                <li>Configure team names and colors in the right panel</li>
                <li>Upload a video file or configure a video source</li>
                <li>Click "Start Processing" to begin detection</li>
                <li>Watch real-time detection and possession tracking</li>
              </ol>
            </div>
          </div>

          {/* Right Column - Stats and Config */}
          <div className="space-y-6">
            <StatsDashboard
              possessionStats={possessionStats}
              metrics={metrics}
              teamAName={teamAName}
              teamBName={teamBName}
              teamAColor={teamAColor}
              teamBColor={teamBColor}
            />

            <ConfigPanel
              onUpdateTeams={handleUpdateTeams}
              onUpdateVideo={handleUpdateVideo}
              onUploadVideo={handleUploadVideo}
              onStart={handleStart}
              onStop={handleStop}
              isProcessing={isProcessing}
            />
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="bg-white mt-12 border-t border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-6 sm:px-6 lg:px-8">
          <div className="text-center text-sm text-gray-500">
            <p>Built with React Router 7, Go, GoCV, and YOLOv8</p>
            <p className="mt-1">Real-time object detection for football games</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
