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
  Upload
} from 'lucide-react';

interface DesktopMenuProps {
  className?: string;
}

function DesktopMenu({ className }: DesktopMenuProps) {
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
      // 这里可以调用 Tauri 命令来导出数据
      const result = await invoke('export_data');
      console.log('Data exported:', result);
    } catch (error) {
      console.error('Export failed:', error);
    }
  };

  const handleImportData = async () => {
    try {
      // 这里可以调用 Tauri 命令来导入数据
      const result = await invoke('import_data');
      console.log('Data imported:', result);
    } catch (error) {
      console.error('Import failed:', error);
    }
  };

  return (
    <div className={twMerge([
      "flex items-center gap-2 px-4 py-2 bg-background/80 backdrop-blur-sm border-b border-primary/20",
      className
    ])}>
      {/* 应用标题 */}
      <div className="flex items-center gap-2 mr-auto">
        <div className="w-6 h-6 bg-primary/20 rounded-sm flex items-center justify-center">
          <span className="text-xs font-bold text-primary">C</span>
        </div>
        <span className="text-sm font-medium text-foreground">Cetiprobe</span>
      </div>

      {/* 桌面功能按钮 */}
      <div className="flex items-center gap-1">
        <Button
          shape="flat"
          onClick={handleImportData}
          className="h-8 w-8 p-0"
        >
          <Upload className="w-4 h-4" />
        </Button>
        
        <Button
          shape="flat"
          onClick={handleExportData}
          className="h-8 w-8 p-0"
        >
          <Download className="w-4 h-4" />
        </Button>

        <Button
          shape="flat"
          onClick={handleFullscreen}
          className="h-8 w-8 p-0"
        >
          <Monitor className="w-4 h-4" />
        </Button>

        <Button
          shape="flat"
          className="h-8 w-8 p-0"
        >
          <Settings className="w-4 h-4" />
        </Button>
      </div>

      {/* 窗口控制按钮 */}
      <div className="flex items-center gap-1 ml-4">
        <Button
          shape="flat"
          onClick={handleMinimize}
          className="h-8 w-8 p-0 hover:bg-yellow-500/20"
        >
          <Minimize2 className="w-4 h-4" />
        </Button>
        
        <Button
          shape="flat"
          onClick={handleMaximize}
          className="h-8 w-8 p-0 hover:bg-green-500/20"
        >
          <Maximize2 className="w-4 h-4" />
        </Button>
        
        <Button
          shape="flat"
          onClick={handleClose}
          className="h-8 w-8 p-0 hover:bg-red-500/20"
        >
          <X className="w-4 h-4" />
        </Button>
      </div>
    </div>
  );
}

export { DesktopMenu };
