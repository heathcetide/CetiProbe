import { useState, useEffect } from "react";
import { twMerge } from "tailwind-merge";
import { Frame } from "@/components/ui/frame";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { SwitchRoot, SwitchControl, SwitchThumb, SwitchLabel } from "@/components/ui/switch";
import { AlertRoot, AlertTitle, AlertDescription, AlertCloseTrigger } from "@/components/ui/alert";
import { DialogRoot, DialogTrigger, DialogBackdrop, DialogPositioner, DialogContent, DialogTitle, DialogDescription, DialogCloseTrigger } from "@/components/ui/dialog";
import { 
  Play, 
  Square, 
  Settings, 
  Trash2, 
  Wifi,
  WifiOff,
  Activity,
  Clock,
  Eye
} from "lucide-react";

// 模拟数据接口
interface CaptureSession {
  id: string;
  name: string;
  startTime: string;
  endTime?: string;
  status: 'running' | 'paused' | 'stopped';
  requestCount: number;
  responseCount: number;
  errorCount: number;
}

interface NetworkRequest {
  id: string;
  method: string;
  url: string;
  status: number;
  time: string;
  size: number;
  duration: number;
  protocol: string;
  host: string;
  headers?: Record<string, string>;
  requestBody?: string;
  responseBody?: string;
  userAgent?: string;
  contentType?: string;
  cookies?: string;
  referer?: string;
}

function CaptureHome() {
  const [isCapturing, setIsCapturing] = useState(false);
  const [captureSession, setCaptureSession] = useState<CaptureSession | null>(null);
  const [requests, setRequests] = useState<NetworkRequest[]>([]);
  const [filteredRequests, setFilteredRequests] = useState<NetworkRequest[]>([]);
  const [selectedRequest, setSelectedRequest] = useState<NetworkRequest | null>(null);
  const [filterMethod, setFilterMethod] = useState<string>('all');
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [autoScroll, setAutoScroll] = useState(true);

  // 模拟数据
  const mockRequests: NetworkRequest[] = [
    {
      id: '1',
      method: 'GET',
      url: 'https://api.example.com/users',
      status: 200,
      time: '14:32:15',
      size: 1024,
      duration: 150,
      protocol: 'HTTPS',
      host: 'api.example.com',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        'Accept': 'application/json, text/plain, */*',
        'Cache-Control': 'no-cache'
      },
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
      contentType: 'application/json',
      cookies: 'session_id=abc123; user_pref=dark_mode',
      referer: 'https://app.example.com/dashboard',
      responseBody: JSON.stringify({
        users: [
          { id: 1, name: 'John Doe', email: 'john@example.com' },
          { id: 2, name: 'Jane Smith', email: 'jane@example.com' }
        ],
        total: 2,
        page: 1
      }, null, 2)
    },
    {
      id: '2',
      method: 'POST',
      url: 'https://api.example.com/auth/login',
      status: 401,
      time: '14:32:18',
      size: 512,
      duration: 300,
      protocol: 'HTTPS',
      host: 'api.example.com',
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        'Accept': 'application/json',
        'Origin': 'https://app.example.com'
      },
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
      contentType: 'application/json',
      referer: 'https://app.example.com/login',
      requestBody: JSON.stringify({
        username: 'user@example.com',
        password: '••••••••'
      }, null, 2),
      responseBody: JSON.stringify({
        error: 'Invalid credentials',
        code: 'AUTH_FAILED',
        timestamp: '2024-01-15T14:32:18.123Z'
      }, null, 2)
    },
    {
      id: '3',
      method: 'GET',
      url: 'https://cdn.example.com/assets/style.css',
      status: 200,
      time: '14:32:20',
      size: 2048,
      duration: 80,
      protocol: 'HTTPS',
      host: 'cdn.example.com',
      headers: {
        'Content-Type': 'text/css',
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
        'Accept': 'text/css,*/*;q=0.1',
        'Cache-Control': 'max-age=3600'
      },
      userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
      contentType: 'text/css',
      referer: 'https://app.example.com/dashboard',
      responseBody: `/* Main stylesheet */
body {
  font-family: 'Inter', sans-serif;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  margin: 0;
  padding: 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}`
    }
  ];

  useEffect(() => {
    setRequests(mockRequests);
    setFilteredRequests(mockRequests);
  }, []);

  useEffect(() => {
    let filtered = requests;
    
    if (filterMethod !== 'all') {
      filtered = filtered.filter(req => req.method === filterMethod);
    }
    
    if (filterStatus !== 'all') {
      if (filterStatus === 'success') {
        filtered = filtered.filter(req => req.status >= 200 && req.status < 300);
      } else if (filterStatus === 'error') {
        filtered = filtered.filter(req => req.status >= 400);
      }
    }
    
    if (searchQuery) {
      filtered = filtered.filter(req => 
        req.url.toLowerCase().includes(searchQuery.toLowerCase()) ||
        req.host.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }
    
    setFilteredRequests(filtered);
  }, [requests, filterMethod, filterStatus, searchQuery]);

  const startCapture = () => {
    setIsCapturing(true);
    setCaptureSession({
      id: Date.now().toString(),
      name: `Session ${new Date().toLocaleTimeString()}`,
      startTime: new Date().toISOString(),
      status: 'running',
      requestCount: 0,
      responseCount: 0,
      errorCount: 0
    });
  };

  const stopCapture = () => {
    setIsCapturing(false);
    if (captureSession) {
      setCaptureSession({
        ...captureSession,
        status: 'stopped',
        endTime: new Date().toISOString()
      });
    }
  };

  const clearRequests = () => {
    setRequests([]);
    setFilteredRequests([]);
    setSelectedRequest(null);
  };

  const getStatusColor = (status: number) => {
    if (status >= 200 && status < 300) return 'text-green-400';
    if (status >= 300 && status < 400) return 'text-yellow-400';
    if (status >= 400) return 'text-red-400';
    return 'text-gray-400';
  };

  const getMethodColor = (method: string) => {
    switch (method) {
      case 'GET': return 'text-blue-400';
      case 'POST': return 'text-green-400';
      case 'PUT': return 'text-yellow-400';
      case 'DELETE': return 'text-red-400';
      default: return 'text-gray-400';
    }
  };

  return (
    <div className="min-h-screen p-6">
      {/* 头部控制面板 */}
      <div className="mb-6">
        <div
          className={twMerge([
            "relative backdrop-blur-xl p-6",
            "[--color-frame-1-stroke:var(--color-primary)]/50",
            "[--color-frame-1-fill:var(--color-primary)]/20",
            "[--color-frame-2-stroke:var(--color-accent)]",
            "[--color-frame-2-fill:var(--color-accent)]/20",
          ])}
        >
          <Frame
            className="drop-shadow-2xl drop-shadow-primary/50"
            paths={JSON.parse(
              '[{"show":true,"style":{"strokeWidth":"1","stroke":"var(--color-frame-1-stroke)","fill":"var(--color-frame-1-fill)"},"path":[["M","19","0"],["L","100% - 18","0"],["L","100% + 0","0% + 18.5"],["L","100% + 0","100% - 30.15267175572519%"],["L","100% - 17","100% - 7.5"],["L","50% + 17.16417910447761%","100% - 7.5"],["L","50% + 15.298507462686567%","100% - 15.5"],["L","50% - 14.552238805970148%","100% - 15.5"],["L","50% - 16.417910447761194%","100% - 6.5"],["L","0% + 17","100% - 7.5"],["L","0% + 0","100% - 24.5"],["L","0% + 0","50% + 19.84732824427481%"],["L","0% + 9","50% + 17.557251908396946%"],["L","0% + 10","50% - 18.829516539440203%"],["L","0","50% - 21.62849872773537%"],["L","0","0% + 19.5"],["L","19","0"]]},{"show":true,"style":{"strokeWidth":"1","stroke":"var(--color-frame-2-stroke)","fill":"var(--color-frame-2-fill)"},"path":[["M","28","100% - 7.000000000000057"],["L","50% - 16.417910447761194%","100% - 7"],["L","50% - 14.552238805970148%","100% - 15.5"],["L","50% + 15.298507462686567%","100% - 15.5"],["L","50% + 17.16417910447761%","100% - 7.5"],["L","100% - 26","100% - 7.5"],["L","100% - 33","100% + 0"],["L","50% + 16.23134328358209%","100% - 1.1368683772161605"],["L","50% + 14.552238805970148%","100% - 8"],["L","50% - 13.619402985074627%","100% - 8"],["L","50% - 15.111940298507463%","100% + 0"],["L","33","100% + 0"],["L","28","100% - 7"]]}]'
            )}
          />
          
          <div className="relative flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                {isCapturing ? (
                  <Wifi className="size-6 text-green-400 animate-pulse" />
                ) : (
                  <WifiOff className="size-6 text-gray-400" />
                )}
                <h1 className="text-2xl font-bold text-shadow-lg text-shadow-primary font-inter">
                  Network Capture
                </h1>
              </div>
              
              {captureSession && (
                <div className="flex items-center gap-4 text-sm opacity-70">
                  <div className="flex items-center gap-1">
                    <Clock className="size-4" />
                    <span>{captureSession.startTime}</span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Activity className="size-4" />
                    <span>{captureSession.requestCount} requests</span>
                  </div>
                </div>
              )}
            </div>

            <div className="flex items-center gap-3">
              <SwitchRoot checked={autoScroll} onCheckedChange={(details) => setAutoScroll(details.checked)}>
                <SwitchControl>
                  <SwitchThumb />
                </SwitchControl>
                <SwitchLabel>Auto Scroll</SwitchLabel>
              </SwitchRoot>

              <Button
                variant={isCapturing ? "destructive" : "success"}
                onClick={isCapturing ? stopCapture : startCapture}
                className="flex items-center gap-2"
              >
                {isCapturing ? (
                  <>
                    <Square className="size-4" />
                    Stop Capture
                  </>
                ) : (
                  <>
                    <Play className="size-4" />
                    Start Capture
                  </>
                )}
              </Button>

              <Button
                variant="secondary"
                onClick={() => {}}
                className="flex items-center gap-2"
              >
                <Settings className="size-4" />
                Settings
              </Button>

              <Button
                variant="destructive"
                onClick={clearRequests}
                className="flex items-center gap-2"
              >
                <Trash2 className="size-4" />
                Clear
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* 请求列表 */}
        <div className="lg:col-span-2">
          <div
            className={twMerge([
              "relative backdrop-blur-xl h-96",
              "[--color-frame-1-stroke:var(--color-primary)]/50",
              "[--color-frame-1-fill:var(--color-primary)]/20",
            ])}
          >
            <Frame
              className="drop-shadow-2xl drop-shadow-primary/50"
              paths={JSON.parse(
                '[{"show":true,"style":{"strokeWidth":"1","stroke":"var(--color-frame-1-stroke)","fill":"var(--color-frame-1-fill)"},"path":[["M","19","0"],["L","100% - 18","0"],["L","100% + 0","0% + 18.5"],["L","100% + 0","100% - 30.15267175572519%"],["L","100% - 17","100% - 7.5"],["L","50% + 17.16417910447761%","100% - 7.5"],["L","50% + 15.298507462686567%","100% - 15.5"],["L","50% - 14.552238805970148%","100% - 15.5"],["L","50% - 16.417910447761194%","100% - 6.5"],["L","0% + 17","100% - 7.5"],["L","0% + 0","100% - 24.5"],["L","0% + 0","50% + 19.84732824427481%"],["L","0% + 9","50% + 17.557251908396946%"],["L","0% + 10","50% - 18.829516539440203%"],["L","0","50% - 21.62849872773537%"],["L","0","0% + 19.5"],["L","19","0"]]},{"show":true,"style":{"strokeWidth":"1","stroke":"var(--color-frame-2-stroke)","fill":"var(--color-frame-2-fill)"},"path":[["M","28","100% - 7.000000000000057"],["L","50% - 16.417910447761194%","100% - 7"],["L","50% - 14.552238805970148%","100% - 15.5"],["L","50% + 15.298507462686567%","100% - 15.5"],["L","50% + 17.16417910447761%","100% - 7.5"],["L","100% - 26","100% - 7.5"],["L","100% - 33","100% + 0"],["L","50% + 16.23134328358209%","100% - 1.1368683772161605"],["L","50% + 14.552238805970148%","100% - 8"],["L","50% - 13.619402985074627%","100% - 8"],["L","50% - 15.111940298507463%","100% + 0"],["L","33","100% + 0"],["L","28","100% - 7"]]}]'
              )}
            />
            
            <div className="relative p-6">
              {/* 过滤器 */}
              <div className="flex items-center gap-4 mb-4">
                <Input
                  placeholder="搜索请求..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="flex-1"
                />
                
                <select
                  value={filterMethod}
                  onChange={(e) => setFilterMethod(e.target.value)}
                  className="px-3 py-2 bg-background/50 border border-primary/30 rounded"
                >
                  <option value="all">All Methods</option>
                  <option value="GET">GET</option>
                  <option value="POST">POST</option>
                  <option value="PUT">PUT</option>
                  <option value="DELETE">DELETE</option>
                </select>

                <select
                  value={filterStatus}
                  onChange={(e) => setFilterStatus(e.target.value)}
                  className="px-3 py-2 bg-background/50 border border-primary/30 rounded"
                >
                  <option value="all">All Status</option>
                  <option value="success">Success (2xx)</option>
                  <option value="error">Error (4xx/5xx)</option>
                </select>
              </div>

              {/* 请求列表 */}
              <div className="space-y-1 max-h-80 overflow-y-auto scrollbar-ultra-thin">
                {filteredRequests.map((request) => (
                  <div
                    key={request.id}
                    className={twMerge([
                      "p-4 rounded-lg transition-all duration-200 hover:bg-primary/10 border",
                      selectedRequest?.id === request.id 
                        ? "bg-primary/20 border-primary/50 shadow-lg shadow-primary/20" 
                        : "border-primary/20 hover:border-primary/40"
                    ])}
                  >
                    <div className="flex items-center justify-between mb-2">
                      <div 
                        className="flex items-center gap-3 cursor-pointer flex-1"
                        onClick={() => setSelectedRequest(request)}
                      >
                        <span className={twMerge(
                          "font-mono text-xs font-bold px-2 py-1 rounded",
                          getMethodColor(request.method),
                          "bg-current/10"
                        )}>
                          {request.method}
                        </span>
                        <span className="text-xs font-mono opacity-70">{request.time}</span>
                        <span className={twMerge(
                          "text-xs font-bold px-2 py-1 rounded",
                          getStatusColor(request.status),
                          "bg-current/10"
                        )}>
                          {request.status}
                        </span>
                      </div>
                      <div className="flex items-center gap-3">
                        <div className="flex items-center gap-2 text-xs opacity-70">
                          <span className="font-mono">{request.duration}ms</span>
                          <span className="font-mono">{request.size}B</span>
                        </div>
                        <DialogRoot>
                          <DialogTrigger asChild>
                            <div
                              className="p-2 h-8 w-8 rounded border border-primary/30 hover:border-primary/50 hover:bg-primary/10 transition-all flex items-center justify-center cursor-pointer"
                            >
                              <Eye className="size-3" />
                            </div>
                          </DialogTrigger>
                          <DialogBackdrop />
                          <DialogPositioner>
                            <DialogContent className="max-w-4xl max-h-[80vh] flex flex-col overflow-y-auto">
                              <div className="flex-shrink-0">
                                <DialogTitle className="flex items-center gap-2">
                                  <span className={twMerge("font-mono px-2 py-1 rounded text-xs", getMethodColor(request.method), "bg-current/10")}>
                                    {request.method}
                                  </span>
                                  <span className="truncate">{request.url}</span>
                                </DialogTitle>
                                <DialogDescription className="text-sm opacity-70 mb-4 font-inter">
                                  {request.host} • {request.time} • {request.duration}ms • {request.size}B
                                </DialogDescription>
                              </div>
                              
                              <div className="flex-1 overflow-y-auto scrollbar-ultra-thin space-y-6 pr-2">
                                {/* 基本信息 */}
                                <div>
                                  <h4 className="text-sm font-bold mb-2 text-primary font-inter">基本信息</h4>
                                  <div className="grid grid-cols-2 gap-4 text-xs">
                                    <div>
                                      <span className="opacity-70 font-inter">协议:</span>
                                      <span className="ml-2 font-jetbrains">{request.protocol}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">状态码:</span>
                                      <span className={twMerge("ml-2 font-bold font-jetbrains", getStatusColor(request.status))}>
                                        {request.status}
                                      </span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">响应时间:</span>
                                      <span className="ml-2 font-jetbrains">{request.duration}ms</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">大小:</span>
                                      <span className="ml-2 font-jetbrains">{request.size}B</span>
                                    </div>
                                  </div>
                                </div>

                                {/* 请求头 */}
                                {request.headers && (
                                  <div>
                                    <h4 className="text-sm font-bold mb-2 text-primary font-inter">请求头</h4>
                                    <div className="bg-background/30 rounded p-3 max-h-32 overflow-y-auto scrollbar-ultra-thin">
                                      <pre className="text-xs font-jetbrains whitespace-pre-wrap">
                                        {Object.entries(request.headers).map(([key, value]) => 
                                          `${key}: ${value}`
                                        ).join('\n')}
                                      </pre>
                                    </div>
                                  </div>
                                )}

                                {/* 请求体 */}
                                {request.requestBody && (
                                  <div>
                                    <h4 className="text-sm font-bold mb-2 text-primary font-inter">请求体</h4>
                                    <div className="bg-background/30 rounded p-3 max-h-32 overflow-y-auto scrollbar-ultra-thin">
                                      <pre className="text-xs font-jetbrains whitespace-pre-wrap">
                                        {request.requestBody}
                                      </pre>
                                    </div>
                                  </div>
                                )}

                                {/* 响应体 */}
                                {request.responseBody && (
                                  <div>
                                    <h4 className="text-sm font-bold mb-2 text-primary font-inter">响应体</h4>
                                    <div className="bg-background/30 rounded p-3 max-h-32 overflow-y-auto scrollbar-ultra-thin">
                                      <pre className="text-xs font-jetbrains whitespace-pre-wrap">
                                        {request.responseBody}
                                      </pre>
                                    </div>
                                  </div>
                                )}

                                {/* 其他信息 */}
                                <div className="grid grid-cols-1 gap-3 text-xs">
                                  {request.userAgent && (
                                    <div>
                                      <span className="opacity-70 font-inter">User-Agent:</span>
                                      <div className="mt-1 p-2 bg-background/30 rounded font-jetbrains text-xs break-all">
                                        {request.userAgent}
                                      </div>
                                    </div>
                                  )}
                                  {request.referer && (
                                    <div>
                                      <span className="opacity-70 font-inter">Referer:</span>
                                      <div className="mt-1 p-2 bg-background/30 rounded font-jetbrains text-xs break-all">
                                        {request.referer}
                                      </div>
                                    </div>
                                  )}
                                  {request.cookies && (
                                    <div>
                                      <span className="opacity-70 font-inter">Cookies:</span>
                                      <div className="mt-1 p-2 bg-background/30 rounded font-jetbrains text-xs break-all">
                                        {request.cookies}
                                      </div>
                                    </div>
                                  )}
                                </div>
                              </div>
                              
                              <DialogCloseTrigger />
                            </DialogContent>
                          </DialogPositioner>
                        </DialogRoot>
                      </div>
                    </div>
                    <div 
                      className="text-sm font-mono truncate mb-1 cursor-pointer"
                      onClick={() => setSelectedRequest(request)}
                    >
                      {request.url}
                    </div>
                    <div 
                      className="text-xs opacity-50 font-mono cursor-pointer"
                      onClick={() => setSelectedRequest(request)}
                    >
                      {request.host}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* 统计面板 */}
        <div className="space-y-6">
          {/* 实时统计 */}
          <div
            className={twMerge([
              "relative backdrop-blur-xl p-6",
              "[--color-frame-1-stroke:var(--color-primary)]/50",
              "[--color-frame-1-fill:var(--color-primary)]/20",
            ])}
          >
            <Frame
              className="drop-shadow-2xl drop-shadow-primary/50"
              paths={JSON.parse(
                '[{"show":true,"style":{"strokeWidth":"1","stroke":"var(--color-frame-1-stroke)","fill":"var(--color-frame-1-fill)"},"path":[["M","19","0"],["L","100% - 18","0"],["L","100% + 0","0% + 18.5"],["L","100% + 0","100% - 30.15267175572519%"],["L","100% - 17","100% - 7.5"],["L","50% + 17.16417910447761%","100% - 7.5"],["L","50% + 15.298507462686567%","100% - 15.5"],["L","50% - 14.552238805970148%","100% - 15.5"],["L","50% - 16.417910447761194%","100% - 6.5"],["L","0% + 17","100% - 7.5"],["L","0% + 0","100% - 24.5"],["L","0% + 0","50% + 19.84732824427481%"],["L","0% + 9","50% + 17.557251908396946%"],["L","0% + 10","50% - 18.829516539440203%"],["L","0","50% - 21.62849872773537%"],["L","0","0% + 19.5"],["L","19","0"]]},{"show":true,"style":{"strokeWidth":"1","stroke":"var(--color-frame-2-stroke)","fill":"var(--color-frame-2-fill)"},"path":[["M","28","100% - 7.000000000000057"],["L","50% - 16.417910447761194%","100% - 7"],["L","50% - 14.552238805970148%","100% - 15.5"],["L","50% + 15.298507462686567%","100% - 15.5"],["L","50% + 17.16417910447761%","100% - 7.5"],["L","100% - 26","100% - 7.5"],["L","100% - 33","100% + 0"],["L","50% + 16.23134328358209%","100% - 1.1368683772161605"],["L","50% + 14.552238805970148%","100% - 8"],["L","50% - 13.619402985074627%","100% - 8"],["L","50% - 15.111940298507463%","100% + 0"],["L","33","100% + 0"],["L","28","100% - 7"]]}]'
              )}
            />
            
            <div className="relative">
              <h3 className="text-lg font-bold text-shadow-lg text-shadow-primary mb-4 font-inter">
                实时统计
              </h3>
              
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm opacity-70 font-inter">总请求数</span>
                  <span className="text-xl font-bold text-primary font-jetbrains">{requests.length}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm opacity-70 font-inter">成功请求</span>
                  <span className="text-xl font-bold text-green-400 font-jetbrains">
                    {requests.filter(r => r.status >= 200 && r.status < 300).length}
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm opacity-70 font-inter">错误请求</span>
                  <span className="text-xl font-bold text-red-400 font-jetbrains">
                    {requests.filter(r => r.status >= 400).length}
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm opacity-70 font-inter">平均响应时间</span>
                  <span className="text-xl font-bold text-accent font-jetbrains">
                    {requests.length > 0 ? Math.round(requests.reduce((sum, r) => sum + r.duration, 0) / requests.length) : 0}ms
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* 请求详情 */}
          {selectedRequest && (
            <div
              className={twMerge([
                "relative backdrop-blur-xl p-6",
                "[--color-frame-1-stroke:var(--color-accent)]/50",
                "[--color-frame-1-fill:var(--color-accent)]/20",
              ])}
            >
              <Frame
                className="drop-shadow-2xl drop-shadow-accent/50"
                paths={JSON.parse(
                  '[{"show":true,"style":{"strokeWidth":"1","stroke":"var(--color-frame-1-stroke)","fill":"var(--color-frame-1-fill)"},"path":[["M","19","0"],["L","100% - 18","0"],["L","100% + 0","0% + 18.5"],["L","100% + 0","100% - 30.15267175572519%"],["L","100% - 17","100% - 7.5"],["L","50% + 17.16417910447761%","100% - 7.5"],["L","50% + 15.298507462686567%","100% - 15.5"],["L","50% - 14.552238805970148%","100% - 15.5"],["L","50% - 16.417910447761194%","100% - 6.5"],["L","0% + 17","100% - 7.5"],["L","0% + 0","100% - 24.5"],["L","0% + 0","50% + 19.84732824427481%"],["L","0% + 9","50% + 17.557251908396946%"],["L","0% + 10","50% - 18.829516539440203%"],["L","0","50% - 21.62849872773537%"],["L","0","0% + 19.5"],["L","19","0"]]},{"show":true,"style":{"strokeWidth":"1","stroke":"var(--color-frame-2-stroke)","fill":"var(--color-frame-2-fill)"},"path":[["M","28","100% - 7.000000000000057"],["L","50% - 16.417910447761194%","100% - 7"],["L","50% - 14.552238805970148%","100% - 15.5"],["L","50% + 15.298507462686567%","100% - 15.5"],["L","50% + 17.16417910447761%","100% - 7.5"],["L","100% - 26","100% - 7.5"],["L","100% - 33","100% + 0"],["L","50% + 16.23134328358209%","100% - 1.1368683772161605"],["L","50% + 14.552238805970148%","100% - 8"],["L","50% - 13.619402985074627%","100% - 8"],["L","50% - 15.111940298507463%","100% + 0"],["L","33","100% + 0"],["L","28","100% - 7"]]}]'
                )}
              />
              
              <div className="relative">
                <h3 className="text-lg font-bold text-shadow-lg text-shadow-accent mb-4 font-inter">
                  请求详情
                </h3>
                
                <div className="space-y-3 text-sm">
                  <div>
                    <span className="opacity-70 font-inter">方法:</span>
                    <span className={twMerge("ml-2 font-jetbrains font-bold", getMethodColor(selectedRequest.method))}>
                      {selectedRequest.method}
                    </span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">状态码:</span>
                    <span className={twMerge("ml-2 font-bold font-jetbrains", getStatusColor(selectedRequest.status))}>
                      {selectedRequest.status}
                    </span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">协议:</span>
                    <span className="ml-2 font-jetbrains">{selectedRequest.protocol}</span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">主机:</span>
                    <span className="ml-2 font-jetbrains">{selectedRequest.host}</span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">响应时间:</span>
                    <span className="ml-2 font-jetbrains">{selectedRequest.duration}ms</span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">大小:</span>
                    <span className="ml-2 font-jetbrains">{selectedRequest.size}B</span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">时间:</span>
                    <span className="ml-2 font-jetbrains">{selectedRequest.time}</span>
                  </div>
                </div>
                
                <div className="mt-4">
                  <span className="opacity-70 text-sm font-inter">URL:</span>
                  <div className="mt-1 p-2 bg-background/30 rounded text-xs break-all font-jetbrains">
                    {selectedRequest.url}
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* 警告提示 */}
      {!isCapturing && requests.length === 0 && (
        <AlertRoot className="mt-6">
          <AlertTitle>准备开始抓包</AlertTitle>
          <AlertDescription>
            点击"Start Capture"按钮开始捕获网络请求。确保已正确配置代理设置。
          </AlertDescription>
          <AlertCloseTrigger />
        </AlertRoot>
      )}
    </div>
  );
}

export default CaptureHome;
