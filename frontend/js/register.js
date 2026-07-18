// ============================================================
// JAGAPILAR — Registration Page Logic
// Step wizard, file upload, form validation, API submission
// ============================================================

/**
 * Navigate between steps in the registration wizard
 */
function goToStep(step) {
    const step1 = document.getElementById('step-1');
    const step2 = document.getElementById('step-2');
    const dot2 = document.getElementById('dot-2');
    const lineFill = document.getElementById('line-fill');

    if (step === 1) {
        step2.classList.remove('active-step');
        step2.classList.add('hidden-step');
        setTimeout(() => {
            step1.classList.remove('hidden-step');
            step1.classList.add('active-step');
        }, 300);

        dot2.classList.replace('bg-primary', 'bg-surface-container-high');
        dot2.classList.replace('text-white', 'text-on-surface-variant');
        lineFill.style.width = '0%';
    } else if (step === 2) {
        // Validate step 1 fields before proceeding
        const name = document.getElementById('school-name').value.trim();
        const city = document.getElementById('school-city').value.trim();
        const grade = document.getElementById('school-grade').value;
        const principal = document.getElementById('school-principal').value.trim();
        const classes = document.getElementById('school-classes').value;

        if (!name || !city || !grade || !principal || !classes) {
            showToast('Mohon lengkapi semua field terlebih dahulu', 'warning');
            return;
        }

        step1.classList.remove('active-step');
        step1.classList.add('hidden-step');
        setTimeout(() => {
            step2.classList.remove('hidden-step');
            step2.classList.add('active-step');
        }, 300);

        dot2.classList.replace('bg-surface-container-high', 'bg-primary');
        dot2.classList.replace('text-on-surface-variant', 'text-white');
        lineFill.style.width = '100%';
    }
}

/**
 * Handle file upload preview
 */
function handleFile(event) {
    const file = event.target.files[0];
    if (!file) return;

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
        showToast('Ukuran file maksimal 5MB', 'error');
        event.target.value = '';
        return;
    }

    // Validate file type
    const allowedTypes = ['application/pdf', 'image/jpeg', 'image/png', 'image/jpg'];
    if (!allowedTypes.includes(file.type)) {
        showToast('Format file harus PDF atau Image (JPG/PNG)', 'error');
        event.target.value = '';
        return;
    }

    const preview = document.getElementById('file-preview');
    const nameText = document.getElementById('file-name-text');
    preview.classList.remove('hidden');
    nameText.innerText = file.name;

    // Success animation
    preview.classList.add('scale-105');
    setTimeout(() => preview.classList.remove('scale-105'), 200);

    showToast('Berkas berhasil dipilih', 'success');
}

/**
 * Submit school registration to the backend API
 */
async function submitRegistration() {
    const schoolData = {
        name: document.getElementById('school-name').value.trim(),
        city: document.getElementById('school-city').value.trim(),
        grade_level: document.getElementById('school-grade').value,
        principal_name: document.getElementById('school-principal').value.trim(),
        total_classes: parseInt(document.getElementById('school-classes').value, 10)
    };

    // Validate all fields
    if (!schoolData.name || !schoolData.city || !schoolData.grade_level || !schoolData.principal_name || !schoolData.total_classes) {
        showToast('Mohon lengkapi semua data sekolah', 'warning');
        return;
    }

    try {
        const result = await apiCall('/schools', {
            method: 'POST',
            body: JSON.stringify(schoolData),
        });

        showToast('Sekolah berhasil didaftarkan! 🎉', 'success', 5000);

        // Handle file upload if present
        const fileInput = document.getElementById('file-upload');
        if (fileInput.files.length > 0) {
            const formData = new FormData();
            formData.append('file', fileInput.files[0]);
            formData.append('school_id', result.id);

            await fetch(`${API_BASE}/schools/upload`, {
                method: 'POST',
                body: formData,
            });
        }

        // Redirect to school dashboard after delay
        setTimeout(() => {
            window.location.href = `school-dashboard.html?id=${result.id}`;
        }, 1500);

    } catch (error) {
        console.error('Registration failed:', error);
    }
}
