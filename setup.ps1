# TruthWeave One-Click Setup (Windows/PowerShell)

Write-Host "🚀 Starting TruthWeave Infrastructure Setup..." -ForegroundColor Cyan

# 1. Check Docker
if (-not (Get-Command "docker" -ErrorAction SilentlyContinue)) {
    Write-Error "Docker is not installed or not in PATH. Please install Docker Desktop and restart your terminal."
    exit 1
}

# 2. Start Services
Write-Host "🐳 Starting Docker Containers..." -ForegroundColor Yellow
docker-compose up -d
if ($LASTEXITCODE -ne 0) {
    Write-Host "docker-compose failed. Trying 'docker compose'..." -ForegroundColor DarkYellow
    docker compose up -d
}

# 3. Wait for Database
Write-Host "⏳ Waiting for Database to be ready (10s)..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

# 4. Run Migrations (via Container)
Write-Host "📦 Applying Database Migrations..." -ForegroundColor Yellow

# Copy schema to container momentarily to execute it
docker cp manual_schema.sql impartial-postgres-1:/tmp/schema.sql
if ($LASTEXITCODE -ne 0) {
    # Try guessing a different container name if default fails
    docker cp impartial_postgres_1:/tmp/schema.sql
}

# Execute psql inside container
docker exec -i impartial-postgres-1 psql -U user -d truthweave -f /tmp/schema.sql
if ($LASTEXITCODE -ne 0) {
    Write-Host "Attempting backup container name..."
    docker exec -i impartial-postgres-1 psql -U user -d truthweave -f /tmp/schema.sql
}

Write-Host "✅ Infrastructure Ready!" -ForegroundColor Green
Write-Host "---------------------------------------------------"
Write-Host "Next Steps:"
Write-Host "1. Open Terminal 1: go run cmd/worker/main.go"
Write-Host "2. Open Terminal 2: go run cmd/api/main.go"
Write-Host "3. Run Android App in Emulator (ensure NetworkModule has 10.0.2.2)"
