// ============================================================
// JAGAPILAR — Teacher Dashboard Logic
// Add student via Neuro ID and display list
// ============================================================

let students = [
    // { id: 'NP-A3F7-K2M9', name: 'Budi (Data Mock)', year: 2012, status: 'pending' }
];

document.addEventListener('DOMContentLoaded', () => {
    renderStudents();

    const checkbox = document.getElementById('teacher-confirm');
    const btnAdd = document.getElementById('btn-add-student');
    if(checkbox && btnAdd) {
        checkbox.addEventListener('change', (e) => {
            btnAdd.disabled = !e.target.checked;
        });
    }
});

function openModal() {
    document.getElementById('add-student-modal').classList.remove('hidden');
    document.getElementById('add-student-form').reset();
    document.getElementById('btn-add-student').disabled = true;
}

function closeModal() {
    document.getElementById('add-student-modal').classList.add('hidden');
}

function addStudent(e) {
    e.preventDefault();
    const neuroId = document.getElementById('student-neuro-id').value.trim().toUpperCase();

    // Regex check basic format
    if (!/^NP-[A-Z0-9]{4}-[A-Z0-9]{4}$/.test(neuroId)) {
        showToast('Format Neuro ID tidak valid. Harusnya NP-XXXX-XXXX', 'error');
        return;
    }

    if(students.find(s => s.id === neuroId)) {
        showToast('Neuro ID ini sudah ada di dashboard Anda.', 'warning');
        return;
    }

    // Simulate API fetch to validate Neuro ID
    setTimeout(() => {
        // Mock data
        students.push({
            id: neuroId,
            name: 'Anak (Data Tersembunyi)',
            year: 2012,
            status: 'pending'
        });

        closeModal();
        renderStudents();
        showToast('Murid berhasil ditambahkan!', 'success');
    }, 500);
}

function renderStudents() {
    const list = document.getElementById('students-list');
    list.innerHTML = '';

    if (students.length === 0) {
        list.innerHTML = `
            <div class="text-center py-16 bg-surface-container-lowest rounded-xl border border-outline-variant">
                <span class="material-symbols-outlined text-6xl text-surface-variant mb-4">school</span>
                <h3 class="font-headline-md text-on-surface-variant">Belum ada murid</h3>
                <p class="text-sm mt-2">Dapatkan Neuro ID dari orang tua dan klik "Tambah Murid".</p>
            </div>
        `;
        return;
    }

    students.forEach(student => {
        let statusBadge = '';
        let actionBtn = '';

        if (student.status === 'pending') {
            statusBadge = `<span class="bg-surface-variant text-on-surface text-xs font-bold px-2 py-1 rounded">MENUNGGU ASESMEN GURU</span>`;
            actionBtn = `<a href="assessment.html?role=teacher&neuro_id=${student.id}" class="bg-primary text-white font-label-caps px-4 py-2 rounded flex items-center gap-1 hover:bg-teal-deep"><span class="material-symbols-outlined text-sm">edit_document</span> ISI PILAR GURU</a>`;
        } else if (student.status === 'teacher_done') {
            statusBadge = `<span class="bg-zone-yellow/20 text-zone-yellow text-xs font-bold px-2 py-1 rounded">MENUNGGU ORANG TUA/ANAK</span>`;
            actionBtn = `<button disabled class="bg-surface-variant text-on-surface/50 font-label-caps px-4 py-2 rounded flex items-center gap-1"><span class="material-symbols-outlined text-sm">check</span> PILAR ANDA SELESAI</button>`;
        } else {
            statusBadge = `<span class="bg-zone-green/20 text-zone-green text-xs font-bold px-2 py-1 rounded">SELESAI 3 PILAR</span>`;
            actionBtn = `<button class="bg-secondary text-white font-label-caps px-4 py-2 rounded flex items-center gap-1"><span class="material-symbols-outlined text-sm">visibility</span> LIHAT HASIL 360°</button>`;
        }

        const card = document.createElement('div');
        card.className = "bg-surface-container-lowest p-6 rounded-xl border border-outline-variant shadow-sm flex flex-col md:flex-row md:items-center justify-between gap-6";
        card.innerHTML = `
            <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                    <h3 class="font-headline-md text-primary font-mono tracking-widest">${student.id}</h3>
                    ${statusBadge}
                </div>
                <div class="text-sm text-on-surface-variant mb-4">
                    Inisial / Nama: ${student.name}
                </div>
            </div>
            <div class="flex flex-col items-end gap-3 min-w-[200px]">
                ${actionBtn}
            </div>
        `;
        list.appendChild(card);
    });
}
