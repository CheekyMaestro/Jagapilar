// ============================================================
// JAGAPILAR — Shared Tailwind Configuration
// Design tokens extracted from the monolith HTML
// ============================================================

tailwind.config = {
    darkMode: "class",
    theme: {
        extend: {
            colors: {
                "tertiary-container": "#a10075",
                "surface-bright": "#f9f9fd",
                "tertiary-fixed-dim": "#ffaed8",
                "on-surface-variant": "#41474e",
                "on-primary-fixed-variant": "#0f4b70",
                "surface-tint": "#306289",
                "surface-variant": "#e2e2e6",
                "on-error-container": "#93000a",
                "primary-fixed-dim": "#9bccf7",
                "primary": "#004064",
                "on-background": "#191c1e",
                "zone-yellow": "#F5C363",
                "on-primary-fixed": "#001d31",
                "surface-dim": "#d9dadd",
                "on-surface": "#191c1e",
                "outline": "#72787f",
                "primary-container": "#23587e",
                "tertiary": "#780056",
                "on-tertiary-fixed-variant": "#890063",
                "on-secondary-container": "#007169",
                "surface": "#f9f9fd",
                "error": "#ba1a1a",
                "zone-green": "#178754",
                "background": "#f9f9fd",
                "surface-container-lowest": "#ffffff",
                "on-primary": "#ffffff",
                "inverse-surface": "#2e3133",
                "sky-blue": "#00A9E6",
                "teal-deep": "#10868B",
                "on-error": "#ffffff",
                "inverse-on-surface": "#f0f0f4",
                "surface-container-low": "#f3f3f7",
                "on-tertiary-container": "#ffb2da",
                "secondary-fixed": "#59faec",
                "surface-container-highest": "#e2e2e6",
                "on-tertiary": "#ffffff",
                "nav-blue": "#479BCE",
                "zone-red": "#D93A3F",
                "primary-fixed": "#cce5ff",
                "on-secondary-fixed": "#00201d",
                "pink-bright": "#FF50C9",
                "surface-container-high": "#e7e8eb",
                "inverse-primary": "#9bccf7",
                "error-container": "#ffdad6",
                "secondary-container": "#59faec",
                "on-secondary-fixed-variant": "#00504a",
                "teal-light": "#32A8AD",
                "surface-container": "#edeef1",
                "outline-variant": "#c1c7cf",
                "secondary": "#006a63",
                "on-tertiary-fixed": "#3c0029",
                "secondary-fixed-dim": "#2eddcf",
                "tertiary-fixed": "#ffd8e9",
                "on-secondary": "#ffffff",
                "on-primary-container": "#9dcef9"
            },
            borderRadius: {
                DEFAULT: "0.25rem",
                lg: "0.5rem",
                xl: "0.75rem",
                xxl: "1.5rem",
                full: "9999px"
            },
            spacing: {
                "stack-lg": "2rem",
                "stack-sm": "0.5rem",
                "margin-mobile": "1.5rem",
                "container-max": "1280px",
                "stack-md": "1rem",
                "margin-desktop": "5rem",
                gutter: "1rem"
            },
            fontFamily: {
                "body-md": ["Inter"],
                "nav-link": ["Inter"],
                "headline-lg-mobile": ["Poppins"],
                "label-caps": ["Poppins"],
                "headline-md": ["Poppins"],
                "headline-lg": ["Poppins"],
                "body-lg": ["Inter"]
            },
            fontSize: {
                "body-md": ["16px", { lineHeight: "1.5", fontWeight: "400" }],
                "nav-link": ["16px", { lineHeight: "1", fontWeight: "500" }],
                "headline-lg-mobile": ["32px", { lineHeight: "1.2", fontWeight: "800" }],
                "label-caps": ["14px", { lineHeight: "1", letterSpacing: "0.1em", fontWeight: "700" }],
                "headline-md": ["24px", { lineHeight: "1.3", fontWeight: "700" }],
                "headline-lg": ["48px", { lineHeight: "1.2", letterSpacing: "-0.02em", fontWeight: "800" }],
                "body-lg": ["18px", { lineHeight: "1.6", fontWeight: "400" }]
            }
        }
    }
};
