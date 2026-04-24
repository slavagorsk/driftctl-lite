package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/driftctl-lite/internal/tfstate"
)

func main() {
	statePath := flag.String("state", "terraform.tfstate", "path to Terraform state file")
	resType := flag.String("type", "", "filter output by resource type (optional)")
	flag.Parse()

	state, err := tfstate.ParseFile(*statePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	resources := state.Resources
	if *resType != "" {
		resources = tfstate.FindByType(state, *resType)
	}

	if len(resources) == 0 {
		fmt.Println("No resources found.")
		return
	}

	fmt.Printf("%-40s %-20s %s\n", "RESOURCE", "ID", "PROVIDER")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, r := range resources {
		id, _ := tfstate.GetAttribute(r, "id")
		key := tfstate.ResourceKey{Type: r.Type, Name: r.Name}
		fmt.Printf("%-40s %-20s %s\n", key, id, r.Provider)
	}
}
