package proxy

import (
	"fmt"
	"os/exec"
	"runtime"
)

// InstallCert 尝试自动安装证书到系统信任区
func InstallCert(certPath string) error {
	switch runtime.GOOS {
	case "darwin":
		return installCertMac(certPath)
	case "windows":
		return installCertWindows(certPath)
	case "linux":
		return installCertLinux(certPath)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// installCertMac macOS使用security命令安装
func installCertMac(certPath string) error {
	cmd := exec.Command("security", "add-trusted-cert",
		"-d", "-r", "trustRoot",
		"-k", "/Library/Keychains/System.keychain",
		certPath)
	return cmd.Run()
}

// installCertWindows Windows使用certutil安装
func installCertWindows(certPath string) error {
	cmd := exec.Command("certutil", "-addstore", "-f", "ROOT", certPath)
	return cmd.Run()
}

// installCertLinux Linux复制到系统证书目录并更新
func installCertLinux(certPath string) error {
	// Ubuntu/Debian
	cmd := exec.Command("sh", "-c",
		fmt.Sprintf("cp %s /usr/local/share/ca-certificates/probe-ca.crt && update-ca-certificates", certPath))
	err := cmd.Run()
	if err == nil {
		return nil
	}

	// CentOS/RHEL
	cmd = exec.Command("sh", "-c",
		fmt.Sprintf("cp %s /etc/pki/ca-trust/source/anchors/probe-ca.crt && update-ca-trust", certPath))
	return cmd.Run()
}

// GetInstallInstructions 获取各平台手动安装指引
func GetInstallInstructions(osType string) map[string]interface{} {
	instructions := map[string]map[string]interface{}{
		"darwin": {
			"title": "macOS 手动安装步骤",
			"steps": []string{
				"方法1（推荐）：",
				"1. 双击下载的 proxy_root_ca.pem 文件",
				"2. 在弹出的钥匙串访问中，选择\"登录\"钥匙串（不是系统）",
				"3. 找到 CetiProbe Root CA 证书，双击打开",
				"4. 展开\"信任\"部分，将\"使用此证书时\"设置为\"始终信任\"",
				"5. 关闭窗口，输入密码确认",
				"",
				"方法2（如果方法1失败）：",
				"1. 打开终端，运行下方命令",
				"2. 输入管理员密码",
				"3. 重启浏览器",
			},
			"command": "sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain proxy_root_ca.pem",
		},
		"windows": {
			"title": "Windows 手动安装步骤",
			"steps": []string{
				"1. 右键点击下载的证书文件，选择\"安装证书\"",
				"2. 选择\"本地计算机\"，点击下一步",
				"3. 选择\"将所有证书都放入下列存储\"",
				"4. 点击\"浏览\"，选择\"受信任的根证书颁发机构\"",
				"5. 点击完成，在安全警告中选择\"是\"",
			},
			"command": "certutil -addstore -f ROOT proxy_root_ca.pem",
		},
		"linux": {
			"title": "Linux 手动安装步骤",
			"steps": []string{
				"Ubuntu/Debian:",
				"  sudo cp proxy_root_ca.pem /usr/local/share/ca-certificates/probe-ca.crt",
				"  sudo update-ca-certificates",
				"",
				"CentOS/RHEL:",
				"  sudo cp proxy_root_ca.pem /etc/pki/ca-trust/source/anchors/probe-ca.crt",
				"  sudo update-ca-trust",
			},
			"command": "",
		},
	}

	if inst, ok := instructions[osType]; ok {
		return inst
	}
	return instructions["darwin"]
}
