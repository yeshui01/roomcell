@echo off
::TABLE_NAME=log_online
:: 读取输入表名
echo please enter database table_name:
set /p TABLE_NAME=
echo is need generate log code(yes/no)?(only first use, don't repeated gen log code !!!):
set /p LOG_CODE=
:: 项目路径配置
set PROJECT_ROOT=E:\\tyh_work_card\\server_project\\roomcell
set SRC_PBORM=%PROJECT_ROOT%\protos\orm\
set TAR_PBORM=%PROJECT_ROOT%\pkg\

:: 生成
ormtools.exe -tableName=%TABLE_NAME% -logCode=%LOG_CODE%
timeout /T 1 /NOBREAK
protoc --proto_path=%SRC_PBORM% --go_out=%TAR_PBORM% %SRC_PBORM%\%TABLE_NAME%_ormpb.proto
::gen_logcode.exe -tableName=%TABLE_NAME% -logCode=%LOG_CODE%
pause