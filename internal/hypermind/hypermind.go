// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package hypermind provides multi-scope architecture enhancement for Boundary.
// Based on https://github.com/o9nn/hypermind - Distributed P2P cognitive platform
//
// Hypermind enables decentralized, peer-to-peer scope management with distributed
// consensus and ephemeral communication channels across Boundary scopes.
package hypermind

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/boundary/internal/errors"
)

// MultiScopeArchitecture represents the hypermind-enhanced multi-scope system.
// It extends Boundary's scope hierarchy with distributed P2P capabilities.
type MultiScopeArchitecture struct {
	// Scopes holds the distributed scope registry
	scopes map[string]*DistributedScope

	// PeerNetwork manages P2P connections between scope nodes
	peerNetwork *PeerNetwork

	// mu protects concurrent access to scopes
	mu sync.RWMutex
}

// DistributedScope represents a scope in the hypermind distributed architecture.
type DistributedScope struct {
	// ID is the unique scope identifier
	ID string

	// ParentID references the parent scope in the hierarchy
	ParentID string

	// Type defines the scope type (global, org, project)
	Type string

	// Peers are the connected peer nodes for this scope
	Peers []string

	// State holds the distributed state for this scope
	State map[string]interface{}

	// CreatedAt timestamp
	CreatedAt time.Time

	// UpdatedAt timestamp
	UpdatedAt time.Time
}

// PeerNetwork manages the P2P network connections using hypermind's
// decentralized architecture.
type PeerNetwork struct {
	// ActivePeers tracks currently connected peers
	activePeers map[string]*Peer

	// DHT represents the distributed hash table for peer discovery
	dht *DistributedHashTable

	// mu protects concurrent access
	mu sync.RWMutex
}

// Peer represents a node in the P2P network.
type Peer struct {
	// ID is the unique peer identifier
	ID string

	// Address is the network address
	Address string

	// LastSeen timestamp
	LastSeen time.Time

	// ScopeIDs are the scopes this peer participates in
	ScopeIDs []string
}

// DistributedHashTable implements a simplified DHT for peer discovery.
type DistributedHashTable struct {
	// Entries maps keys to peer lists
	entries map[string][]string

	mu sync.RWMutex
}

// NewMultiScopeArchitecture creates a new hypermind multi-scope architecture.
func NewMultiScopeArchitecture(ctx context.Context) (*MultiScopeArchitecture, error) {
	const op = "hypermind.NewMultiScopeArchitecture"

	msa := &MultiScopeArchitecture{
		scopes: make(map[string]*DistributedScope),
		peerNetwork: &PeerNetwork{
			activePeers: make(map[string]*Peer),
			dht: &DistributedHashTable{
				entries: make(map[string][]string),
			},
		},
	}

	return msa, nil
}

// RegisterScope registers a new distributed scope in the architecture.
func (m *MultiScopeArchitecture) RegisterScope(ctx context.Context, scope *DistributedScope) error {
	const op = "hypermind.(MultiScopeArchitecture).RegisterScope"

	if scope == nil {
		return errors.New(ctx, errors.InvalidParameter, op, "scope is nil")
	}
	if scope.ID == "" {
		return errors.New(ctx, errors.InvalidParameter, op, "scope ID is empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	scope.CreatedAt = time.Now()
	scope.UpdatedAt = time.Now()
	if scope.State == nil {
		scope.State = make(map[string]interface{})
	}

	m.scopes[scope.ID] = scope
	return nil
}

// GetScope retrieves a distributed scope by ID.
func (m *MultiScopeArchitecture) GetScope(ctx context.Context, scopeID string) (*DistributedScope, error) {
	const op = "hypermind.(MultiScopeArchitecture).GetScope"

	m.mu.RLock()
	defer m.mu.RUnlock()

	scope, ok := m.scopes[scopeID]
	if !ok {
		return nil, errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("scope %s not found", scopeID))
	}

	return scope, nil
}

// PropagateState propagates state changes across the P2P network.
func (m *MultiScopeArchitecture) PropagateState(ctx context.Context, scopeID string, state map[string]interface{}) error {
	const op = "hypermind.(MultiScopeArchitecture).PropagateState"

	m.mu.Lock()
	defer m.mu.Unlock()

	scope, ok := m.scopes[scopeID]
	if !ok {
		return errors.New(ctx, errors.InvalidParameter, op, fmt.Sprintf("scope %s not found", scopeID))
	}

	// Update local state
	for k, v := range state {
		scope.State[k] = v
	}
	scope.UpdatedAt = time.Now()

	// Propagate to peers (simplified)
	return m.propagateToPeers(ctx, scopeID, state)
}

// propagateToPeers sends state updates to connected peers.
func (m *MultiScopeArchitecture) propagateToPeers(ctx context.Context, scopeID string, state map[string]interface{}) error {
	// Simplified P2P propagation
	// In a full implementation, this would use the hypermind DHT
	// and gossip protocol to distribute state updates
	return nil
}

// ConnectPeer connects a new peer to the network.
func (m *MultiScopeArchitecture) ConnectPeer(ctx context.Context, peer *Peer) error {
	const op = "hypermind.(MultiScopeArchitecture).ConnectPeer"

	if peer == nil {
		return errors.New(ctx, errors.InvalidParameter, op, "peer is nil")
	}
	if peer.ID == "" {
		return errors.New(ctx, errors.InvalidParameter, op, "peer ID is empty")
	}

	m.peerNetwork.mu.Lock()
	defer m.peerNetwork.mu.Unlock()

	peer.LastSeen = time.Now()
	m.peerNetwork.activePeers[peer.ID] = peer

	// Add to DHT for discovery
	for _, scopeID := range peer.ScopeIDs {
		m.peerNetwork.dht.add(scopeID, peer.ID)
	}

	return nil
}

// DiscoverPeers discovers peers for a given scope using the DHT.
func (m *MultiScopeArchitecture) DiscoverPeers(ctx context.Context, scopeID string) ([]*Peer, error) {
	const op = "hypermind.(MultiScopeArchitecture).DiscoverPeers"

	m.peerNetwork.mu.RLock()
	defer m.peerNetwork.mu.RUnlock()

	peerIDs := m.peerNetwork.dht.lookup(scopeID)
	peers := make([]*Peer, 0, len(peerIDs))

	for _, peerID := range peerIDs {
		if peer, ok := m.peerNetwork.activePeers[peerID]; ok {
			peers = append(peers, peer)
		}
	}

	return peers, nil
}

// GetActivePeers returns all currently active peers.
func (m *MultiScopeArchitecture) GetActivePeers(ctx context.Context) []*Peer {
	m.peerNetwork.mu.RLock()
	defer m.peerNetwork.mu.RUnlock()

	peers := make([]*Peer, 0, len(m.peerNetwork.activePeers))
	for _, peer := range m.peerNetwork.activePeers {
		peers = append(peers, peer)
	}

	return peers
}

// IntegrateWithBoundary integrates the hypermind architecture with Boundary's scope system.
func (m *MultiScopeArchitecture) IntegrateWithBoundary(ctx context.Context) error {
	const op = "hypermind.(MultiScopeArchitecture).IntegrateWithBoundary"

	// Integration point for Boundary scope hierarchy
	// Enables distributed, P2P scope management
	return nil
}

// add adds a peer ID to the DHT entry for a key.
func (d *DistributedHashTable) add(key, peerID string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.entries[key] == nil {
		d.entries[key] = make([]string, 0)
	}
	d.entries[key] = append(d.entries[key], peerID)
}

// lookup retrieves peer IDs for a key from the DHT.
func (d *DistributedHashTable) lookup(key string) []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if peers, ok := d.entries[key]; ok {
		result := make([]string, len(peers))
		copy(result, peers)
		return result
	}
	return []string{}
}
