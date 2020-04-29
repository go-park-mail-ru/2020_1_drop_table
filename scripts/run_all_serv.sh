cd ../internal/app/main || exit
go run start.go &
cd ../../microservices/staff/main || exit
go run start.go &
cd ../../survey/main || exit
go run start.go &
