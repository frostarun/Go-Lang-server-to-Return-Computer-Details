# Go-Lang-server-to-Return-Computer-Details
This is a Go Lang server as a service to Return Computer Details in XML for a REST - GET call .  I have used Kardianos and Gopsutil for running as a service and getting cpu information.

# How to use :: 

go run getusage.go -service=install

go run getusage.go -service=start

localhost:8083/getusage
