set APP_CONFIG_ROOT=E:/tyh_work_card/server_project/app_config
start /b .\bin\myth_zone.exe --configPath=%APP_CONFIG_ROOT%/zone > .\\logs\\zone.log
timeout /T 1 /NOBREAK
start /b .\bin\myth_archive.exe --configPath=%APP_CONFIG_ROOT%//archive > .\\logs\\archive.log
timeout /T 1 /NOBREAK
start /b .\bin\myth_dbserver.exe --configPath=%APP_CONFIG_ROOT%//dbserver > .\\logs\\dbserver.log
timeout /T 1 /NOBREAK
start /b .\bin\myth_dblog.exe --configPath=%APP_CONFIG_ROOT%//dblog > .\\logs\\dblog.log
timeout /T 1 /NOBREAK
start /b .\bin\myth_social.exe --configPath=%APP_CONFIG_ROOT%//social > .\\logs\\social.log
timeout /T 1 /NOBREAK
start /b .\bin\myth_rank.exe --configPath=%APP_CONFIG_ROOT%//rank > .\\logs\\rank.log
timeout /T 1 /NOBREAK
start /b .\bin\myth_game.exe --configPath=%APP_CONFIG_ROOT%//game > .\\logs\\game.log
timeout /T 1 /NOBREAK
start /b .\bin\myth_gateway.exe --configPath=%APP_CONFIG_ROOT%//gateway > .\\logs\\gateway.log
timeout /T 1 /NOBREAK

echo "服务器启动完成"
pause