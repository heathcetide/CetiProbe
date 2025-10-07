class ProbeApp {
    constructor() {
        this.packets = [];
        this.filteredPackets = [];
        this.currentFilter = {};
        this.autoRefresh = true;
        this.refreshInterval = null;
        
        this.init();
    }
    
    init() {
        this.bindEvents();
        this.loadStats();
        this.loadPackets();
        this.startAutoRefresh();
    }
    
    bindEvents() {
        // 过滤按钮
        document.getElementById('filter-btn').addEventListener('click', () => {
            this.applyFilter();
        });
        
        // 重置过滤按钮
        document.getElementById('reset-filter-btn').addEventListener('click', () => {
            this.resetFilter();
        });
        
        // 清空数据按钮
        document.getElementById('clear-btn').addEventListener('click', () => {
            this.clearData();
        });
        
        // 导出数据按钮
        document.getElementById('export-btn').addEventListener('click', () => {
            this.exportData();
        });
        
        // 回车键应用过滤
        document.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && e.target.closest('.filter-row')) {
                this.applyFilter();
            }
        });
    }
    
    async loadStats() {
        try {
            const response = await fetch('/api/stats');
            const stats = await response.json();
            
            document.getElementById('total-packets').textContent = stats.total_packets || 0;
            document.getElementById('http-packets').textContent = stats.http_packets || 0;
            document.getElementById('https-packets').textContent = stats.https_packets || 0;
            document.getElementById('unique-ips').textContent = stats.unique_ips || 0;
        } catch (error) {
            console.error('加载统计信息失败:', error);
        }
    }
    
    async loadPackets() {
        try {
            const response = await fetch('/api/packets?limit=1000');
            const data = await response.json();
            
            this.packets = data.packets || [];
            this.filteredPackets = [...this.packets];
            this.renderPackets();
        } catch (error) {
            console.error('加载数据包失败:', error);
            this.showError('加载数据包失败: ' + error.message);
        }
    }
    
    async applyFilter() {
        const filter = {
            src_ip: document.getElementById('src-ip').value.trim(),
            dst_ip: document.getElementById('dst-ip').value.trim(),
            port: document.getElementById('port').value ? parseInt(document.getElementById('port').value) : 0,
            protocol: document.getElementById('protocol').value,
            http_method: document.getElementById('http-method').value,
            search_text: document.getElementById('search-text').value.trim()
        };
        
        // 移除空值
        Object.keys(filter).forEach(key => {
            if (filter[key] === '' || filter[key] === 0) {
                delete filter[key];
            }
        });
        
        this.currentFilter = filter;
        
        try {
            const response = await fetch('/api/packets/filter', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(filter)
            });
            
            const data = await response.json();
            this.filteredPackets = data.packets || [];
            this.renderPackets();
        } catch (error) {
            console.error('应用过滤失败:', error);
            this.showError('应用过滤失败: ' + error.message);
        }
    }
    
    resetFilter() {
        document.getElementById('src-ip').value = '';
        document.getElementById('dst-ip').value = '';
        document.getElementById('port').value = '';
        document.getElementById('protocol').value = '';
        document.getElementById('http-method').value = '';
        document.getElementById('search-text').value = '';
        
        this.currentFilter = {};
        this.filteredPackets = [...this.packets];
        this.renderPackets();
    }
    
    async clearData() {
        if (!confirm('确定要清空所有数据吗？此操作不可撤销。')) {
            return;
        }
        
        try {
            const response = await fetch('/api/clear', {
                method: 'POST'
            });
            
            if (response.ok) {
                this.packets = [];
                this.filteredPackets = [];
                this.renderPackets();
                this.loadStats();
                this.showMessage('数据已清空');
            } else {
                throw new Error('清空数据失败');
            }
        } catch (error) {
            console.error('清空数据失败:', error);
            this.showError('清空数据失败: ' + error.message);
        }
    }
    
    async exportData() {
        try {
            const response = await fetch('/api/export?format=json');
            
            if (response.ok) {
                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = `probe-packets-${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.json`;
                document.body.appendChild(a);
                a.click();
                document.body.removeChild(a);
                window.URL.revokeObjectURL(url);
            } else {
                throw new Error('导出失败');
            }
        } catch (error) {
            console.error('导出数据失败:', error);
            this.showError('导出数据失败: ' + error.message);
        }
    }
    
    renderPackets() {
        const tbody = document.getElementById('packets-tbody');
        
        if (this.filteredPackets.length === 0) {
            tbody.innerHTML = '<tr><td colspan="9" class="empty">暂无数据包</td></tr>';
            return;
        }
        
        tbody.innerHTML = this.filteredPackets.map(packet => {
            const time = new Date(packet.timestamp).toLocaleString('zh-CN');
            const srcAddr = `${packet.src_ip}:${packet.src_port}`;
            const dstAddr = `${packet.dst_ip}:${packet.dst_port}`;
            const protocolClass = `protocol-${packet.protocol.toLowerCase()}`;
            const methodClass = packet.http_method ? `http-${packet.http_method.toLowerCase()}` : '';
            const statusClass = packet.http_status ? this.getStatusClass(packet.http_status) : '';
            
            return `
                <tr>
                    <td class="time-cell">${time}</td>
                    <td class="ip-cell">${srcAddr}</td>
                    <td class="ip-cell">${dstAddr}</td>
                    <td class="protocol-cell">
                        <span class="protocol-tcp ${protocolClass}">${packet.protocol}</span>
                    </td>
                    <td class="length-cell">${packet.length}</td>
                    <td class="http-method ${methodClass}">${packet.http_method || '-'}</td>
                    <td class="url-cell" title="${packet.http_url || ''}">${packet.http_url || '-'}</td>
                    <td class="status-cell ${statusClass}">${packet.http_status || '-'}</td>
                    <td class="user-agent-cell" title="${packet.user_agent || ''}">${packet.user_agent || '-'}</td>
                </tr>
            `;
        }).join('');
    }
    
    getStatusClass(status) {
        const code = parseInt(status);
        if (code >= 200 && code < 300) return 'status-200';
        if (code >= 300 && code < 400) return 'status-300';
        if (code >= 400 && code < 500) return 'status-400';
        if (code >= 500) return 'status-500';
        return '';
    }
    
    startAutoRefresh() {
        this.refreshInterval = setInterval(() => {
            this.loadStats();
            this.loadPackets();
        }, 2000); // 每2秒刷新一次
    }
    
    stopAutoRefresh() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
        }
    }
    
    showMessage(message) {
        // 简单的消息提示，可以后续改进为更好的UI
        alert(message);
    }
    
    showError(message) {
        console.error(message);
        // 简单的错误提示，可以后续改进为更好的UI
        alert('错误: ' + message);
    }
}

// 页面加载完成后初始化应用
document.addEventListener('DOMContentLoaded', () => {
    new ProbeApp();
});
