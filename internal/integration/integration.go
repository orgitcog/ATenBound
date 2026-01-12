// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package integration provides the unified integration layer for
// Tensor Logic, Hypermind, and ATenSpace frameworks in Boundary.
//
// This package demonstrates how all three frameworks work together:
// - Tensor Logic: Provides tensor equation framework for all variables
// - Hypermind: Enables distributed multi-scope P2P architecture
// - ATenSpace: Defines Space through Boundary domain model with tensor operations
package integration

import (
	"context"

	"github.com/hashicorp/boundary/internal/atenspace"
	"github.com/hashicorp/boundary/internal/errors"
	"github.com/hashicorp/boundary/internal/hypermind"
	"github.com/hashicorp/boundary/internal/tensorlogic"
)

// UnifiedFramework integrates all three frameworks into a cohesive system.
type UnifiedFramework struct {
	// TensorLogic provides the tensor equation framework
	TensorLogic *tensorlogic.Framework

	// Hypermind provides multi-scope P2P architecture
	Hypermind *hypermind.MultiScopeArchitecture

	// ATenSpace provides the Space defined by Boundary domain model
	ATenSpace *atenspace.Space
}

// NewUnifiedFramework creates a new integrated framework instance.
func NewUnifiedFramework(ctx context.Context) (*UnifiedFramework, error) {
	const op = "integration.NewUnifiedFramework"

	// Initialize Tensor Logic framework
	tl, err := tensorlogic.NewFramework(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op, errors.WithMsg("failed to initialize tensor logic"))
	}

	// Initialize Hypermind multi-scope architecture
	hm, err := hypermind.NewMultiScopeArchitecture(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op, errors.WithMsg("failed to initialize hypermind"))
	}

	// Initialize ATenSpace
	as, err := atenspace.NewSpace(ctx)
	if err != nil {
		return nil, errors.Wrap(ctx, err, op, errors.WithMsg("failed to initialize atenspace"))
	}

	uf := &UnifiedFramework{
		TensorLogic: tl,
		Hypermind:   hm,
		ATenSpace:   as,
	}

	return uf, nil
}

// IntegrateWithBoundary integrates all three frameworks with Boundary's domain model.
// This is the key integration point where all frameworks work together:
// 1. Tensor Logic: All Boundary variables use tensor equations
// 2. Hypermind: Scopes become distributed P2P entities
// 3. ATenSpace: Boundary domain defines the Space with tensor operations
func (u *UnifiedFramework) IntegrateWithBoundary(ctx context.Context) error {
	const op = "integration.(UnifiedFramework).IntegrateWithBoundary"

	// Integrate Tensor Logic with Boundary variables
	if err := u.TensorLogic.IntegrateWithBoundary(ctx); err != nil {
		return errors.Wrap(ctx, err, op, errors.WithMsg("tensor logic integration failed"))
	}

	// Integrate Hypermind with Boundary scope system
	if err := u.Hypermind.IntegrateWithBoundary(ctx); err != nil {
		return errors.Wrap(ctx, err, op, errors.WithMsg("hypermind integration failed"))
	}

	// Integrate ATenSpace where Space is defined by Boundary
	if err := u.ATenSpace.IntegrateWithBoundary(ctx); err != nil {
		return errors.Wrap(ctx, err, op, errors.WithMsg("atenspace integration failed"))
	}

	return nil
}

// CreateBoundaryScope creates a scope that integrates all three frameworks.
// This demonstrates the unified architecture in action:
// - The scope is represented as a tensor variable (Tensor Logic)
// - The scope participates in P2P network (Hypermind)
// - The scope is an atom in the Space (ATenSpace)
func (u *UnifiedFramework) CreateBoundaryScope(ctx context.Context, scopeID, scopeType string) error {
	const op = "integration.(UnifiedFramework).CreateBoundaryScope"

	// Create tensor variable for the scope (Tensor Logic)
	scopeVar := &tensorlogic.Variable{
		Name:    scopeID,
		Indices: []string{"entity", "property"},
		Type:    tensorlogic.HybridType,
	}
	if err := u.TensorLogic.RegisterVariable(ctx, scopeVar); err != nil {
		return errors.Wrap(ctx, err, op)
	}

	// Create distributed scope (Hypermind)
	distScope := &hypermind.DistributedScope{
		ID:   scopeID,
		Type: scopeType,
	}
	if err := u.Hypermind.RegisterScope(ctx, distScope); err != nil {
		return errors.Wrap(ctx, err, op)
	}

	// Create atom in Space (ATenSpace)
	atom := &atenspace.Atom{
		ID:   scopeID,
		Type: atenspace.AggregateAtom,
		Name: scopeID,
	}
	if err := u.ATenSpace.AddAtom(ctx, atom); err != nil {
		return errors.Wrap(ctx, err, op)
	}

	// Attach tensor to atom
	tensor := &atenspace.Tensor{
		ID:     scopeID + "_tensor",
		Shape:  []int{10, 10},
		Data:   make([]float64, 100),
		DType:  "float64",
		Device: "cpu",
	}
	if err := u.ATenSpace.AttachTensor(ctx, scopeID, tensor); err != nil {
		return errors.Wrap(ctx, err, op)
	}

	return nil
}

// QueryScope demonstrates querying across all three frameworks.
func (u *UnifiedFramework) QueryScope(ctx context.Context, scopeID string) (*ScopeInfo, error) {
	const op = "integration.(UnifiedFramework).QueryScope"

	info := &ScopeInfo{
		ID: scopeID,
	}

	// Get tensor representation (Tensor Logic)
	if tensorVar, err := u.TensorLogic.Evaluate(ctx, scopeID); err == nil {
		info.TensorVariable = tensorVar
	}

	// Get distributed scope info (Hypermind)
	if distScope, err := u.Hypermind.GetScope(ctx, scopeID); err == nil {
		info.DistributedScope = distScope
	}

	// Get atom representation (ATenSpace)
	if atom, err := u.ATenSpace.GetAtom(ctx, scopeID); err == nil {
		info.Atom = atom
	}

	return info, nil
}

// ScopeInfo aggregates information from all three frameworks.
type ScopeInfo struct {
	ID               string
	TensorVariable   *tensorlogic.Variable
	DistributedScope *hypermind.DistributedScope
	Atom             *atenspace.Atom
}

// DefineDomainBoundary creates a boundary that spans all frameworks.
func (u *UnifiedFramework) DefineDomainBoundary(ctx context.Context, boundaryID, boundaryType string, atomIDs []string) error {
	const op = "integration.(UnifiedFramework).DefineDomainBoundary"

	// Define boundary in ATenSpace (where Space is defined by Boundary)
	boundary := &atenspace.DomainBoundary{
		ID:      boundaryID,
		Name:    boundaryID,
		Type:    atenspace.BoundaryType(boundaryType),
		AtomIDs: atomIDs,
	}
	if err := u.ATenSpace.DefineBoundary(ctx, boundary); err != nil {
		return errors.Wrap(ctx, err, op)
	}

	return nil
}

// PropagateState demonstrates state propagation across frameworks.
func (u *UnifiedFramework) PropagateState(ctx context.Context, scopeID string, state map[string]interface{}) error {
	const op = "integration.(UnifiedFramework).PropagateState"

	// Propagate through Hypermind P2P network
	if err := u.Hypermind.PropagateState(ctx, scopeID, state); err != nil {
		return errors.Wrap(ctx, err, op)
	}

	// Update atom attributes in ATenSpace
	atom, err := u.ATenSpace.GetAtom(ctx, scopeID)
	if err != nil {
		return errors.Wrap(ctx, err, op)
	}

	for k, v := range state {
		atom.Attributes[k] = v
	}

	return nil
}
