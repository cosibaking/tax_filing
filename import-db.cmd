@echo off
setlocal
cd /d "%~dp0"
echo 已迁移到 scripts\import-mysql.ps1，正在调用...
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0scripts\import-mysql.ps1" -Import %*
endlocal
