package main

import (
	"fmt"

	"github.com/spf13/cobra"

	podsvc "kubeapp/internal/pods"
)

func newPodsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pods",
		Short: "Manage pods",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List pods",
		RunE:  runPodsList,
	}
	listCmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List across all namespaces")

	cmd.AddCommand(listCmd)
	cmd.AddCommand(&cobra.Command{
		Use:   "describe <name>",
		Short: "Describe a pod",
		Args:  cobra.ExactArgs(1),
		RunE:  runPodsDescribe,
	})

	return cmd
}

func runPodsList(cmd *cobra.Command, args []string) error {
	ns := namespace
	if allNamespaces {
		ns = ""
	}

	svc := podsvc.NewService(client)
	pods, err := svc.List(cmd.Context(), ns)
	if err != nil {
		return err
	}

	if len(pods) == 0 {
		fmt.Println("No pods found.")
		return nil
	}

	fmt.Printf("%-52s %-15s %-12s %-8s\n", "NAME", "NAMESPACE", "STATUS", "READY")
	fmt.Printf("%-52s %-15s %-12s %-8s\n", "----", "---------", "------", "-----")
	for _, p := range pods {
		ready := 0
		for _, cs := range p.Status.ContainerStatuses {
			if cs.Ready {
				ready++
			}
		}
		fmt.Printf("%-52s %-15s %-12s %d/%d\n",
			p.Name, p.Namespace, string(p.Status.Phase),
			ready, len(p.Spec.Containers))
	}
	return nil
}

func runPodsDescribe(cmd *cobra.Command, args []string) error {
	svc := podsvc.NewService(client)
	pod, err := svc.Get(cmd.Context(), namespace, args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Name:       %s\n", pod.Name)
	fmt.Printf("Namespace:  %s\n", pod.Namespace)
	fmt.Printf("Node:       %s\n", pod.Spec.NodeName)
	fmt.Printf("Status:     %s\n", pod.Status.Phase)
	fmt.Printf("Pod IP:     %s\n", pod.Status.PodIP)
	if pod.Status.StartTime != nil {
		fmt.Printf("Started:    %s\n", pod.Status.StartTime)
	}

	if len(pod.Labels) > 0 {
		fmt.Println("Labels:")
		for k, v := range pod.Labels {
			fmt.Printf("  %s=%s\n", k, v)
		}
	}

	fmt.Println("Containers:")
	for _, c := range pod.Spec.Containers {
		fmt.Printf("  Name:   %s\n", c.Name)
		fmt.Printf("  Image:  %s\n", c.Image)
		if len(c.Ports) > 0 {
			fmt.Printf("  Ports:  ")
			for i, p := range c.Ports {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%d/%s", p.ContainerPort, p.Protocol)
			}
			fmt.Println()
		}
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.Name == c.Name {
				fmt.Printf("  Ready:    %v\n", cs.Ready)
				fmt.Printf("  Restarts: %d\n", cs.RestartCount)
			}
		}
	}

	if len(pod.Status.Conditions) > 0 {
		fmt.Println("Conditions:")
		for _, cond := range pod.Status.Conditions {
			fmt.Printf("  %-24s %s\n", cond.Type, cond.Status)
		}
	}

	return nil
}
