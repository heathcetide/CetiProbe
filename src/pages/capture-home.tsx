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

// 抓包数据接口
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

// 数据包元信息
interface PacketMetadata {
  capture_time: string;
  data_size: number;
  wire_length: number;
  capture_length: number;
  truncated: boolean;
  interface_index: number;
  capture_length_ok: boolean;
}

// 链路层信息
interface LinkLayer {
  timestamp: string;
  src_mac: string;
  dst_mac: string;
  eth_type: string;
  length: number;
}

// 网络层信息
interface NetworkLayer {
  timestamp: string;
  ip_version: number;
  src_ip: string;
  dst_ip: string;
  protocol: string;
  length: number;
  ttl: number;
  ihl: number;
  identifier: number;
  flags: number;
  checksum: number;
  is_src_loopback: boolean;
  is_dst_loopback: boolean;
  is_src_link_local: boolean;
  is_dst_link_local: boolean;
  is_src_ip_valid: boolean;
  is_dst_ip_valid: boolean;
}

// 传输层信息
interface TransportLayer {
  timestamp: string;
  src_port: number;
  dst_port: number;
  protocol: string;
  seq_number: number;
  ack_number: number;
  window_size: number;
  checksum: number;
  is_psh: boolean;
  is_ack: boolean;
}

// 应用层信息
interface ApplicationLayer {
  timestamp: string;
  payload: string; // base64编码
}

// 错误层信息
interface ErrorLayer {
  timestamp: string;
  error: string;
  layer: string;
  fatal: boolean;
  code: number;
}

// 完整的数据包信息
interface PacketInfo {
  id: number;
  metadata: PacketMetadata;
  linkLayer: LinkLayer;
  networkLayer: NetworkLayer;
  transportLayer: TransportLayer;
  applicationLayer: ApplicationLayer;
  errorLayer: ErrorLayer | null;
}

function CaptureHome() {
  const [isCapturing, setIsCapturing] = useState(false);
  const [captureSession, setCaptureSession] = useState<CaptureSession | null>(null);
  const [packets, setPackets] = useState<PacketInfo[]>([]);
  const [filteredPackets, setFilteredPackets] = useState<PacketInfo[]>([]);
  const [selectedPacket, setSelectedPacket] = useState<PacketInfo | null>(null);
  const [filterProtocol, setFilterProtocol] = useState<string>('all');
  const [filterIP, setFilterIP] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [autoScroll, setAutoScroll] = useState(true);

  // 模拟抓包数据
  const mockPackets: PacketInfo[] = [
    {
      id: 0,
      metadata: {
        capture_time: "2025-10-09T16:49:42.595378+08:00",
        data_size: 1460,
        wire_length: 1460,
        capture_length: 1460,
        truncated: false,
        interface_index: 13,
        capture_length_ok: true
      },
      linkLayer: {
        timestamp: "2025-10-09T16:49:42.599304+08:00",
        src_mac: "10:5f:02:b0:2d:a8",
        dst_mac: "18:4a:53:24:af:87",
        eth_type: "IPv4",
        length: 14
      },
      networkLayer: {
        timestamp: "2025-10-09T16:49:42.599305+08:00",
        ip_version: 4,
        src_ip: "183.47.124.77",
        dst_ip: "192.168.1.21",
        protocol: "TCP",
        length: 1446,
        ttl: 52,
        ihl: 5,
        identifier: 62183,
        flags: 2,
        checksum: 22832,
        is_src_loopback: false,
        is_dst_loopback: false,
        is_src_link_local: false,
        is_dst_link_local: false,
        is_src_ip_valid: true,
        is_dst_ip_valid: true
      },
      transportLayer: {
        timestamp: "2025-10-09T16:49:42.599306+08:00",
        src_port: 80,
        dst_port: 49793,
        protocol: "TCP",
        seq_number: 1019211179,
        ack_number: 1649748342,
        window_size: 501,
        checksum: 6324,
        is_psh: true,
        is_ack: true
      },
      applicationLayer: {
        timestamp: "2025-10-09T16:49:42.599307+08:00",
        payload: "266FC/4cEKemsLU2YJUixgdQxJfU2samqyFTYT8/ehVE0lcICmzxBggyKTKPn0h6ABPaF7vXyC9p1kbl061p94SM9eBddHZI5uf73WbVpVuSz0Voq8WTDGAytNR1uhvGRbzwQfn3i3GgpxY1sxoJIHFnHuS2QzEc3tlChdXmSzmcWLq+4qHPymJGHqj5SVRAuY3x/ejY3kaaLJIjf4Rxg2YRCmeqe2U7v6ZdhW+07VEQivwLxgtINNJAPZtyZFJge1M0P1R8OASEQLUxp46a/pU11A0I5fxxk6yzckd+ykTbCW0s8MS5OvFL7LRTdx5sMryrDA2u3978hA2Ni6XpNNVKnwhjXZ4mLIUKCh2SIrd9+j46AZGKfAaQTphLPTYzyBBzFp6MDLz/NJgHA/RUVLW4rjVOjejJElA7yE59MXZcJx2qGKKjRQAJ5SUXw+9EOm+Fge9wxXUMs4uWKyNQMbww6Mkwb7+v5mgxSan8kXvQmmkcGLJGkKhoeNUIqx2p+BRsdZTSwd8FYFq35DogLtxLQYpqkQf8EXSqJcAYzgwweFPCuCVJnElWNM93fdu2W07m+QhaZyLH62+tDuiCxrf3fUb8WaZrajKb9AocF8xiRb+L2YV5AMDgOzrpLs/QxhDoMUNMgC2RmnfO0UbP6LBGQwmUGS1Cp8rwzz0gSx4ZEQLdI6Nk94fcFuig/hZqN20moEI2JUqd17sA0x/5LN/Kyv38cOKTKSKes9cVqZUODX83bcAG0xG9pDDGWe24i3Sq5S4p8l3j8RT0Ghq7R0RGPuGPcwPI6bjvyHyv8VZIHdjUdEG/x/UEJ4ih/cmAY+HjxRjgn706Mkvpg9IUH9PKkcEBvNtA0rBgF6pOg4FjHGYV4z8Ng6v/Sr3QldQZlzV250N4dZmldqQuhpFfNQGZuk/lj2XbHr1Prvwf+RvfT9OXos9RgSKlLprHd+GUUw0R5JRQKvkjTkRCeyzYcTIGrldK6ng56UuBXxnVLJ0m2agD3d4pWJ2OxTdKX4O7P7wdiHSY9yLt4SObgbqL768TXg5vrRoQ/w5Re3kgc9GYA2+pWcHlx9emdo8neGz9t5djEvlU5sb5+XQ6DFGhsJgR/fNc710BXhjbmfH2xsDMhsmS1Rj7XzvyJigQjzgqaoUqAGXtKoccW/WHvUUiXAykR9r7+qCzh3ukg6Rm+pKpNj8CR0Y4mXvu5p8cn4BIuJWl9RjfMNfRYRd8Bcx0HqfbVQh17g0EppTFOmMEQBRKhow6KrIyLk/9aaY8eMjOjGe3CjvRTG064wWGHx6sBffLCKx+n2Oe5WAV2GspLkUmLSycGCptuv0CP9ogNTwZ/gzXay69s96rjVg4itl5R478GKdB/CijYbEDSgHiPwJwd19ZcJrwY9k+y6pwC+XAJAUjcolt0/Bl1Gp+oaRCsAhq+/ey0RcuqWGtdd2ggsz+wjmhpw5pmm4zGcLt5/4frU5JCEjmRbwR4swDvqasvocmtz4VwhKG6nr7ntV1yoRph8EkuhOBEfLQPNFKeijsQlb/LtHbvd6/H/sb4+MwNwAJHeQp8HEnZA7ibbRr8SdbJfAyHf9G+Mlz1UAlCu966jbyPStFfUz7zGnyhLge0I/pdjc3xKTEJ2q4tsAPvXam7J56iqVNZWnWPxVse1AjBQSZ/6c5qIrgnis38zBf0b0kcJkDZIxR+52EeUM+iNNKpQ2T3yiYGuFFAUW9rbBcbU1pqYRPTCL5BHycLAZmNdZKwdGSx0WAiMqOgn2FSrPhsYCMisB1uDzVYItwGvGSsGksRqDD9Kr6Jz5xhqgsCt5RA34U2nQoCHrPcgOMNGX4TBXxBAAXgwBFNL6hKmPke2vYQm4MkQp53dQGfBk="
      },
      errorLayer: null
    }
  ];

  useEffect(() => {
    setPackets(mockPackets);
    setFilteredPackets(mockPackets);
  }, []);

  useEffect(() => {
    let filtered = packets;
    
    if (filterProtocol !== 'all') {
      filtered = filtered.filter(packet => packet.networkLayer.protocol === filterProtocol);
    }
    
    if (filterIP !== 'all') {
      filtered = filtered.filter(packet => 
        packet.networkLayer.src_ip.includes(filterIP) ||
        packet.networkLayer.dst_ip.includes(filterIP)
      );
    }
    
    if (searchQuery) {
      filtered = filtered.filter(packet => 
        packet.networkLayer.src_ip.toLowerCase().includes(searchQuery.toLowerCase()) ||
        packet.networkLayer.dst_ip.toLowerCase().includes(searchQuery.toLowerCase()) ||
        packet.transportLayer.src_port.toString().includes(searchQuery) ||
        packet.transportLayer.dst_port.toString().includes(searchQuery)
      );
    }
    
    setFilteredPackets(filtered);
  }, [packets, filterProtocol, filterIP, searchQuery]);

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

  const clearPackets = () => {
    setPackets([]);
    setFilteredPackets([]);
    setSelectedPacket(null);
  };

  const getProtocolColor = (protocol: string) => {
    switch (protocol) {
      case 'TCP': return 'text-blue-400';
      case 'UDP': return 'text-green-400';
      case 'ICMP': return 'text-yellow-400';
      case 'HTTP': return 'text-purple-400';
      case 'HTTPS': return 'text-pink-400';
      default: return 'text-gray-400';
    }
  };

  const getPortColor = (port: number) => {
    if (port === 80 || port === 443) return 'text-blue-400';
    if (port === 22 || port === 23) return 'text-green-400';
    if (port === 21 || port === 25) return 'text-yellow-400';
    if (port >= 1024) return 'text-purple-400';
    return 'text-gray-400';
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
                onClick={clearPackets}
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
                  value={filterProtocol}
                  onChange={(e) => setFilterProtocol(e.target.value)}
                  className="px-3 py-2 bg-background/50 border border-primary/30 rounded"
                >
                  <option value="all">All Protocols</option>
                  <option value="TCP">TCP</option>
                  <option value="UDP">UDP</option>
                  <option value="ICMP">ICMP</option>
                </select>

                <select
                  value={filterIP}
                  onChange={(e) => setFilterIP(e.target.value)}
                  className="px-3 py-2 bg-background/50 border border-primary/30 rounded"
                >
                  <option value="all">All IPs</option>
                  <option value="192.168">Local Network</option>
                  <option value="10.">Private Network</option>
                  <option value="172.">Private Network</option>
                </select>
              </div>

              {/* 数据包列表 */}
              <div className="space-y-1 max-h-80 overflow-y-auto scrollbar-ultra-thin">
                {filteredPackets.map((packet) => (
                  <div
                    key={packet.id}
                    className={twMerge([
                      "p-4 rounded-lg transition-all duration-200 hover:bg-primary/10 border",
                      selectedPacket?.id === packet.id 
                        ? "bg-primary/20 border-primary/50 shadow-lg shadow-primary/20" 
                        : "border-primary/20 hover:border-primary/40"
                    ])}
                  >
                    <div className="flex items-center justify-between mb-2">
                      <div 
                        className="flex items-center gap-3 cursor-pointer flex-1"
                        onClick={() => setSelectedPacket(packet)}
                      >
                        <span className={twMerge(
                          "font-mono text-xs font-bold px-2 py-1 rounded",
                          getProtocolColor(packet.networkLayer.protocol),
                          "bg-current/10"
                        )}>
                          {packet.networkLayer.protocol}
                        </span>
                        <span className="text-xs font-mono opacity-70">
                          {new Date(packet.metadata.capture_time).toLocaleTimeString()}
                        </span>
                        <span className={twMerge(
                          "text-xs font-bold px-2 py-1 rounded",
                          getPortColor(packet.transportLayer.src_port),
                          "bg-current/10"
                        )}>
                          {packet.transportLayer.src_port}
                        </span>
                      </div>
                      <div className="flex items-center gap-3">
                        <div className="flex items-center gap-2 text-xs opacity-70">
                          <span className="font-mono">{packet.metadata.data_size}B</span>
                          <span className="font-mono">TTL:{packet.networkLayer.ttl}</span>
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
                                  <span className={twMerge("font-mono px-2 py-1 rounded text-xs", getProtocolColor(packet.networkLayer.protocol), "bg-current/10")}>
                                    {packet.networkLayer.protocol}
                                  </span>
                                  <span className="truncate">{packet.networkLayer.src_ip}:{packet.transportLayer.src_port} → {packet.networkLayer.dst_ip}:{packet.transportLayer.dst_port}</span>
                                </DialogTitle>
                                <DialogDescription className="text-sm opacity-70 mb-4 font-inter">
                                  {packet.networkLayer.src_ip} → {packet.networkLayer.dst_ip} • {new Date(packet.metadata.capture_time).toLocaleString()} • {packet.metadata.data_size}B
                                </DialogDescription>
                              </div>
                              
                              <div className="flex-1 overflow-y-auto scrollbar-ultra-thin space-y-6 pr-2">
                                {/* 元数据信息 */}
                                <div>
                                  <h4 className="text-sm font-bold mb-2 text-primary font-inter">数据包元信息</h4>
                                  <div className="grid grid-cols-2 gap-4 text-xs">
                                    <div>
                                      <span className="opacity-70 font-inter">捕获时间:</span>
                                      <span className="ml-2 font-jetbrains">{new Date(packet.metadata.capture_time).toLocaleString()}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">数据大小:</span>
                                      <span className="ml-2 font-jetbrains">{packet.metadata.data_size}B</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">接口索引:</span>
                                      <span className="ml-2 font-jetbrains">{packet.metadata.interface_index}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">是否截断:</span>
                                      <span className="ml-2 font-jetbrains">{packet.metadata.truncated ? '是' : '否'}</span>
                                    </div>
                                  </div>
                                </div>

                                {/* 网络层信息 */}
                                <div>
                                  <h4 className="text-sm font-bold mb-2 text-primary font-inter">网络层信息</h4>
                                  <div className="grid grid-cols-2 gap-4 text-xs">
                                    <div>
                                      <span className="opacity-70 font-inter">源IP:</span>
                                      <span className="ml-2 font-jetbrains">{packet.networkLayer.src_ip}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">目标IP:</span>
                                      <span className="ml-2 font-jetbrains">{packet.networkLayer.dst_ip}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">协议:</span>
                                      <span className="ml-2 font-jetbrains">{packet.networkLayer.protocol}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">TTL:</span>
                                      <span className="ml-2 font-jetbrains">{packet.networkLayer.ttl}</span>
                                    </div>
                                  </div>
                                </div>

                                {/* 传输层信息 */}
                                <div>
                                  <h4 className="text-sm font-bold mb-2 text-primary font-inter">传输层信息</h4>
                                  <div className="grid grid-cols-2 gap-4 text-xs">
                                    <div>
                                      <span className="opacity-70 font-inter">源端口:</span>
                                      <span className="ml-2 font-jetbrains">{packet.transportLayer.src_port}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">目标端口:</span>
                                      <span className="ml-2 font-jetbrains">{packet.transportLayer.dst_port}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">序列号:</span>
                                      <span className="ml-2 font-jetbrains">{packet.transportLayer.seq_number}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">确认号:</span>
                                      <span className="ml-2 font-jetbrains">{packet.transportLayer.ack_number}</span>
                                    </div>
                                  </div>
                                </div>

                                {/* 应用层数据 */}
                                {packet.applicationLayer.payload && (
                                  <div>
                                    <h4 className="text-sm font-bold mb-2 text-primary font-inter">应用层载荷</h4>
                                    <div className="bg-background/30 rounded p-3 max-h-32 overflow-y-auto scrollbar-ultra-thin">
                                      <pre className="text-xs font-jetbrains whitespace-pre-wrap">
                                        {packet.applicationLayer.payload}
                                      </pre>
                                    </div>
                                  </div>
                                )}

                                {/* 链路层信息 */}
                                <div>
                                  <h4 className="text-sm font-bold mb-2 text-primary font-inter">链路层信息</h4>
                                  <div className="grid grid-cols-2 gap-4 text-xs">
                                    <div>
                                      <span className="opacity-70 font-inter">源MAC:</span>
                                      <span className="ml-2 font-jetbrains">{packet.linkLayer.src_mac}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">目标MAC:</span>
                                      <span className="ml-2 font-jetbrains">{packet.linkLayer.dst_mac}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">以太网类型:</span>
                                      <span className="ml-2 font-jetbrains">{packet.linkLayer.eth_type}</span>
                                    </div>
                                    <div>
                                      <span className="opacity-70 font-inter">长度:</span>
                                      <span className="ml-2 font-jetbrains">{packet.linkLayer.length}B</span>
                                    </div>
                                  </div>
                                </div>

                                {/* 错误层信息 */}
                                {packet.errorLayer && (
                                  <div>
                                    <h4 className="text-sm font-bold mb-2 text-primary font-inter">错误信息</h4>
                                    <div className="grid grid-cols-2 gap-4 text-xs">
                                      <div>
                                        <span className="opacity-70 font-inter">错误类型:</span>
                                        <span className="ml-2 font-jetbrains">{packet.errorLayer.error}</span>
                                      </div>
                                      <div>
                                        <span className="opacity-70 font-inter">错误层:</span>
                                        <span className="ml-2 font-jetbrains">{packet.errorLayer.layer}</span>
                                      </div>
                                      <div>
                                        <span className="opacity-70 font-inter">是否致命:</span>
                                        <span className="ml-2 font-jetbrains">{packet.errorLayer.fatal ? '是' : '否'}</span>
                                      </div>
                                      <div>
                                        <span className="opacity-70 font-inter">错误代码:</span>
                                        <span className="ml-2 font-jetbrains">{packet.errorLayer.code}</span>
                                      </div>
                                    </div>
                                  </div>
                                )}
                              </div>
                              
                              <DialogCloseTrigger />
                            </DialogContent>
                          </DialogPositioner>
                        </DialogRoot>
                      </div>
                    </div>
                    <div 
                      className="text-sm font-mono truncate mb-1 cursor-pointer"
                      onClick={() => setSelectedPacket(packet)}
                    >
                      {packet.networkLayer.src_ip}:{packet.transportLayer.src_port} → {packet.networkLayer.dst_ip}:{packet.transportLayer.dst_port}
                    </div>
                    <div 
                      className="text-xs opacity-50 font-mono cursor-pointer"
                      onClick={() => setSelectedPacket(packet)}
                    >
                      {packet.linkLayer.src_mac} → {packet.linkLayer.dst_mac}
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
                  <span className="text-sm opacity-70 font-inter">总数据包数</span>
                  <span className="text-xl font-bold text-primary font-jetbrains">{packets.length}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm opacity-70 font-inter">TCP包数</span>
                  <span className="text-xl font-bold text-blue-400 font-jetbrains">
                    {packets.filter(p => p.networkLayer.protocol === 'TCP').length}
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm opacity-70 font-inter">UDP包数</span>
                  <span className="text-xl font-bold text-green-400 font-jetbrains">
                    {packets.filter(p => p.networkLayer.protocol === 'UDP').length}
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm opacity-70 font-inter">平均包大小</span>
                  <span className="text-xl font-bold text-accent font-jetbrains">
                    {packets.length > 0 ? Math.round(packets.reduce((sum, p) => sum + p.metadata.data_size, 0) / packets.length) : 0}B
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* 数据包详情 */}
          {selectedPacket && (
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
                  数据包详情
                </h3>
                
                <div className="space-y-3 text-sm">
                  <div>
                    <span className="opacity-70 font-inter">协议:</span>
                    <span className={twMerge("ml-2 font-jetbrains font-bold", getProtocolColor(selectedPacket.networkLayer.protocol))}>
                      {selectedPacket.networkLayer.protocol}
                    </span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">源IP:</span>
                    <span className="ml-2 font-jetbrains">{selectedPacket.networkLayer.src_ip}</span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">目标IP:</span>
                    <span className="ml-2 font-jetbrains">{selectedPacket.networkLayer.dst_ip}</span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">源端口:</span>
                    <span className="ml-2 font-jetbrains">{selectedPacket.transportLayer.src_port}</span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">目标端口:</span>
                    <span className="ml-2 font-jetbrains">{selectedPacket.transportLayer.dst_port}</span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">数据大小:</span>
                    <span className="ml-2 font-jetbrains">{selectedPacket.metadata.data_size}B</span>
                  </div>
                  <div>
                    <span className="opacity-70 font-inter">捕获时间:</span>
                    <span className="ml-2 font-jetbrains">{new Date(selectedPacket.metadata.capture_time).toLocaleString()}</span>
                  </div>
                </div>
                
                <div className="mt-4">
                  <span className="opacity-70 text-sm font-inter">连接信息:</span>
                  <div className="mt-1 p-2 bg-background/30 rounded text-xs break-all font-jetbrains">
                    {selectedPacket.networkLayer.src_ip}:{selectedPacket.transportLayer.src_port} → {selectedPacket.networkLayer.dst_ip}:{selectedPacket.transportLayer.dst_port}
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* 警告提示 */}
      {!isCapturing && packets.length === 0 && (
        <AlertRoot className="mt-6">
          <AlertTitle>准备开始抓包</AlertTitle>
          <AlertDescription>
            点击"Start Capture"按钮开始捕获网络数据包。确保已正确配置网络接口。
          </AlertDescription>
          <AlertCloseTrigger />
        </AlertRoot>
      )}
    </div>
  );
}

export default CaptureHome;
