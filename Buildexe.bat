:: 账号服务器
go build -o ./bin/roomcell_roommgr.exe ./cmd/hall_roommgr/main.go
go build -o ./bin/roomcell_room.exe ./cmd/hall_room/main.go
go build -o ./bin/roomcell_account.exe ./cmd/account/main.go
go build -o ./bin/roomcell_hall.exe ./cmd/hall_server/main.go
go build -o ./bin/roomcell_gate.exe ./cmd/hall_gate/main.go
go build -o ./bin/roomcell_data.exe ./cmd/hall_data/main.go