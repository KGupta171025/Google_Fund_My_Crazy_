import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "METRO-SYNAPSE (PRANA-GRID) | Google Gemini 'Fund My Crazy'",
  description:
    "Autonomous Multimodal Urban Nervous System & Kinetic Bio-Infrastructure. A ₹1 Crore Google Gemini Initiative Moonshot.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased min-h-screen bg-[#080c14] text-slate-100 selection:bg-cyan-500/30 selection:text-cyan-200">
        {children}
      </body>
    </html>
  );
}
