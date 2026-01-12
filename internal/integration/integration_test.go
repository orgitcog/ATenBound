// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package integration

import (
	"context"
	"testing"

	"github.com/hashicorp/boundary/internal/atenspace"
	"github.com/hashicorp/boundary/internal/hypermind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUnifiedFramework(t *testing.T) {
	ctx := context.Background()

	t.Run("creates unified framework successfully", func(t *testing.T) {
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)
		require.NotNil(t, uf)
		assert.NotNil(t, uf.TensorLogic)
		assert.NotNil(t, uf.Hypermind)
		assert.NotNil(t, uf.ATenSpace)
	})
}

func TestUnifiedFramework_IntegrateWithBoundary(t *testing.T) {
	ctx := context.Background()

	t.Run("integration succeeds for all frameworks", func(t *testing.T) {
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)

		err = uf.IntegrateWithBoundary(ctx)
		assert.NoError(t, err)
	})
}

func TestUnifiedFramework_CreateBoundaryScope(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		scopeID   string
		scopeType string
		wantErr   bool
	}{
		{
			name:      "create global scope",
			scopeID:   "global",
			scopeType: "global",
			wantErr:   false,
		},
		{
			name:      "create org scope",
			scopeID:   "org-123",
			scopeType: "org",
			wantErr:   false,
		},
		{
			name:      "create project scope",
			scopeID:   "project-456",
			scopeType: "project",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uf, err := NewUnifiedFramework(ctx)
			require.NoError(t, err)

			err = uf.CreateBoundaryScope(ctx, tt.scopeID, tt.scopeType)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify scope exists in all three frameworks
				// 1. Tensor Logic
				_, err := uf.TensorLogic.Evaluate(ctx, tt.scopeID)
				assert.NoError(t, err)

				// 2. Hypermind
				_, err = uf.Hypermind.GetScope(ctx, tt.scopeID)
				assert.NoError(t, err)

				// 3. ATenSpace
				atom, err := uf.ATenSpace.GetAtom(ctx, tt.scopeID)
				assert.NoError(t, err)
				assert.NotNil(t, atom)

				// Verify tensor is attached
				_, err = uf.ATenSpace.GetTensor(ctx, tt.scopeID)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnifiedFramework_QueryScope(t *testing.T) {
	ctx := context.Background()

	t.Run("query scope across all frameworks", func(t *testing.T) {
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)

		scopeID := "test-scope"
		err = uf.CreateBoundaryScope(ctx, scopeID, "org")
		require.NoError(t, err)

		info, err := uf.QueryScope(ctx, scopeID)
		require.NoError(t, err)
		require.NotNil(t, info)

		assert.Equal(t, scopeID, info.ID)
		assert.NotNil(t, info.TensorVariable)
		assert.NotNil(t, info.DistributedScope)
		assert.NotNil(t, info.Atom)
	})

	t.Run("query non-existent scope", func(t *testing.T) {
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)

		info, err := uf.QueryScope(ctx, "nonexistent")
		require.NoError(t, err)
		require.NotNil(t, info)

		// Should return info with nil components
		assert.Equal(t, "nonexistent", info.ID)
		assert.Nil(t, info.TensorVariable)
		assert.Nil(t, info.DistributedScope)
		assert.Nil(t, info.Atom)
	})
}

func TestUnifiedFramework_DefineDomainBoundary(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		boundaryID   string
		boundaryType string
		atomIDs      []string
		wantErr      bool
	}{
		{
			name:         "define transactional boundary",
			boundaryID:   "boundary-1",
			boundaryType: "transactional",
			atomIDs:      []string{"atom-1", "atom-2"},
			wantErr:      false,
		},
		{
			name:         "define scope boundary",
			boundaryID:   "boundary-2",
			boundaryType: "scope",
			atomIDs:      []string{"scope-1"},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uf, err := NewUnifiedFramework(ctx)
			require.NoError(t, err)

			err = uf.DefineDomainBoundary(ctx, tt.boundaryID, tt.boundaryType, tt.atomIDs)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify boundary exists in ATenSpace
				boundaries := uf.ATenSpace.GetBoundaries(ctx)
				found := false
				for _, b := range boundaries {
					if b.ID == tt.boundaryID {
						found = true
						assert.Equal(t, tt.boundaryType, string(b.Type))
						assert.Equal(t, tt.atomIDs, b.AtomIDs)
						break
					}
				}
				assert.True(t, found, "boundary not found")
			}
		})
	}
}

func TestUnifiedFramework_PropagateState(t *testing.T) {
	ctx := context.Background()

	t.Run("propagate state successfully", func(t *testing.T) {
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)

		scopeID := "test-scope"
		err = uf.CreateBoundaryScope(ctx, scopeID, "org")
		require.NoError(t, err)

		state := map[string]interface{}{
			"status":  "active",
			"version": 1,
		}

		err = uf.PropagateState(ctx, scopeID, state)
		require.NoError(t, err)

		// Verify state in Hypermind
		distScope, err := uf.Hypermind.GetScope(ctx, scopeID)
		require.NoError(t, err)
		assert.Equal(t, "active", distScope.State["status"])
		assert.Equal(t, 1, distScope.State["version"])

		// Verify state in ATenSpace
		atom, err := uf.ATenSpace.GetAtom(ctx, scopeID)
		require.NoError(t, err)
		assert.Equal(t, "active", atom.Attributes["status"])
		assert.Equal(t, 1, atom.Attributes["version"])
	})

	t.Run("error on non-existent scope", func(t *testing.T) {
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)

		state := map[string]interface{}{"key": "value"}
		err = uf.PropagateState(ctx, "nonexistent", state)
		require.Error(t, err)
	})
}

func TestUnifiedFramework_ComplexScenario(t *testing.T) {
	ctx := context.Background()

	t.Run("full integration scenario", func(t *testing.T) {
		// Create unified framework
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)

		// Integrate with Boundary
		err = uf.IntegrateWithBoundary(ctx)
		require.NoError(t, err)

		// Create scope hierarchy
		globalScope := "global"
		orgScope := "org-1"
		projectScope := "project-1"

		require.NoError(t, uf.CreateBoundaryScope(ctx, globalScope, "global"))
		require.NoError(t, uf.CreateBoundaryScope(ctx, orgScope, "org"))
		require.NoError(t, uf.CreateBoundaryScope(ctx, projectScope, "project"))

		// Define domain boundary
		err = uf.DefineDomainBoundary(ctx, "org-boundary", "scope", []string{orgScope, projectScope})
		require.NoError(t, err)

		// Propagate state
		state := map[string]interface{}{
			"environment": "production",
			"region":      "us-west",
		}
		err = uf.PropagateState(ctx, orgScope, state)
		require.NoError(t, err)

		// Query scope
		info, err := uf.QueryScope(ctx, orgScope)
		require.NoError(t, err)
		require.NotNil(t, info)

		// Verify all components
		assert.Equal(t, orgScope, info.ID)
		assert.NotNil(t, info.TensorVariable)
		assert.NotNil(t, info.DistributedScope)
		assert.NotNil(t, info.Atom)

		// Verify state propagation
		assert.Equal(t, "production", info.DistributedScope.State["environment"])
		assert.Equal(t, "production", info.Atom.Attributes["environment"])

		// Query boundary
		atoms, err := uf.ATenSpace.QueryByBoundary(ctx, "org-boundary")
		require.NoError(t, err)
		assert.Equal(t, 2, len(atoms))
	})
}

func TestUnifiedFramework_TensorLogicIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("tensor operations on boundary scopes", func(t *testing.T) {
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)

		// Create scopes
		scope1 := "scope-1"
		scope2 := "scope-2"
		require.NoError(t, uf.CreateBoundaryScope(ctx, scope1, "org"))
		require.NoError(t, uf.CreateBoundaryScope(ctx, scope2, "org"))

		// Perform tensor operations
		v1, err := uf.TensorLogic.Evaluate(ctx, scope1)
		require.NoError(t, err)

		v2, err := uf.TensorLogic.Evaluate(ctx, scope2)
		require.NoError(t, err)

		// Join operation
		result, err := uf.TensorLogic.Join(ctx, v1, v2)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestUnifiedFramework_HypermindIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("peer network for scopes", func(t *testing.T) {
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)

		// Create scope
		scopeID := "distributed-scope"
		require.NoError(t, uf.CreateBoundaryScope(ctx, scopeID, "org"))

		// Connect peers to the scope
		peer1 := &hypermind.Peer{
			ID:       "peer-1",
			Address:  "192.168.1.1:8080",
			ScopeIDs: []string{scopeID},
		}
		peer2 := &hypermind.Peer{
			ID:       "peer-2",
			Address:  "192.168.1.2:8080",
			ScopeIDs: []string{scopeID},
		}

		require.NoError(t, uf.Hypermind.ConnectPeer(ctx, peer1))
		require.NoError(t, uf.Hypermind.ConnectPeer(ctx, peer2))

		// Discover peers
		peers, err := uf.Hypermind.DiscoverPeers(ctx, scopeID)
		require.NoError(t, err)
		assert.Equal(t, 2, len(peers))
	})
}

func TestUnifiedFramework_ATenSpaceIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("hypergraph representation of boundary domain", func(t *testing.T) {
		uf, err := NewUnifiedFramework(ctx)
		require.NoError(t, err)

		// Create scopes
		parent := "parent-scope"
		child := "child-scope"
		require.NoError(t, uf.CreateBoundaryScope(ctx, parent, "org"))
		require.NoError(t, uf.CreateBoundaryScope(ctx, child, "project"))

		// Create link between scopes
		link := &atenspace.Link{
			ID:       "parent-child-link",
			Type:     atenspace.ScopeLink,
			Source:   parent,
			Target:   child,
			Strength: 1.0,
		}
		require.NoError(t, uf.ATenSpace.AddLink(ctx, link))

		// Query links
		links := uf.ATenSpace.GetLinksForAtom(ctx, parent)
		assert.Equal(t, 1, len(links))
		assert.Equal(t, child, links[0].Target)
	})
}

func TestScopeInfo_Structure(t *testing.T) {
	info := &ScopeInfo{
		ID: "test-scope",
	}

	assert.Equal(t, "test-scope", info.ID)
	assert.Nil(t, info.TensorVariable)
	assert.Nil(t, info.DistributedScope)
	assert.Nil(t, info.Atom)
}
