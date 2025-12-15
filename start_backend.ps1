# Script de Pornire TruthWeave Antigravity
Write-Host "🚀 Inițializare Secvență de Lansare Antigravity..." -ForegroundColor Cyan

# 1. Pornire Infrastructură Docker
Write-Host "docker-compose up -d..." -ForegroundColor Yellow
docker-compose up -d

Write-Host "⏳ Așteptăm 15 secunde pentru inițializarea bazelor de date..." -ForegroundColor DarkGray
Start-Sleep -Seconds 15

# 2. Pornire Worker (Background)
Write-Host "👷 Pornire Worker (Fereastră Nouă)..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-Command", "& { go run ./cmd/worker }"

# 3. Pornire API (Foreground)
Write-Host "🌐 Pornire API Server..." -ForegroundColor Green
go run ./cmd/api

# Notă: Dacă API-ul eșuează, worker-ul va continua să ruleze în fundal.
