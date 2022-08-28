timeout /T 1 /NOBREAK
taskkill /im myth_gateway.exe /t /f
timeout /T 1 /NOBREAK
taskkill /im myth_game.exe /t /f
timeout /T 1 /NOBREAK
taskkill /im myth_social.exe /t /f
timeout /T 1 /NOBREAK
taskkill /im myth_rank.exe /t /f
timeout /T 1 /NOBREAK
taskkill /im myth_dblog.exe /t /f
timeout /T 1 /NOBREAK
taskkill /im myth_dbserver.exe /t /f
timeout /T 1 /NOBREAK
taskkill /im myth_archive.exe /t /f
timeout /T 1 /NOBREAK
taskkill /im myth_zone.exe /t /f
timeout /T 1 /NOBREAK
echo "服务器停止"
:: 停止
pause