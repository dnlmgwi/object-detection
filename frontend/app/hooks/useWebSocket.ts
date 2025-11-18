import { useEffect, useRef, useState, useCallback } from 'react';

export interface WSMessage {
  type: 'detection' | 'possession' | 'metrics' | 'error';
  data: any;
}

export interface DetectionResult {
  frameId: number;
  timestamp: number;
  players: Player[];
  ball: Ball | null;
  possession: Possession | null;
}

export interface Player {
  id: string;
  bbox: BoundingBox;
  team: string;
  confidence: number;
  colorMatch?: number;
}

export interface Ball {
  bbox: BoundingBox;
  confidence: number;
  velocity?: { x: number; y: number };
}

export interface BoundingBox {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface Possession {
  team: string;
  playerId: string;
  duration: number;
  distance: number;
}

export interface PossessionStats {
  teamA: number;
  teamB: number;
  currentHolder: string;
  lastChange: number;
  totalDuration: number;
}

export interface PerformanceMetrics {
  fps: number;
  latency: number;
  cpuUsage?: number;
  memoryUsage?: number;
}

interface UseWebSocketOptions {
  url: string;
  onDetection?: (data: DetectionResult) => void;
  onPossession?: (data: PossessionStats) => void;
  onMetrics?: (data: PerformanceMetrics) => void;
  onError?: (error: string) => void;
  reconnectAttempts?: number;
  reconnectInterval?: number;
}

export const useWebSocket = ({
  url,
  onDetection,
  onPossession,
  onMetrics,
  onError,
  reconnectAttempts = 3,
  reconnectInterval = 2000,
}: UseWebSocketOptions) => {
  const [isConnected, setIsConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState<WSMessage | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectCount = useRef(0);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>();

  const connect = useCallback(() => {
    try {
      const ws = new WebSocket(url);

      ws.onopen = () => {
        console.log('WebSocket connected');
        setIsConnected(true);
        reconnectCount.current = 0;
      };

      ws.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data);
          setLastMessage(message);

          switch (message.type) {
            case 'detection':
              onDetection?.(message.data as DetectionResult);
              break;
            case 'possession':
              onPossession?.(message.data as PossessionStats);
              break;
            case 'metrics':
              onMetrics?.(message.data as PerformanceMetrics);
              break;
            case 'error':
              onError?.(message.data.message || 'Unknown error');
              break;
          }
        } catch (err) {
          console.error('Failed to parse WebSocket message:', err);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        setIsConnected(false);
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected');
        setIsConnected(false);

        // Attempt reconnection
        if (reconnectCount.current < reconnectAttempts) {
          reconnectCount.current++;
          console.log(
            `Attempting to reconnect... (${reconnectCount.current}/${reconnectAttempts})`
          );

          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, reconnectInterval);
        } else {
          console.error('Max reconnection attempts reached');
          onError?.('Connection lost. Please refresh the page.');
        }
      };

      wsRef.current = ws;
    } catch (err) {
      console.error('Failed to create WebSocket:', err);
      setIsConnected(false);
    }
  }, [url, onDetection, onPossession, onMetrics, onError, reconnectAttempts, reconnectInterval]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
    }

    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }

    setIsConnected(false);
  }, []);

  useEffect(() => {
    connect();

    return () => {
      disconnect();
    };
  }, [connect, disconnect]);

  return {
    isConnected,
    lastMessage,
    disconnect,
    reconnect: connect,
  };
};
