package main
import (
	"middleware/packages/authentication"
	"middleware/packages/encryption"
	"middleware/packages/logmon"
	"fmt"
	"net/http"
)

func costumerEndpointHandler(w http.ResponseWriter, r  *http.Request){

}

func main(){
	http.HandleFunc("/costumer", auth.RoleMiddleware("costumer")(costumerEndpointHandler))
	}