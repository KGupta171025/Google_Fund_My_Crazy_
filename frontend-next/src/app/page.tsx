"use client";

import React, { useEffect, useState } from "react";
import {
  Sparkles,
  Zap,
  Activity,
  Github,
  Globe,
  Radio,
  Cpu,
  Shield,
  Layers,
  HeartPulse,
} from "lucide-react";
import { CityCanvas3D } from "@/components/CityCanvas3D";
import { TelemetryHUD } from "@/components/TelemetryHUD";
import { GreenCorridorSim } from "@/components/GreenCorridorSim";
import { GeminiSentinelChat } from "@/components/GeminiSentinelChat";
import { NeuralModelVisualizer } from "@/components/NeuralModelVisualizer";
import { ImpactEconomics } from "@/components/ImpactEconomics";
import { fetchCityState, MOCK_INITIAL_STATE } from "@/lib/api";
import { UrbanGridState } from "@/types";

export default function Home() {
  const [gridState, setGridState] = useState<UrbanGridState>(MOCK_INITIAL_STATE);
  const [isCorridorActive, setIsCorridorActive] = useState<boolean>(false);
  const [activeTab, setActiveTab] = useState<"twin" | "ai" | "impact">("twin");

  useEffect(() => {
    // Initial fetch
    fetchCityState().then((s) => setGridState(s));

    // Periodic poll / state update
    const interval = setInterval(async () => {
      const updated = await fetchCityState();
      setGridState((prev) => ({
        ...updated,
        total_kinetic_mwh: prev.total_kinetic_mwh + 0.00025,
        total_solar_mwh: prev.total_solar_mwh + 0.00015,
        total_co2_offset_tonnes: prev.total_co2_offset_tonnes + 0.0002,
        economic_savings_inr: prev.economic_savings_inr + 1.8,
      }));
    }, 1200);

    return () => clearInterval(interval);
  }, []);

  return (
    <main className="min-h-screen bg-[#080c14] text-slate-100 p-4 md:p-6 lg:p-8 space-y-6 max-w-[1600px] mx-auto">
      {/* 1. Header & FMC Initiative Banner */}
      <header className="flex flex-col lg:flex-row items-start lg:items-center justify-between gap-4 p-5 rounded-2xl glass-panel border border-slate-700/80 shadow-2xl">
        <div className="flex items-center gap-3.5">
          <div className="p-3 rounded-xl bg-gradient-to-br from-cyan-500/20 via-purple-500/20 to-pink-500/20 border border-cyan-500/40">
            <Zap className="w-7 h-7 text-cyan-400 animate-pulse" />
          </div>
          <div>
            <div className="flex items-center gap-2 flex-wrap">
              <h1 className="text-xl md:text-2xl font-black tracking-tight text-white glow-text-cyan font-mono">
                METRO-SYNAPSE (PRANA-GRID)
              </h1>
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-mono font-bold bg-gradient-to-r from-amber-500/20 to-rose-500/20 text-amber-300 border border-amber-500/40">
                ₹1 CRORE GOOGLE GEMINI MOONSHOT
              </span>
            </div>
            <p className="text-xs font-mono text-slate-400 mt-1">
              Autonomous Multimodal Urban Nervous System & Kinetic Bio-Infrastructure
            </p>
          </div>
        </div>

        {/* Tech Stack Badges */}
        <div className="flex items-center gap-2 flex-wrap">
          <span className="px-2.5 py-1 rounded-lg text-[10px] font-mono bg-cyan-950/80 text-cyan-300 border border-cyan-800">
            ⚡ Go Core (100k+ ev/s)
          </span>
          <span className="px-2.5 py-1 rounded-lg text-[10px] font-mono bg-purple-950/80 text-purple-300 border border-purple-800">
            🧠 Python AI (From Scratch)
          </span>
          <span className="px-2.5 py-1 rounded-lg text-[10px] font-mono bg-indigo-950/80 text-indigo-300 border border-indigo-800">
            🏎️ C++ SIMD Kernel
          </span>
          <span className="px-2.5 py-1 rounded-lg text-[10px] font-mono bg-emerald-950/80 text-emerald-300 border border-emerald-800">
            ✨ Gemini 1.5/2.0 Native
          </span>
          <a
            href="https://github.com/KGupta171025/Google_Fund_My_Crazy_"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-1.5 px-3 py-1 rounded-lg text-[11px] font-mono bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 transition-all ml-2"
          >
            <Github className="w-3.5 h-3.5" />
            <span>GitHub</span>
          </a>
        </div>
      </header>

      {/* 2. Real-time Telemetry HUD */}
      <section>
        <TelemetryHUD state={gridState} />
      </section>

      {/* 3. Navigation View Switcher */}
      <div className="flex items-center gap-2 border-b border-slate-800 pb-2">
        <button
          onClick={() => setActiveTab("twin")}
          className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-mono transition-all ${
            activeTab === "twin"
              ? "bg-cyan-500/20 text-cyan-300 border border-cyan-500/50 font-bold shadow-md"
              : "text-slate-400 hover:text-slate-200"
          }`}
        >
          <Radio className="w-4 h-4 text-cyan-400" />
          3D Cyber-City Twin & Emergency Dispatch
        </button>

        <button
          onClick={() => setActiveTab("ai")}
          className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-mono transition-all ${
            activeTab === "ai"
              ? "bg-purple-500/20 text-purple-300 border border-purple-500/50 font-bold shadow-md"
              : "text-slate-400 hover:text-slate-200"
          }`}
        >
          <Cpu className="w-4 h-4 text-purple-400" />
          Custom Scratch AI & Gemini Sentinel
        </button>

        <button
          onClick={() => setActiveTab("impact")}
          className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-mono transition-all ${
            activeTab === "impact"
              ? "bg-amber-500/20 text-amber-300 border border-amber-500/50 font-bold shadow-md"
              : "text-slate-400 hover:text-slate-200"
          }`}
        >
          <Layers className="w-4 h-4 text-amber-400" />
          FMC Criteria & Unit Economics
        </button>
      </div>

      {/* 4. Tab 1: 3D Twin + Emergency Dispatch */}
      {activeTab === "twin" && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <CityCanvas3D
              isGreenCorridorActive={isCorridorActive}
              activeVehiclesCount={gridState.active_vehicles}
              kineticMWh={gridState.total_kinetic_mwh}
            />
          </div>
          <div className="lg:col-span-1">
            <GreenCorridorSim
              isCorridorActive={isCorridorActive}
              setIsCorridorActive={setIsCorridorActive}
            />
          </div>
        </div>
      )}

      {/* 5. Tab 2: Custom Scratch Neural Network + Gemini Sentinel */}
      {activeTab === "ai" && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <NeuralModelVisualizer />
          <GeminiSentinelChat />
        </div>
      )}

      {/* 6. Tab 3: FMC Evaluation Criteria & Economics */}
      {activeTab === "impact" && (
        <div>
          <ImpactEconomics />
        </div>
      )}

      {/* 7. Bottom Unified Strip: 4 Key Arterials Condition Status */}
      <section className="p-4 rounded-xl glass-panel border border-slate-800">
        <div className="text-xs font-mono text-slate-400 mb-3 flex items-center justify-between">
          <span>Active Urban Smart Arterial Network Health:</span>
          <span className="text-cyan-400">● Spatial Geohash Index: Precision 7</span>
        </div>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
          {gridState.segments.map((s) => (
            <div
              key={s.id}
              className={`p-3 rounded-lg border text-xs font-mono transition-all ${
                s.is_green_corridor || isCorridorActive
                  ? "bg-emerald-950/40 border-emerald-500/50 text-emerald-200"
                  : "bg-slate-900/60 border-slate-800/80 text-slate-300"
              }`}
            >
              <div className="font-semibold text-slate-100 text-[11px] truncate mb-1">
                {s.name}
              </div>
              <div className="flex items-center justify-between text-[10px] text-slate-400">
                <span>Health: <strong className="text-emerald-400">{s.condition_score.toFixed(1)}%</strong></span>
                <span>Speed: <strong className="text-cyan-400">{s.avg_speed_kmh} km/h</strong></span>
              </div>
              <div className="flex items-center justify-between text-[10px] text-slate-400 mt-1">
                <span>Kinetic: <strong className="text-amber-400">{s.kinetic_mwh_today.toFixed(1)} MWh</strong></span>
                <span className={s.is_green_corridor || isCorridorActive ? "text-emerald-400 font-bold" : "text-slate-500"}>
                  {s.is_green_corridor || isCorridorActive ? "CORRIDOR ACTIVE" : "OPTIMAL"}
                </span>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* 8. Footer */}
      <footer className="pt-4 border-t border-slate-800 text-center text-xs font-mono text-slate-500 flex flex-col sm:flex-row items-center justify-between gap-2">
        <div>
          METRO-SYNAPSE • Google Gemini "Fund My Crazy" Student Innovation Challenge
        </div>
        <div className="flex items-center gap-4">
          <a
            href="https://fundmycrazy.com/"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-cyan-400 transition-colors"
          >
            fundmycrazy.com
          </a>
          <span>•</span>
          <a
            href="https://github.com/KGupta171025/Google_Fund_My_Crazy_"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-cyan-400 transition-colors"
          >
            github.com/KGupta171025/Google_Fund_My_Crazy_
          </a>
        </div>
      </footer>
    </main>
  );
}
