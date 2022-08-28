::打包客户端
@echo on
@REM F:\MyDevelopHome\go_path\bin/fyne.exe package --name cellclient.exe --exe G:\work\dangwan\fs_project\roomcell\bin  --src G:\work\dangwan\fs_project\roomcell\cmd\cellclient --icon icon.png
@REM F:\MyDevelopHome\go_path\bin/fyne.exe package --name cellclient.exe  --src G:\work\dangwan\fs_project\roomcell\cmd\cellclient --icon icon.png
@REM timeout /T 1
@REM DEL ./bin/cellclient.exe
@REM MOVE ./cmd/cellclient/cellclient.exe ./bin/
set PROJECTROOT=G://work//dangwan//fs_project//roomcell
copy %PROJECTROOT%/cmd/cellclient/cellclient.exe %PROJECTROOT%/bin//cellclient.exe
@REM del %PROJECT_ROOT%//cmd//cellclient//cellclient.exe