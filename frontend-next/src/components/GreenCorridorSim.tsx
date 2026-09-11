"use client";

import React, { useState } from "react";
import { Siren, Hospital, Clock, CheckCircle2, AlertTriangle, ArrowRight, ShieldCheck } from "lucide-react";
import { dispatchGreenCorridor } from "@/lib/api";
import { CorridorDispatchResponse } from "@/types";

interface GreenCorridorSimProps {
  isCorridorActive: boolean;
  setIsCorridorActive: (val: boolean) => void;
  onDispatchSuccess?: (res: CorridorDispatchResponse) => void;
}

export const GreenCorridorSim: React.FC<GreenCorridorSimProps> = ({
  isCorridorActive,
  setIsCorridorActive,
  onDispatchSuccess,
}) => {
  const [vehicleType, setVehicleType] = useState<string>("AMBULANCE_CRITICAL");
  const [selectedHospital, setSelectedHospital] = useState<string>("Manipal Institute of Cardiac Sciences");
  const [isDispatching, setIsDispatching] = useState<boolean>(false);
  const [dispatchResult, setDispatchResult] = useState<CorridorDispatchResponse | null>(null);

  const hospitals = [
    { name: "Manipal Institute of Cardiac Sciences", dist: "5.5 km", baseEta: "20.1 min", clearedEta: "4.5 min" },
    { name: "Narayana Institute of Neurosciences", dist: "8.2 km", baseEta: "31.4 min", clearedEta: "6.8 min" },
    { name: "Apollo Emergency Trauma Centre", dist: "4.1 km", baseEta: "16.8 min", clearedEta: "3.2 min" },
  ];

  const handleDispatch = async () => {
    setIsDispatching(true);
    try {
      const res = await dispatchGreenCorridor({
        vehicle_type: vehicleType,
        priority: 1,
        start_segment_id: "seg-mgroad-indiranagar",
        target_hospital: selectedHospital,
      });
      setDispatchResult(res);
      setIsCorridorActive(true);
      if (onDispatchSuccess) onDispatchSuccess(res);
    } catch (e) {
      console.error(e);
    } finally {
      setIsDispatching(false);
    }
  };

  const handleClear = () => {
    setIsCorridorActive(false);
    setDispatchResult(null);
  };

  return (
    <div className="p-5 rounded-2xl glass-panel border border-slate-700/80 shadow-2xl h-full flex flex-col justify-between">
      <div>
        {/* Header */}
        <div className="flex items-center justify-between pb-3 border-b border-slate-800">
          <div className="flex items-center gap-2.5">
            <div className="p-2 rounded-lg bg-red-500/20 text-red-400 border border-red-500/40">
              <Siren className="w-5 h-5 animate-pulse" />
            </div>
            <div>
              <h3 className="text-sm font-semibold text-white tracking-wide">
                Lifeline Green-Corridor AI Dispatch
              </h3>
              <p className="text-[11px] font-mono text-slate-400">
                Fluid wave signal preemption for emergency vehicles
              </p>
            </div>
          </div>

          <span
            className={`px-2.5 py-1 rounded-full text-[10px] font-mono font-semibold uppercase ${
              isCorridorActive
                ? "bg-emerald-500/20 text-emerald-300 border border-emerald-500/50 animate-pulse"
                : "bg-slate-800 text-slate-400 border border-slate-700"
            }`}
          >
            {isCorridorActive ? "WAVE CLEARANCE ACTIVE" : "STANDBY"}
          </span>
        </div>

        {/* Form Controls */}
        <div className="space-y-4 mt-4">
          <div>
            <label className="text-xs font-mono text-slate-300 mb-1.5 block">
              1. Emergency Priority Classification:
            </label>
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => setVehicleType("AMBULANCE_CRITICAL")}
                className={`py-2 px-3 rounded-lg text-xs font-mono text-left transition-all ${
                  vehicleType === "AMBULANCE_CRITICAL"
                    ? "bg-red-500/20 border border-red-500 text-red-200 font-semibold"
                    : "bg-slate-900/60 border border-slate-800 text-slate-400 hover:border-slate-700"
                }`}
              >
                🚨 Code Red: Cardiac / Trauma
              </button>

              <button
                type="button"
                onClick={() => setVehicleType("ORGAN_TRANSIT")}
                className={`py-2 px-3 rounded-lg text-xs font-mono text-left transition-all ${
                  vehicleType === "ORGAN_TRANSIT"
                    ? "bg-purple-500/20 border border-purple-500 text-purple-200 font-semibold"
                    : "bg-slate-900/60 border border-slate-800 text-slate-400 hover:border-slate-700"
                }`}
              >
                🫀 Live Organ Transplant
              </button>
            </div>
          </div>

          <div>
            <label className="text-xs font-mono text-slate-300 mb-1.5 block">
              2. Target Emergency Trauma Destination:
            </label>
            <div className="space-y-1.5">
              {hospitals.map((h) => (
                <div
                  key={h.name}
                  onClick={() => setSelectedHospital(h.name)}
                  className={`p-2.5 rounded-lg text-xs font-mono cursor-pointer transition-all flex items-center justify-between ${
                    selectedHospital === h.name
                      ? "bg-cyan-950/60 border border-cyan-500/60 text-cyan-200"
                      : "bg-slate-900/40 border border-slate-800/80 text-slate-400 hover:bg-slate-900"
                  }`}
                >
                  <div className="flex items-center gap-2">
                    <Hospital className="w-3.5 h-3.5 text-cyan-400" />
                    <span className="font-medium text-[11px] text-slate-200">{h.name}</span>
                  </div>
                  <div className="flex items-center gap-3 text-[10px]">
                    <span className="text-slate-500">{h.dist}</span>
                    <span className="text-red-400 line-through">{h.baseEta}</span>
                    <span className="text-emerald-400 font-semibold">{h.clearedEta}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Action / Result Area */}
      <div className="pt-4 border-t border-slate-800 mt-4">
        {isCorridorActive ? (
          <div className="space-y-3">
            <div className="p-3 rounded-xl bg-emerald-950/40 border border-emerald-500/40 text-emerald-200 text-xs font-mono flex items-center justify-between">
              <div className="flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                <div>
                  <div className="font-semibold text-emerald-300">Phase Wave Synchronized</div>
                  <div className="text-[10px] text-slate-400">
                    Saved: ~15.6 mins (+45% Trauma Survival Chance)
                  </div>
                </div>
              </div>
              <button
                onClick={handleClear}
                className="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 text-[10px] border border-slate-700"
              >
                Reset Grid
              </button>
            </div>
          </div>
        ) : (
          <button
            onClick={handleDispatch}
            disabled={isDispatching}
            className="w-full py-3 px-4 rounded-xl bg-gradient-to-r from-red-600 via-rose-600 to-amber-600 hover:from-red-500 hover:to-amber-500 text-white font-mono text-xs font-bold uppercase tracking-wider shadow-lg shadow-red-950/50 flex items-center justify-center gap-2 transition-all active:scale-[0.98]"
          >
            <Siren className="w-4 h-4 animate-bounce" />
            {isDispatching ? "Synchronizing Traffic Wave..." : "Engage Lifeline Green Corridor"}
          </button>
        )}
      </div>
    </div>
  );
};
