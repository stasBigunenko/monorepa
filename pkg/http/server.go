package httpserver

import (
	"log"
	itemsservice "monorepa/service/items"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	itemsservice itemsservice.ItemsService
}

func New() HTTPServer {
	return HTTPServer{}
}

func (hs HTTPServer) GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/items", hs.ListItems).Methods("GET")

	router.Use(hs.authMiddleware)

	return router
}

func (hs HTTPServer) ListItems(w http.ResponseWriter, req *http.Request) {
	// resp, err := hs.client.ListItems(ctx, &pb.ListItemsRequest{})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }

	// var items Item
	// json.Unmarshal(resp.Body, &items)
}

func (hs HTTPServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Auth starting...")

		next.ServeHTTP(w, r)
	})
}
