const API = '/kv';
let currentPage = 1;
const pageSize = 10;

async function loadKeys(page = 1) {
    const tbody = document.querySelector('#kvTable tbody');
    tbody.innerHTML = '';
    const res = await fetch(`${API}?page=${page}&size=${pageSize}`);
    if (!res.ok) return;
    const result = await res.json();
    const { data, page: current, total } = result;
    currentPage = current;
  
    for (let key of data) {
      const valRes = await fetch(`${API}/${key}`);
      if (!valRes.ok) continue;
      const valObj = await valRes.json();
      const value = valObj.Value ?? '';
      const ttl = valObj.TTL ?? 0;
  
      const row = document.createElement('tr');
      row.innerHTML = `
      <td>
        <span class="key-view">${key}</span>
        <input class="key-edit" value="${key}" style="display:none;" readonly>
      </td>
      <td>
        <input class="value-field" type="text" value="${value}" style="background:#f0f0f0;" readonly>
      </td>
      <td>
        <input class="ttl-field" type="number" value="${ttl}" min="0" style="width:60px; background:#f0f0f0;" readonly>
      </td>
      <td>
        <button class="btn-edit" onclick="enableEdit(this)">编辑</button>
        <button class="btn-save" style="display:none;" onclick="saveEdit(this, '${key}')">保存</button>
        <button class="btn-cancel" style="display:none;" onclick="cancelEdit(this)">取消</button>
        <button onclick="deleteKey('${key}')">删除</button>
      </td>
      `;   
      tbody.appendChild(row);
    }
  
    updatePagination(total);
}  

async function deleteKey(key) {
  await fetch(`${API}/${key}`, { method: 'DELETE' });
  loadKeys(currentPage);
}

document.getElementById('addForm').addEventListener('submit', async e => {
  e.preventDefault();
  const key = document.getElementById('key').value;
  const value = document.getElementById('value').value;
  const ttl = parseInt(document.getElementById('ttl').value);
  await fetch(`${API}/${key}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ value, ttl })
  });
  e.target.reset();
  loadKeys(currentPage);
});

// 修改后的 enableEdit 函数
function enableEdit(btn) {
  const row = btn.closest('tr');
  // 启用可编辑
  const valueInput = row.querySelector('.value-field');
  const ttlInput = row.querySelector('.ttl-field');

  valueInput.readOnly = false;
  ttlInput.readOnly = false;
  valueInput.style.background = '#ffffff';
  ttlInput.style.background = '#ffffff';

  // 显示“保存”和“取消”按钮，隐藏“编辑”
  row.querySelector('.btn-edit').style.display = 'none';
  row.querySelector('.btn-save').style.display = 'inline-block';
  row.querySelector('.btn-cancel').style.display = 'inline-block';
}

// 修改后的 cancelEdit 函数
function cancelEdit(btn) {
  const row = btn.closest('tr');

  const valueInput = row.querySelector('.value-field');
  const ttlInput = row.querySelector('.ttl-field');

  // 恢复为只读和灰色背景
  valueInput.readOnly = true;
  ttlInput.readOnly = true;
  valueInput.style.background = '#f0f0f0';
  ttlInput.style.background = '#f0f0f0';

  // 隐藏“保存”和“取消”，显示“编辑”
  row.querySelector('.btn-edit').style.display = 'inline-block';
  row.querySelector('.btn-save').style.display = 'none';
  row.querySelector('.btn-cancel').style.display = 'none';

  // 恢复原值（可选）
  loadKeys(currentPage);
}

async function saveEdit(btn, oldKey) {
  const row = btn.closest('tr');
  const keyInput = row.querySelector('.key-edit');
  const newKey = keyInput ? keyInput.value.trim() : oldKey;

  const valueInput = row.querySelector('.value-field');
  const ttlInput = row.querySelector('.ttl-field');

  const newValue = valueInput ? valueInput.value.trim() : '';
  const newTTL = ttlInput ? parseInt(ttlInput.value) : 0;

  if (!newKey) return alert('Key不能为空');

  const payload = {
    value: newValue,
    ttl: newTTL
  };

  if (newKey !== oldKey) {
    await fetch(`${API}/${oldKey}`, { method: 'DELETE' });
  }

  await fetch(`${API}/${newKey}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  });

  loadKeys(currentPage);
}


function updatePagination(total) {
  const pagination = document.getElementById('pagination');
  const totalPages = Math.ceil(total / pageSize);
  pagination.innerHTML = '';

  if (totalPages <= 1) return;

  if (currentPage > 1) {
    const prev = document.createElement('button');
    prev.textContent = '⬅️ 上一页';
    prev.onclick = () => loadKeys(currentPage - 1);
    pagination.appendChild(prev);
  }

  const info = document.createElement('span');
  info.textContent = ` 第 ${currentPage} 页 / 共 ${totalPages} 页 `;
  pagination.appendChild(info);

  if (currentPage < totalPages) {
    const next = document.createElement('button');
    next.textContent = '➡️ 下一页';
    next.onclick = () => loadKeys(currentPage + 1);
    pagination.appendChild(next);
  }
}

async function exportData() {
    const res = await fetch(`${API}?page=1&size=10000`); // 假设不会超过 10000 条
    if (!res.ok) return;
    const result = await res.json();
    const keys = result.data;
  
    const exportList = [];
  
    for (let key of keys) {
      const valRes = await fetch(`${API}/${key}`);
      if (!valRes.ok) continue;
      const valObj = await valRes.json();
      exportList.push({
        key,
        value: valObj.value ?? '',
        ttl: valObj.ttl ?? 0
      });
    }
  
    const blob = new Blob([JSON.stringify(exportList, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'kv_export.json';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
} 
  
async function importData(event) {
    const file = event.target.files[0];
    if (!file) return;
  
    const text = await file.text();
    let entries;
    try {
      entries = JSON.parse(text);
    } catch (err) {
      alert('导入的文件格式错误，请检查 JSON 格式');
      return;
    }
  
    if (!Array.isArray(entries)) {
      alert('导入数据应为 JSON 数组格式');
      return;
    }
  
    for (const item of entries) {
      if (!item.key || item.value === undefined) continue;
      await fetch(`${API}/${item.key}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ value: item.value, ttl: item.ttl || 0 })
      });
    }
  
    alert('导入完成');
    loadKeys(currentPage);
}  
  

loadKeys();
