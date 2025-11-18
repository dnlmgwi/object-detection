import { useEffect, useRef } from 'react';
import type { DetectionResult, Player, Ball } from '~/hooks/useWebSocket';

interface VideoPlayerProps {
  detection: DetectionResult | null;
  teamAColor: string;
  teamBColor: string;
}

export const VideoPlayer = ({ detection, teamAColor, teamBColor }: VideoPlayerProps) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (!detection || !canvasRef.current) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Draw players
    detection.players.forEach((player) => {
      drawPlayer(ctx, player, teamAColor, teamBColor);
    });

    // Draw ball
    if (detection.ball) {
      drawBall(ctx, detection.ball);
    }

    // Draw possession indicator
    if (detection.possession) {
      drawPossessionIndicator(ctx, detection.possession, canvas.width, canvas.height);
    }
  }, [detection, teamAColor, teamBColor]);

  return (
    <div className="relative w-full aspect-video bg-gray-900 rounded-lg overflow-hidden">
      <canvas
        ref={canvasRef}
        width={1280}
        height={720}
        className="w-full h-full"
      />
      {!detection && (
        <div className="absolute inset-0 flex items-center justify-center text-white">
          <div className="text-center">
            <div className="text-xl mb-2">No video feed</div>
            <div className="text-sm text-gray-400">
              Configure and start video processing to see live detection
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

function drawPlayer(
  ctx: CanvasRenderingContext2D,
  player: Player,
  teamAColor: string,
  teamBColor: string
) {
  const { bbox, team, confidence } = player;

  // Determine color
  let color = '#808080'; // Gray for neutral
  if (team === 'teamA') color = teamAColor;
  if (team === 'teamB') color = teamBColor;

  // Draw bounding box
  ctx.strokeStyle = color;
  ctx.lineWidth = 3;
  ctx.strokeRect(bbox.x, bbox.y, bbox.width, bbox.height);

  // Draw label background
  const label = `${team} ${(confidence * 100).toFixed(0)}%`;
  ctx.font = '14px sans-serif';
  const metrics = ctx.measureText(label);
  const padding = 4;

  ctx.fillStyle = color;
  ctx.fillRect(
    bbox.x,
    bbox.y - 20,
    metrics.width + padding * 2,
    20
  );

  // Draw label text
  ctx.fillStyle = '#FFFFFF';
  ctx.fillText(label, bbox.x + padding, bbox.y - 6);
}

function drawBall(ctx: CanvasRenderingContext2D, ball: Ball) {
  const { bbox } = ball;
  const centerX = bbox.x + bbox.width / 2;
  const centerY = bbox.y + bbox.height / 2;
  const radius = Math.max(bbox.width, bbox.height) / 2;

  // Draw circle
  ctx.strokeStyle = '#00FF00';
  ctx.lineWidth = 3;
  ctx.beginPath();
  ctx.arc(centerX, centerY, radius, 0, 2 * Math.PI);
  ctx.stroke();

  // Draw center dot
  ctx.fillStyle = '#00FF00';
  ctx.beginPath();
  ctx.arc(centerX, centerY, 3, 0, 2 * Math.PI);
  ctx.fill();
}

function drawPossessionIndicator(
  ctx: CanvasRenderingContext2D,
  possession: any,
  width: number,
  height: number
) {
  const text = `Possession: ${possession.team} (${possession.duration.toFixed(1)}s)`;
  ctx.font = 'bold 20px sans-serif';
  const metrics = ctx.measureText(text);

  const x = width - metrics.width - 20;
  const y = 30;

  // Background
  ctx.fillStyle = 'rgba(0, 0, 0, 0.7)';
  ctx.fillRect(x - 10, y - 20, metrics.width + 20, 30);

  // Text
  ctx.fillStyle = '#FFFFFF';
  ctx.fillText(text, x, y);
}
