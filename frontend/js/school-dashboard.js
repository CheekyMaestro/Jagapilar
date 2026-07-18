// ============================================================
// JAGAPILAR — School Dashboard Logic
// Displays aggregate data only if minimum sample (>=5) is met.
// ============================================================

const MIN_SAMPLE_SIZE = 5;

// Mock API Data
const MOCK_DATA = {
    schoolName: "SDN 01 Percontohan",
    totalLinkedStudents: 12,
    totalCompletedAssessments: 8, // Changed to 8 to show aggregate data (>=5). Change to <5 to test lockdown.
    zones: {
        hijau: 3,
        kuning: 4,
        merah: 1
    },
    classes: [
        { name: "Kelas 4A", total: 4, hijau: 1, kuning: 2, merah: 1 },
        { name: "Kelas 5B", total: 4, hijau: 2, kuning: 2, merah: 0 }
    ]
};

document.addEventListener('DOMContentLoaded', () => {
    // In real app, fetch from API. We use mock data here.
    renderDashboard(MOCK_DATA);
});

function renderDashboard(data) {
    document.getElementById('school-name').textContent = data.schoolName;
    document.getElementById('stat-total').textContent = data.totalLinkedStudents;
    
    document.getElementById('completed-count').textContent = data.totalCompletedAssessments;
    const progressPct = Math.min((data.totalCompletedAssessments / MIN_SAMPLE_SIZE) * 100, 100);
    document.getElementById('completed-bar').style.width = `${progressPct}%`;

    if (data.totalCompletedAssessments < MIN_SAMPLE_SIZE) {
        // Lockdown mode
        document.getElementById('insufficient-data-alert').classList.remove('hidden');
        document.getElementById('aggregated-data-view').classList.add('hidden');
        
        document.getElementById('stat-hijau').textContent = '-';
        document.getElementById('stat-kuning').textContent = '-';
        document.getElementById('stat-merah').textContent = '-';
    } else {
        // Show Aggregate
        document.getElementById('insufficient-data-alert').classList.add('hidden');
        document.getElementById('aggregated-data-view').classList.remove('hidden');
        
        // Populate Zone Stats
        document.getElementById('stat-hijau').textContent = data.zones.hijau;
        document.getElementById('stat-kuning').textContent = data.zones.kuning;
        document.getElementById('stat-merah').textContent = data.zones.merah;

        // Render Class Distribution
        const classList = document.getElementById('class-distribution-list');
        classList.innerHTML = '';
        
        data.classes.forEach(cls => {
            const row = document.createElement('div');
            row.className = "p-4 bg-surface-container-low rounded-lg border border-outline-variant/50";
            
            const pctHijau = Math.round((cls.hijau / cls.total) * 100) || 0;
            const pctKuning = Math.round((cls.kuning / cls.total) * 100) || 0;
            const pctMerah = Math.round((cls.merah / cls.total) * 100) || 0;

            row.innerHTML = `
                <div class="flex justify-between items-center mb-3">
                    <h4 class="font-bold text-on-surface">${cls.name} <span class="text-xs text-on-surface-variant font-normal">(${cls.total} murid selesai)</span></h4>
                </div>
                <div class="w-full h-4 flex rounded-full overflow-hidden bg-surface-variant shadow-inner">
                    <div class="bg-zone-green h-full" style="width: ${pctHijau}%" title="Hijau: ${pctHijau}%"></div>
                    <div class="bg-zone-yellow h-full" style="width: ${pctKuning}%" title="Kuning: ${pctKuning}%"></div>
                    <div class="bg-zone-red h-full" style="width: ${pctMerah}%" title="Merah: ${pctMerah}%"></div>
                </div>
                <div class="flex gap-4 mt-2 text-xs font-label-caps text-on-surface-variant">
                    <span class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-zone-green"></div> ${pctHijau}% Hijau</span>
                    <span class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-zone-yellow"></div> ${pctKuning}% Kuning</span>
                    <span class="flex items-center gap-1"><div class="w-2 h-2 rounded-full bg-zone-red"></div> ${pctMerah}% Merah</span>
                </div>
            `;
            classList.appendChild(row);
        });
    }
}
