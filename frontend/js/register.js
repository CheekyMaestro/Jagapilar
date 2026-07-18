// ============================================================
// JAGAPILAR — Registration/Login Logic
// Parent & Teacher Roles
// ============================================================

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
        // Style active tab
        tabParent.classList.replace('text-on-surface-variant', 'text-primary');
        tabParent.classList.replace('border-transparent', 'border-primary');
        tabParent.classList.add('bg-primary/5');
        tabParent.classList.remove('hover:bg-surface-container-low');

        // Style inactive tab
        tabTeacher.classList.replace('text-primary', 'text-on-surface-variant');
        tabTeacher.classList.replace('border-primary', 'border-transparent');
        tabTeacher.classList.remove('bg-primary/5');
        tabTeacher.classList.add('hover:bg-surface-container-low');

        schoolFieldContainer.classList.add('hidden');
        userRoleInput.value = 'parent';
    } else if (role === 'teacher') {
        // Style active tab
        tabTeacher.classList.replace('text-on-surface-variant', 'text-primary');
        tabTeacher.classList.replace('border-transparent', 'border-primary');
        tabTeacher.classList.add('bg-primary/5');
        tabTeacher.classList.remove('hover:bg-surface-container-low');

        // Style inactive tab
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
        role,
        name,
        contact,
        password,
        school: role === 'teacher' ? school : null
    };

    // Simulate API call and redirect
    showToast('Berhasil masuk! Mengarahkan ke dashboard...', 'success');
    
    setTimeout(() => {
        if (role === 'parent') {
            window.location.href = 'dashboard-parent.html';
        } else {
            window.location.href = 'dashboard-teacher.html';
        }
    }, 1500);
}
