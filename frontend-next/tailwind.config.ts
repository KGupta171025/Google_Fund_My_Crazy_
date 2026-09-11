import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "#080c14",
        foreground: "#f1f5f9",
        cyber: {
          cyan: "#00f0ff",
          green: "#00ff88",
          amber: "#ffaa00",
          rose: "#ff2a6d",
          purple: "#7b2cbf",
          navy: "#0a1128",
          card: "#0d1829",
          border: "#1e293b",
        },
      },
      fontFamily: {
        mono: ["ui-monospace", "SFMono-Regular", "Menlo", "Monaco", "Consolas", "monospace"],
      },
      animation: {
        "pulse-fast": "pulse 1.2s cubic-bezier(0.4, 0, 0.6, 1) infinite",
        "glow-cyan": "glowCyan 2s ease-in-out infinite alternate",
      },
      keyframes: {
        glowCyan: {
          "0%": { boxShadow: "0 0 10px rgba(0, 240, 255, 0.3)" },
          "100%": { boxShadow: "0 0 25px rgba(0, 240, 255, 0.8)" },
        },
      },
    },
  },
  plugins: [],
};
export default config;
