cd /d %~dp0
@echo off
powershell -Command "Start-Process -Wait PowerShell -Verb RunAs -ArgumentList 'Set-ExecutionPolicy RemoteSigned'"