# Force cleanup of Android build artifacts and caches

Write-Host "🧹 Cleaning Android Project..." -ForegroundColor Yellow

# Kill running Gradle daemons (just in case)
Stop-Process -Name "java" -ErrorAction SilentlyContinue

# Define paths
$androidDir = "c:\Users\Andrei\.gemini\antigravity\scratch\IMPARTIAL\android"

# Remove .gradle directory
if (Test-Path "$androidDir\.gradle") {
    Remove-Item "$androidDir\.gradle" -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "Removed .gradle cache"
}

# Remove .idea directory (Project settings)
if (Test-Path "$androidDir\.idea") {
    Remove-Item "$androidDir\.idea" -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "Removed .idea settings"
}

# Remove build directory
if (Test-Path "$androidDir\build") {
    Remove-Item "$androidDir\build" -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "Removed root build directory"
}

# Remove app/build directory
if (Test-Path "$androidDir\app\build") {
    Remove-Item "$androidDir\app\build" -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "Removed app build directory"
}

Write-Host "✨ Clean complete! Now RESTART Android Studio." -ForegroundColor Green
