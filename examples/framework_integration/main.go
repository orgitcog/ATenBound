// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package main provides an example of using the integrated framework
// with Tensor Logic, Hypermind, and ATenSpace in Boundary.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/boundary/internal/atenspace"
	"github.com/hashicorp/boundary/internal/hypermind"
	"github.com/hashicorp/boundary/internal/integration"
	"github.com/hashicorp/boundary/internal/tensorlogic"
)

func main() {
	ctx := context.Background()

	// Example 1: Using Individual Frameworks
	fmt.Println("=== Example 1: Individual Framework Usage ===\n")
	tensorLogicExample(ctx)
	hypermindExample(ctx)
	atenSpaceExample(ctx)

	// Example 2: Unified Integration
	fmt.Println("\n=== Example 2: Unified Integration ===\n")
	unifiedIntegrationExample(ctx)
}

// tensorLogicExample demonstrates Tensor Logic framework usage
func tensorLogicExample(ctx context.Context) {
	fmt.Println("--- Tensor Logic Framework ---")

	// Create framework
	framework, err := tensorlogic.NewFramework(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Register variables
	scope := &tensorlogic.Variable{
		Name:    "organization_scope",
		Indices: []string{"users", "permissions"},
		Shape:   []int{100, 50},
		Type:    tensorlogic.HybridType,
		Data:    make([]float64, 5000),
	}
	framework.RegisterVariable(ctx, scope)

	users := &tensorlogic.Variable{
		Name:    "users",
		Indices: []string{"id", "attributes"},
		Shape:   []int{100, 10},
		Type:    tensorlogic.SymbolicType,
		Data:    make([]float64, 1000),
	}
	framework.RegisterVariable(ctx, users)

	// Define an equation
	eq := &tensorlogic.TensorEquation{
		Left: tensorlogic.Variable{
			Name:    "access_matrix",
			Indices: []string{"users", "resources"},
		},
		Right:     "users_ij * permissions_jk",
		Operation: "join",
	}
	framework.DefineEquation(ctx, eq)

	// Evaluate
	result, _ := framework.Evaluate(ctx, "organization_scope")
	fmt.Printf("Evaluated tensor variable: %s (type: %s)\n", result.Name, result.Type)
	fmt.Printf("Variables registered: %d\n", len(framework.Variables))
	fmt.Printf("Equations defined: %d\n\n", len(framework.Equations))
}

// hypermindExample demonstrates Hypermind multi-scope architecture
func hypermindExample(ctx context.Context) {
	fmt.Println("--- Hypermind Multi-Scope Architecture ---")

	// Create architecture
	msa, err := hypermind.NewMultiScopeArchitecture(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Register scopes
	globalScope := &hypermind.DistributedScope{
		ID:   "global",
		Type: "global",
	}
	msa.RegisterScope(ctx, globalScope)

	orgScope := &hypermind.DistributedScope{
		ID:       "org-acme",
		ParentID: "global",
		Type:     "org",
	}
	msa.RegisterScope(ctx, orgScope)

	projectScope := &hypermind.DistributedScope{
		ID:       "project-alpha",
		ParentID: "org-acme",
		Type:     "project",
	}
	msa.RegisterScope(ctx, projectScope)

	// Connect peers
	peer1 := &hypermind.Peer{
		ID:       "peer-us-west",
		Address:  "192.168.1.10:8080",
		ScopeIDs: []string{"org-acme"},
	}
	msa.ConnectPeer(ctx, peer1)

	peer2 := &hypermind.Peer{
		ID:       "peer-us-east",
		Address:  "192.168.1.20:8080",
		ScopeIDs: []string{"org-acme", "project-alpha"},
	}
	msa.ConnectPeer(ctx, peer2)

	// Propagate state
	state := map[string]interface{}{
		"status":      "active",
		"region":      "us",
		"environment": "production",
	}
	msa.PropagateState(ctx, "org-acme", state)

	// Discover peers
	peers, _ := msa.DiscoverPeers(ctx, "org-acme")
	fmt.Printf("Distributed scopes: 3 (global, org, project)\n")
	fmt.Printf("Connected peers: %d\n", len(msa.GetActivePeers(ctx)))
	fmt.Printf("Peers for org-acme: %d\n\n", len(peers))
}

// atenSpaceExample demonstrates ATenSpace integration
func atenSpaceExample(ctx context.Context) {
	fmt.Println("--- ATenSpace (Space defined by Boundary) ---")

	// Create space
	space, err := atenspace.NewSpace(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Add atoms (Boundary entities)
	globalAtom := &atenspace.Atom{
		ID:   "global",
		Type: atenspace.AggregateAtom,
		Name: "Global Scope",
	}
	space.AddAtom(ctx, globalAtom)

	orgAtom := &atenspace.Atom{
		ID:   "org-acme",
		Type: atenspace.AggregateAtom,
		Name: "ACME Organization",
	}
	space.AddAtom(ctx, orgAtom)

	userAtom := &atenspace.Atom{
		ID:   "user-alice",
		Type: atenspace.EntityAtom,
		Name: "Alice",
	}
	space.AddAtom(ctx, userAtom)

	// Create links
	scopeLink := &atenspace.Link{
		ID:       "link-global-org",
		Type:     atenspace.ScopeLink,
		Source:   "global",
		Target:   "org-acme",
		Strength: 1.0,
	}
	space.AddLink(ctx, scopeLink)

	memberLink := &atenspace.Link{
		ID:       "link-user-org",
		Type:     atenspace.MembershipLink,
		Source:   "org-acme",
		Target:   "user-alice",
		Strength: 0.9,
	}
	space.AddLink(ctx, memberLink)

	// Attach tensors
	orgTensor := &atenspace.Tensor{
		ID:     "tensor-org",
		Shape:  []int{10, 10},
		Data:   make([]float64, 100),
		DType:  "float64",
		Device: "cpu",
	}
	space.AttachTensor(ctx, "org-acme", orgTensor)

	// Define boundary (Space is defined by Boundary domain model)
	boundary := &atenspace.DomainBoundary{
		ID:      "boundary-org",
		Name:    "Organization Boundary",
		Type:    atenspace.ScopeBoundary,
		AtomIDs: []string{"org-acme", "user-alice"},
	}
	space.DefineBoundary(ctx, boundary)

	// Query
	atoms, _ := space.QueryByBoundary(ctx, "boundary-org")
	links := space.GetLinksForAtom(ctx, "org-acme")

	fmt.Printf("Atoms in space: 3 (global, org, user)\n")
	fmt.Printf("Links in space: 2 (scope, membership)\n")
	fmt.Printf("Domain boundaries: %d\n", len(space.GetBoundaries(ctx)))
	fmt.Printf("Atoms in org boundary: %d\n", len(atoms))
	fmt.Printf("Links for org-acme: %d\n\n", len(links))
}

// unifiedIntegrationExample demonstrates the unified framework
func unifiedIntegrationExample(ctx context.Context) {
	fmt.Println("--- Unified Framework Integration ---")

	// Create unified framework
	uf, err := integration.NewUnifiedFramework(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Integrate with Boundary
	if err := uf.IntegrateWithBoundary(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Integrated Tensor Logic, Hypermind, and ATenSpace with Boundary")

	// Create boundary scopes (integrated across all frameworks)
	scopes := []struct {
		id       string
		scopeType string
	}{
		{"global", "global"},
		{"org-acme", "org"},
		{"project-alpha", "project"},
	}

	for _, s := range scopes {
		if err := uf.CreateBoundaryScope(ctx, s.id, s.scopeType); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("✓ Created scope '%s' across all frameworks\n", s.id)
	}

	// Define domain boundary
	if err := uf.DefineDomainBoundary(ctx, "org-boundary", "scope",
		[]string{"org-acme", "project-alpha"}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Defined domain boundary")

	// Propagate state across frameworks
	state := map[string]interface{}{
		"status":      "active",
		"environment": "production",
		"region":      "us-west",
	}
	if err := uf.PropagateState(ctx, "org-acme", state); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Propagated state across P2P network")

	// Query scope across all frameworks
	info, err := uf.QueryScope(ctx, "org-acme")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n--- Scope Query Results ---")
	fmt.Printf("Scope ID: %s\n", info.ID)

	if info.TensorVariable != nil {
		fmt.Printf("✓ Tensor Logic: Variable '%s' (type: %s)\n",
			info.TensorVariable.Name, info.TensorVariable.Type)
	}

	if info.DistributedScope != nil {
		fmt.Printf("✓ Hypermind: Distributed scope (type: %s, state entries: %d)\n",
			info.DistributedScope.Type, len(info.DistributedScope.State))
	}

	if info.Atom != nil {
		fmt.Printf("✓ ATenSpace: Atom (type: %s, has tensor: %v)\n",
			info.Atom.Type, info.Atom.TensorID != "")
	}

	fmt.Println("\n✓ All frameworks successfully integrated and working together!")
}
