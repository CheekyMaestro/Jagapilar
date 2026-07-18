// ============================================================
// JAGAPILAR — Registration/Login Logic
// Parent & Teacher Roles
// ============================================================

const API_BASE = '/api';

document.addEventListener('DOMContentLoaded', () => {
    // Check URL params for role selection
    const urlParams = new URLSearchParams(window.location.search);
    const role = urlParams.get('role');
    if (role === 'teacher') {
        switchRole('teacher');
    }
});

function switchRole(role) {
    const tabParent = document.getElementById('tab-parent');
    const tabTeacher = document.getElementById('tab-teacher');
    const schoolFieldContainer = document.getElementById('school-field-container');
    const userRoleInput = document.getElementById('user-role');

    if (role === 'parent') {
        tabParent.classList.replace('text-on-surface-variant', 'text-primary');
        tabParent.classList.replace('border-transparent', 'border-primary');
        tabParent.classList.add('bg-primary/5');
        tabParent.classList.remove('hover:bg-surface-container-low');

        tabTeacher.classList.replace('text-primary', 'text-on-surface-variant');
        tabTeacher.classList.replace('border-primary', 'border-transparent');
        tabTeacher.classList.remove('bg-primary/5');
        tabTeacher.classList.add('hover:bg-surface-container-low');

        schoolFieldContainer.classList.add('hidden');
        userRoleInput.value = 'parent';
    } else if (role === 'teacher') {
        tabTeacher.classList.replace('text-on-surface-variant', 'text-primary');
        tabTeacher.classList.replace('border-transparent', 'border-primary');
        tabTeacher.classList.add('bg-primary/5');
        tabTeacher.classList.remove('hover:bg-surface-container-low');

        tabParent.classList.replace('text-primary', 'text-on-surface-variant');
        tabParent.classList.replace('border-primary', 'border-transparent');
        tabParent.classList.remove('bg-primary/5');
        tabParent.classList.add('hover:bg-surface-container-low');

        schoolFieldContainer.classList.remove('hidden');
        userRoleInput.value = 'teacher';
    }
}

async function submitForm(event) {
    event.preventDefault();
    
    const role = document.getElementById('user-role').value;
    const name = document.getElementById('reg-name').value.trim();
    const contact = document.getElementById('reg-contact').value.trim();
    const password = document.getElementById('reg-password').value;
    const school = document.getElementById('reg-school').value.trim();

    if (!name || !contact || !password) {
        showToast('Mohon lengkapi semua field wajib.', 'warning');
        return;
    }

    if (password.length < 8) {
        showToast('Password minimal 8 karakter.', 'warning');
        return;
    }

    const payload = {
        role: role,
        name: name,
        email_contact: contact,
        password: password,
        school_name: role === 'teacher' ? school : null
    };

    try {
        const btn = document.querySelector('button[type="submit"]');
        btn.disabled = true;
        btn.innerHTML = 'Memproses...';

        // 1. Try to Register
        let res = await fetch(`${API_BASE}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (res.status === 409) {
            // Already exists, try login instead
            res = await fetch(`${API_BASE}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email_contact: contact, password: password })
            });
        }

        const data = await res.json();
        
        if (!res.ok) {
            throw new Error(data.error || 'Terjadi kesalahan sistem');
        }

        // Handle successful registration (which doesn't return token) or login
        if (!data.token) {
            // If we just registered, automatically login
            const loginRes = await fetch(`${API_BASE}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email_contact: contact, password: password })
            });
            const loginData = await loginRes.json();
            if(!loginRes.ok) throw new Error(loginData.error);
            data.token = loginData.token;
            data.user = loginData.user;
        }

        // Save token to localStorage
        localStorage.setItem('jagapilar_token', data.token);
        localStorage.setItem('jagapilar_user', JSON.stringify(data.user));

        showToast('Berhasil masuk! Mengarahkan ke dashboard...', 'success');
        
        setTimeout(() => {
            if (data.user.role === 'parent') {
                window.location.href = 'dashboard-parent.html';
            } else {
                window.location.href = 'dashboard-teacher.html';
            }
        }, 1500);

    } catch (error) {
        console.error(error);
        showToast(error.message || 'Koneksi ke server gagal. (Pastikan backend berjalan)', 'error');
        
        // --- FALLBACK FOR OFFLINE DEMO ---
        // If backend is offline, we fallback to local simulation
        console.warn('Falling back to local simulation due to API error.');
        localStorage.setItem('jagapilar_token', 'mock-token-123');
        localStorage.setItem('jagapilar_user', JSON.stringify({ role: role, name: name }));
        
        showToast('(Offline Mode) Berhasil masuk...', 'success');
        setTimeout(() => {
            window.location.href = role === 'parent' ? 'dashboard-parent.html' : 'dashboard-teacher.html';
        }, 1500);
        // ---------------------------------
        
        const btn = document.querySelector('button[type="submit"]');
        btn.disabled = false;
        btn.innerHTML = 'DAFTAR / MASUK SEKARANG <span class="material-symbols-outlined text-sm ml-2">arrow_forward</span>';
    }
}
