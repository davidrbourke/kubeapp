package main

import (
	"fmt"

	"github.com/spf13/cobra"

	svcsvc "kubeapp/internal/services"
)

func newServicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "services",
		Short:   "Manage services",
		Aliases: []string{"svc"},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List services",
		RunE:  runServicesList,
	}
	listCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List across all namespaces")

	cmd.AddCommand(listCmd)
	cmd.AddCommand(&cobra.Command{
		Use:   "describe <name>",
		Short: "Describe a service",
		Args:  cobra.ExactArgs(1),
		RunE:  runServicesDescribe,
	})

	return cmd
}

func runServicesList(cmd *cobra.Command, args []string) error {
	ns := namespace
	if allNamespaces {
		ns = ""
	}

	svc := svcsvc.NewService(client)
	services, err := svc.List(cmd.Context(), ns)
	if err != nil {
		return err
	}

	if len(services) == 0 {
		fmt.Println("No services found.")
		return nil
	}

	fmt.Printf("%-40s %-15s %-12s %-16s %-16s\n", "NAME", "NAMESPACE", "TYPE", "CLUSTER-IP", "EXTERNAL-IP")
	fmt.Printf("%-40s %-15s %-12s %-16s %-16s\n", "----", "---------", "----", "----------", "-----------")
	for _, s := range services {
		externalIP := "<none>"
		if len(s.Status.LoadBalancer.Ingress) > 0 {
			externalIP = s.Status.LoadBalancer.Ingress[0].IP
			if externalIP == "" {
				externalIP = s.Status.LoadBalancer.Ingress[0].Hostname
			}
		} else if len(s.Spec.ExternalIPs) > 0 {
			externalIP = s.Spec.ExternalIPs[0]
		}
		fmt.Printf("%-40s %-15s %-12s %-16s %-16s\n",
			s.Name, s.Namespace, string(s.Spec.Type), s.Spec.ClusterIP, externalIP)
	}
	return nil
}

func runServicesDescribe(cmd *cobra.Command, args []string) error {
	svc := svcsvc.NewService(client)
	service, err := svc.Get(cmd.Context(), namespace, args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Name:       %s\n", service.Name)
	fmt.Printf("Namespace:  %s\n", service.Namespace)
	fmt.Printf("Type:       %s\n", service.Spec.Type)
	fmt.Printf("ClusterIP:  %s\n", service.Spec.ClusterIP)

	if len(service.Spec.ExternalIPs) > 0 {
		fmt.Printf("ExternalIPs: %v\n", service.Spec.ExternalIPs)
	}

	if len(service.Status.LoadBalancer.Ingress) > 0 {
		fmt.Println("LoadBalancer Ingress:")
		for _, ing := range service.Status.LoadBalancer.Ingress {
			if ing.IP != "" {
				fmt.Printf("  IP:       %s\n", ing.IP)
			}
			if ing.Hostname != "" {
				fmt.Printf("  Hostname: %s\n", ing.Hostname)
			}
		}
	}

	if len(service.Spec.Ports) > 0 {
		fmt.Println("Ports:")
		for _, p := range service.Spec.Ports {
			name := p.Name
			if name == "" {
				name = "<unnamed>"
			}
			fmt.Printf("  %-16s %d:%s/%s\n", name, p.Port, p.TargetPort.String(), p.Protocol)
		}
	}

	if len(service.Spec.Selector) > 0 {
		fmt.Println("Selector:")
		for k, v := range service.Spec.Selector {
			fmt.Printf("  %s=%s\n", k, v)
		}
	}

	if len(service.Labels) > 0 {
		fmt.Println("Labels:")
		for k, v := range service.Labels {
			fmt.Printf("  %s=%s\n", k, v)
		}
	}

	return nil
}
