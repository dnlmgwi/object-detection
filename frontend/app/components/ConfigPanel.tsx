import { useState } from 'react';

interface Team {
  id: string;
  name: string;
  primaryColor: { r: number; g: number; b: number };
  secondaryColor: { r: number; g: number; b: number };
}

interface VideoConfig {
  source: 'file' | 'webcam' | 'rtsp';
  path: string;
  fps: number;
}

interface ConfigPanelProps {
  onUpdateTeams: (teamA: Team, teamB: Team) => Promise<void>;
  onUpdateVideo: (video: VideoConfig) => Promise<void>;
  onUploadVideo: (file: File) => Promise<string>;
  onStart: () => Promise<void>;
  onStop: () => Promise<void>;
  isProcessing: boolean;
}

export const ConfigPanel = ({
  onUpdateTeams,
  onUpdateVideo,
  onUploadVideo,
  onStart,
  onStop,
  isProcessing,
}: ConfigPanelProps) => {
  const [teamAName, setTeamAName] = useState('Team A');
  const [teamAColor, setTeamAColor] = useState('#FF0000');
  const [teamBName, setTeamBName] = useState('Team B');
  const [teamBColor, setTeamBColor] = useState('#0000FF');

  const [videoSource, setVideoSource] = useState<'file' | 'webcam' | 'rtsp'>('file');
  const [videoPath, setVideoPath] = useState('');
  const [fps, setFps] = useState(15);
  const [uploading, setUploading] = useState(false);

  const hexToRgb = (hex: string) => {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result
      ? {
          r: parseInt(result[1], 16),
          g: parseInt(result[2], 16),
          b: parseInt(result[3], 16),
        }
      : { r: 0, g: 0, b: 0 };
  };

  const handleSaveTeams = async () => {
    const teamA: Team = {
      id: 'teamA',
      name: teamAName,
      primaryColor: hexToRgb(teamAColor),
      secondaryColor: { r: 255, g: 255, b: 255 },
    };

    const teamB: Team = {
      id: 'teamB',
      name: teamBName,
      primaryColor: hexToRgb(teamBColor),
      secondaryColor: { r: 255, g: 255, b: 255 },
    };

    await onUpdateTeams(teamA, teamB);
  };

  const handleSaveVideo = async () => {
    const video: VideoConfig = {
      source: videoSource,
      path: videoPath,
      fps,
    };

    await onUpdateVideo(video);
  };

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    try {
      const filePath = await onUploadVideo(file);
      setVideoPath(filePath);
      setVideoSource('file');
    } catch (error) {
      console.error('Upload failed:', error);
      alert('Failed to upload video');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Team Configuration */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold mb-4">Team Configuration</h2>

        <div className="space-y-4">
          {/* Team A */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Team A Name
            </label>
            <input
              type="text"
              value={teamAName}
              onChange={(e) => setTeamAName(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Team A Color
            </label>
            <div className="flex items-center space-x-2">
              <input
                type="color"
                value={teamAColor}
                onChange={(e) => setTeamAColor(e.target.value)}
                className="h-10 w-20 border border-gray-300 rounded cursor-pointer"
              />
              <input
                type="text"
                value={teamAColor}
                onChange={(e) => setTeamAColor(e.target.value)}
                className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          {/* Team B */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Team B Name
            </label>
            <input
              type="text"
              value={teamBName}
              onChange={(e) => setTeamBName(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Team B Color
            </label>
            <div className="flex items-center space-x-2">
              <input
                type="color"
                value={teamBColor}
                onChange={(e) => setTeamBColor(e.target.value)}
                className="h-10 w-20 border border-gray-300 rounded cursor-pointer"
              />
              <input
                type="text"
                value={teamBColor}
                onChange={(e) => setTeamBColor(e.target.value)}
                className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          <button
            onClick={handleSaveTeams}
            className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Save Teams
          </button>
        </div>
      </div>

      {/* Video Source Configuration */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold mb-4">Video Source</h2>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Source Type
            </label>
            <select
              value={videoSource}
              onChange={(e) => setVideoSource(e.target.value as any)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="file">Video File</option>
              <option value="webcam">Webcam</option>
              <option value="rtsp">RTSP Stream</option>
            </select>
          </div>

          {videoSource === 'file' && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Upload Video File
              </label>
              <input
                type="file"
                accept="video/*"
                onChange={handleFileUpload}
                disabled={uploading}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              {uploading && (
                <div className="mt-2 text-sm text-blue-600">Uploading...</div>
              )}
              {videoPath && (
                <div className="mt-2 text-sm text-green-600">
                  File: {videoPath}
                </div>
              )}
            </div>
          )}

          {videoSource === 'webcam' && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Webcam Device ID
              </label>
              <input
                type="number"
                value={videoPath}
                onChange={(e) => setVideoPath(e.target.value)}
                placeholder="0"
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          )}

          {videoSource === 'rtsp' && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                RTSP URL
              </label>
              <input
                type="text"
                value={videoPath}
                onChange={(e) => setVideoPath(e.target.value)}
                placeholder="rtsp://..."
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Processing FPS
            </label>
            <input
              type="number"
              value={fps}
              onChange={(e) => setFps(parseInt(e.target.value) || 15)}
              min="1"
              max="30"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <div className="mt-1 text-xs text-gray-500">
              Higher FPS = more accurate but higher CPU usage
            </div>
          </div>

          <button
            onClick={handleSaveVideo}
            className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Save Video Config
          </button>
        </div>
      </div>

      {/* Control */}
      <div className="bg-white rounded-lg shadow-md p-6">
        <h2 className="text-xl font-bold mb-4">Control</h2>

        <div className="space-y-3">
          {!isProcessing ? (
            <button
              onClick={onStart}
              disabled={!videoPath}
              className="w-full px-4 py-3 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:bg-gray-300 disabled:cursor-not-allowed font-semibold"
            >
              ▶ Start Processing
            </button>
          ) : (
            <button
              onClick={onStop}
              className="w-full px-4 py-3 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 font-semibold"
            >
              ⏹ Stop Processing
            </button>
          )}

          {!videoPath && (
            <div className="text-sm text-orange-600 text-center">
              ⚠️ Please configure and save a video source first
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
