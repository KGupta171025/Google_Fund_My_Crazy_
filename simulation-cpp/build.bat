@echo off
echo Building Metro-Synapse C++ Simulation Engine with g++ -O3 ...
if not exist bin mkdir bin
g++ -O3 -std=c++17 src/kinetic_kernel.cpp src/green_corridor.cpp src/main_sim.cpp -o bin/metro_sim.exe
if %ERRORLEVEL% EQU 0 (
    echo Compilation Successful! Running simulation benchmark:
    bin\metro_sim.exe
) else (
    echo Compilation Failed.
)
