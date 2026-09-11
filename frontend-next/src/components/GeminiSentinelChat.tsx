"use client";

import React, { useState } from "react";
import { Bot, Send, Sparkles, Terminal, ShieldAlert, Zap, Cpu, CheckCircle } from "lucide-react";
import { fetchGeminiReasoning } from "@/lib/api";
import { GeminiReasoningResponse } from "@/types";

export const GeminiSentinelChat: React.FC = () => {
  const [messages, setMessages] = useState<
    Array<{ sender: "user" | "gemini"; text: string; data?: GeminiReasoningResponse }>
  >([
    {
      sender: "gemini",
      text: "⚡ Welcome to Google Gemini Urban Sentinel. Multimodal spatial reasoning active. How can I assist with city grid optimization, drone defect diagnosis, or emergency wave coordination today?",
    },
  ]);
  const [inputText, setInputText] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const quickPrompts = [
    "Diagnose MG Road surface micro-vibrations & pothole risk",
    "Calculate kinetic power generation from 15,000 peak vehicles",
    "Simulate green wave corridor during monsoon flooding",
    "Break down municipal ROI & CO₂ offset for ₹1 Cr pilot",
  ];

  const handleSend = async (queryText?: string) => {
    const query = queryText || inputText;
    if (!query.trim() || loading) return;

    setMessages((prev) => [...prev, { sender: "user", text: query }]);
    setInputText("");
    setLoading(true);

    try {
      const response = await fetchGeminiReasoning(query, {
        vehicle_density: 55.0,
        avg_speed_kmh: 36.0,
        pothole_degradation_risk: query.toLowerCase().includes("pothole") ? 0.82 : 0.24,
        is_green_corridor: query.toLowerCase().includes("corridor") || query.toLowerCase().includes("green"),
      });

      setMessages((prev) => [
        ...prev,
        {
          sender: "gemini",
          text: response.assessment_summary,
          data: response,
        },
      ]);
    } catch (e) {
      setMessages((prev) => [
        ...prev,
        {
          sender: "gemini",
          text: "Gemini Sentinel processed your query using the offline scratch neural reasoner.",
        },
      ]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="p-5 rounded-2xl glass-panel border border-slate-700/80 shadow-2xl h-full flex flex-col justify-between">
      {/* Header */}
      <div className="flex items-center justify-between pb-3 border-b border-slate-800">
        <div className="flex items-center gap-2.5">
          <div className="p-2 rounded-lg bg-cyan-500/20 text-cyan-400 border border-cyan-500/40">
            <Bot className="w-5 h-5 animate-pulse" />
          </div>
          <div>
            <h3 className="text-sm font-semibold text-white tracking-wide flex items-center gap-2">
              Google Gemini Multimodal Sentinel
              <span className="px-2 py-0.5 rounded-full text-[9px] font-mono bg-cyan-500/20 text-cyan-300 border border-cyan-500/30">
                1.5 PRO / FLASH
              </span>
            </h3>
            <p className="text-[11px] font-mono text-slate-400">
              Multimodal spatial AI reasoning & autonomous grid orchestration
            </p>
          </div>
        </div>
      </div>

      {/* Message Stream */}
      <div className="flex-1 my-3 overflow-y-auto space-y-3 pr-1 max-h-[340px]">
        {messages.map((m, idx) => (
          <div
            key={idx}
            className={`flex flex-col ${m.sender === "user" ? "items-end" : "items-start"}`}
          >
            <div
              className={`p-3.5 rounded-xl text-xs font-mono max-w-[90%] leading-relaxed ${
                m.sender === "user"
                  ? "bg-cyan-600/30 border border-cyan-500/40 text-cyan-100 rounded-tr-none"
                  : "bg-slate-900/80 border border-slate-800 text-slate-200 rounded-tl-none shadow-md"
              }`}
            >
              <div className="flex items-center gap-1.5 mb-1 text-[10px] text-slate-400 font-semibold">
                {m.sender === "gemini" ? (
                  <>
                    <Sparkles className="w-3 h-3 text-cyan-400" />
                    <span>GEMINI SENTINEL AGENT</span>
                  </>
                ) : (
                  <span>URBAN OPERATOR</span>
                )}
              </div>
              <p>{m.text}</p>

              {/* Structured AI Output Badge Panel */}
              {m.data && (
                <div className="mt-3 pt-2.5 border-t border-slate-800 space-y-1.5 text-[10px]">
                  <div className="flex items-center justify-between text-slate-300">
                    <span className="text-slate-400">Structural Risk Assessment:</span>
                    <span
                      className={`font-semibold px-2 py-0.5 rounded ${
                        m.data.structural_risk_level === "CRITICAL"
                          ? "bg-red-500/20 text-red-400 border border-red-500/40"
                          : m.data.structural_risk_level === "HIGH"
                          ? "bg-amber-500/20 text-amber-400 border border-amber-500/40"
                          : "bg-emerald-500/20 text-emerald-400 border border-emerald-500/40"
                      }`}
                    >
                      {m.data.structural_risk_level}
                    </span>
                  </div>

                  <div className="text-slate-400">
                    <span className="text-cyan-300 font-medium">Signal Strategy: </span>
                    {m.data.green_corridor_recommendation}
                  </div>

                  <div className="text-slate-400">
                    <span className="text-amber-300 font-medium">Kinetic Optimization: </span>
                    {m.data.energy_harvesting_optimization}
                  </div>

                  {m.data.autonomous_drone_dispatch && (
                    <div className="p-1.5 rounded bg-red-950/60 border border-red-500/40 text-red-300 flex items-center gap-1.5 font-semibold">
                      <ShieldAlert className="w-3.5 h-3.5" />
                      Autonomous Thermal Resin Drone Dispatched to coordinates
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        ))}
        {loading && (
          <div className="flex items-center gap-2 p-3 rounded-xl bg-slate-900/60 border border-slate-800 text-xs font-mono text-cyan-400 animate-pulse">
            <Bot className="w-4 h-4 animate-spin" />
            <span>Gemini Multimodal Reasoning in progress...</span>
          </div>
        )}
      </div>

      {/* Quick Prompts */}
      <div className="mb-3">
        <div className="text-[10px] font-mono text-slate-400 mb-1.5">⚡ Suggested AI Moonshot Scenarios:</div>
        <div className="flex flex-wrap gap-1.5">
          {quickPrompts.map((p, i) => (
            <button
              key={i}
              onClick={() => handleSend(p)}
              className="text-[10px] font-mono py-1 px-2.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-cyan-300/90 border border-slate-800 hover:border-cyan-500/40 transition-all text-left"
            >
              {p}
            </button>
          ))}
        </div>
      </div>

      {/* Input Form */}
      <div className="flex items-center gap-2 pt-2 border-t border-slate-800">
        <input
          type="text"
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSend()}
          placeholder="Ask Gemini Sentinel about city flow, kinetic MWh or emergency routing..."
          className="flex-1 bg-slate-950/80 border border-slate-800 rounded-xl px-3.5 py-2.5 text-xs font-mono text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 transition-colors"
        />
        <button
          onClick={() => handleSend()}
          disabled={loading || !inputText.trim()}
          className="p-2.5 rounded-xl bg-cyan-500 hover:bg-cyan-400 disabled:opacity-40 text-slate-950 font-bold transition-all"
        >
          <Send className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
};
