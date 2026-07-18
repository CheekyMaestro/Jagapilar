// ============================================================
// JAGAPILAR Kids — Assessment Logic
// ============================================================

const STUDENT_ITEMS = [
    { code: 'M1', text: 'Aku jadi gelisah atau nggak enak rasanya, kalau lama nggak pegang HP/gadget.', isReverse: false, construct: 'Craving / withdrawal subjektif' },
    { code: 'M2', text: 'Awalnya aku cuma mau main sebentar, tapi ujung-ujungnya main gadget lama banget.', isReverse: false, construct: 'Loss of control' },
    { code: 'M3', text: 'Aku main gadget biar bisa lupa, kalau lagi sedih, bosan, atau marah.', isReverse: false, construct: 'Experiential avoidance' },
    { code: 'M4', text: 'Aku masih suka main di luar rumah atau ngobrol langsung sama teman dan keluarga, dibanding main gadget terus-terusan.', isReverse: true, construct: 'Restorative engagement' },
    { code: 'M5', text: 'Aku suka kepikiran gadget terus, padahal lagi belajar atau ngerjain PR.', isReverse: false, construct: 'Intrusive thoughts' },
    { code: 'M6', text: 'Kalau lagi ada masalah, aku lebih milih main gadget sendirian daripada cerita ke orang tua atau guru.', isReverse: false, construct: 'Self-isolation' },
];

let currentToken = null;
let currentSessionId = null;
let currentQuestionIndex = 0;
let answers = {}; // { 'M1': 3, 'M2': 5, ... }

document.addEventListener('DOMContentLoaded', async () => {
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token');
    
    if (token) {
        currentToken = token;
        try {
            // Validate Token
            const res = await apiCall('/auth/validate-token', {
                method: 'POST',
                body: JSON.stringify({ token: token })
            });
            
            if (res && res.valid && res.role === 'student') {
                // Create session
                const sessionRes = await fetch(`${API_BASE}/assessment/sessions`, {
                    method: 'POST',
                    headers: { 'Authorization': `Bearer ${token}` }
                });
                
                if(sessionRes.ok) {
                    const session = await sessionRes.json();
                    currentSessionId = session.id;
                    
                    // Show questions
                    document.getElementById('loading-section').classList.add('hidden');
                    document.getElementById('question-section').classList.remove('hidden');
                    renderQuestion();
                } else {
                    showError('Gagal membuat sesi asesmen.');
                }
            } else {
                showError('Link asesmen tidak valid atau bukan untuk murid.');
            }
        } catch(e) {
            showError('Gagal memvalidasi link asesmen.');
        }
    } else {
        // Demo mode
        document.getElementById('loading-section').classList.add('hidden');
        document.getElementById('question-section').classList.remove('hidden');
        renderQuestion();
    }
});

function showError(msg) {
    document.getElementById('loading-section').innerHTML = `
        <div class="text-error font-headline-md mb-2">Oops!</div>
        <p class="text-on-surface-variant">${msg}</p>
        <button onclick="window.location.href='index.html'" class="mt-6 text-nav-link font-nav-link text-primary hover:text-nav-blue underline">Kembali ke Beranda</button>
    `;
}

function renderQuestion() {
    const item = STUDENT_ITEMS[currentQuestionIndex];
    const total = STUDENT_ITEMS.length;

    // Update Progress
    document.getElementById('level-indicator').textContent = `LEVEL ${currentQuestionIndex + 1} DARI ${total}`;
    document.getElementById('nav-progress-text').textContent = `${currentQuestionIndex + 1}/${total}`;
    
    // Update Stars
    let starsHtml = '';
    for(let i=0; i<total; i++) {
        if(i <= currentQuestionIndex) {
            starsHtml += `<span class="material-symbols-outlined text-3xl md:text-4xl star-active">stars</span>`;
        } else {
            starsHtml += `<span class="material-symbols-outlined text-3xl md:text-4xl star-inactive">stars</span>`;
        }
    }
    document.getElementById('stars-container').innerHTML = starsHtml;

    // Update Text
    document.getElementById('question-text').textContent = `"${item.text}"`;

    // Reset emojis
    document.querySelectorAll('.emoji-btn').forEach(btn => {
        btn.classList.remove('active');
    });

    // Re-apply if already answered
    const existingAnswer = answers[item.code];
    if (existingAnswer) {
        // Find the button (nth child)
        const btn = document.querySelectorAll('.emoji-btn')[existingAnswer - 1];
        if (btn) btn.classList.add('active');
    }

    // Back button visibility
    const backBtn = document.getElementById('btn-back');
    if (currentQuestionIndex === 0) {
        backBtn.classList.add('opacity-50', 'pointer-events-none');
    } else {
        backBtn.classList.remove('opacity-50', 'pointer-events-none');
    }

    // Button state
    updateNextButtonState();

    // Trigger Buddy Animation
    const buddy = document.getElementById('jaga-buddy');
    buddy.classList.remove('animate-bounce');
    void buddy.offsetWidth;
    buddy.classList.add('animate-bounce');
}

function selectEmoji(element, value) {
    const item = STUDENT_ITEMS[currentQuestionIndex];
    answers[item.code] = value;
    
    document.querySelectorAll('.emoji-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    element.classList.add('active');

    // Progressive save
    if (currentSessionId && currentToken) {
        fetch(`${API_BASE}/assessment/responses`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify({
                session_id: currentSessionId,
                item_code: item.code,
                raw_score: value
            })
        }).catch(e => console.log('Progressive save fail:', e));
    }

    updateNextButtonState();
    
    // Auto advance after small delay
    if (currentQuestionIndex < STUDENT_ITEMS.length - 1) {
        setTimeout(() => nextQuestion(), 500);
    }
}

function updateNextButtonState() {
    const item = STUDENT_ITEMS[currentQuestionIndex];
    const btn = document.getElementById('next-btn');
    if (answers[item.code]) {
        btn.classList.remove('opacity-50', 'cursor-not-allowed');
        btn.disabled = false;
        
        if (currentQuestionIndex === STUDENT_ITEMS.length - 1) {
            btn.textContent = 'SELESAI';
        } else {
            btn.textContent = 'LANJUT';
        }
    } else {
        btn.classList.add('opacity-50', 'cursor-not-allowed');
        btn.disabled = true;
    }
}

function nextQuestion() {
    const item = STUDENT_ITEMS[currentQuestionIndex];
    if (!answers[item.code]) return;

    if (currentQuestionIndex < STUDENT_ITEMS.length - 1) {
        currentQuestionIndex++;
        renderQuestion();
    } else {
        submitAssessment();
    }
}

function prevQuestion() {
    if (currentQuestionIndex > 0) {
        currentQuestionIndex--;
        renderQuestion();
    }
}

function submitAssessment() {
    let totalScore = 0;
    const itemScores = [];

    for (const item of STUDENT_ITEMS) {
        const rawScore = answers[item.code];
        const effectiveScore = item.isReverse ? (6 - rawScore) : rawScore;
        totalScore += effectiveScore;
        itemScores.push({
            code: item.code,
            rawScore: rawScore
        });
    }

    // Submit to API
    if (currentToken) {
        fetch(`${API_BASE}/assessment/submit`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${currentToken}`
            },
            body: JSON.stringify({
                pillar: 'student',
                total_score: totalScore,
                zone: 'unknown', // backend will calculate and ignore this anyway for single pillar
                items: itemScores.map(s => ({
                    item_code: s.code,
                    raw_score: s.rawScore,
                }))
            })
        }).catch(e => console.log(e));
    }

    // Show Success Screen
    document.getElementById('question-section').classList.add('hidden');
    document.getElementById('success-section').classList.remove('hidden');
    
    // Hide bottom nav
    document.getElementById('bottom-nav').classList.add('translate-y-full');
    
    // Trigger Confetti
    createConfetti();
    setInterval(createConfetti, 5000);
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Confetti Logic
const colors = ['#00A9E6', '#FF50C9', '#F5C363', '#178754', '#004064'];
function createConfetti() {
    const container = document.getElementById('confetti-container');
    if(!container) return;
    
    for (let i = 0; i < 60; i++) {
        const confetti = document.createElement('div');
        confetti.className = 'confetti-piece animate-confetti';
        
        const color = colors[Math.floor(Math.random() * colors.length)];
        const x = (Math.random() - 0.5) * 800 + 'px';
        const y = (Math.random() - 0.5) * 800 - 100 + 'px';
        const r = Math.random() * 720 + 'deg';
        const delay = Math.random() * 2 + 's';

        confetti.style.backgroundColor = color;
        confetti.style.setProperty('--x', x);
        confetti.style.setProperty('--y', y);
        confetti.style.setProperty('--r', r);
        confetti.style.animationDelay = delay;
        
        if (Math.random() > 0.5) confetti.style.borderRadius = '50%';
        container.appendChild(confetti);

        setTimeout(() => {
            confetti.remove();
        }, 4500);
    }
}
