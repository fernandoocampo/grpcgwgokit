package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"strings"

	grpcadapter "github.com/fernandoocampo/grpcgwgokit/internal/handler"
	"github.com/fernandoocampo/grpcgwgokit/internal/service"
	pb "github.com/fernandoocampo/grpcgwgokit/pkg/proto/grpcgwgokit/pb"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

func main() {

	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
		gRPCAddr = flag.String("grpc", ":8081", "gRPC listen address")
	)
	flag.Parse()

	serviceUser := service.NewBasicEcho()
	// finder endpoint
	endpoints := service.NewEndpoints(serviceUser)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var g group.Group

	// gRPC
	{
		handler := grpcadapter.NewGRPCServer(endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterYourServiceServer(gRPCServer, handler)
		g.Add(func() error {
			log.Println("msg", "starting grpc server", "grpc:", *gRPCAddr)
			listener, err := net.Listen("tcp", *gRPCAddr)
			if err != nil {
				return err
			}
			return gRPCServer.Serve(listener)
		}, func(err error) {
			log.Println("transport", "grpc", "error", err.Error())
		})
	}

	// http rest
	{
		router := mux.NewRouter()
		mux := runtime.NewServeMux()
		dialOpts := []grpc.DialOption{grpc.WithInsecure()}
		g.Add(func() error {
			err := pb.RegisterYourServiceHandlerFromEndpoint(ctx, mux, *gRPCAddr, dialOpts)
			if err != nil {
				return err
			}
			router.Handle("/", mux)
			log.Println("msg", "starting http server", "http:", *httpAddr)
			return http.ListenAndServe(*httpAddr, allowCORS(mux))
		}, func(err error) {
			log.Println("transport", "http", "error", err.Error())
		})
	}

	log.Println(g.Run())
}

// allowCORS allows Cross Origin Resoruce Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	glog.Infof("preflight request for %s", r.URL.Path)
	return
}
