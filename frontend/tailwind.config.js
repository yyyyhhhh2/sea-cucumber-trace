/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Plus Jakarta Sans", "ui-sans-serif", "system-ui", "sans-serif"],
        display: ["Outfit", "Plus Jakarta Sans", "sans-serif"],
      },
      colors: {
        abyss: { 950: "#020617", 900: "#0f172a", 800: "#1e293b" },
        sea: { 400: "#22d3ee", 500: "#06b6d4", 600: "#0891b2" },
      },
      boxShadow: {
        glow: "0 0 80px -20px rgba(34, 211, 238, 0.45)",
        card: "0 25px 50px -12px rgba(2, 6, 23, 0.65)",
      },
      backgroundImage: {
        "hero-mesh":
          "radial-gradient(ellipse 80% 60% at 50% -10%, rgba(34,211,238,0.22), transparent 55%), radial-gradient(ellipse 60% 50% at 100% 0%, rgba(6,182,212,0.12), transparent 45%)",
        "card-shine":
          "linear-gradient(135deg, rgba(255,255,255,0.12) 0%, transparent 42%, rgba(34,211,238,0.08) 100%)",
      },
      keyframes: {
        "fade-up": {
          "0%": { opacity: "0", transform: "translateY(14px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
        "fade-in": {
          "0%": { opacity: "0" },
          "100%": { opacity: "1" },
        },
        float: {
          "0%, 100%": { transform: "translateY(0)" },
          "50%": { transform: "translateY(-6px)" },
        },
        drift: {
          "0%": { transform: "translate3d(0, 0, 0) scale(1)" },
          "50%": { transform: "translate3d(28px, -18px, 0) scale(1.06)" },
          "100%": { transform: "translate3d(0, 0, 0) scale(1)" },
        },
        sway: {
          "0%, 100%": { transform: "translateX(0) translateY(0)" },
          "50%": { transform: "translateX(18px) translateY(-10px)" },
        },
        shimmer: {
          "0%": { backgroundPosition: "0% 50%" },
          "100%": { backgroundPosition: "200% 50%" },
        },
        "spin-slow": {
          "0%": { transform: "rotate(0deg)" },
          "100%": { transform: "rotate(360deg)" },
        },
        pulseSoft: {
          "0%, 100%": { opacity: "0.45" },
          "50%": { opacity: "0.85" },
        },
      },
      animation: {
        "fade-up": "fade-up 0.55s ease-out both",
        "fade-in": "fade-in 0.9s ease-out both",
        float: "float 5s ease-in-out infinite",
        drift: "drift 12s ease-in-out infinite",
        sway: "sway 8s ease-in-out infinite",
        shimmer: "shimmer 8s linear infinite",
        "spin-slow": "spin-slow 26s linear infinite",
        "pulse-soft": "pulseSoft 2s ease-in-out infinite",
      },
    },
  },
  plugins: [],
};
