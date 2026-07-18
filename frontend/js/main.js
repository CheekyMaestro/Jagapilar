// ============================================================
// JAGAPILAR — Shared JavaScript Utilities
// Navbar scroll, button interactions, toast notifications
// ============================================================

/**
 * Navbar scroll effect — shrinks on scroll
 */
function initNavbarScroll() {
    const nav = document.querySelector('nav');
    if (!nav) return;

    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            nav.classList.add('shadow-lg', 'h-[80px]');
            nav.classList.remove('h-[90px]');
        } else {
            nav.classList.remove('shadow-lg', 'h-[80px]');
            nav.classList.add('h-[90px]');
        }
    });
}

/**
 * Button micro-interactions — press effect
 */
function initButtonInteractions() {
    document.querySelectorAll('button').forEach(btn => {
        btn.addEventListener('mousedown', () => {
            btn.style.transform = 'scale(0.95)';
        });
        btn.addEventListener('mouseup', () => {
            btn.style.transform = 'scale(1)';
        });
        btn.addEventListener('mouseleave', () => {
            btn.style.transform = 'scale(1)';
        });
    });
}

/**
 * Mobile Navigation Toggle
 */
function initMobileNav() {
    const menuBtn = document.getElementById('mobile-menu-btn');
    const mobileMenu = document.getElementById('mobile-menu');
    if (!menuBtn || !mobileMenu) return;

    menuBtn.addEventListener('click', () => {
        mobileMenu.classList.toggle('hidden');
    });
}

/**
 * Toast notification system
 */
function showToast(message, type = 'success', duration = 3000) {
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    toast.textContent = message;
    document.body.appendChild(toast);

    setTimeout(() => {
        toast.style.animation = 'slideInRight 0.3s ease-out reverse';
        setTimeout(() => toast.remove(), 300);
    }, duration);
}

/**
 * API Helper — fetch wrapper with error handling
 */
const API_BASE = window.location.origin + '/api';

async function apiCall(endpoint, options = {}) {
    const url = `${API_BASE}${endpoint}`;
    const defaultHeaders = {
        'Content-Type': 'application/json',
    };

    // Add auth token if available
    const token = localStorage.getItem('jagapilar_token');
    if (token) {
        defaultHeaders['Authorization'] = `Bearer ${token}`;
    }

    try {
        const response = await fetch(url, {
            headers: { ...defaultHeaders, ...options.headers },
            ...options,
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || `HTTP ${response.status}`);
        }

        return data;
    } catch (error) {
        console.error(`API Error [${endpoint}]:`, error);
        showToast(error.message || 'Terjadi kesalahan', 'error');
        throw error;
    }
}

/**
 * Input label micro-interaction — highlight on focus
 */
function initInputInteractions() {
    document.querySelectorAll('input, select, textarea').forEach(element => {
        element.addEventListener('focus', () => {
            const label = element.parentElement?.querySelector('label');
            if (label) label.classList.add('text-primary');
        });
        element.addEventListener('blur', () => {
            const label = element.parentElement?.querySelector('label');
            if (label) label.classList.remove('text-primary');
        });
    });
}

/**
 * Initialize all shared behaviors
 */
document.addEventListener('DOMContentLoaded', () => {
    initNavbarScroll();
    initButtonInteractions();
    initMobileNav();
    initInputInteractions();
});
