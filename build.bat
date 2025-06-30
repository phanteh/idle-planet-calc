@echo off
echo Building Windows Package
fyne package -os windows
echo Building Android Package
fyne package -os android
echo Done
