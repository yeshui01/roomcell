@echo off

set PROJECT_ROOT=E:\tyh_work_card\server_project\roomcell

set PROTO_ROOT_PATH=%PROJECT_ROOT%\protos\

set SRC_CLIENT=%PROJECT_ROOT%\protos\proto_client
@REM set TAR_CLIENT=%PROJECT_ROOT%\pkg\pb\
set TAR_CLIENT_PKG_ROOT=E:\tyh_work_card\server_project\\

set SRC_SERVER=%PROJECT_ROOT%\protos\proto_server
set TAR_SERVER=%PROJECT_ROOT%\pkg\pb\

::set SRC_BATTLE=%PROJECT_ROOT%\proto\battle
::set TAR_BATTLE=%PROJECT_ROOT%\pkg\pb\
::protoc --proto_path=%SRC_BATTLE% --go_out=%TAR_BATTLE% %SRC_BATTLE%\battle2.proto
::protoc --proto_path=%SRC_BATTLE% --go_out=%TAR_BATTLE% %SRC_BATTLE%\battle2_client.proto

:: 客户端
@REM protoc --proto_path=%SRC_CLIENT% --go_out=%TAR_CLIENT_PKG_ROOT%  %SRC_CLIENT%\c_player.proto
@REM protoc --proto_path=%SRC_CLIENT% --go_out=%TAR_CLIENT_PKG_ROOT%  %SRC_CLIENT%\c_room_draw_guess.proto
protoc --proto_path=%SRC_CLIENT% --go_out=%TAR_CLIENT_PKG_ROOT%  %SRC_CLIENT%\c_room_undercover.proto
protoc --proto_path=%SRC_CLIENT% --go_out=%TAR_CLIENT_PKG_ROOT%  %SRC_CLIENT%\c_room.proto

:: 服务器
@REM protoc --proto_path=%SRC_SERVER% --go_out=%TAR_SERVER% %SRC_SERVER%\s_frame.proto
@REM protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_player.proto
@REM protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_db.proto
@REM protoc --proto_path=%SRC_SERVER% --proto_path=%PROJECT_ROOT% --proto_path=%PROTO_ROOT_PATH% --go_out=%TAR_SERVER% %SRC_SERVER%\s_room.proto
pause