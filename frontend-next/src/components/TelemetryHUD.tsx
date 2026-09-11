"use client";

import React from "react";
import { Zap, Sun, ShieldCheck, HeartPulse, DollarSign, CloudRain, Cpu } from "lucide-react";
import { UrbanGridState } from "@/types";

interface TelemetryHUDProps {
  state: UrbanGridState;
}

export const TelemetryHUD: React.FC<TelemetryHUDProps> = ({ state }) => {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3.5 w-full">
      {/* 1. Kinetic Power MWh */}
      <div className="p-4 rounded-xl glass-panel border border-cyan-500/30 hover:border-cyan-400 transition-all shadow-lg group">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[11px] font-mono uppercase tracking-wider text-slate-400">
            Kinetic Harvest
          </span>
          <Zap className="w-4 h-4 text-cyan-400 group-hover:scale-110 transition-transform" />
        </div>
        <div className="text-2xl font-bold font-mono text-cyan-300">
          {state.total_kinetic_mwh.toFixed(2)}
          <span className="text-xs font-normal text-cyan-400/80 ml-1">MWh</span>
        </div>
        <div className="text-[10px] font-mono text-emerald-400 mt-1 flex items-center gap-1">
          <span>▲ +18.4%</span>
          <span className="text-slate-500">rolling friction</span>
        </div>
      </div>

      {/* 2. Solar PV MWh */}
      <div className="p-4 rounded-xl glass-panel border border-amber-500/30 hover:border-amber-400 transition-all shadow-lg group">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[11px] font-mono uppercase tracking-wider text-slate-400">
            Solar Surface
          </span>
          <Sun className="w-4 h-4 text-amber-400 group-hover:scale-110 transition-transform" />
        </div>
        <div className="text-2xl font-bold font-mono text-amber-300">
          {state.total_solar_mwh.toFixed(2)}
          <span className="text-xs font-normal text-amber-400/80 ml-1">MWh</span>
        </div>
        <div className="text-[10px] font-mono text-amber-400 mt-1 flex items-center gap-1">
          <span>780 W/m²</span>
          <span className="text-slate-500">irradiance</span>
        </div>
      </div>

      {/* 3. Total Clean Power */}
      <div className="p-4 rounded-xl glass-panel border border-purple-500/30 hover:border-purple-400 transition-all shadow-lg group">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[11px] font-mono uppercase tracking-wider text-slate-400">
            Instantaneous Grid
          </span>
          <Cpu className="w-4 h-4 text-purple-400 group-hover:scale-110 transition-transform" />
        </div>
        <div className="text-2xl font-bold font-mono text-purple-300">
          {state.total_power_generated_mw.toFixed(2)}
          <span className="text-xs font-normal text-purple-400/80 ml-1">MW</span>
        </div>
        <div className="text-[10px] font-mono text-purple-400 mt-1 flex items-center gap-1">
          <span>Decentralized</span>
          <span className="text-slate-500">micro-grid feed</span>
        </div>
      </div>

      {/* 4. Golden Hour Response Time */}
      <div className="p-4 rounded-xl glass-panel border border-emerald-500/30 hover:border-emerald-400 transition-all shadow-lg group">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[11px] font-mono uppercase tracking-wider text-slate-400">
            Ambulance Lifeline
          </span>
          <HeartPulse className="w-4 h-4 text-emerald-400 group-hover:scale-110 transition-transform" />
        </div>
        <div className="text-2xl font-bold font-mono text-emerald-300">
          {state.avg_response_time_min.toFixed(1)}
          <span className="text-xs font-normal text-emerald-400/80 ml-1">min</span>
        </div>
        <div className="text-[10px] font-mono text-emerald-400 mt-1 flex items-center gap-1">
          <span className="line-through text-slate-500">28.5m</span>
          <span>▼ -85% delay</span>
        </div>
      </div>

      {/* 5. CO2 Offset */}
      <div className="p-4 rounded-xl glass-panel border border-teal-500/30 hover:border-teal-400 transition-all shadow-lg group">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[11px] font-mono uppercase tracking-wider text-slate-400">
            CO₂ Avoided
          </span>
          <CloudRain className="w-4 h-4 text-teal-400 group-hover:scale-110 transition-transform" />
        </div>
        <div className="text-2xl font-bold font-mono text-teal-300">
          {state.total_co2_offset_tonnes.toFixed(1)}
          <span className="text-xs font-normal text-teal-400/80 ml-1">Tonnes</span>
        </div>
        <div className="text-[10px] font-mono text-teal-400 mt-1 flex items-center gap-1">
          <span>Net Carbon</span>
          <span className="text-slate-500">Negative Road</span>
        </div>
      </div>

      {/* 6. Economic Value Created */}
      <div className="p-4 rounded-xl glass-panel border border-rose-500/30 hover:border-rose-400 transition-all shadow-lg group">
        <div className="flex items-center justify-between mb-2">
          <span className="text-[11px] font-mono uppercase tracking-wider text-slate-400">
            Value Generated
          </span>
          <DollarSign className="w-4 h-4 text-rose-400 group-hover:scale-110 transition-transform" />
        </div>
        <div className="text-2xl font-bold font-mono text-rose-300">
          ₹{(state.economic_savings_inr / 100000).toFixed(2)}
          <span className="text-xs font-normal text-rose-400/80 ml-1">Lakh</span>
        </div>
        <div className="text-[10px] font-mono text-rose-400 mt-1 flex items-center gap-1">
          <span>Direct ROI</span>
          <span className="text-slate-500">& power feed</span>
        </div>
      </div>
    </div>
  );
};
