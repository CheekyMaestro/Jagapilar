// ============================================================
// JAGAPILAR — Assessment Page Logic (Neuro ID Flow)
// Questionnaire items (18 items from PRD), Likert selection,
// progressive save, client-side scoring, API integration
// ============================================================

// --- 18 Assessment Items from PRD (3 pillars × 6 items) ---
const ASSESSMENT_ITEMS = {
    parent: [
        { code: 'O1', text: 'Anak saya tantrum/marah berlebihan saat gadget diambil atau dibatasi waktunya.', isReverse: false, construct: 'Dopamine withdrawal' },
        { code: 'O2', text: 'Anak saya terus meminta tambahan waktu bermain gadget meski sudah diberi batas waktu yang jelas.', isReverse: false, construct: 'Toleransi & eskalasi' },
        { code: 'O3', text: 'Anak saya sulit tidur atau tidur larut malam karena bermain gadget.', isReverse: false, construct: 'Gangguan tidur' },
        { code: 'O4', text: 'Anak saya masih bisa beraktivitas tanpa gadget selama lebih dari 60 menit tanpa rewel.', isReverse: true, construct: 'Delayed gratification' },
        { code: 'O5', text: 'Anak saya lebih memilih main gadget sendirian daripada bermain/mengobrol dengan keluarga atau teman.', isReverse: false, construct: 'Experiential avoidance' },
        { code: 'O6', text: 'Saya sendiri sering bermain gadget saat sedang menemani atau mengawasi anak.', isReverse: false, construct: 'Parental digital neglect' },
    ],
    teacher: [
        { code: 'G1', text: 'Anak ini mampu fokus mengerjakan tugas/mendengarkan penjelasan >15 menit tanpa teralihkan.', isReverse: true, construct: 'Sustained attention' },
        { code: 'G2', text: 'Anak ini tampak lelah, mengantuk, atau kurang bersemangat di kelas pada pagi hari.', isReverse: false, construct: 'Indikasi gangguan tidur' },
        { code: 'G3', text: 'Anak ini menunjukkan impulsivitas — sulit menunggu giliran, memotong pembicaraan, bertindak tanpa berpikir.', isReverse: false, construct: 'Impulsivitas / ADHD-like' },
        { code: 'G4', text: 'Anak ini kesulitan mengingat instruksi atau materi yang baru saja disampaikan.', isReverse: false, construct: 'Working memory' },
        { code: 'G5', text: 'Anak ini menarik diri dari interaksi sosial dengan teman sebaya saat jam istirahat.', isReverse: false, construct: 'Isolasi sosial' },
        { code: 'G6', text: 'Prestasi/performa akademik anak ini menurun dibanding awal semester tanpa sebab akademik yang jelas.', isReverse: false, construct: 'Penurunan performa akademik' },
    ],
    // student role here is fallback, though they use assessment-kids.html
    student: [
        { code: 'M1', text: 'Aku jadi gelisah atau nggak enak rasanya, kalau lama nggak pegang HP/gadget.', isReverse: false, construct: 'Craving / withdrawal subjektif' },
        { code: 'M2', text: 'Awalnya aku cuma mau main sebentar, tapi ujung-ujungnya main gadget lama banget.', isReverse: false, construct: 'Loss of control' },
        { code: 'M3', text: 'Aku main gadget biar bisa lupa, kalau lagi sedih, bosan, atau marah.', isReverse: false, construct: 'Experiential avoidance' },
        { code: 'M4', text: 'Aku masih suka main di luar rumah atau ngobrol langsung sama teman dan keluarga, dibanding main gadget terus-terusan.', isReverse: true, construct: 'Restorative engagement' },
        { code: 'M5', text: 'Aku suka kepikiran gadget terus, padahal lagi belajar atau ngerjain PR.', isReverse: false, construct: 'Intrusive thoughts' },
        { code: 'M6', text: 'Kalau lagi ada masalah, aku lebih milih main gadget sendirian daripada cerita ke orang tua atau guru.', isReverse: false, construct: 'Self-isolation' },
    ]
};

const ROLE_LABELS = {
    parent: { label: 'Orang Tua', icon: 'family_restroom' },
    teacher: { label: 'Guru', icon: 'school' },
    student: { label: 'Anak / Murid', icon: 'child_care' }
};

const LIKERT_LABELS = ['Tidak Pernah', 'Jarang', 'Kadang-kadang', 'Sering', 'Selalu'];

// --- State ---
let currentRole = null;
let currentNeuroId = null;
let currentQuestionIndex = 0;
let answers = {}; // { 'O1': 3, 'O2': 5, ... }

// --- Initialization ---
document.addEventListener('DOMContentLoaded', () => {
    // Check URL params for role and neuro_id
    const params = new URLSearchParams(window.location.search);
    const role = params.get('role');
    const neuroId = params.get('neuro_id');
    
    if (role && ASSESSMENT_ITEMS[role]) {
        if (role === 'student') {
            window.location.href = `assessment-kids.html?neuro_id=${neuroId || ''}`;
            return;
        }

        currentNeuroId = neuroId; // Can be null if they somehow got here without one
        selectRole(role);
        
        // Hide the intro selection manually
        document.getElementById('role-selection').classList.add('hidden');
        document.getElementById('assessment-form').classList.remove('hidden');
    }
});

/**
 * Select a role and show the assessment form
 */
function selectRole(role) {
    if (role === 'student') {
        window.location.href = 'assessment-kids.html';
        return;
    }
    
    currentRole = role;
    currentQuestionIndex = 0;
    answers = {};

    const roleInfo = ROLE_LABELS[role];

    // Update UI
    document.getElementById('role-selection').classList.add('hidden');
    document.getElementById('assessment-form').classList.remove('hidden');
    document.getElementById('results-section').classList.add('hidden');

    document.getElementById('role-icon').textContent = roleInfo.icon;
    document.getElementById('role-label').textContent = roleInfo.label;

    // Render first question
    renderQuestion();
}

/**
 * Reset role selection
 */
function resetRole() {
    currentRole = null;
    currentQuestionIndex = 0;
    answers = {};

    document.getElementById('role-selection').classList.remove('hidden');
    document.getElementById('assessment-form').classList.add('hidden');
    document.getElementById('results-section').classList.add('hidden');
}

/**
 * Render the current question
 */
function renderQuestion() {
    const items = ASSESSMENT_ITEMS[currentRole];
    const item = items[currentQuestionIndex];
    const total = items.length;

    // Update question text
    document.getElementById('q-code').textContent = item.code;
    document.getElementById('q-construct').textContent = item.construct;
    document.getElementById('q-text').textContent = item.text;

    // Update progress
    const pct = Math.round(((currentQuestionIndex + 1) / total) * 100);
    document.getElementById('current-q').textContent = currentQuestionIndex + 1;
    document.getElementById('progress-pct').textContent = pct + '%';
    document.getElementById('progress-bar').style.width = pct + '%';

    // Update navigation buttons
    document.getElementById('btn-prev').disabled = currentQuestionIndex === 0;

    if (currentQuestionIndex === total - 1) {
        document.getElementById('btn-next-text').textContent = 'LIHAT HASIL';
        document.getElementById('btn-next-icon').textContent = 'check_circle';
    } else {
        document.getElementById('btn-next-text').textContent = 'BERIKUTNYA';
        document.getElementById('btn-next-icon').textContent = 'chevron_right';
    }

    // Highlight selected answer if exists
    const selectedValue = answers[item.code];
    document.querySelectorAll('.likert-option').forEach(btn => {
        const val = parseInt(btn.getAttribute('data-value'));
        if (val === selectedValue) {
            btn.classList.add('selected', 'border-primary', 'bg-primary', 'text-white');
            btn.classList.remove('border-outline-variant/50', 'bg-surface-container-low');
            btn.querySelector('div').classList.add('border-white', 'text-white');
            btn.querySelector('div').classList.remove('border-outline-variant', 'text-on-surface-variant');
            btn.querySelector('span').classList.add('text-white');
            btn.querySelector('span').classList.remove('text-on-surface');
        } else {
            btn.classList.remove('selected', 'border-primary', 'bg-primary', 'text-white');
            btn.classList.add('border-outline-variant/50', 'bg-surface-container-low');
            btn.querySelector('div').classList.remove('border-white', 'text-white');
            btn.querySelector('div').classList.add('border-outline-variant', 'text-on-surface-variant');
            btn.querySelector('span').classList.remove('text-white');
            btn.querySelector('span').classList.add('text-on-surface');
        }
    });

    // Add reverse-scoring indicator
    if (item.isReverse) {
        document.getElementById('q-construct').textContent += ' ↺ (reverse)';
    }

    // Fade-in animation
    const container = document.getElementById('question-container');
    container.classList.remove('fade-in-up');
    void container.offsetWidth; // trigger reflow
    container.classList.add('fade-in-up');
}

/**
 * Select an answer for the current question
 */
function selectAnswer(value) {
    const items = ASSESSMENT_ITEMS[currentRole];
    const item = items[currentQuestionIndex];
    answers[item.code] = value;

    // Simulate progressive save
    console.log(`Saved item ${item.code} with value ${value} for Neuro ID ${currentNeuroId}`);

    // Update button styles
    renderQuestion();

    // Auto-advance after brief delay (unless last question)
    if (currentQuestionIndex < items.length - 1) {
        setTimeout(() => nextQuestion(), 400);
    }
}

/**
 * Go to next question or submit
 */
function nextQuestion() {
    const items = ASSESSMENT_ITEMS[currentRole];
    const item = items[currentQuestionIndex];

    // Validate current answer
    if (!answers[item.code]) {
        showToast('Silakan pilih jawaban terlebih dahulu', 'warning');
        return;
    }

    if (currentQuestionIndex < items.length - 1) {
        currentQuestionIndex++;
        renderQuestion();
    } else {
        // Submit and show results
        submitAssessment();
    }
}

/**
 * Go to previous question
 */
function prevQuestion() {
    if (currentQuestionIndex > 0) {
        currentQuestionIndex--;
        renderQuestion();
    }
}

/**
 * Calculate score and show results
 */
function submitAssessment() {
    const items = ASSESSMENT_ITEMS[currentRole];

    // Check all questions answered
    for (const item of items) {
        if (!answers[item.code]) {
            showToast(`Pertanyaan ${item.code} belum dijawab`, 'warning');
            return;
        }
    }

    // Calculate pillar score with reverse-scoring
    let totalScore = 0;
    const itemScores = [];

    for (const item of items) {
        const rawScore = answers[item.code];
        const effectiveScore = item.isReverse ? (6 - rawScore) : rawScore;
        totalScore += effectiveScore;
        itemScores.push({
            code: item.code,
            text: item.text,
            rawScore,
            effectiveScore,
            isReverse: item.isReverse,
            construct: item.construct
        });
    }

    // Classify zone
    let zone, zoneLabel, zoneColor, zoneBg, explanation;
    if (totalScore <= 13.9) {
        zone = 'hijau';
        zoneLabel = 'ZONA HIJAU: AMAN';
        zoneColor = 'text-zone-green';
        zoneBg = 'bg-zone-green';
        explanation = {
            title: 'Apa Artinya?',
            text: 'Aktivitas kognitif anak normal dan interaksi sosial sehat. Tetap lakukan pemantauan rutin sebagai langkah preventif.',
            bgClass: 'bg-zone-green/10 border-zone-green/20',
            titleColor: 'text-zone-green'
        };
    } else if (totalScore <= 21.9) {
        zone = 'kuning';
        zoneLabel = 'ZONA KUNING: WASPADA';
        zoneColor = 'text-zone-yellow';
        zoneBg = 'bg-zone-yellow';
        explanation = {
            title: 'Perlu Perhatian',
            text: 'Perlu pengawasan durasi dan konten digital harian. Pertimbangkan untuk membatasi waktu layar dan meningkatkan aktivitas alternatif.',
            bgClass: 'bg-zone-yellow/10 border-zone-yellow/20',
            titleColor: 'text-zone-yellow'
        };
    } else {
        zone = 'merah';
        zoneLabel = 'ZONA MERAH: PERLU TINDAKAN';
        zoneColor = 'text-zone-red';
        zoneBg = 'bg-zone-red';
        explanation = {
            title: 'Tindakan Segera Diperlukan',
            text: 'Segera konsultasikan dengan psikolog atau dokter anak. Skor menunjukkan risiko tinggi gangguan neurokognitif yang memerlukan intervensi profesional.',
            bgClass: 'bg-zone-red/10 border-zone-red/20',
            titleColor: 'text-zone-red'
        };
    }

    // Show results
    document.getElementById('assessment-form').classList.add('hidden');
    document.getElementById('results-section').classList.remove('hidden');

    // Update result elements
    document.getElementById('result-score').textContent = totalScore;
    document.getElementById('result-score').className = `text-6xl font-extrabold font-headline-lg mb-2 ${zoneColor}`;
    document.getElementById('result-icon-bg').className = `w-20 h-20 rounded-full mx-auto mb-6 flex items-center justify-center ${zoneBg}`;
    document.getElementById('result-zone-badge').textContent = zoneLabel;
    document.getElementById('result-zone-badge').className = `zone-badge inline-block px-6 py-2 rounded-full font-label-caps text-label-caps shadow-lg text-white ${zoneBg}`;
    document.getElementById('result-role-label').textContent = `Pilar ${ROLE_LABELS[currentRole].label}`;

    // Explanation
    document.getElementById('result-explanation').innerHTML = `
        <div class="p-4 rounded-xl ${explanation.bgClass} border">
            <h4 class="font-bold ${explanation.titleColor} mb-1">${explanation.title}</h4>
            <p class="text-sm text-on-surface-variant">${explanation.text}</p>
        </div>
    `;

    // Per-item breakdown
    const itemsHtml = itemScores.map(s => `
        <div class="flex items-center justify-between p-3 rounded-lg bg-surface-container-low">
            <div class="flex items-center gap-3">
                <span class="px-2 py-1 rounded bg-primary/10 text-primary text-xs font-bold">${s.code}</span>
                <span class="text-sm text-on-surface-variant truncate max-w-[300px]">${s.construct}</span>
                ${s.isReverse ? '<span class="text-xs text-nav-blue">↺ reverse</span>' : ''}
            </div>
            <div class="flex items-center gap-2">
                <span class="text-sm text-on-surface-variant">${LIKERT_LABELS[s.rawScore - 1]}</span>
                <span class="font-bold ${s.effectiveScore >= 4 ? 'text-zone-red' : s.effectiveScore >= 3 ? 'text-zone-yellow' : 'text-zone-green'}">${s.effectiveScore}</span>
            </div>
        </div>
    `).join('');
    document.getElementById('result-items').innerHTML = itemsHtml;

    // Scroll to top
    window.scrollTo({ top: 0, behavior: 'smooth' });

    console.log(`Submitted for ${currentNeuroId}`, itemScores);
}

/**
 * Reset assessment and start over
 */
function resetAssessment() {
    currentQuestionIndex = 0;
    answers = {};
    document.getElementById('results-section').classList.add('hidden');
    
    if (currentNeuroId) {
        document.getElementById('assessment-form').classList.remove('hidden');
        renderQuestion();
    } else {
        document.getElementById('role-selection').classList.remove('hidden');
    }
    window.scrollTo({ top: 0, behavior: 'smooth' });
}
