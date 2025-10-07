class ProbeApp {
    constructor() {
        this.packets = [];
        this.filteredPackets = [];
        this.currentFilter = {};
        this.autoRefresh = true;
        this.refreshInterval = null;
        this.isFilterActive = false;
        
        this.init();
    }
    
    init() {
        this.bindEvents();
        this.setDefaultFilter();
        this.loadStats();
        this.loadPackets();
        this.startAutoRefresh();
    }
    
    setDefaultFilter() {
        // 设置默认过滤条件
        document.getElementById('dst-ip').value = '139.155.132.244';
        this.currentFilter = {
            dst_ip: '139.155.132.244'
        };
        this.isFilterActive = true;
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
            
            // 调试信息：打印前几个数据包的信息
            if (this.packets.length > 0) {
                console.log('调试: 加载的数据包信息:');
                this.packets.slice(0, 3).forEach((packet, index) => {
                    console.log(`数据包 ${index + 1}:`, {
                        src_ip: packet.src_ip,
                        dst_ip: packet.dst_ip,
                        domain: packet.domain,
                        host: packet.host,
                        http_method: packet.http_method,
                        http_url: packet.http_url,
                        path: packet.path
                    });
                });
            }
            
            // 如果当前有过滤条件，重新应用过滤
            if (Object.keys(this.currentFilter).length > 0) {
                await this.applyCurrentFilter();
            } else {
                this.filteredPackets = [...this.packets];
                this.renderPackets();
            }
        } catch (error) {
            console.error('加载数据包失败:', error);
            this.showError('加载数据包失败: ' + error.message);
        }
    }
    
    async applyFilter() {
        const filter = this.getCurrentFilterFromUI();
        this.currentFilter = filter;
        this.isFilterActive = Object.keys(filter).length > 0;
        await this.applyCurrentFilter();
        this.updateFilterStatus();
    }
    
    getCurrentFilterFromUI() {
        const filter = {
            src_ip: document.getElementById('src-ip').value.trim(),
            dst_ip: document.getElementById('dst-ip').value.trim(),
            port: document.getElementById('port').value ? parseInt(document.getElementById('port').value) : 0,
            protocol: document.getElementById('protocol').value,
            http_method: document.getElementById('http-method').value,
            host: document.getElementById('host').value.trim(),
            domain: document.getElementById('domain').value.trim(),
            path: document.getElementById('path').value.trim(),
            user_agent: document.getElementById('user-agent').value.trim(),
            content_type: document.getElementById('content-type').value.trim(),
            referer: document.getElementById('referer').value.trim(),
            server: document.getElementById('server').value.trim(),
            search_text: document.getElementById('search-text').value.trim()
        };
        
        // 移除空值
        Object.keys(filter).forEach(key => {
            if (filter[key] === '' || filter[key] === 0) {
                delete filter[key];
            }
        });
        
        return filter;
    }
    
    async applyCurrentFilter() {
        if (Object.keys(this.currentFilter).length === 0) {
            this.filteredPackets = [...this.packets];
            this.renderPackets();
            return;
        }
        
        try {
            const response = await fetch('/api/packets/filter', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(this.currentFilter)
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
        document.getElementById('host').value = '';
        document.getElementById('domain').value = '';
        document.getElementById('path').value = '';
        document.getElementById('user-agent').value = '';
        document.getElementById('content-type').value = '';
        document.getElementById('referer').value = '';
        document.getElementById('server').value = '';
        document.getElementById('search-text').value = '';
        
        this.currentFilter = {};
        this.isFilterActive = false;
        this.filteredPackets = [...this.packets];
        this.renderPackets();
        this.updateFilterStatus();
    }
    
    updateFilterStatus() {
        const statusElement = document.getElementById('filter-status');
        if (!statusElement) {
            // 创建状态指示器
            this.createFilterStatusIndicator();
            return;
        }
        
        if (this.isFilterActive) {
            const filterCount = Object.keys(this.currentFilter).length;
            statusElement.innerHTML = `🔍 过滤已激活 (${filterCount} 个条件) | 显示 ${this.filteredPackets.length} 个数据包`;
            statusElement.className = 'filter-status active';
        } else {
            statusElement.innerHTML = `📊 显示所有数据包 (${this.filteredPackets.length} 个)`;
            statusElement.className = 'filter-status inactive';
        }
    }
    
    createFilterStatusIndicator() {
        const filtersDiv = document.querySelector('.filters');
        if (filtersDiv) {
            const statusDiv = document.createElement('div');
            statusDiv.id = 'filter-status';
            statusDiv.className = 'filter-status';
            statusDiv.innerHTML = '📊 显示所有数据包';
            filtersDiv.appendChild(statusDiv);
        }
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
            tbody.innerHTML = '<tr><td colspan="11" class="empty">暂无数据包</td></tr>';
            return;
        }
        
        tbody.innerHTML = this.filteredPackets.map((packet, index) => {
            const time = new Date(packet.timestamp).toLocaleString('zh-CN');
            const srcAddr = `${packet.src_ip}:${packet.src_port}`;
            const dstAddr = `${packet.dst_ip}:${packet.dst_port}`;
            const protocolClass = `protocol-${packet.protocol.toLowerCase()}`;
            const methodClass = packet.http_method ? `http-${packet.http_method.toLowerCase()}` : '';
            const statusClass = packet.http_status ? this.getStatusClass(packet.http_status) : '';
            
            return `
                <tr class="packet-row" data-index="${index}">
                    <td class="time-cell">${time}</td>
                    <td class="ip-cell">${srcAddr}</td>
                    <td class="ip-cell">${dstAddr}</td>
                    <td class="protocol-cell">
                        <span class="protocol-tcp ${protocolClass}">${packet.protocol}</span>
                    </td>
                    <td class="length-cell">${packet.length}</td>
                    <td class="domain-cell" title="${packet.domain || ''}">${packet.domain || '-'}</td>
                    <td class="path-cell" title="${packet.path || ''}">${packet.path || '-'}</td>
                    <td class="http-method ${methodClass}">${packet.http_method || '-'}</td>
                    <td class="status-cell ${statusClass}">${packet.http_status || '-'}</td>
                    <td class="user-agent-cell" title="${packet.user_agent || ''}">${packet.user_agent || '-'}</td>
                    <td class="action-cell">
                        <button class="btn-detail" onclick="window.app.showPacketDetail(${index})">详情</button>
                    </td>
                </tr>
            `;
        }).join('');
        
        // 渲染完成后更新过滤状态
        this.updateFilterStatus();
        
        // 绑定详情按钮事件
        this.bindDetailButtons();
    }
    
    bindDetailButtons() {
        // 移除之前的事件监听器
        document.removeEventListener('click', this.handleDetailClick);
        
        // 绑定新的事件监听器
        this.handleDetailClick = (e) => {
            console.log('点击事件:', e.target);
            if (e.target.classList.contains('btn-detail')) {
                console.log('点击了详情按钮');
                const index = parseInt(e.target.getAttribute('data-index'));
                console.log('按钮索引:', index);
                this.showPacketDetail(index);
            }
        };
        
        document.addEventListener('click', this.handleDetailClick);
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
    
    showPacketDetail(index) {
        console.log('点击详情按钮，索引:', index);
        console.log('当前过滤数据包数量:', this.filteredPackets.length);
        
        const packet = this.filteredPackets[index];
        if (!packet) {
            console.error('数据包不存在，索引:', index);
            this.showError('数据包不存在');
            return;
        }
        
        console.log('显示数据包详情:', packet);
        
        // 创建详情模态框
        this.createDetailModal(packet);
    }
    
    createDetailModal(packet) {
        // 移除已存在的模态框
        const existingModal = document.getElementById('packet-detail-modal');
        if (existingModal) {
            existingModal.remove();
        }
        
        // 创建模态框
        const modal = document.createElement('div');
        modal.id = 'packet-detail-modal';
        modal.className = 'modal';
        modal.style.display = 'block';
        
        const time = new Date(packet.timestamp).toLocaleString('zh-CN');
        
        modal.innerHTML = `
            <div class="modal-content">
                <div class="modal-header">
                    <h3>数据包详情</h3>
                    <span class="close" onclick="this.closest('.modal').remove()">&times;</span>
                </div>
                <div class="modal-body">
                    <div class="detail-section">
                        <h4>基本信息</h4>
                        <table class="detail-table">
                            <tr><td>时间戳</td><td>${time}</td></tr>
                            <tr><td>源IP</td><td>${packet.src_ip}</td></tr>
                            <tr><td>目标IP</td><td>${packet.dst_ip}</td></tr>
                            <tr><td>源端口</td><td>${packet.src_port}</td></tr>
                            <tr><td>目标端口</td><td>${packet.dst_port}</td></tr>
                            <tr><td>协议</td><td>${packet.protocol}</td></tr>
                            <tr><td>长度</td><td>${packet.length} 字节</td></tr>
                        </table>
                    </div>
                    
                    <div class="detail-section">
                        <h4>HTTP信息</h4>
                        <table class="detail-table">
                            <tr><td>HTTP方法</td><td>${packet.http_method || '-'}</td></tr>
                            <tr><td>HTTP状态</td><td>${packet.http_status || '-'}</td></tr>
                            <tr><td>完整URL</td><td>${packet.http_url || '-'}</td></tr>
                            <tr><td>域名</td><td>${packet.domain || '-'}</td></tr>
                            <tr><td>主机名</td><td>${packet.host || '-'}</td></tr>
                            <tr><td>路径</td><td>${packet.path || '-'}</td></tr>
                            <tr><td>查询参数</td><td>${packet.query || '-'}</td></tr>
                        </table>
                    </div>
                    
                    <div class="detail-section">
                        <h4>HTTP头部</h4>
                        <table class="detail-table">
                            <tr><td>User-Agent</td><td>${packet.user_agent || '-'}</td></tr>
                            <tr><td>Content-Type</td><td>${packet.content_type || '-'}</td></tr>
                            <tr><td>Referer</td><td>${packet.referer || '-'}</td></tr>
                            <tr><td>Server</td><td>${packet.server || '-'}</td></tr>
                            <tr><td>Accept</td><td>${packet.accept || '-'}</td></tr>
                            <tr><td>Accept-Language</td><td>${packet.accept_language || '-'}</td></tr>
                            <tr><td>Accept-Encoding</td><td>${packet.accept_encoding || '-'}</td></tr>
                            <tr><td>Connection</td><td>${packet.connection || '-'}</td></tr>
                            <tr><td>Cache-Control</td><td>${packet.cache_control || '-'}</td></tr>
                            <tr><td>Authorization</td><td>${packet.authorization || '-'}</td></tr>
                            <tr><td>Cookie</td><td>${packet.cookie || '-'}</td></tr>
                            <tr><td>Set-Cookie</td><td>${packet.set_cookie || '-'}</td></tr>
                        </table>
                    </div>
                    
                    <div class="detail-section">
                        <h4>响应信息</h4>
                        <table class="detail-table">
                            <tr><td>Location</td><td>${packet.location || '-'}</td></tr>
                            <tr><td>Last-Modified</td><td>${packet.last_modified || '-'}</td></tr>
                            <tr><td>ETag</td><td>${packet.etag || '-'}</td></tr>
                            <tr><td>Expires</td><td>${packet.expires || '-'}</td></tr>
                            <tr><td>Date</td><td>${packet.date || '-'}</td></tr>
                            <tr><td>Content-Length</td><td>${packet.content_length || '-'}</td></tr>
                            <tr><td>Transfer-Encoding</td><td>${packet.transfer_encoding || '-'}</td></tr>
                        </table>
                    </div>
                    
                    <div class="detail-section">
                        <h4>代理信息</h4>
                        <table class="detail-table">
                            <tr><td>X-Forwarded-For</td><td>${packet.x_forwarded_for || '-'}</td></tr>
                            <tr><td>X-Real-IP</td><td>${packet.x_real_ip || '-'}</td></tr>
                            <tr><td>X-Requested-With</td><td>${packet.x_requested_with || '-'}</td></tr>
                            <tr><td>Via</td><td>${packet.via || '-'}</td></tr>
                        </table>
                    </div>
                    
                    ${packet.payload && packet.payload.length > 0 ? `
                    <div class="detail-section">
                        <h4>原始数据</h4>
                        <pre class="payload-content">${this.escapeHtml(packet.payload)}</pre>
                    </div>
                    ` : ''}
                </div>
                <div class="modal-footer">
                    <button class="btn btn-secondary" onclick="this.closest('.modal').remove()">关闭</button>
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
        
        // 绑定关闭事件
        const closeBtn = modal.querySelector('.close');
        if (closeBtn) {
            closeBtn.addEventListener('click', () => {
                modal.remove();
            });
        }
        
        // 点击背景关闭模态框
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });
        
        // ESC键关闭模态框
        const handleEsc = (e) => {
            if (e.key === 'Escape') {
                modal.remove();
                document.removeEventListener('keydown', handleEsc);
            }
        };
        document.addEventListener('keydown', handleEsc);
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// 页面加载完成后初始化应用
let app;
document.addEventListener('DOMContentLoaded', () => {
    app = new ProbeApp();
    window.app = app; // 设置为全局变量
});
