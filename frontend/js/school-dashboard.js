document.addEventListener('DOMContentLoaded', () => {
    // Check for school ID in URL or localStorage
    const urlParams = new URLSearchParams(window.location.search);
    let schoolId = urlParams.get('id');
    
    if (!schoolId) {
        schoolId = localStorage.getItem('jagapilar_school_id');
        if (!schoolId) {
            window.location.href = 'register.html';
            return;
        }
        // Update URL to make it shareable
        window.history.replaceState({}, '', `?id=${schoolId}`);
    } else {
        localStorage.setItem('jagapilar_school_id', schoolId);
    }

    loadDashboard(schoolId);

    // Logout
    document.getElementById('logoutBtn').addEventListener('click', () => {
        localStorage.removeItem('jagapilar_school_id');
        window.location.href = 'index.html';
    });

    // Add Student Form
    document.getElementById('addStudentForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const btn = document.getElementById('submitStudentBtn');
        const originalText = btn.innerHTML;
        btn.innerHTML = '<div class="spinner w-5 h-5 border-2"></div> Mendaftarkan...';
        btn.disabled = true;
        
        try {
            const res = await apiCall('/api/children/register-full', {
                method: 'POST',
                body: JSON.stringify({
                    school_id: schoolId,
                    grade: parseInt(document.getElementById('studentGrade').value),
                    birth_year: parseInt(document.getElementById('studentBirthYear').value)
                })
            });
            
            showToast('Murid berhasil didaftarkan! Link otomatis dibuat.', 'success');
            document.getElementById('addStudentModal').classList.add('hidden');
            document.getElementById('addStudentForm').reset();
            
            // Reload dashboard
            loadDashboard(schoolId);
        } catch (error) {
            showToast(error.message || 'Gagal mendaftarkan murid', 'error');
        } finally {
            btn.innerHTML = originalText;
            btn.disabled = false;
        }
    });

    // Search
    document.getElementById('searchInput').addEventListener('input', (e) => {
        const term = e.target.value.toLowerCase();
        const rows = document.querySelectorAll('#studentsTableBody tr');
        rows.forEach(row => {
            if (row.classList.contains('empty-row')) return;
            const code = row.querySelector('.anon-code').textContent.toLowerCase();
            row.style.display = code.includes(term) ? '' : 'none';
        });
    });
});

async function loadDashboard(schoolId) {
    try {
        const data = await apiCall(`/api/schools/${schoolId}/dashboard`);
        
        // Update UI Info
        document.getElementById('schoolNameDisplay').textContent = data.school.name;
        document.getElementById('schoolInfoDisplay').textContent = `${data.school.city} • Kepala Sekolah: ${data.school.principal_name}`;
        
        // Stats
        const totalMurid = data.children ? data.children.length : 0;
        document.getElementById('statTotal').textContent = totalMurid;
        
        animateValue('statHijau', 0, data.total_hijau || 0, 1000);
        animateValue('statKuning', 0, data.total_kuning || 0, 1000);
        animateValue('statMerah', 0, data.total_merah || 0, 1000);

        if (data.total_review > 0) {
            document.getElementById('reviewAlert').classList.remove('hidden');
        } else {
            document.getElementById('reviewAlert').classList.add('hidden');
        }

        renderTable(data.children || []);
        
    } catch (error) {
        showToast('Gagal memuat dashboard', 'error');
        document.getElementById('schoolNameDisplay').textContent = 'Error Memuat Data';
    }
}

function renderTable(children) {
    const tbody = document.getElementById('studentsTableBody');
    tbody.innerHTML = '';
    
    if (children.length === 0) {
        tbody.innerHTML = `<tr class="empty-row"><td colspan="6" class="px-6 py-10 text-center text-slate-500">Belum ada murid terdaftar. Klik "Daftarkan Murid" untuk memulai.</td></tr>`;
        return;
    }

    const baseUrl = window.location.origin + window.location.pathname.replace('school-dashboard.html', 'assessment.html');

    children.forEach(item => {
        const tr = document.createElement('tr');
        tr.className = 'hover:bg-slate-50 transition-colors border-b border-slate-100 last:border-0';
        
        let zoneBadge = `<span class="px-2.5 py-1 rounded-full text-xs font-medium bg-slate-100 text-slate-600">Belum Selesai</span>`;
        if (item.composite) {
            if (item.composite.needs_manual_review) {
                zoneBadge = `<span class="px-2.5 py-1 rounded-full text-xs font-medium bg-orange-100 text-orange-700 flex items-center gap-1 w-max"><i class="ph ph-warning"></i> Cek Manual</span>`;
            } else if (item.composite.composite_zone === 'hijau') {
                zoneBadge = `<span class="px-2.5 py-1 rounded-full text-xs font-medium bg-green-100 text-green-700">Zona Hijau</span>`;
            } else if (item.composite.composite_zone === 'kuning') {
                zoneBadge = `<span class="px-2.5 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-700">Zona Kuning</span>`;
            } else {
                zoneBadge = `<span class="px-2.5 py-1 rounded-full text-xs font-medium bg-red-100 text-red-700">Zona Merah</span>`;
            }
        }

        tr.innerHTML = `
            <td class="px-6 py-4 font-outfit font-semibold text-slate-800 anon-code">${item.child.anon_code}</td>
            <td class="px-6 py-4 text-slate-600">Kelas ${item.child.grade}</td>
            <td class="px-6 py-4">${zoneBadge}</td>
            <td class="px-6 py-4 text-center">${createCopyButton(baseUrl, item.parent_token, 'parent')}</td>
            <td class="px-6 py-4 text-center">${createCopyButton(baseUrl, item.teacher_token, 'teacher')}</td>
            <td class="px-6 py-4 text-center">${createCopyButton(baseUrl, item.student_token, 'student')}</td>
        `;
        tbody.appendChild(tr);
    });
}

function createCopyButton(baseUrl, token, role) {
    if (!token) return '<span class="text-slate-300">-</span>';
    
    let colorClass = '';
    let label = '';
    switch(role) {
        case 'parent': colorClass = 'text-blue-600 bg-blue-50 hover:bg-blue-100'; label = 'Ortu'; break;
        case 'teacher': colorClass = 'text-emerald-600 bg-emerald-50 hover:bg-emerald-100'; label = 'Guru'; break;
        case 'student': colorClass = 'text-purple-600 bg-purple-50 hover:bg-purple-100'; label = 'Anak'; break;
    }
    
    const url = `${baseUrl}?token=${token}`;
    
    return `
        <button onclick="copyToClipboard('${url}')" class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors flex items-center gap-1.5 mx-auto ${colorClass}" title="Salin Link Asesmen">
            <i class="ph ph-link"></i> Salin Link
        </button>
    `;
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showToast('Link berhasil disalin ke clipboard!', 'success');
    }).catch(err => {
        console.error('Failed to copy: ', err);
        showToast('Gagal menyalin link', 'error');
    });
}

function animateValue(id, start, end, duration) {
    if (start === end) {
        document.getElementById(id).innerHTML = end;
        return;
    }
    let range = end - start;
    let current = start;
    let increment = end > start ? 1 : -1;
    let stepTime = Math.abs(Math.floor(duration / range));
    let obj = document.getElementById(id);
    let timer = setInterval(function() {
        current += increment;
        obj.innerHTML = current;
        if (current == end) {
            clearInterval(timer);
        }
    }, stepTime);
}
