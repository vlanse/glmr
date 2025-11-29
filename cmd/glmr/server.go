package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/swaggest/swgui/v5emb"
	"github.com/vlanse/glmr/internal"
	"github.com/vlanse/glmr/internal/util/swagger"
	"google.golang.org/grpc"
)

var (
	grpcServerEndpoint = "0.0.0.0:10002"
	httpServerEndpoint = "localhost:8082"
)

func runGrpcServer(srv *grpc.Server) error {
	addr := grpcServerEndpoint
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("open web-ui port on %s: %w", addr, err)
	}

	go func() {
		log.Fatal(srv.Serve(listener))
	}()

	return nil
}

func serveGrpcGateway(mux *runtime.ServeMux) error {
	sw := internal.SwaggerContent

	swaggerMerger := swagger.NewMerger(internal.Name)
	_ = fs.WalkDir(sw, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, "swagger.json") {
			data, err := sw.ReadFile(path)
			if err != nil {
				return nil
			}
			_ = swaggerMerger.AddFile(bytes.NewReader(data))
		}
		return nil
	})

	swaggerData, _ := swaggerMerger.Content()
	if err := mux.HandlePath(
		http.MethodGet,
		"/swagger/swagger.json",
		func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(swaggerData)
		},
	); err != nil {
		return err
	}

	h := v5emb.New(
		internal.Name,
		fmt.Sprintf("http://%s/swagger/swagger.json", httpServerEndpoint),
		"/docs/",
	)
	if err := mux.HandlePath(
		http.MethodGet,
		"/docs/{content}",
		func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			h.ServeHTTP(w, r)
		},
	); err != nil {
		return err
	}
	return nil
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func serveFrontend(mux *runtime.ServeMux) error {
	content, err := fs.Sub(internal.FrontendContent, "ui/dist")
	if err != nil {
		log.Fatalln(err)
	}

	if err = mux.HandlePath(
		http.MethodGet, "/{content}",
		func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			http.FileServer(http.FS(content)).ServeHTTP(w, r)
		},
	); err != nil {
		return err
	}

	if err = mux.HandlePath(
		http.MethodGet, "/assets/{content}",
		func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			http.FileServer(http.FS(content)).ServeHTTP(w, r)
		},
	); err != nil {
		return err
	}
	return nil
}

func runServer(srv *grpc.Server, mux *runtime.ServeMux) error {
	if err := runGrpcServer(srv); err != nil {
		return err
	}

	if err := serveGrpcGateway(mux); err != nil {
		return err
	}

	if err := serveFrontend(mux); err != nil {
		return err
	}

	return http.ListenAndServe(httpServerEndpoint, corsMiddleware(mux))
}
