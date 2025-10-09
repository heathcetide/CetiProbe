// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use tauri::{Manager, State};
use serde::{Deserialize, Serialize};
use std::process::{Command, Stdio};
use std::path::Path;

#[derive(Debug, Serialize, Deserialize)]
struct AppState {
    theme: String,
    window_title: String,
}

// Learn more about Tauri commands at https://tauri.app/v2/guides/features/command
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

#[tauri::command]
fn get_app_info() -> serde_json::Value {
    serde_json::json!({
        "name": "Cetiprobe",
        "version": "0.0.0",
        "description": "A powerful network analysis and packet capture tool"
    })
}

#[tauri::command]
fn set_theme(theme: &str, _state: State<AppState>) -> Result<(), String> {
    println!("Setting theme to: {}", theme);
    // Here you could implement theme switching logic
    Ok(())
}

#[tauri::command]
fn get_theme(state: State<AppState>) -> String {
    state.theme.clone()
}

#[tauri::command]
async fn export_data() -> Result<String, String> {
    // Implement data export logic here
    Ok("Data exported successfully".to_string())
}

#[tauri::command]
async fn import_data() -> Result<String, String> {
    // Implement data import logic here
    Ok("Data imported successfully".to_string())
}

#[tauri::command]
async fn check_backend_status() -> Result<bool, String> {
    // 检查后端服务是否运行
    match reqwest::get("http://localhost:8081").await {
        Ok(response) => Ok(response.status().is_success()),
        Err(_) => Ok(false),
    }
}

fn start_backend_server() {
    // 检查 Go 是否安装
    let go_available = Command::new("go")
        .arg("version")
        .output()
        .is_ok();
    
    if !go_available {
        println!("Warning: Go is not installed or not in PATH. Backend server will not start.");
        return;
    }
    
    // 检查 server 目录是否存在
    let server_path = Path::new("../server");
    if !server_path.exists() {
        println!("Warning: Server directory not found. Backend server will not start.");
        return;
    }
    
    // 启动 Go 后端服务
    let mut child = match Command::new("go")
        .arg("run")
        .arg("cmd/main.go")
        .current_dir("../server")
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
    {
        Ok(child) => {
            println!("Go backend server started on port 8080");
            child
        }
        Err(e) => {
            println!("Failed to start Go backend server: {}", e);
            return;
        }
    };
    
    // 在后台运行，不等待进程结束
    std::thread::spawn(move || {
        let _ = child.wait();
    });
}

fn main() {
        let app_state = AppState {
            theme: "dark".to_string(),
            window_title: "Cetiprobe".to_string(),
        };

    tauri::Builder::default()
        .manage(app_state)
        .invoke_handler(tauri::generate_handler![
            greet,
            get_app_info,
            set_theme,
            get_theme,
            export_data,
            import_data,
            check_backend_status
        ])
        .setup(|app| {
            let window = app.get_webview_window("main").unwrap();
            
            // Set window properties
            window.set_title("Cetiprobe").unwrap();
            
            // 启动 Go 后端服务
            start_backend_server();
            
            println!("Cetiprobe application started!");
            
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
