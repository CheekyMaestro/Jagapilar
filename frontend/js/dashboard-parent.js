// ============================================================
// JAGAPILAR — Parent Dashboard Logic
// Handle generation of Neuro ID and rendering child cards
// ============================================================

// Mock state
let children = [
    // { id: 'NP-A3F7-K2M9', name: 'Budi', year: 2012, gender: 'L', status: 'pending' }
];

document.addEventListener('DOMContentLoaded', () => {
    renderChildren();

    // Checkbox toggle state
    const checkbox = document.getElementById('consent-checkbox');
    const btnGen = document.getElementById('btn-generate');
    if(checkbox && btnGen) {
        checkbox.addEventListener('change', (e) => {
            btnGen.disabled = !e.target.checked;
        });
    }
});

function openModal() {
    document.getElementById('add-child-modal').classList.remove('hidden');
    document.getElementById('add-child-form').classList.remove('hidden');
    document.getElementById('result-container').classList.add('hidden');
    document.getElementById('add-child-form').reset();
    document.getElementById('btn-generate').disabled = true;
}

function closeModal() {
    document.getElementById('add-child-modal').classList.add('hidden');
    renderChildren(); // Re-render when closing
}

function generateRandomNeuroID() {
    const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'; // No 0,O,1,I,L
    let p1 = '';
    let p2 = '';
    for(let i=0; i<4; i++) {
        p1 += chars.charAt(Math.floor(Math.random() * chars.length));
        p2 += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return `NP-${p1}-${p2}`;
}

const API_BASE = '/api';

async function generateNeuroID(e) {
    e.preventDefault();
    const name = document.getElementById('child-name').value;
    const year = document.getElementById('child-year').value;
    const gender = document.getElementById('child-gender').value;

    const token = localStorage.getItem('jagapilar_token');
    let newId = '';

    try {
        const res = await fetch(`${API_BASE}/children`, {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ name, birth_year: parseInt(year), gender })
        });

        if (!res.ok) throw new Error('API Error');
        const data = await res.json();
        newId = data.neuro_id;
    } catch (error) {
        console.warn('API gagal, menggunakan mode offline/mock.', error);
        newId = generateRandomNeuroID();
    }
    
    // Save to mock state
    children.push({
        id: newId,
        name: name,
        year: year,
        gender: gender,
        status: 'pending' // pending | parent_done | complete
    });

    document.getElementById('add-child-form').classList.add('hidden');
    document.getElementById('result-container').classList.remove('hidden');
    document.getElementById('generated-neuro-id').innerText = newId;
    
    showToast('Neuro ID berhasil dibuat!', 'success');
}

function copyNeuroIdFromModal() {
    const text = document.getElementById('generated-neuro-id').innerText;
    navigator.clipboard.writeText(text);
    showToast('Neuro ID disalin ke clipboard!', 'success');
}

function copyNeuroId(id) {
    navigator.clipboard.writeText(id);
    showToast('Neuro ID disalin ke clipboard!', 'success');
}

function shareNeuroId(id, name) {
    const text = `Halo, tolong bantu isi asesmen JAGAPILAR untuk anak saya ${name}. Neuro ID: ${id}. Silakan buka https://jagapilar.id/ dan pilih peran Anda.`;
    const encoded = encodeURIComponent(text);
    window.open(`https://wa.me/?text=${encoded}`, '_blank');
}

function renderChildren() {
    const list = document.getElementById('children-list');
    list.innerHTML = '';

    if (children.length === 0) {
        list.innerHTML = `
            <div class="text-center py-16 bg-surface-container-lowest rounded-xl border border-outline-variant">
                <span class="material-symbols-outlined text-6xl text-surface-variant mb-4">child_care</span>
                <h3 class="font-headline-md text-on-surface-variant">Belum ada anak</h3>
                <p class="text-sm mt-2">Klik tombol "Tambah Anak" untuk memulai.</p>
            </div>
        `;
        return;
    }

    children.forEach(child => {
        let statusBadge = '';
        let actionBtn = '';

        if (child.status === 'pending') {
            statusBadge = `<span class="bg-surface-variant text-on-surface text-xs font-bold px-2 py-1 rounded">MENUNGGU ASESMEN</span>`;
            actionBtn = `<a href="assessment.html?role=parent&neuro_id=${child.id}" class="bg-primary text-white font-label-caps px-4 py-2 rounded flex items-center gap-1 hover:bg-teal-deep"><span class="material-symbols-outlined text-sm">edit_document</span> ISI PILAR ORANG TUA</a>`;
        } else if (child.status === 'parent_done') {
            statusBadge = `<span class="bg-zone-yellow/20 text-zone-yellow text-xs font-bold px-2 py-1 rounded">MENUNGGU GURU & ANAK</span>`;
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
                    <h3 class="font-headline-md text-primary">${child.name}</h3>
                    ${statusBadge}
                </div>
                <div class="flex items-center gap-4 text-sm text-on-surface-variant mb-4">
                    <span class="flex items-center gap-1"><span class="material-symbols-outlined text-sm">calendar_today</span> Lhr: ${child.year}</span>
                    <span class="flex items-center gap-1"><span class="material-symbols-outlined text-sm">wc</span> ${child.gender === 'L' ? 'Laki-laki' : 'Perempuan'}</span>
                </div>
                
                <div class="bg-surface-container-low p-3 rounded-lg flex items-center justify-between max-w-sm">
                    <div class="flex flex-col">
                        <span class="text-xs font-bold text-on-surface-variant">NEURO ID</span>
                        <span class="font-mono text-primary font-bold tracking-widest">${child.id}</span>
                    </div>
                    <div class="flex gap-2">
                        <button onclick="copyNeuroId('${child.id}')" class="p-2 hover:bg-surface-variant rounded-full text-primary" title="Salin">
                            <span class="material-symbols-outlined text-sm">content_copy</span>
                        </button>
                        <button onclick="shareNeuroId('${child.id}', '${child.name}')" class="p-2 hover:bg-surface-variant rounded-full text-zone-green" title="Bagikan via WA">
                            <span class="material-symbols-outlined text-sm">share</span>
                        </button>
                    </div>
                </div>
            </div>
            <div class="flex flex-col items-end gap-3 min-w-[200px]">
                ${actionBtn}
                <button class="text-sm font-label-caps text-on-surface-variant hover:text-primary flex items-center gap-1">
                    <span class="material-symbols-outlined text-sm">domain_add</span> HUBUNGKAN KE SEKOLAH
                </button>
            </div>
        `;
        list.appendChild(card);
    });
}
