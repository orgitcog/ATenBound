// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package atenspace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSpace(t *testing.T) {
	ctx := context.Background()

	t.Run("creates space successfully", func(t *testing.T) {
		s, err := NewSpace(ctx)
		require.NoError(t, err)
		require.NotNil(t, s)
		assert.NotNil(t, s.atoms)
		assert.NotNil(t, s.links)
		assert.NotNil(t, s.tensorStore)
		assert.NotNil(t, s.boundaries)
		assert.Equal(t, 0, len(s.atoms))
		assert.Equal(t, 0, len(s.links))
	})
}

func TestSpace_AddAtom(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*Space, *Atom)
		wantErr bool
		errMsg  string
	}{
		{
			name: "add entity atom",
			setup: func() (*Space, *Atom) {
				s, _ := NewSpace(ctx)
				atom := &Atom{
					ID:   "atom-1",
					Type: EntityAtom,
					Name: "User Entity",
				}
				return s, atom
			},
			wantErr: false,
		},
		{
			name: "add aggregate atom",
			setup: func() (*Space, *Atom) {
				s, _ := NewSpace(ctx)
				atom := &Atom{
					ID:   "atom-2",
					Type: AggregateAtom,
					Name: "Scope Aggregate",
					Attributes: map[string]interface{}{
						"version": 1,
					},
				}
				return s, atom
			},
			wantErr: false,
		},
		{
			name: "add resource atom",
			setup: func() (*Space, *Atom) {
				s, _ := NewSpace(ctx)
				atom := &Atom{
					ID:   "atom-3",
					Type: ResourceAtom,
					Name: "Target Resource",
				}
				return s, atom
			},
			wantErr: false,
		},
		{
			name: "error on nil atom",
			setup: func() (*Space, *Atom) {
				s, _ := NewSpace(ctx)
				return s, nil
			},
			wantErr: true,
			errMsg:  "atom is nil",
		},
		{
			name: "error on empty atom ID",
			setup: func() (*Space, *Atom) {
				s, _ := NewSpace(ctx)
				atom := &Atom{
					ID:   "",
					Type: EntityAtom,
				}
				return s, atom
			},
			wantErr: true,
			errMsg:  "atom ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, atom := tt.setup()
			err := s.AddAtom(ctx, atom)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Contains(t, s.atoms, atom.ID)
				assert.NotZero(t, s.atoms[atom.ID].CreatedAt)
				assert.NotNil(t, s.atoms[atom.ID].Attributes)
			}
		})
	}
}

func TestSpace_AddLink(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*Space, *Link)
		wantErr bool
		errMsg  string
	}{
		{
			name: "add inheritance link",
			setup: func() (*Space, *Link) {
				s, _ := NewSpace(ctx)
				_ = s.AddAtom(ctx, &Atom{ID: "atom-1", Type: EntityAtom})
				_ = s.AddAtom(ctx, &Atom{ID: "atom-2", Type: EntityAtom})
				link := &Link{
					ID:       "link-1",
					Type:     InheritanceLink,
					Source:   "atom-1",
					Target:   "atom-2",
					Strength: 0.9,
				}
				return s, link
			},
			wantErr: false,
		},
		{
			name: "add scope link",
			setup: func() (*Space, *Link) {
				s, _ := NewSpace(ctx)
				_ = s.AddAtom(ctx, &Atom{ID: "parent", Type: AggregateAtom})
				_ = s.AddAtom(ctx, &Atom{ID: "child", Type: AggregateAtom})
				link := &Link{
					ID:       "link-2",
					Type:     ScopeLink,
					Source:   "parent",
					Target:   "child",
					Strength: 1.0,
				}
				return s, link
			},
			wantErr: false,
		},
		{
			name: "error on nil link",
			setup: func() (*Space, *Link) {
				s, _ := NewSpace(ctx)
				return s, nil
			},
			wantErr: true,
			errMsg:  "link is nil",
		},
		{
			name: "error on empty source",
			setup: func() (*Space, *Link) {
				s, _ := NewSpace(ctx)
				link := &Link{
					ID:     "link-3",
					Source: "",
					Target: "atom-1",
				}
				return s, link
			},
			wantErr: true,
			errMsg:  "link source or target is empty",
		},
		{
			name: "error on non-existent source atom",
			setup: func() (*Space, *Link) {
				s, _ := NewSpace(ctx)
				_ = s.AddAtom(ctx, &Atom{ID: "atom-1", Type: EntityAtom})
				link := &Link{
					ID:     "link-4",
					Source: "nonexistent",
					Target: "atom-1",
				}
				return s, link
			},
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, link := tt.setup()
			err := s.AddLink(ctx, link)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Contains(t, s.links, link)
				assert.NotZero(t, link.CreatedAt)
			}
		})
	}
}

func TestSpace_AttachTensor(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*Space, string, *Tensor)
		wantErr bool
		errMsg  string
	}{
		{
			name: "attach tensor successfully",
			setup: func() (*Space, string, *Tensor) {
				s, _ := NewSpace(ctx)
				_ = s.AddAtom(ctx, &Atom{ID: "atom-1", Type: EntityAtom})
				tensor := &Tensor{
					ID:     "tensor-1",
					Shape:  []int{3, 3},
					Data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
					DType:  "float64",
					Device: "cpu",
				}
				return s, "atom-1", tensor
			},
			wantErr: false,
		},
		{
			name: "error on nil tensor",
			setup: func() (*Space, string, *Tensor) {
				s, _ := NewSpace(ctx)
				_ = s.AddAtom(ctx, &Atom{ID: "atom-1", Type: EntityAtom})
				return s, "atom-1", nil
			},
			wantErr: true,
			errMsg:  "tensor is nil",
		},
		{
			name: "error on non-existent atom",
			setup: func() (*Space, string, *Tensor) {
				s, _ := NewSpace(ctx)
				tensor := &Tensor{ID: "tensor-1"}
				return s, "nonexistent", tensor
			},
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, atomID, tensor := tt.setup()
			err := s.AttachTensor(ctx, atomID, tensor)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				atom, _ := s.GetAtom(ctx, atomID)
				assert.Equal(t, tensor.ID, atom.TensorID)
				assert.Contains(t, s.tensorStore, tensor.ID)
			}
		})
	}
}

func TestSpace_DefineBoundary(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*Space, *DomainBoundary)
		wantErr bool
		errMsg  string
	}{
		{
			name: "define transactional boundary",
			setup: func() (*Space, *DomainBoundary) {
				s, _ := NewSpace(ctx)
				boundary := &DomainBoundary{
					ID:      "boundary-1",
					Name:    "Transaction Boundary",
					Type:    TransactionalBoundary,
					AtomIDs: []string{"atom-1", "atom-2"},
				}
				return s, boundary
			},
			wantErr: false,
		},
		{
			name: "define scope boundary",
			setup: func() (*Space, *DomainBoundary) {
				s, _ := NewSpace(ctx)
				boundary := &DomainBoundary{
					ID:      "boundary-2",
					Name:    "Org Scope",
					Type:    ScopeBoundary,
					AtomIDs: []string{"atom-3"},
					Properties: map[string]interface{}{
						"level": "org",
					},
				}
				return s, boundary
			},
			wantErr: false,
		},
		{
			name: "error on nil boundary",
			setup: func() (*Space, *DomainBoundary) {
				s, _ := NewSpace(ctx)
				return s, nil
			},
			wantErr: true,
			errMsg:  "boundary is nil",
		},
		{
			name: "error on empty boundary ID",
			setup: func() (*Space, *DomainBoundary) {
				s, _ := NewSpace(ctx)
				boundary := &DomainBoundary{
					ID:   "",
					Name: "Test",
				}
				return s, boundary
			},
			wantErr: true,
			errMsg:  "boundary ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, boundary := tt.setup()
			err := s.DefineBoundary(ctx, boundary)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Contains(t, s.boundaries, boundary)
				assert.NotNil(t, boundary.Properties)
			}
		})
	}
}

func TestSpace_GetAtom(t *testing.T) {
	ctx := context.Background()

	t.Run("get existing atom", func(t *testing.T) {
		s, _ := NewSpace(ctx)
		atom := &Atom{ID: "atom-1", Type: EntityAtom, Name: "Test"}
		_ = s.AddAtom(ctx, atom)

		result, err := s.GetAtom(ctx, "atom-1")
		require.NoError(t, err)
		assert.Equal(t, "atom-1", result.ID)
		assert.Equal(t, "Test", result.Name)
	})

	t.Run("error on non-existent atom", func(t *testing.T) {
		s, _ := NewSpace(ctx)

		result, err := s.GetAtom(ctx, "nonexistent")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestSpace_GetLinksForAtom(t *testing.T) {
	ctx := context.Background()

	t.Run("get links for atom", func(t *testing.T) {
		s, _ := NewSpace(ctx)
		_ = s.AddAtom(ctx, &Atom{ID: "atom-1", Type: EntityAtom})
		_ = s.AddAtom(ctx, &Atom{ID: "atom-2", Type: EntityAtom})
		_ = s.AddAtom(ctx, &Atom{ID: "atom-3", Type: EntityAtom})

		_ = s.AddLink(ctx, &Link{ID: "link-1", Source: "atom-1", Target: "atom-2", Type: InheritanceLink})
		_ = s.AddLink(ctx, &Link{ID: "link-2", Source: "atom-2", Target: "atom-3", Type: AssociationLink})

		links := s.GetLinksForAtom(ctx, "atom-2")
		assert.Equal(t, 2, len(links))
	})

	t.Run("no links for atom", func(t *testing.T) {
		s, _ := NewSpace(ctx)
		_ = s.AddAtom(ctx, &Atom{ID: "atom-1", Type: EntityAtom})

		links := s.GetLinksForAtom(ctx, "atom-1")
		assert.Equal(t, 0, len(links))
	})
}

func TestSpace_GetTensor(t *testing.T) {
	ctx := context.Background()

	t.Run("get tensor for atom", func(t *testing.T) {
		s, _ := NewSpace(ctx)
		_ = s.AddAtom(ctx, &Atom{ID: "atom-1", Type: EntityAtom})
		tensor := &Tensor{
			ID:    "tensor-1",
			Shape: []int{2, 2},
			Data:  []float64{1, 2, 3, 4},
		}
		_ = s.AttachTensor(ctx, "atom-1", tensor)

		result, err := s.GetTensor(ctx, "atom-1")
		require.NoError(t, err)
		assert.Equal(t, "tensor-1", result.ID)
		assert.Equal(t, []int{2, 2}, result.Shape)
	})

	t.Run("error on atom without tensor", func(t *testing.T) {
		s, _ := NewSpace(ctx)
		_ = s.AddAtom(ctx, &Atom{ID: "atom-1", Type: EntityAtom})

		result, err := s.GetTensor(ctx, "atom-1")
		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestSpace_GetBoundaries(t *testing.T) {
	ctx := context.Background()

	t.Run("get all boundaries", func(t *testing.T) {
		s, _ := NewSpace(ctx)
		_ = s.DefineBoundary(ctx, &DomainBoundary{ID: "b1", Type: TransactionalBoundary})
		_ = s.DefineBoundary(ctx, &DomainBoundary{ID: "b2", Type: ScopeBoundary})

		boundaries := s.GetBoundaries(ctx)
		assert.Equal(t, 2, len(boundaries))
	})

	t.Run("no boundaries", func(t *testing.T) {
		s, _ := NewSpace(ctx)

		boundaries := s.GetBoundaries(ctx)
		assert.Equal(t, 0, len(boundaries))
	})
}

func TestSpace_QueryByBoundary(t *testing.T) {
	ctx := context.Background()

	t.Run("query atoms by boundary", func(t *testing.T) {
		s, _ := NewSpace(ctx)
		_ = s.AddAtom(ctx, &Atom{ID: "atom-1", Type: EntityAtom})
		_ = s.AddAtom(ctx, &Atom{ID: "atom-2", Type: EntityAtom})
		_ = s.DefineBoundary(ctx, &DomainBoundary{
			ID:      "boundary-1",
			Type:    TransactionalBoundary,
			AtomIDs: []string{"atom-1", "atom-2"},
		})

		atoms, err := s.QueryByBoundary(ctx, "boundary-1")
		require.NoError(t, err)
		assert.Equal(t, 2, len(atoms))
	})

	t.Run("error on non-existent boundary", func(t *testing.T) {
		s, _ := NewSpace(ctx)

		atoms, err := s.QueryByBoundary(ctx, "nonexistent")
		require.Error(t, err)
		assert.Nil(t, atoms)
	})
}

func TestSpace_IntegrateWithBoundary(t *testing.T) {
	ctx := context.Background()

	t.Run("integration succeeds", func(t *testing.T) {
		s, err := NewSpace(ctx)
		require.NoError(t, err)

		err = s.IntegrateWithBoundary(ctx)
		assert.NoError(t, err)
	})
}

func TestAtomTypes(t *testing.T) {
	tests := []struct {
		name     string
		atomType AtomType
		expected string
	}{
		{"entity atom", EntityAtom, "entity"},
		{"aggregate atom", AggregateAtom, "aggregate"},
		{"resource atom", ResourceAtom, "resource"},
		{"relation atom", RelationAtom, "relation"},
		{"concept atom", ConceptAtom, "concept"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.atomType))
		})
	}
}

func TestLinkTypes(t *testing.T) {
	tests := []struct {
		name     string
		linkType LinkType
		expected string
	}{
		{"inheritance link", InheritanceLink, "inheritance"},
		{"membership link", MembershipLink, "membership"},
		{"scope link", ScopeLink, "scope"},
		{"dependency link", DependencyLink, "dependency"},
		{"association link", AssociationLink, "association"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.linkType))
		})
	}
}

func TestBoundaryTypes(t *testing.T) {
	tests := []struct {
		name         string
		boundaryType BoundaryType
		expected     string
	}{
		{"transactional boundary", TransactionalBoundary, "transactional"},
		{"security boundary", SecurityBoundary, "security"},
		{"scope boundary", ScopeBoundary, "scope"},
		{"logical boundary", LogicalBoundary, "logical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.boundaryType))
		})
	}
}

func TestSpace_ComplexScenario(t *testing.T) {
	ctx := context.Background()
	s, err := NewSpace(ctx)
	require.NoError(t, err)

	// Create atom hierarchy
	globalAtom := &Atom{ID: "global", Type: AggregateAtom, Name: "Global Scope"}
	orgAtom := &Atom{ID: "org-1", Type: AggregateAtom, Name: "Organization"}
	projectAtom := &Atom{ID: "project-1", Type: AggregateAtom, Name: "Project"}

	require.NoError(t, s.AddAtom(ctx, globalAtom))
	require.NoError(t, s.AddAtom(ctx, orgAtom))
	require.NoError(t, s.AddAtom(ctx, projectAtom))

	// Create links
	require.NoError(t, s.AddLink(ctx, &Link{
		ID:       "link-1",
		Type:     ScopeLink,
		Source:   "global",
		Target:   "org-1",
		Strength: 1.0,
	}))
	require.NoError(t, s.AddLink(ctx, &Link{
		ID:       "link-2",
		Type:     ScopeLink,
		Source:   "org-1",
		Target:   "project-1",
		Strength: 1.0,
	}))

	// Attach tensors
	tensor := &Tensor{
		ID:     "tensor-1",
		Shape:  []int{10},
		Data:   make([]float64, 10),
		DType:  "float64",
		Device: "cpu",
	}
	require.NoError(t, s.AttachTensor(ctx, "org-1", tensor))

	// Define boundary
	boundary := &DomainBoundary{
		ID:      "org-boundary",
		Name:    "Organization Boundary",
		Type:    ScopeBoundary,
		AtomIDs: []string{"org-1", "project-1"},
	}
	require.NoError(t, s.DefineBoundary(ctx, boundary))

	// Verify
	atoms, err := s.QueryByBoundary(ctx, "org-boundary")
	require.NoError(t, err)
	assert.Equal(t, 2, len(atoms))

	orgLinks := s.GetLinksForAtom(ctx, "org-1")
	assert.Equal(t, 2, len(orgLinks))

	retrievedTensor, err := s.GetTensor(ctx, "org-1")
	require.NoError(t, err)
	assert.Equal(t, "tensor-1", retrievedTensor.ID)
}
