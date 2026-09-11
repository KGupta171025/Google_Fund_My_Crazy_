"use client";

import React, { useState, useEffect } from "react";
import { Cpu, Activity, Zap, ShieldAlert, BarChart3, Sliders } from "lucide-react";
import { fetchAISegmentPrediction } from "@/lib/api";
import { AIPredictionResponse } from "@/types";

export const NeuralModelVisualizer: React.FC = () => {
  const [density, setDensity] = useState<number>(65);
  const [speed, setSpeed] = useState<number>(32);
  const [pressure, setPressure] = useState<number>(180);
  const [temperature, setTemperature] = useState<number>(34);
  const [rainfall, setRainfall] = useState<number>(15);

  const [prediction, setPrediction] = useState<AIPredictionResponse>({
    congestion_probability_30m: 0.742,
    pothole_degradation_risk: 0.618,
    optimal_green_corridor_speed_kmh: 58.0,
    structural_status: "WARNING",
  });

  useEffect(() => {
    const updatePrediction = async () => {
      const pred = await fetchAISegmentPrediction({
        vehicle_density: density,
        avg_speed_kmh: speed,
        pressure_kpa: pressure,
        surface_temp: temperature,
        rainfall_mm: rainfall,
        hour_of_day: 17.5, // 5:30 PM peak
      });
      setPrediction(pred);
    };

    updatePrediction();
  }, [density, speed, pressure, temperature, rainfall]);

  return (
    <div className="p-5 rounded-2xl glass-panel border border-slate-700/80 shadow-2xl space-y-5">
      {/* Header */}
      <div className="flex items-center justify-between pb-3 border-b border-slate-800">
        <div className="flex items-center gap-2.5">
          <div className="p-2 rounded-lg bg-purple-500/20 text-purple-400 border border-purple-500/40">
            <Cpu className="w-5 h-5 animate-pulse" />
          </div>
          <div>
            <h3 className="text-sm font-semibold text-white tracking-wide flex items-center gap-2">
              Custom Scratch Neural Network Engine (Python)
              <span className="px-2 py-0.5 rounded-full text-[9px] font-mono bg-purple-500/20 text-purple-300 border border-purple-500/30">
                100% PURE MATH • ZERO BLACKBOX
              </span>
            </h3>
            <p className="text-[11px] font-mono text-slate-400">
              Forward/Backward Autograd • Adam Optimizer • Physical Spatiotemporal Tensor Inference
            </p>
          </div>
        </div>

        <div className="text-[10px] font-mono text-emerald-400 bg-emerald-950/50 px-2.5 py-1 rounded border border-emerald-500/30">
          ● INFERENCE: &lt; 0.4ms
        </div>
      </div>

      {/* Grid: Inputs vs Network Topology & Outputs */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
        {/* Sliders Input Panel */}
        <div className="p-4 rounded-xl bg-slate-900/60 border border-slate-800 space-y-3.5">
          <div className="flex items-center gap-1.5 text-xs font-mono text-slate-300 mb-1">
            <Sliders className="w-3.5 h-3.5 text-purple-400" />
            <span>Spatiotemporal Physical Inputs:</span>
          </div>

          <div>
            <div className="flex justify-between text-xs font-mono text-slate-400 mb-1">
              <span>Vehicle Density:</span>
              <span className="text-cyan-300 font-bold">{density} veh/km</span>
            </div>
            <input
              type="range"
              min="5"
              max="100"
              value={density}
              onChange={(e) => setDensity(Number(e.target.value))}
              className="w-full h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-cyan-400"
            />
          </div>

          <div>
            <div className="flex justify-between text-xs font-mono text-slate-400 mb-1">
              <span>Average Speed:</span>
              <span className="text-cyan-300 font-bold">{speed} km/h</span>
            </div>
            <input
              type="range"
              min="5"
              max="120"
              value={speed}
              onChange={(e) => setSpeed(Number(e.target.value))}
              className="w-full h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-cyan-400"
            />
          </div>

          <div>
            <div className="flex justify-between text-xs font-mono text-slate-400 mb-1">
              <span>Piezo Dynamic Pressure:</span>
              <span className="text-purple-300 font-bold">{pressure} kPa</span>
            </div>
            <input
              type="range"
              min="30"
              max="300"
              value={pressure}
              onChange={(e) => setPressure(Number(e.target.value))}
              className="w-full h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-purple-400"
            />
          </div>

          <div>
            <div className="flex justify-between text-xs font-mono text-slate-400 mb-1">
              <span>Rainfall / Water Ingress:</span>
              <span className="text-teal-300 font-bold">{rainfall} mm/h</span>
            </div>
            <input
              type="range"
              min="0"
              max="80"
              value={rainfall}
              onChange={(e) => setRainfall(Number(e.target.value))}
              className="w-full h-1.5 bg-slate-800 rounded-lg appearance-none cursor-pointer accent-teal-400"
            />
          </div>
        </div>

        {/* Neural Network Architecture & Predictions */}
        <div className="p-4 rounded-xl bg-slate-900/60 border border-slate-800 flex flex-col justify-between space-y-4">
          <div>
            <div className="text-xs font-mono text-slate-300 mb-2 flex items-center justify-between">
              <span>Topology: 6 → Dense(32) → ReLU → Dense(16) → 3</span>
              <span className="text-[10px] text-purple-400 font-bold">He Normal</span>
            </div>

            {/* Visual Node Layer Representation */}
            <div className="grid grid-cols-4 gap-2 text-center text-[9px] font-mono py-2 bg-slate-950/70 rounded-lg border border-slate-800/80 mb-3">
              <div className="text-cyan-400">INPUT (6)</div>
              <div className="text-purple-400">DENSE (32)</div>
              <div className="text-indigo-400">DENSE (16)</div>
              <div className="text-emerald-400">OUTPUT (3)</div>
            </div>

            {/* Neural Predictions */}
            <div className="space-y-2.5 font-mono text-xs">
              <div className="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
                <span className="text-slate-400">30-Min Congestion Probability:</span>
                <span className="font-bold text-amber-400">
                  {(prediction.congestion_probability_30m * 100).toFixed(1)}%
                </span>
              </div>

              <div className="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
                <span className="text-slate-400">Pothole Structural Fatigue Risk:</span>
                <span
                  className={`font-bold ${
                    prediction.pothole_degradation_risk > 0.6
                      ? "text-red-400"
                      : "text-emerald-400"
                  }`}
                >
                  {(prediction.pothole_degradation_risk * 100).toFixed(1)}% (
                  {prediction.structural_status})
                </span>
              </div>

              <div className="p-2.5 rounded-lg bg-slate-950 border border-slate-800 flex items-center justify-between">
                <span className="text-slate-400">Optimal Green Wave Velocity:</span>
                <span className="font-bold text-emerald-400">
                  {prediction.optimal_green_corridor_speed_kmh} km/h
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
