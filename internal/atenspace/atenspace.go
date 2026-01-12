// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package atenspace provides ATenSpace integration for Boundary.
// Based on https://github.com/o9nn/ATenSpace - ATen Tensors + OpenCog AtomSpace
//
// ATenSpace combines ATen tensor operations (from PyTorch) with OpenCog's AtomSpace
// (hypergraph knowledge base) to create a unified framework where "Space" is defined
// by the Boundary domain model.
package atenspace

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/boundary/internal/errors"
)

// Space represents the ATenSpace where the "Space" is defined by Boundary's domain model.
// It combines tensor computation with hypergraph knowledge representation.
type Space struct {
	// Atoms are the nodes in the hypergraph (Boundary domain entities)
	atoms map[string]*Atom

	// Links are the edges in the hypergraph (relationships between entities)
	links []*Link

	// TensorStore maps atoms to their tensor representations
	tensorStore map[string]*Tensor

	// Boundaries define the domain boundaries (from Boundary domain model)
	boundaries []*DomainBoundary

	// mu protects concurrent access
	mu sync.RWMutex
}

// Atom represents a node in the AtomSpace hypergraph.
// In ATenSpace, atoms correspond to Boundary domain entities.
type Atom struct {
	// ID is the unique atom identifier
	ID string

	// Type is the atom type (Entity, Aggregate, Resource, etc.)
	Type AtomType

	// Name is the human-readable name
	Name string

	// Attributes hold additional properties
	Attributes map[string]interface{}

	// TensorID references the associated tensor representation
	TensorID string

	// CreatedAt timestamp
	CreatedAt time.Time
}

// AtomType defines the type of atom in the space.
type AtomType string

const (
	// EntityAtom represents Boundary entities
	EntityAtom AtomType = "entity"

	// AggregateAtom represents Boundary aggregates
	AggregateAtom AtomType = "aggregate"

	// ResourceAtom represents Boundary resources
	ResourceAtom AtomType = "resource"

	// RelationAtom represents relationships
	RelationAtom AtomType = "relation"

	// ConceptAtom represents abstract concepts
	ConceptAtom AtomType = "concept"
)

// Link represents an edge in the AtomSpace hypergraph.
type Link struct {
	// ID is the unique link identifier
	ID string

	// Type is the link type
	Type LinkType

	// Source is the source atom ID
	Source string

	// Target is the target atom ID
	Target string

	// Strength represents the link strength (0.0 to 1.0)
	Strength float64

	// CreatedAt timestamp
	CreatedAt time.Time
}

// LinkType defines the type of link between atoms.
type LinkType string

const (
	// InheritanceLink represents is-a relationships
	InheritanceLink LinkType = "inheritance"

	// MembershipLink represents part-of relationships
	MembershipLink LinkType = "membership"

	// ScopeLink represents scope containment
	ScopeLink LinkType = "scope"

	// DependencyLink represents dependencies
	DependencyLink LinkType = "dependency"

	// AssociationLink represents general associations
	AssociationLink LinkType = "association"
)

// Tensor represents the ATen tensor associated with an atom.
type Tensor struct {
	// ID is the unique tensor identifier
	ID string

	// Shape defines the tensor dimensions
	Shape []int

	// Data holds the tensor data (flattened)
	Data []float64

	// DType is the data type
	DType string

	// Device specifies where the tensor is stored (cpu, cuda, etc.)
	Device string
}

// DomainBoundary defines a boundary in the Boundary domain model.
// This is the key integration point where "Space" is defined by "Boundary".
type DomainBoundary struct {
	// ID is the boundary identifier
	ID string

	// Name is the boundary name
	Name string

	// Type is the boundary type (transactional, security, scope, etc.)
	Type BoundaryType

	// AtomIDs are the atoms within this boundary
	AtomIDs []string

	// Properties define boundary-specific properties
	Properties map[string]interface{}
}

// BoundaryType defines the type of domain boundary.
type BoundaryType string

const (
	// TransactionalBoundary represents transactional consistency boundaries
	TransactionalBoundary BoundaryType = "transactional"

	// SecurityBoundary represents security boundaries
	SecurityBoundary BoundaryType = "security"

	// ScopeBoundary represents scope boundaries
	ScopeBoundary BoundaryType = "scope"

	// LogicalBoundary represents logical domain boundaries
	LogicalBoundary BoundaryType = "logical"
)

// NewSpace creates a new ATenSpace instance.
func NewSpace(ctx context.Context) (*Space, error) {
	const op = "atenspace.NewSpace"

	s := &Space{
		atoms:       make(map[string]*Atom),
		links:       make([]*Link, 0),
		tensorStore: make(map[string]*Tensor),
		boundaries:  make([]*DomainBoundary, 0),
	}

	return s, nil
}

// AddAtom adds a new atom to the space.
func (s *Space) AddAtom(ctx context.Context, atom *Atom) error {
	const op = "atenspace.(Space).AddAtom"

	if atom == nil {
		return errors.New(ctx, errors.InvalidParameter, op, "atom is nil")
	}
	if atom.ID == "" {
		return errors.New(ctx, errors.InvalidParameter, op, "atom ID is empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	atom.CreatedAt = time.Now()
	if atom.Attributes == nil {
		atom.Attributes = make(map[string]interface{})
	}

	s.atoms[atom.ID] = atom
	return nil
}

// AddLink adds a new link between atoms in the space.
func (s *Space) AddLink(ctx context.Context, link *Link) error {
	const op = "atenspace.(Space).AddLink"

	if link == nil {
		return errors.New(ctx, errors.InvalidParameter, op, "link is nil")
	}
	if link.Source == "" || link.Target == "" {
		return errors.New(ctx, errors.InvalidParameter, op, "link source or target is empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify source and target atoms exist
	if _, ok := s.atoms[link.Source]; !ok {
		return errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("source atom %s not found", link.Source))
	}
	if _, ok := s.atoms[link.Target]; !ok {
		return errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("target atom %s not found", link.Target))
	}

	link.CreatedAt = time.Now()
	s.links = append(s.links, link)
	return nil
}

// AttachTensor attaches an ATen tensor to an atom.
func (s *Space) AttachTensor(ctx context.Context, atomID string, tensor *Tensor) error {
	const op = "atenspace.(Space).AttachTensor"

	if tensor == nil {
		return errors.New(ctx, errors.InvalidParameter, op, "tensor is nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	atom, ok := s.atoms[atomID]
	if !ok {
		return errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("atom %s not found", atomID))
	}

	atom.TensorID = tensor.ID
	s.tensorStore[tensor.ID] = tensor
	return nil
}

// DefineBoundary defines a new domain boundary in the space.
// This is where "Space" is defined by "Boundary" domain model.
func (s *Space) DefineBoundary(ctx context.Context, boundary *DomainBoundary) error {
	const op = "atenspace.(Space).DefineBoundary"

	if boundary == nil {
		return errors.New(ctx, errors.InvalidParameter, op, "boundary is nil")
	}
	if boundary.ID == "" {
		return errors.New(ctx, errors.InvalidParameter, op, "boundary ID is empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if boundary.Properties == nil {
		boundary.Properties = make(map[string]interface{})
	}

	s.boundaries = append(s.boundaries, boundary)
	return nil
}

// GetAtom retrieves an atom by ID.
func (s *Space) GetAtom(ctx context.Context, atomID string) (*Atom, error) {
	const op = "atenspace.(Space).GetAtom"

	s.mu.RLock()
	defer s.mu.RUnlock()

	atom, ok := s.atoms[atomID]
	if !ok {
		return nil, errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("atom %s not found", atomID))
	}

	return atom, nil
}

// GetLinksForAtom retrieves all links connected to an atom.
func (s *Space) GetLinksForAtom(ctx context.Context, atomID string) []*Link {
	s.mu.RLock()
	defer s.mu.RUnlock()

	links := make([]*Link, 0)
	for _, link := range s.links {
		if link.Source == atomID || link.Target == atomID {
			links = append(links, link)
		}
	}

	return links
}

// GetTensor retrieves the tensor for an atom.
func (s *Space) GetTensor(ctx context.Context, atomID string) (*Tensor, error) {
	const op = "atenspace.(Space).GetTensor"

	s.mu.RLock()
	defer s.mu.RUnlock()

	atom, ok := s.atoms[atomID]
	if !ok {
		return nil, errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("atom %s not found", atomID))
	}

	if atom.TensorID == "" {
		return nil, errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("atom %s has no tensor", atomID))
	}

	tensor, ok := s.tensorStore[atom.TensorID]
	if !ok {
		return nil, errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("tensor %s not found", atom.TensorID))
	}

	return tensor, nil
}

// GetBoundaries retrieves all domain boundaries in the space.
func (s *Space) GetBoundaries(ctx context.Context) []*DomainBoundary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	boundaries := make([]*DomainBoundary, len(s.boundaries))
	copy(boundaries, s.boundaries)
	return boundaries
}

// QueryByBoundary queries atoms within a specific domain boundary.
func (s *Space) QueryByBoundary(ctx context.Context, boundaryID string) ([]*Atom, error) {
	const op = "atenspace.(Space).QueryByBoundary"

	s.mu.RLock()
	defer s.mu.RUnlock()

	var boundary *DomainBoundary
	for _, b := range s.boundaries {
		if b.ID == boundaryID {
			boundary = b
			break
		}
	}

	if boundary == nil {
		return nil, errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("boundary %s not found", boundaryID))
	}

	atoms := make([]*Atom, 0, len(boundary.AtomIDs))
	for _, atomID := range boundary.AtomIDs {
		if atom, ok := s.atoms[atomID]; ok {
			atoms = append(atoms, atom)
		}
	}

	return atoms, nil
}

// IntegrateWithBoundary integrates ATenSpace with Boundary's domain model.
// This establishes "Space" as defined by "Boundary".
func (s *Space) IntegrateWithBoundary(ctx context.Context) error {
	const op = "atenspace.(Space).IntegrateWithBoundary"

	// Integration point where Boundary domain model defines the Space
	// All Boundary entities, aggregates, and resources are represented
	// as atoms in the hypergraph with tensor representations
	return nil
}
