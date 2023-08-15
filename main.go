// main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gogorush/simple_auth/auth"
	"github.com/gogorush/simple_auth/utils"
)

var (
	fileConfig string // config file path
)

func init() {
	flag.StringVar(&fileConfig, "f", "default.yml", "-f: config file path")
}

func main() {

    flag.Parse()

	utils.InitConfig(fileConfig)

    fmt.Println(utils.GetConfig())
    auth.JwtKey = []byte(utils.GetConfig().JwtKey)

	http.HandleFunc("/create-user", auth.HandleCreateUser)
	http.HandleFunc("/delete-user", auth.HandleDeleteUser)
	http.HandleFunc("/create-role", auth.HandleCreateRole)
	http.HandleFunc("/delete-role", auth.HandleDeleteRole)
	http.HandleFunc("/add-role-to-user", auth.HandleAddRoleToUser)
	http.HandleFunc("/authenticate", auth.HandleAuthenticate)
	http.HandleFunc("/invalidate-token", auth.HandleInvalidateToken)
	http.HandleFunc("/check-role", auth.HandleCheckRole)
	http.HandleFunc("/get-all-roles", auth.HandleGetAllRoles)


	// Load the HTTPS certificate and key
	//cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	//if err != nil {
	//	log.Fatalf("Server failed to start: %v", err)

	//}
	//tlsConfig := &tls.Config{
	//	Certificates: []tls.Certificate{cert},
	//}

	server := &http.Server{
		Addr: utils.GetConfig().Address,
		//TLSConfig: tlsConfig,
	}

	log.Printf("Starting server on https://localhost%v", server.Addr)
	//log.Fatal(server.ListenAndServeTLS("", ""))
	log.Fatal(server.ListenAndServe())
}
