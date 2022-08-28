set APP_CONFIG_ROOT=./
start /b ./roomcell_account.exe --configPath=%APP_CONFIG_ROOT% > .\\logs\\account_server.log

timeout /T 1 /NOBREAK
start /b ./roomcell_roommgr.exe --configPath=%APP_CONFIG_ROOT% > .\\logs\\hall_roommgr.log

timeout /T 1 /NOBREAK
start /b ./roomcell_room.exe --configPath=%APP_CONFIG_ROOT% > .\\logs\\hall_room.log

timeout /T 1 /NOBREAK
start /b ./roomcell_hall.exe --configPath=%APP_CONFIG_ROOT% > .\\logs\\hall_server.log

timeout /T 1 /NOBREAK
start /b ./roomcell_data.exe --configPath=%APP_CONFIG_ROOT% > .\\logs\\hall_data.log

timeout /T 1 /NOBREAK
start /b ./roomcell_gate.exe --configPath=%APP_CONFIG_ROOT% > .\\logs\\hall_gate.log