"use client";

import React, { useEffect, useRef, useState } from "react";
import * as THREE from "three";
import { Sparkles, Activity, ShieldAlert, Zap, Layers, RefreshCw } from "lucide-react";

interface CityCanvas3DProps {
  isGreenCorridorActive: boolean;
  activeVehiclesCount: number;
  kineticMWh: number;
}

export const CityCanvas3D: React.FC<CityCanvas3DProps> = ({
  isGreenCorridorActive,
  activeVehiclesCount,
  kineticMWh,
}) => {
  const mountRef = useRef<HTMLDivElement>(null);
  const [viewMode, setViewMode] = useState<"kinetic" | "corridor" | "heatmap">("kinetic");
  const [cameraPreset, setCameraPreset] = useState<"aerial" | "arterial" | "birds_eye">("aerial");

  useEffect(() => {
    const container = mountRef.current;
    if (!container) return;

    // 1. Scene Setup
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x080c14);
    scene.fog = new THREE.FogExp2(0x080c14, 0.015);

    // 2. Camera Setup
    const camera = new THREE.PerspectiveCamera(
      55,
      container.clientWidth / container.clientHeight,
      0.1,
      1000
    );
    camera.position.set(45, 38, 55);
    camera.lookAt(0, 0, 0);

    // 3. Renderer Setup
    const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
    renderer.setSize(container.clientWidth, container.clientHeight);
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    renderer.toneMapping = THREE.ACESFilmicToneMapping;
    renderer.toneMappingExposure = 1.2;
    container.appendChild(renderer.domElement);

    // 4. Lighting
    const ambientLight = new THREE.AmbientLight(0x1a2639, 1.8);
    scene.add(ambientLight);

    const dirLight = new THREE.DirectionalLight(0x00f0ff, 2.2);
    dirLight.position.set(30, 50, 20);
    scene.add(dirLight);

    const greenLight = new THREE.PointLight(0x00ff88, 3.0, 60);
    greenLight.position.set(0, 5, 0);
    scene.add(greenLight);

    // 5. Ground Plane & Kinetic Grid
    const gridHelper = new THREE.GridHelper(120, 60, 0x00f0ff, 0x112233);
    gridHelper.position.y = -0.05;
    scene.add(gridHelper);

    // 6. Procedural Futuristic City Buildings
    const buildingGroup = new THREE.Group();
    const buildingGeom = new THREE.BoxGeometry(1, 1, 1);
    const buildingMaterials = [
      new THREE.MeshStandardMaterial({ color: 0x0d1b2a, roughness: 0.3, metalness: 0.8 }),
      new THREE.MeshStandardMaterial({ color: 0x1b263b, roughness: 0.4, metalness: 0.7 }),
      new THREE.MeshStandardMaterial({ color: 0x0f172a, roughness: 0.2, metalness: 0.9 }),
    ];

    const gridSize = 7;
    const spacing = 12;
    for (let x = -gridSize; x <= gridSize; x++) {
      for (let z = -gridSize; z <= gridSize; z++) {
        // Leave main cross corridors open for roads
        if (Math.abs(x) <= 1 || Math.abs(z) <= 1 || (x % 3 === 0 && z % 3 === 0)) {
          continue;
        }

        const height = 4 + Math.random() * 22;
        const mat = buildingMaterials[Math.floor(Math.random() * buildingMaterials.length)];
        const building = new THREE.Mesh(buildingGeom, mat);
        building.scale.set(3.8 + Math.random() * 2, height, 3.8 + Math.random() * 2);
        building.position.set(x * spacing + (Math.random() - 0.5) * 2, height / 2, z * spacing + (Math.random() - 0.5) * 2);
        buildingGroup.add(building);

        // Holographic beacon on top of tallest towers
        if (height > 18) {
          const beaconGeom = new THREE.SphereGeometry(0.3, 8, 8);
          const beaconMat = new THREE.MeshBasicMaterial({ color: 0x00f0ff });
          const beacon = new THREE.Mesh(beaconGeom, beaconMat);
          beacon.position.set(building.position.x, height + 0.4, building.position.z);
          buildingGroup.add(beacon);
        }
      }
    }
    scene.add(buildingGroup);

    // 7. Kinetic Road Grid Lines
    const roadGroup = new THREE.Group();
    const roadMatStandard = new THREE.MeshBasicMaterial({ color: 0x00f0ff });
    const roadMatCorridor = new THREE.MeshBasicMaterial({ color: 0x00ff88 });

    // Main Arterial Cross Roads
    const roadPlaneH = new THREE.Mesh(new THREE.PlaneGeometry(120, 6), new THREE.MeshBasicMaterial({ color: 0x07111e }));
    roadPlaneH.rotation.x = -Math.PI / 2;
    roadPlaneH.position.y = 0.01;
    roadGroup.add(roadPlaneH);

    const roadPlaneV = new THREE.Mesh(new THREE.PlaneGeometry(6, 120), new THREE.MeshBasicMaterial({ color: 0x07111e }));
    roadPlaneV.rotation.x = -Math.PI / 2;
    roadPlaneV.position.y = 0.01;
    roadGroup.add(roadPlaneV);

    scene.add(roadGroup);

    // 8. Vehicle Swarms (Kinetic energy particle flows)
    const vehicleCount = 180;
    const vehicleGeom = new THREE.BoxGeometry(0.8, 0.4, 1.4);
    const vehicleMat = new THREE.MeshBasicMaterial({ color: 0x00f0ff });
    const vehicles: { mesh: THREE.Mesh; axis: "x" | "z"; speed: number; dir: number }[] = [];

    for (let i = 0; i < vehicleCount; i++) {
      const axis = Math.random() > 0.5 ? "x" : "z";
      const dir = Math.random() > 0.5 ? 1 : -1;
      const mesh = new THREE.Mesh(vehicleGeom, vehicleMat.clone());

      if (axis === "x") {
        mesh.position.set((Math.random() - 0.5) * 110, 0.25, (dir > 0 ? 1.5 : -1.5));
      } else {
        mesh.position.set((dir > 0 ? 1.5 : -1.5), 0.25, (Math.random() - 0.5) * 110);
        mesh.rotation.y = Math.PI / 2;
      }

      scene.add(mesh);
      vehicles.push({
        mesh,
        axis,
        speed: 0.3 + Math.random() * 0.4,
        dir,
      });
    }

    // 9. Emergency Ambulance Unit
    const ambGeom = new THREE.BoxGeometry(1.2, 0.8, 2.2);
    const ambMat = new THREE.MeshStandardMaterial({
      color: 0xff0044,
      emissive: 0xff0044,
      emissiveIntensity: 1.5,
    });
    const ambulance = new THREE.Mesh(ambGeom, ambMat);
    ambulance.position.set(0, 0.4, -50);
    scene.add(ambulance);

    // Ambulance Siren Point Light
    const sirenLight = new THREE.PointLight(0xff0044, 4.0, 20);
    ambulance.add(sirenLight);

    // 10. Interactive Mouse / Orbit Drag Handlers
    let isDragging = false;
    let prevMouseX = 0;
    let prevMouseY = 0;
    let targetRotationY = 0;
    let targetRotationX = 0.3;

    const onMouseDown = (e: MouseEvent) => {
      isDragging = true;
      prevMouseX = e.clientX;
      prevMouseY = e.clientY;
    };

    const onMouseMove = (e: MouseEvent) => {
      if (!isDragging) return;
      const deltaX = e.clientX - prevMouseX;
      const deltaY = e.clientY - prevMouseY;
      prevMouseX = e.clientX;
      prevMouseY = e.clientY;

      targetRotationY += deltaX * 0.005;
      targetRotationX = Math.max(0.1, Math.min(1.2, targetRotationX + deltaY * 0.005));
    };

    const onMouseUp = () => {
      isDragging = false;
    };

    const onWheel = (e: WheelEvent) => {
      e.preventDefault();
      const zoomFactor = e.deltaY * 0.05;
      const currentDist = camera.position.length();
      const newDist = Math.max(25, Math.min(130, currentDist + zoomFactor));
      camera.position.setLength(newDist);
    };

    const domElem = renderer.domElement;
    domElem.addEventListener("mousedown", onMouseDown);
    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
    domElem.addEventListener("wheel", onWheel, { passive: false });

    // 11. Animation Loop
    let animationFrameId: number;
    let clock = new THREE.Clock();

    const animate = () => {
      animationFrameId = requestAnimationFrame(animate);
      const elapsedTime = clock.getElapsedTime();

      // Camera Orbit Interpolation
      const radius = camera.position.length();
      camera.position.x = radius * Math.sin(targetRotationY) * Math.cos(targetRotationX);
      camera.position.y = radius * Math.sin(targetRotationX);
      camera.position.z = radius * Math.cos(targetRotationY) * Math.cos(targetRotationX);
      camera.lookAt(0, 0, 0);

      // Animate Vehicles
      vehicles.forEach((v) => {
        if (v.axis === "x") {
          v.mesh.position.x += v.speed * v.dir;
          if (v.mesh.position.x > 55) v.mesh.position.x = -55;
          if (v.mesh.position.x < -55) v.mesh.position.x = 55;
        } else {
          v.mesh.position.z += v.speed * v.dir;
          if (v.mesh.position.z > 55) v.mesh.position.z = -55;
          if (v.mesh.position.z < -55) v.mesh.position.z = 55;
        }

        // Color update based on Green Corridor
        const isNearAmbulance = Math.abs(v.mesh.position.z - ambulance.position.z) < 15;
        const mat = v.mesh.material as THREE.MeshBasicMaterial;

        if (isGreenCorridorActive && isNearAmbulance) {
          mat.color.setHex(0x00ff88); // Cleared wave
        } else if (viewMode === "heatmap") {
          mat.color.setHex(0xffaa00);
        } else {
          mat.color.setHex(0x00f0ff);
        }
      });

      // Animate Ambulance
      if (isGreenCorridorActive) {
        ambulance.visible = true;
        ambulance.position.z += 0.65; // High speed along corridor
        if (ambulance.position.z > 55) {
          ambulance.position.z = -55;
        }
        sirenLight.intensity = 2.0 + Math.sin(elapsedTime * 15) * 2.0;
        greenLight.intensity = 3.5 + Math.sin(elapsedTime * 8) * 1.5;
        greenLight.position.set(ambulance.position.x, 3, ambulance.position.z);
      } else {
        ambulance.visible = false;
        greenLight.intensity = 1.0;
      }

      renderer.render(scene, camera);
    };

    animate();

    // 12. Responsive Resize
    const handleResize = () => {
      if (!container) return;
      camera.aspect = container.clientWidth / container.clientHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(container.clientWidth, container.clientHeight);
    };
    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
      domElem.removeEventListener("mousedown", onMouseDown);
      window.removeEventListener("mousemove", onMouseMove);
      window.removeEventListener("mouseup", onMouseUp);
      domElem.removeEventListener("wheel", onWheel);
      cancelAnimationFrame(animationFrameId);
      if (container && renderer.domElement) {
        container.removeChild(renderer.domElement);
      }
      renderer.dispose();
    };
  }, [isGreenCorridorActive, viewMode]);

  return (
    <div className="relative w-full h-[520px] rounded-2xl overflow-hidden glass-panel border border-slate-700/60 shadow-2xl">
      {/* 3D WebGL Canvas */}
      <div ref={mountRef} className="w-full h-full cursor-grab active:cursor-grabbing" />

      {/* Top Overlay Badge */}
      <div className="absolute top-4 left-4 flex items-center gap-3 pointer-events-none">
        <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-900/90 border border-cyan-500/40 backdrop-blur-md">
          <span className="w-2.5 h-2.5 rounded-full bg-cyan-400 animate-ping" />
          <span className="text-xs font-mono font-semibold tracking-wider text-cyan-300">
            3D DIGITAL TWIN • LIVE TELEMETRY
          </span>
        </div>
        {isGreenCorridorActive && (
          <div className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-emerald-950/90 border border-emerald-500/60 backdrop-blur-md animate-pulse">
            <span className="w-2.5 h-2.5 rounded-full bg-emerald-400" />
            <span className="text-xs font-mono font-semibold tracking-wider text-emerald-300">
              🚑 GREEN CORRIDOR WAVE ACTIVE
            </span>
          </div>
        )}
      </div>

      {/* Mode Switcher Pill */}
      <div className="absolute bottom-4 left-4 flex items-center gap-1.5 p-1 rounded-xl bg-slate-900/90 border border-slate-700/80 backdrop-blur-md">
        <button
          onClick={() => setViewMode("kinetic")}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-mono transition-all ${
            viewMode === "kinetic"
              ? "bg-cyan-500/20 text-cyan-300 border border-cyan-500/40 font-semibold"
              : "text-slate-400 hover:text-slate-200"
          }`}
        >
          <Zap className="w-3.5 h-3.5 text-cyan-400" />
          Kinetic Energy
        </button>

        <button
          onClick={() => setViewMode("corridor")}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-mono transition-all ${
            viewMode === "corridor"
              ? "bg-emerald-500/20 text-emerald-300 border border-emerald-500/40 font-semibold"
              : "text-slate-400 hover:text-slate-200"
          }`}
        >
          <Activity className="w-3.5 h-3.5 text-emerald-400" />
          Lifeline Corridor
        </button>

        <button
          onClick={() => setViewMode("heatmap")}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-mono transition-all ${
            viewMode === "heatmap"
              ? "bg-amber-500/20 text-amber-300 border border-amber-500/40 font-semibold"
              : "text-slate-400 hover:text-slate-200"
          }`}
        >
          <Layers className="w-3.5 h-3.5 text-amber-400" />
          Heatmap AQI
        </button>
      </div>

      {/* Floating Canvas Controls */}
      <div className="absolute bottom-4 right-4 text-[11px] font-mono text-slate-400 bg-slate-900/80 px-3 py-1.5 rounded-lg border border-slate-800">
        Left Click + Drag to Orbit • Scroll to Zoom
      </div>
    </div>
  );
};
