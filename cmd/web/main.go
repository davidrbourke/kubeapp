package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"k8s.io/client-go/util/homedir"

	k8sclient "kubeapp/internal/k8s"
	podsvc "kubeapp/internal/pods"
	svcsvc "kubeapp/internal/services"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	defaultKubeconfig := ""
	if home := homedir.HomeDir(); home != "" {
		defaultKubeconfig = filepath.Join(home, ".kube", "config")
	}

	port := flag.Int("port", 8080, "HTTP server port")
	kubeconfig := flag.String("kubeconfig", defaultKubeconfig, "path to kubeconfig file")
	namespace := flag.String("namespace", "", "default namespace filter (empty = all namespaces)")
	flag.Parse()

	client, err := k8sclient.NewClient(*kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to cluster: %v\n", err)
		os.Exit(1)
	}

	tmpls, err := parseTemplates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse templates: %v\n", err)
		os.Exit(1)
	}

	app := &App{
		podSvc:    podsvc.NewService(client),
		svcSvc:    svcsvc.NewService(client),
		tmpls:     tmpls,
		namespace: *namespace,
	}

	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not prepare static files: %v\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticSub)))
	mux.HandleFunc("GET /", app.indexHandler)
	mux.HandleFunc("GET /pods/{namespace}/{name}", app.podDetailHandler)
	mux.HandleFunc("GET /services/{namespace}/{name}", app.serviceDetailHandler)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("kubeapp-web listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func parseTemplates() (map[string]*template.Template, error) {
	pages := []string{"index.html", "pod.html", "service.html"}
	tmpls := make(map[string]*template.Template, len(pages))
	for _, page := range pages {
		t, err := template.New("").Funcs(templateFuncs()).ParseFS(templateFS,
			"templates/layout.html",
			"templates/"+page,
		)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", page, err)
		}
		tmpls[page] = t
	}
	return tmpls, nil
}
