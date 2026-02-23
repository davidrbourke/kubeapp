package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"

	corev1 "k8s.io/api/core/v1"

	podsvc "kubeapp/internal/pods"
	svcsvc "kubeapp/internal/services"
)

// App holds shared dependencies for HTTP handlers.
type App struct {
	podSvc    *podsvc.Service
	svcSvc    *svcsvc.Service
	tmpls     map[string]*template.Template
	namespace string // default namespace filter; empty = all namespaces
}

type indexData struct {
	Pods     []corev1.Pod
	Services []corev1.Service
	NS       string
}

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("namespace")
	if ns == "" {
		ns = a.namespace
	}

	var (
		pods   []corev1.Pod
		svcs   []corev1.Service
		podErr error
		svcErr error
		wg     sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		pods, podErr = a.podSvc.List(r.Context(), ns)
	}()
	go func() {
		defer wg.Done()
		svcs, svcErr = a.svcSvc.List(r.Context(), ns)
	}()
	wg.Wait()

	if podErr != nil {
		http.Error(w, fmt.Sprintf("listing pods: %v", podErr), http.StatusInternalServerError)
		return
	}
	if svcErr != nil {
		http.Error(w, fmt.Sprintf("listing services: %v", svcErr), http.StatusInternalServerError)
		return
	}

	data := indexData{Pods: pods, Services: svcs, NS: ns}
	if err := a.tmpls["index.html"].ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) podDetailHandler(w http.ResponseWriter, r *http.Request) {
	ns := r.PathValue("namespace")
	name := r.PathValue("name")

	pod, err := a.podSvc.Get(r.Context(), ns, name)
	if err != nil {
		http.Error(w, fmt.Sprintf("pod not found: %v", err), http.StatusNotFound)
		return
	}

	if err := a.tmpls["pod.html"].ExecuteTemplate(w, "layout", pod); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) podRestartHandler(w http.ResponseWriter, r *http.Request) {
	ns := r.PathValue("namespace")
	name := r.PathValue("name")

	if err := a.podSvc.Delete(r.Context(), ns, name); err != nil {
		http.Error(w, fmt.Sprintf("could not restart pod: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) serviceDetailHandler(w http.ResponseWriter, r *http.Request) {
	ns := r.PathValue("namespace")
	name := r.PathValue("name")

	svc, err := a.svcSvc.Get(r.Context(), ns, name)
	if err != nil {
		http.Error(w, fmt.Sprintf("service not found: %v", err), http.StatusNotFound)
		return
	}

	if err := a.tmpls["service.html"].ExecuteTemplate(w, "layout", svc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"phaseBadge": func(phase corev1.PodPhase) string {
			switch phase {
			case corev1.PodRunning:
				return "badge-running"
			case corev1.PodPending:
				return "badge-pending"
			case corev1.PodFailed:
				return "badge-failed"
			case corev1.PodSucceeded:
				return "badge-succeeded"
			default:
				return "badge-unknown"
			}
		},
		"readyCount": func(pod corev1.Pod) int {
			n := 0
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.Ready {
					n++
				}
			}
			return n
		},
		"externalIP": func(svc corev1.Service) string {
			if len(svc.Status.LoadBalancer.Ingress) > 0 {
				if ip := svc.Status.LoadBalancer.Ingress[0].IP; ip != "" {
					return ip
				}
				return svc.Status.LoadBalancer.Ingress[0].Hostname
			}
			if len(svc.Spec.ExternalIPs) > 0 {
				return svc.Spec.ExternalIPs[0]
			}
			return "<none>"
		},
		"containerStatus": func(pod *corev1.Pod, name string) *corev1.ContainerStatus {
			for i := range pod.Status.ContainerStatuses {
				if pod.Status.ContainerStatuses[i].Name == name {
					return &pod.Status.ContainerStatuses[i]
				}
			}
			return nil
		},
		"portList": func(ports []corev1.ContainerPort) string {
			if len(ports) == 0 {
				return ""
			}
			parts := make([]string, len(ports))
			for i, p := range ports {
				parts[i] = fmt.Sprintf("%d/%s", p.ContainerPort, p.Protocol)
			}
			return strings.Join(parts, ", ")
		},
	}
}
