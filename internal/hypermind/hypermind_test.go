// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package hypermind

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMultiScopeArchitecture(t *testing.T) {
	ctx := context.Background()

	t.Run("creates architecture successfully", func(t *testing.T) {
		msa, err := NewMultiScopeArchitecture(ctx)
		require.NoError(t, err)
		require.NotNil(t, msa)
		assert.NotNil(t, msa.scopes)
		assert.NotNil(t, msa.peerNetwork)
		assert.Equal(t, 0, len(msa.scopes))
	})
}

func TestMultiScopeArchitecture_RegisterScope(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*MultiScopeArchitecture, *DistributedScope)
		wantErr bool
		errMsg  string
	}{
		{
			name: "register global scope",
			setup: func() (*MultiScopeArchitecture, *DistributedScope) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				scope := &DistributedScope{
					ID:       "global",
					ParentID: "",
					Type:     "global",
					Peers:    []string{},
				}
				return msa, scope
			},
			wantErr: false,
		},
		{
			name: "register org scope",
			setup: func() (*MultiScopeArchitecture, *DistributedScope) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				scope := &DistributedScope{
					ID:       "org-1",
					ParentID: "global",
					Type:     "org",
					Peers:    []string{"peer1", "peer2"},
				}
				return msa, scope
			},
			wantErr: false,
		},
		{
			name: "register project scope",
			setup: func() (*MultiScopeArchitecture, *DistributedScope) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				scope := &DistributedScope{
					ID:       "project-1",
					ParentID: "org-1",
					Type:     "project",
					Peers:    []string{"peer3"},
				}
				return msa, scope
			},
			wantErr: false,
		},
		{
			name: "error on nil scope",
			setup: func() (*MultiScopeArchitecture, *DistributedScope) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				return msa, nil
			},
			wantErr: true,
			errMsg:  "scope is nil",
		},
		{
			name: "error on empty scope ID",
			setup: func() (*MultiScopeArchitecture, *DistributedScope) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				scope := &DistributedScope{
					ID:   "",
					Type: "org",
				}
				return msa, scope
			},
			wantErr: true,
			errMsg:  "scope ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msa, scope := tt.setup()
			err := msa.RegisterScope(ctx, scope)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Contains(t, msa.scopes, scope.ID)
				assert.NotZero(t, msa.scopes[scope.ID].CreatedAt)
				assert.NotZero(t, msa.scopes[scope.ID].UpdatedAt)
				assert.NotNil(t, msa.scopes[scope.ID].State)
			}
		})
	}
}

func TestMultiScopeArchitecture_GetScope(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*MultiScopeArchitecture, string)
		wantErr bool
		errMsg  string
	}{
		{
			name: "get existing scope",
			setup: func() (*MultiScopeArchitecture, string) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				scope := &DistributedScope{
					ID:   "test-scope",
					Type: "org",
				}
				_ = msa.RegisterScope(ctx, scope)
				return msa, "test-scope"
			},
			wantErr: false,
		},
		{
			name: "error on non-existent scope",
			setup: func() (*MultiScopeArchitecture, string) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				return msa, "nonexistent"
			},
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msa, scopeID := tt.setup()
			scope, err := msa.GetScope(ctx, scopeID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, scope)
			} else {
				require.NoError(t, err)
				require.NotNil(t, scope)
				assert.Equal(t, scopeID, scope.ID)
			}
		})
	}
}

func TestMultiScopeArchitecture_PropagateState(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*MultiScopeArchitecture, string, map[string]interface{})
		wantErr bool
		errMsg  string
	}{
		{
			name: "propagate state successfully",
			setup: func() (*MultiScopeArchitecture, string, map[string]interface{}) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				scope := &DistributedScope{
					ID:   "test-scope",
					Type: "org",
				}
				_ = msa.RegisterScope(ctx, scope)
				state := map[string]interface{}{
					"key1": "value1",
					"key2": 42,
				}
				return msa, "test-scope", state
			},
			wantErr: false,
		},
		{
			name: "error on non-existent scope",
			setup: func() (*MultiScopeArchitecture, string, map[string]interface{}) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				state := map[string]interface{}{"key": "value"}
				return msa, "nonexistent", state
			},
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msa, scopeID, state := tt.setup()
			oldTime := time.Now().Add(-1 * time.Second)
			
			err := msa.PropagateState(ctx, scopeID, state)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				scope, _ := msa.GetScope(ctx, scopeID)
				for k, v := range state {
					assert.Equal(t, v, scope.State[k])
				}
				assert.True(t, scope.UpdatedAt.After(oldTime))
			}
		})
	}
}

func TestMultiScopeArchitecture_ConnectPeer(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*MultiScopeArchitecture, *Peer)
		wantErr bool
		errMsg  string
	}{
		{
			name: "connect peer successfully",
			setup: func() (*MultiScopeArchitecture, *Peer) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				peer := &Peer{
					ID:       "peer-1",
					Address:  "192.168.1.1:8080",
					ScopeIDs: []string{"scope-1", "scope-2"},
				}
				return msa, peer
			},
			wantErr: false,
		},
		{
			name: "error on nil peer",
			setup: func() (*MultiScopeArchitecture, *Peer) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				return msa, nil
			},
			wantErr: true,
			errMsg:  "peer is nil",
		},
		{
			name: "error on empty peer ID",
			setup: func() (*MultiScopeArchitecture, *Peer) {
				msa, _ := NewMultiScopeArchitecture(ctx)
				peer := &Peer{
					ID:      "",
					Address: "192.168.1.1:8080",
				}
				return msa, peer
			},
			wantErr: true,
			errMsg:  "peer ID is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msa, peer := tt.setup()
			err := msa.ConnectPeer(ctx, peer)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Contains(t, msa.peerNetwork.activePeers, peer.ID)
				assert.NotZero(t, msa.peerNetwork.activePeers[peer.ID].LastSeen)
			}
		})
	}
}

func TestMultiScopeArchitecture_DiscoverPeers(t *testing.T) {
	ctx := context.Background()

	t.Run("discover peers for scope", func(t *testing.T) {
		msa, _ := NewMultiScopeArchitecture(ctx)
		
		// Connect peers
		peer1 := &Peer{
			ID:       "peer-1",
			Address:  "addr1",
			ScopeIDs: []string{"scope-1"},
		}
		peer2 := &Peer{
			ID:       "peer-2",
			Address:  "addr2",
			ScopeIDs: []string{"scope-1", "scope-2"},
		}
		
		_ = msa.ConnectPeer(ctx, peer1)
		_ = msa.ConnectPeer(ctx, peer2)
		
		// Discover peers for scope-1
		peers, err := msa.DiscoverPeers(ctx, "scope-1")
		require.NoError(t, err)
		assert.Equal(t, 2, len(peers))
	})

	t.Run("discover peers for scope with no peers", func(t *testing.T) {
		msa, _ := NewMultiScopeArchitecture(ctx)
		
		peers, err := msa.DiscoverPeers(ctx, "empty-scope")
		require.NoError(t, err)
		assert.Equal(t, 0, len(peers))
	})
}

func TestMultiScopeArchitecture_GetActivePeers(t *testing.T) {
	ctx := context.Background()

	t.Run("get all active peers", func(t *testing.T) {
		msa, _ := NewMultiScopeArchitecture(ctx)
		
		peers := []*Peer{
			{ID: "peer-1", Address: "addr1", ScopeIDs: []string{"scope-1"}},
			{ID: "peer-2", Address: "addr2", ScopeIDs: []string{"scope-2"}},
			{ID: "peer-3", Address: "addr3", ScopeIDs: []string{"scope-3"}},
		}
		
		for _, p := range peers {
			_ = msa.ConnectPeer(ctx, p)
		}
		
		activePeers := msa.GetActivePeers(ctx)
		assert.Equal(t, 3, len(activePeers))
	})

	t.Run("no active peers", func(t *testing.T) {
		msa, _ := NewMultiScopeArchitecture(ctx)
		
		activePeers := msa.GetActivePeers(ctx)
		assert.Equal(t, 0, len(activePeers))
	})
}

func TestMultiScopeArchitecture_IntegrateWithBoundary(t *testing.T) {
	ctx := context.Background()

	t.Run("integration succeeds", func(t *testing.T) {
		msa, err := NewMultiScopeArchitecture(ctx)
		require.NoError(t, err)

		err = msa.IntegrateWithBoundary(ctx)
		assert.NoError(t, err)
	})
}

func TestDistributedScope_Creation(t *testing.T) {
	scope := &DistributedScope{
		ID:       "test-scope",
		ParentID: "parent-scope",
		Type:     "project",
		Peers:    []string{"peer1", "peer2"},
		State:    map[string]interface{}{"key": "value"},
	}

	assert.Equal(t, "test-scope", scope.ID)
	assert.Equal(t, "parent-scope", scope.ParentID)
	assert.Equal(t, "project", scope.Type)
	assert.Equal(t, 2, len(scope.Peers))
	assert.NotNil(t, scope.State)
}

func TestPeer_Creation(t *testing.T) {
	peer := &Peer{
		ID:       "peer-123",
		Address:  "192.168.1.100:8080",
		ScopeIDs: []string{"scope-1", "scope-2", "scope-3"},
	}

	assert.Equal(t, "peer-123", peer.ID)
	assert.Equal(t, "192.168.1.100:8080", peer.Address)
	assert.Equal(t, 3, len(peer.ScopeIDs))
}

func TestDistributedHashTable_AddAndLookup(t *testing.T) {
	dht := &DistributedHashTable{
		entries: make(map[string][]string),
	}

	t.Run("add and lookup single peer", func(t *testing.T) {
		dht.add("key1", "peer1")
		peers := dht.lookup("key1")
		assert.Equal(t, 1, len(peers))
		assert.Contains(t, peers, "peer1")
	})

	t.Run("add multiple peers to same key", func(t *testing.T) {
		dht.add("key2", "peer1")
		dht.add("key2", "peer2")
		dht.add("key2", "peer3")
		peers := dht.lookup("key2")
		assert.Equal(t, 3, len(peers))
	})

	t.Run("lookup non-existent key", func(t *testing.T) {
		peers := dht.lookup("nonexistent")
		assert.Equal(t, 0, len(peers))
	})
}

func TestPeerNetwork_Creation(t *testing.T) {
	pn := &PeerNetwork{
		activePeers: make(map[string]*Peer),
		dht: &DistributedHashTable{
			entries: make(map[string][]string),
		},
	}

	assert.NotNil(t, pn.activePeers)
	assert.NotNil(t, pn.dht)
	assert.Equal(t, 0, len(pn.activePeers))
}

func TestMultiScopeArchitecture_ComplexScenario(t *testing.T) {
	ctx := context.Background()
	msa, err := NewMultiScopeArchitecture(ctx)
	require.NoError(t, err)

	// Register scope hierarchy
	globalScope := &DistributedScope{ID: "global", Type: "global"}
	orgScope := &DistributedScope{ID: "org-1", ParentID: "global", Type: "org"}
	projectScope := &DistributedScope{ID: "project-1", ParentID: "org-1", Type: "project"}

	require.NoError(t, msa.RegisterScope(ctx, globalScope))
	require.NoError(t, msa.RegisterScope(ctx, orgScope))
	require.NoError(t, msa.RegisterScope(ctx, projectScope))

	// Connect peers
	peer1 := &Peer{ID: "peer-1", Address: "addr1", ScopeIDs: []string{"org-1"}}
	peer2 := &Peer{ID: "peer-2", Address: "addr2", ScopeIDs: []string{"project-1"}}

	require.NoError(t, msa.ConnectPeer(ctx, peer1))
	require.NoError(t, msa.ConnectPeer(ctx, peer2))

	// Propagate state
	state := map[string]interface{}{"status": "active"}
	require.NoError(t, msa.PropagateState(ctx, "org-1", state))

	// Verify
	scope, err := msa.GetScope(ctx, "org-1")
	require.NoError(t, err)
	assert.Equal(t, "active", scope.State["status"])

	activePeers := msa.GetActivePeers(ctx)
	assert.Equal(t, 2, len(activePeers))
}
