import { useState, useEffect } from 'react';
import { invoke } from '@tauri-apps/api/core';
import { getCurrentWebviewWindow } from '@tauri-apps/api/webviewWindow';

interface DesktopState {
  isMaximized: boolean;
  isFullscreen: boolean;
  theme: string;
  windowTitle: string;
}

export function useDesktop() {
  const [state, setState] = useState<DesktopState>({
    isMaximized: false,
    isFullscreen: false,
    theme: 'dark',
    windowTitle: 'Cosmic UI Desktop'
  });

  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initializeDesktop = async () => {
      try {
        // 获取应用信息
        const appInfo = await invoke('get_app_info');
        console.log('App info:', appInfo);

        // 获取当前主题
        const theme = await invoke('get_theme');
        
        // 检查窗口状态
        const window = getCurrentWebviewWindow();
        const isMaximized = await window.isMaximized();
        const isFullscreen = await window.isFullscreen();

        setState({
          isMaximized,
          isFullscreen,
          theme: theme as string,
          windowTitle: 'Cosmic UI Desktop'
        });

        setIsLoading(false);
      } catch (error) {
        console.error('Failed to initialize desktop:', error);
        setIsLoading(false);
      }
    };

    initializeDesktop();
  }, []);

  const setTheme = async (theme: string) => {
    try {
      await invoke('set_theme', { theme });
      setState(prev => ({ ...prev, theme }));
    } catch (error) {
      console.error('Failed to set theme:', error);
    }
  };

  const toggleMaximize = async () => {
    try {
      const window = getCurrentWebviewWindow();
      const isMaximized = await window.isMaximized();
      
      if (isMaximized) {
        await window.unmaximize();
      } else {
        await window.maximize();
      }
      
      setState(prev => ({ ...prev, isMaximized: !isMaximized }));
    } catch (error) {
      console.error('Failed to toggle maximize:', error);
    }
  };

  const toggleFullscreen = async () => {
    try {
      const window = getCurrentWebviewWindow();
      const isFullscreen = await window.isFullscreen();
      
      await window.setFullscreen(!isFullscreen);
      setState(prev => ({ ...prev, isFullscreen: !isFullscreen }));
    } catch (error) {
      console.error('Failed to toggle fullscreen:', error);
    }
  };

  const minimize = async () => {
    try {
      const window = getCurrentWebviewWindow();
      await window.minimize();
    } catch (error) {
      console.error('Failed to minimize:', error);
    }
  };

  const close = async () => {
    try {
      const window = getCurrentWebviewWindow();
      await window.close();
    } catch (error) {
      console.error('Failed to close:', error);
    }
  };

  return {
    ...state,
    isLoading,
    setTheme,
    toggleMaximize,
    toggleFullscreen,
    minimize,
    close
  };
}
