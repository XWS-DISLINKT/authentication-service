module authentication_service

go 1.18

replace github.com/XWS-DISLINKT/dislinkt/common => ../common

require (
	github.com/XWS-DISLINKT/dislinkt/common v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	google.golang.org/grpc v1.45.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220317150908-0efb43f6373e // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)
