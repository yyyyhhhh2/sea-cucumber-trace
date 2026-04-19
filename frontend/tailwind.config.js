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
        float: {
          "0%, 100%": { transform: "translateY(0)" },
          "50%": { transform: "translateY(-6px)" },
        },
        pulseSoft: {
          "0%, 100%": { opacity: "0.45" },
          "50%": { opacity: "0.85" },
        },
      },
      animation: {
        "fade-up": "fade-up 0.55s ease-out both",
        float: "float 5s ease-in-out infinite",
        "pulse-soft": "pulseSoft 2s ease-in-out infinite",
      },
    },
  },
  plugins: [],
};
