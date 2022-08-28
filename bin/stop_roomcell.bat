timeout /T 1 /NOBREAK
taskkill /im roomcell_account.exe /t /f

timeout /T 1 /NOBREAK
taskkill /im roomcell_gate.exe /t /f

timeout /T 1 /NOBREAK
taskkill /im roomcell_data.exe /t /f

timeout /T 1 /NOBREAK
taskkill /im roomcell_hall.exe /t /f

timeout /T 1 /NOBREAK
taskkill /im roomcell_room.exe /t /f

timeout /T 1 /NOBREAK
taskkill /im roomcell_roommgr.exe /t /f
echo "[stop roomcell server finish!!!!]"
:: 停止
pause