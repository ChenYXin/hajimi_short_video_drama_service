// 管理后台 JavaScript

$(document).ready(function() {
    // 侧边栏切换
    $('#sidebarCollapse').on('click', function() {
        $('#sidebar').toggleClass('active');
        $('#content').toggleClass('active');
    });
    
    // 自动隐藏警告消息
    $('.alert').each(function() {
        const alert = $(this);
        if (alert.hasClass('alert-success') || alert.hasClass('alert-info')) {
            setTimeout(function() {
                alert.fadeOut();
            }, 5000);
        }
    });
    
    // 初始化工具提示
    $('[data-bs-toggle="tooltip"]').tooltip();
    
    // 初始化弹出框
    $('[data-bs-toggle="popover"]').popover();
    
    // 表格行点击效果
    $('.table-hover tbody tr').on('click', function(e) {
        if (!$(e.target).closest('.btn, .btn-group, a').length) {
            $(this).toggleClass('table-active');
        }
    });
    
    // 全选/取消全选
    $('#selectAll').on('change', function() {
        const isChecked = $(this).is(':checked');
        $('.row-checkbox').prop('checked', isChecked);
        updateBatchToolbar();
    });
    
    // 行选择
    $(document).on('change', '.row-checkbox', function() {
        updateBatchToolbar();
        
        const totalRows = $('.row-checkbox').length;
        const checkedRows = $('.row-checkbox:checked').length;
        
        $('#selectAll').prop('indeterminate', checkedRows > 0 && checkedRows < totalRows);
        $('#selectAll').prop('checked', checkedRows === totalRows);
    });
    
    // 搜索表单自动提交
    let searchTimeout;
    $('input[name="search"]').on('input', function() {
        clearTimeout(searchTimeout);
        const form = $(this).closest('form');
        searchTimeout = setTimeout(function() {
            form.submit();
        }, 500);
    });
    
    // 文件拖拽上传
    $('.upload-area').on('dragover', function(e) {
        e.preventDefault();
        $(this).addClass('dragover');
    });
    
    $('.upload-area').on('dragleave', function(e) {
        e.preventDefault();
        $(this).removeClass('dragover');
    });
    
    $('.upload-area').on('drop', function(e) {
        e.preventDefault();
        $(this).removeClass('dragover');
        
        const files = e.originalEvent.dataTransfer.files;
        if (files.length > 0) {
            handleFileUpload(files[0], $(this));
        }
    });
    
    // 确认删除对话框
    $('.btn-delete').on('click', function(e) {
        e.preventDefault();
        const url = $(this).attr('href') || $(this).data('url');
        const title = $(this).data('title') || '此项目';
        
        if (confirm('确定要删除 "' + title + '" 吗？此操作不可恢复。')) {
            if ($(this).attr('href')) {
                window.location.href = url;
            } else {
                // AJAX 删除
                performDelete(url);
            }
        }
    });
    
    // 状态切换
    $('.status-toggle').on('change', function() {
        const id = $(this).data('id');
        const status = $(this).is(':checked') ? 'active' : 'inactive';
        const url = $(this).data('url');
        
        updateStatus(id, status, url);
    });
    
    // 表单验证
    $('form[data-validate="true"]').on('submit', function(e) {
        if (!validateForm($(this))) {
            e.preventDefault();
        }
    });
    
    // 图片预览
    $('input[type="file"][accept*="image"]').on('change', function() {
        const file = this.files[0];
        const preview = $(this).data('preview');
        
        if (file && preview) {
            const reader = new FileReader();
            reader.onload = function(e) {
                $(preview).attr('src', e.target.result).show();
            };
            reader.readAsDataURL(file);
        }
    });
    
    // 数据表格排序
    $('.sortable th').on('click', function() {
        const column = $(this).data('column');
        const currentSort = $(this).data('sort') || 'asc';
        const newSort = currentSort === 'asc' ? 'desc' : 'asc';
        
        // 更新 URL 参数
        const url = new URL(window.location);
        url.searchParams.set('sort', column);
        url.searchParams.set('order', newSort);
        window.location.href = url.toString();
    });
    
    // 自动保存草稿
    let autoSaveTimeout;
    $('form[data-autosave="true"] input, form[data-autosave="true"] textarea').on('input', function() {
        clearTimeout(autoSaveTimeout);
        const form = $(this).closest('form');
        
        autoSaveTimeout = setTimeout(function() {
            saveDraft(form);
        }, 2000);
    });
});

// 更新批量操作工具栏
function updateBatchToolbar() {
    const checkedCount = $('.row-checkbox:checked').length;
    
    if (checkedCount > 0) {
        $('#batchToolbar').show();
        $('#selectedCount').text(checkedCount);
    } else {
        $('#batchToolbar').hide();
    }
}

// 文件上传处理
function handleFileUpload(file, uploadArea) {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('type', uploadArea.data('type') || 'others');
    
    // 显示上传进度
    uploadArea.html(`
        <div class="upload-progress">
            <i class="fas fa-spinner fa-spin fa-2x text-primary"></i>
            <p class="text-primary mt-2">上传中...</p>
            <div class="progress mt-2">
                <div class="progress-bar" role="progressbar" style="width: 0%"></div>
            </div>
        </div>
    `);
    
    $.ajax({
        url: '/api/upload',
        type: 'POST',
        data: formData,
        processData: false,
        contentType: false,
        headers: {
            'Authorization': 'Bearer ' + getAuthToken()
        },
        xhr: function() {
            const xhr = new window.XMLHttpRequest();
            xhr.upload.addEventListener('progress', function(e) {
                if (e.lengthComputable) {
                    const percentComplete = (e.loaded / e.total) * 100;
                    uploadArea.find('.progress-bar').css('width', percentComplete + '%');
                }
            });
            return xhr;
        },
        success: function(response) {
            if (response.success) {
                const fileUrl = response.data.url;
                const fileName = response.data.filename;
                
                // 更新上传区域显示
                if (file.type.startsWith('image/')) {
                    uploadArea.html(`
                        <img src="${fileUrl}" alt="预览" class="img-fluid rounded">
                        <p class="text-success mt-2"><i class="fas fa-check-circle"></i> 上传成功</p>
                    `);
                } else {
                    uploadArea.html(`
                        <i class="fas fa-file fa-2x text-success"></i>
                        <p class="text-success mt-2"><i class="fas fa-check-circle"></i> ${fileName}</p>
                    `);
                }
                
                // 设置隐藏字段值
                const hiddenInput = uploadArea.siblings('input[type="hidden"]');
                if (hiddenInput.length) {
                    hiddenInput.val(fileUrl);
                }
                
                // 触发自定义事件
                uploadArea.trigger('fileUploaded', [fileUrl, fileName]);
            } else {
                showUploadError(uploadArea, response.message);
            }
        },
        error: function() {
            showUploadError(uploadArea, '上传失败，请重试');
        }
    });
}

// 显示上传错误
function showUploadError(uploadArea, message) {
    uploadArea.html(`
        <i class="fas fa-exclamation-triangle fa-2x text-danger"></i>
        <p class="text-danger mt-2">${message}</p>
        <button type="button" class="btn btn-outline-primary btn-sm mt-2" onclick="$(this).parent().trigger('click')">
            重新上传
        </button>
    `);
}

// 执行删除操作
function performDelete(url) {
    $.ajax({
        url: url,
        type: 'DELETE',
        headers: {
            'Authorization': 'Bearer ' + getAuthToken()
        },
        success: function(response) {
            if (response.success) {
                showNotification('删除成功', 'success');
                setTimeout(function() {
                    location.reload();
                }, 1000);
            } else {
                showNotification('删除失败：' + response.message, 'error');
            }
        },
        error: function() {
            showNotification('删除失败，请重试', 'error');
        }
    });
}

// 更新状态
function updateStatus(id, status, url) {
    $.ajax({
        url: url,
        type: 'POST',
        data: JSON.stringify({ status: status }),
        contentType: 'application/json',
        headers: {
            'Authorization': 'Bearer ' + getAuthToken()
        },
        success: function(response) {
            if (response.success) {
                showNotification('状态更新成功', 'success');
            } else {
                showNotification('状态更新失败：' + response.message, 'error');
                // 恢复开关状态
                $('.status-toggle[data-id="' + id + '"]').prop('checked', status === 'inactive');
            }
        },
        error: function() {
            showNotification('状态更新失败，请重试', 'error');
            // 恢复开关状态
            $('.status-toggle[data-id="' + id + '"]').prop('checked', status === 'inactive');
        }
    });
}

// 表单验证
function validateForm(form) {
    let isValid = true;
    
    // 清除之前的错误提示
    form.find('.is-invalid').removeClass('is-invalid');
    form.find('.invalid-feedback').remove();
    
    // 验证必填字段
    form.find('[required]').each(function() {
        const field = $(this);
        const value = field.val().trim();
        
        if (!value) {
            showFieldError(field, '此字段为必填项');
            isValid = false;
        }
    });
    
    // 验证邮箱格式
    form.find('input[type="email"]').each(function() {
        const field = $(this);
        const value = field.val().trim();
        
        if (value && !isValidEmail(value)) {
            showFieldError(field, '请输入有效的邮箱地址');
            isValid = false;
        }
    });
    
    // 验证密码确认
    const password = form.find('input[name="password"]').val();
    const confirmPassword = form.find('input[name="password_confirmation"]').val();
    
    if (password && confirmPassword && password !== confirmPassword) {
        showFieldError(form.find('input[name="password_confirmation"]'), '密码确认不匹配');
        isValid = false;
    }
    
    return isValid;
}

// 显示字段错误
function showFieldError(field, message) {
    field.addClass('is-invalid');
    field.after('<div class="invalid-feedback">' + message + '</div>');
}

// 验证邮箱格式
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// 保存草稿
function saveDraft(form) {
    const formData = form.serialize();
    const draftKey = 'draft_' + form.attr('id') || 'form';
    
    localStorage.setItem(draftKey, formData);
    
    // 显示保存提示
    showNotification('草稿已自动保存', 'info', 2000);
}

// 加载草稿
function loadDraft(formId) {
    const draftKey = 'draft_' + formId;
    const draftData = localStorage.getItem(draftKey);
    
    if (draftData) {
        const params = new URLSearchParams(draftData);
        const form = $('#' + formId);
        
        params.forEach((value, key) => {
            const field = form.find('[name="' + key + '"]');
            if (field.length) {
                field.val(value);
            }
        });
        
        showNotification('已加载草稿内容', 'info', 3000);
    }
}

// 清除草稿
function clearDraft(formId) {
    const draftKey = 'draft_' + formId;
    localStorage.removeItem(draftKey);
}

// 显示通知消息
function showNotification(message, type = 'info', duration = 5000) {
    const alertClass = 'alert-' + (type === 'error' ? 'danger' : type);
    const icon = {
        success: 'fas fa-check-circle',
        error: 'fas fa-exclamation-circle',
        warning: 'fas fa-exclamation-triangle',
        info: 'fas fa-info-circle'
    }[type] || 'fas fa-info-circle';
    
    const notification = $(`
        <div class="alert ${alertClass} alert-dismissible fade show position-fixed" 
             style="top: 20px; right: 20px; z-index: 9999; min-width: 300px;">
            <i class="${icon}"></i> ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    `);
    
    $('body').append(notification);
    
    if (duration > 0) {
        setTimeout(function() {
            notification.fadeOut(function() {
                $(this).remove();
            });
        }, duration);
    }
}

// 获取认证令牌
function getAuthToken() {
    return localStorage.getItem('admin_token') || sessionStorage.getItem('admin_token') || '';
}

// 设置认证令牌
function setAuthToken(token, remember = false) {
    if (remember) {
        localStorage.setItem('admin_token', token);
    } else {
        sessionStorage.setItem('admin_token', token);
    }
}

// 清除认证令牌
function clearAuthToken() {
    localStorage.removeItem('admin_token');
    sessionStorage.removeItem('admin_token');
}

// 检查认证状态
function checkAuthStatus() {
    const token = getAuthToken();
    if (!token) {
        window.location.href = '/admin/login';
        return false;
    }
    return true;
}

// AJAX 请求拦截器
$(document).ajaxSend(function(event, xhr, settings) {
    // 自动添加认证头
    if (!settings.headers || !settings.headers.Authorization) {
        const token = getAuthToken();
        if (token) {
            xhr.setRequestHeader('Authorization', 'Bearer ' + token);
        }
    }
});

// AJAX 错误处理
$(document).ajaxError(function(event, xhr, settings) {
    if (xhr.status === 401) {
        clearAuthToken();
        showNotification('登录已过期，请重新登录', 'warning');
        setTimeout(function() {
            window.location.href = '/admin/login';
        }, 2000);
    } else if (xhr.status === 403) {
        showNotification('权限不足', 'error');
    } else if (xhr.status >= 500) {
        showNotification('服务器错误，请稍后重试', 'error');
    }
});

// 批量操作函数
function batchActivate() {
    const selectedIds = getSelectedIds();
    if (selectedIds.length === 0) {
        showNotification('请先选择要操作的项目', 'warning');
        return;
    }
    
    if (confirm('确定要激活选中的 ' + selectedIds.length + ' 个项目吗？')) {
        performBatchOperation('activate', selectedIds);
    }
}

function batchDeactivate() {
    const selectedIds = getSelectedIds();
    if (selectedIds.length === 0) {
        showNotification('请先选择要操作的项目', 'warning');
        return;
    }
    
    if (confirm('确定要禁用选中的 ' + selectedIds.length + ' 个项目吗？')) {
        performBatchOperation('deactivate', selectedIds);
    }
}

function batchDelete() {
    const selectedIds = getSelectedIds();
    if (selectedIds.length === 0) {
        showNotification('请先选择要操作的项目', 'warning');
        return;
    }
    
    if (confirm('确定要删除选中的 ' + selectedIds.length + ' 个项目吗？此操作不可恢复。')) {
        performBatchOperation('delete', selectedIds);
    }
}

// 获取选中的 ID
function getSelectedIds() {
    const ids = [];
    $('.row-checkbox:checked').each(function() {
        ids.push($(this).val());
    });
    return ids;
}

// 执行批量操作
function performBatchOperation(operation, ids) {
    const url = '/api/admin/batch/' + operation;
    
    $.ajax({
        url: url,
        type: 'POST',
        data: JSON.stringify({ ids: ids }),
        contentType: 'application/json',
        headers: {
            'Authorization': 'Bearer ' + getAuthToken()
        },
        success: function(response) {
            if (response.success) {
                showNotification('批量操作成功', 'success');
                setTimeout(function() {
                    location.reload();
                }, 1000);
            } else {
                showNotification('批量操作失败：' + response.message, 'error');
            }
        },
        error: function() {
            showNotification('批量操作失败，请重试', 'error');
        }
    });
}

// 数据导出
function exportData(format = 'excel') {
    const url = new URL(window.location);
    url.pathname = url.pathname + '/export';
    url.searchParams.set('format', format);
    
    // 创建隐藏的下载链接
    const link = document.createElement('a');
    link.href = url.toString();
    link.download = '';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}

// 打印页面
function printPage() {
    window.print();
}

// 全屏切换
function toggleFullscreen() {
    if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen();
    } else {
        document.exitFullscreen();
    }
}

// 工具函数
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN');
}

function truncateText(text, maxLength) {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

function throttle(func, limit) {
    let inThrottle;
    return function() {
        const args = arguments;
        const context = this;
        if (!inThrottle) {
            func.apply(context, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}