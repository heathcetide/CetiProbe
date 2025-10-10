import { invoke } from '@tauri-apps/api/core';
import { getCurrentWebviewWindow } from '@tauri-apps/api/webviewWindow';
import { Button } from '@/components/ui/button';
import { twMerge } from 'tailwind-merge';
import { 
  Minimize2, 
  Maximize2, 
  X, 
  Monitor, 
  Settings,
  Download,
  Upload,
  Zap,
  Github,
  Home,
  BookOpen,
  Palette,
  Network
} from 'lucide-react';
import { Link } from 'react-router';

interface UnifiedHeaderProps {
  className?: string;
  showMenu: boolean;
  setShowMenu: (show: boolean) => void;
}

function UnifiedHeader({ className, showMenu, setShowMenu }: UnifiedHeaderProps) {
  const handleMinimize = async () => {
    const window = getCurrentWebviewWindow();
    await window.minimize();
  };

  const handleMaximize = async () => {
    const window = getCurrentWebviewWindow();
    const isMaximized = await window.isMaximized();
    if (isMaximized) {
      await window.unmaximize();
    } else {
      await window.maximize();
    }
  };

  const handleClose = async () => {
    const window = getCurrentWebviewWindow();
    await window.close();
  };

  const handleFullscreen = async () => {
    const window = getCurrentWebviewWindow();
    const isFullscreen = await window.isFullscreen();
    await window.setFullscreen(!isFullscreen);
  };

  const handleExportData = async () => {
    try {
      const result = await invoke('export_data');
      console.log('Data exported:', result);
    } catch (error) {
      console.error('Export failed:', error);
    }
  };

  const handleImportData = async () => {
    try {
      const result = await invoke('import_data');
      console.log('Data imported:', result);
    } catch (error) {
      console.error('Import failed:', error);
    }
  };

  return (
    <div className={twMerge([
      "flex items-center justify-between px-6 py-3 bg-background/90 backdrop-blur-md border-b border-primary/20",
      "transition-all duration-300 ease-out",
      className
    ])}>
      {/* 左侧：应用标题和导航 */}
      <div className="flex items-center gap-6">
        {/* 应用标题 */}
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-gradient-to-br from-primary to-accent rounded-lg flex items-center justify-center shadow-lg">
            <span className="text-sm font-bold text-white">C</span>
          </div>
          <div>
            <h1 className="text-lg font-bold text-foreground">Cetiprobe</h1>
            <p className="text-xs text-muted-foreground">Network Analysis Tool</p>
          </div>
        </div>

        {/* 导航菜单 */}
        <nav className="hidden lg:flex items-center gap-1">
          <Link
            to="/"
            className="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-primary/10 transition-all duration-200"
          >
            <Home className="w-4 h-4" />
            Home
          </Link>
          <Link
            to="/docs"
            className="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-primary/10 transition-all duration-200"
          >
            <BookOpen className="w-4 h-4" />
            Docs
          </Link>
          <Link
            to="/docs/frame"
            className="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-primary/10 transition-all duration-200"
          >
            <Palette className="w-4 h-4" />
            Components
          </Link>
          <Link
            to="/capture"
            className="flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-primary/10 transition-all duration-200"
          >
            <Network className="w-4 h-4" />
            Capture
          </Link>
        </nav>
      </div>

      {/* 右侧：功能按钮和窗口控制 */}
      <div className="flex items-center gap-2">
        {/* 功能按钮 */}
        <div className="hidden sm:flex items-center gap-1">
          <Button
            variant="secondary"
            shape="flat"
            onClick={handleImportData}
            className="px-3 py-1 text-xs"
          >
            <Upload className="w-4 h-4" />
          </Button>
          
          <Button
            variant="secondary"
            shape="flat"
            onClick={handleExportData}
            className="px-3 py-1 text-xs"
          >
            <Download className="w-4 h-4" />
          </Button>

          <Button
            variant="secondary"
            shape="flat"
            onClick={handleFullscreen}
            className="px-3 py-1 text-xs"
          >
            <Monitor className="w-4 h-4" />
          </Button>

          <Button
            variant="secondary"
            shape="flat"
            className="px-3 py-1 text-xs"
          >
            <Settings className="w-4 h-4" />
          </Button>
        </div>

        {/* GitHub链接 */}
        <a
          href="https://github.com/rizkimuhammada/cosmic-ui"
          target="_blank"
          rel="noopener noreferrer"
          className="hidden sm:block"
        >
          <Button
            variant="accent"
            shape="flat"
            className="px-3 py-1 text-xs"
          >
            <Github className="w-4 h-4" />
          </Button>
        </a>

        {/* 移动端菜单按钮 */}
        <Button
          variant="default"
          shape="flat"
          onClick={() => setShowMenu(!showMenu)}
          className="lg:hidden px-3 py-1 text-xs"
        >
          <Zap className="w-4 h-4" />
        </Button>

        {/* 窗口控制按钮 */}
        <div className="flex items-center gap-1 ml-2">
          <Button
            variant="secondary"
            shape="flat"
            onClick={handleMinimize}
            className="px-2 py-1 text-xs hover:bg-yellow-500/20"
          >
            <Minimize2 className="w-4 h-4" />
          </Button>
          
          <Button
            variant="secondary"
            shape="flat"
            onClick={handleMaximize}
            className="px-2 py-1 text-xs hover:bg-green-500/20"
          >
            <Maximize2 className="w-4 h-4" />
          </Button>
          
          <Button
            variant="destructive"
            shape="flat"
            onClick={handleClose}
            className="px-2 py-1 text-xs"
          >
            <X className="w-4 h-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}

export { UnifiedHeader };
