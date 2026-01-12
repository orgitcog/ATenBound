# Framework Integration Documentation

This document describes the integration of three advanced frameworks into the ATenBound (Boundary) project:

1. **Tensor Logic** (https://tensor-logic.org/) - The Language of AI
2. **Hypermind** (https://github.com/o9nn/hypermind) - Multi-scope P2P Architecture
3. **ATenSpace** (https://github.com/o9nn/ATenSpace) - ATen Tensors + OpenCog AtomSpace

## Overview

The integration provides a unified framework that enhances Boundary with:
- Tensor equation-based variable system (Tensor Logic)
- Distributed peer-to-peer scope architecture (Hypermind)
- Hypergraph knowledge representation where "Space" is defined by the Boundary domain model (ATenSpace)

## Architecture

### Tensor Logic Integration (`internal/tensorlogic/`)

Tensor Logic is a cutting-edge AI programming paradigm that unifies neural networks and symbolic AI by expressing everything as tensor equations.

**Key Components:**
- `Variable`: Represents tensor logic variables with named indices and shapes
- `TensorEquation`: Defines tensor equations using Einstein summation notation
- `Framework`: Main framework instance managing variables and equations

**Variable Types:**
- `SymbolicType`: Symbolic/logical variables
- `NeuralType`: Neural network variables
- `ProbabilisticType`: Probabilistic inference variables
- `HybridType`: Variables combining multiple types

**Operations:**
- `RegisterVariable`: Register new tensor variables
- `DefineEquation`: Define tensor equations
- `Evaluate`: Evaluate tensor expressions
- `Project`: Tensor projection (reduction along indices)
- `Join`: Tensor join operation (generalized Einstein summation)

**Example:**
```go
ctx := context.Background()
framework, _ := tensorlogic.NewFramework(ctx)

// Register a variable
v := &tensorlogic.Variable{
    Name:    "scope_state",
    Indices: []string{"entity", "property"},
    Shape:   []int{100, 10},
    Type:    tensorlogic.HybridType,
}
framework.RegisterVariable(ctx, v)

// Evaluate
result, _ := framework.Evaluate(ctx, "scope_state")
```

### Hypermind Multi-Scope Architecture (`internal/hypermind/`)

Hypermind provides decentralized, peer-to-peer scope management with distributed consensus.

**Key Components:**
- `MultiScopeArchitecture`: Main architecture managing distributed scopes
- `DistributedScope`: Scope in the distributed architecture
- `PeerNetwork`: P2P network connection management
- `DistributedHashTable`: DHT for peer discovery

**Features:**
- Distributed scope registry
- P2P peer connections between scope nodes
- State propagation across the network
- Peer discovery using DHT

**Example:**
```go
ctx := context.Background()
msa, _ := hypermind.NewMultiScopeArchitecture(ctx)

// Register a distributed scope
scope := &hypermind.DistributedScope{
    ID:       "org-1",
    ParentID: "global",
    Type:     "org",
}
msa.RegisterScope(ctx, scope)

// Connect a peer
peer := &hypermind.Peer{
    ID:       "peer-1",
    Address:  "192.168.1.1:8080",
    ScopeIDs: []string{"org-1"},
}
msa.ConnectPeer(ctx, peer)

// Propagate state
state := map[string]interface{}{"status": "active"}
msa.PropagateState(ctx, "org-1", state)
```

### ATenSpace Integration (`internal/atenspace/`)

ATenSpace combines ATen tensor operations with OpenCog's AtomSpace hypergraph, where "Space" is defined by the Boundary domain model.

**Key Components:**
- `Space`: Main space instance managing atoms and links
- `Atom`: Node in the hypergraph (Boundary domain entities)
- `Link`: Edge in the hypergraph (relationships)
- `Tensor`: ATen tensor representation
- `DomainBoundary`: Boundary definition in the domain model

**Atom Types:**
- `EntityAtom`: Boundary entities
- `AggregateAtom`: Boundary aggregates
- `ResourceAtom`: Boundary resources
- `RelationAtom`: Relationships
- `ConceptAtom`: Abstract concepts

**Link Types:**
- `InheritanceLink`: Is-a relationships
- `MembershipLink`: Part-of relationships
- `ScopeLink`: Scope containment
- `DependencyLink`: Dependencies
- `AssociationLink`: General associations

**Boundary Types:**
- `TransactionalBoundary`: Transactional consistency boundaries
- `SecurityBoundary`: Security boundaries
- `ScopeBoundary`: Scope boundaries
- `LogicalBoundary`: Logical domain boundaries

**Example:**
```go
ctx := context.Background()
space, _ := atenspace.NewSpace(ctx)

// Add an atom
atom := &atenspace.Atom{
    ID:   "scope-1",
    Type: atenspace.AggregateAtom,
    Name: "Organization Scope",
}
space.AddAtom(ctx, atom)

// Attach a tensor
tensor := &atenspace.Tensor{
    ID:     "tensor-1",
    Shape:  []int{10, 10},
    Data:   make([]float64, 100),
    DType:  "float64",
    Device: "cpu",
}
space.AttachTensor(ctx, "scope-1", tensor)

// Define a boundary (Space defined by Boundary)
boundary := &atenspace.DomainBoundary{
    ID:      "org-boundary",
    Name:    "Organization Boundary",
    Type:    atenspace.ScopeBoundary,
    AtomIDs: []string{"scope-1"},
}
space.DefineBoundary(ctx, boundary)
```

## Unified Integration Layer (`internal/integration/`)

The integration package provides a unified framework that ties all three systems together.

**Key Components:**
- `UnifiedFramework`: Integrates all three frameworks
- `ScopeInfo`: Aggregates information from all frameworks

**Example:**
```go
ctx := context.Background()
uf, _ := integration.NewUnifiedFramework(ctx)

// Integrate with Boundary
uf.IntegrateWithBoundary(ctx)

// Create a boundary scope across all frameworks
uf.CreateBoundaryScope(ctx, "org-1", "org")

// Query scope across all frameworks
info, _ := uf.QueryScope(ctx, "org-1")
// info contains: TensorVariable, DistributedScope, and Atom

// Define a domain boundary
uf.DefineDomainBoundary(ctx, "boundary-1", "transactional", 
    []string{"atom-1", "atom-2"})

// Propagate state across all frameworks
state := map[string]interface{}{"status": "active"}
uf.PropagateState(ctx, "org-1", state)
```

## Testing

All components include comprehensive unit tests:

### Tensor Logic Tests (`internal/tensorlogic/tensorlogic_test.go`)
- 22 tests covering all framework operations
- Tests for all variable types and operations
- Error handling and edge cases

### Hypermind Tests (`internal/hypermind/hypermind_test.go`)
- 18 tests covering distributed scope management
- Tests for P2P networking and peer discovery
- State propagation and DHT operations

### ATenSpace Tests (`internal/atenspace/atenspace_test.go`)
- 20 tests covering hypergraph operations
- Tests for all atom and link types
- Boundary definitions and tensor operations

### Integration Tests (`internal/integration/integration_test.go`)
- 11 tests covering unified framework operations
- End-to-end integration scenarios
- Cross-framework interactions

**Running Tests:**
```bash
# Test individual packages
go test ./internal/tensorlogic/...
go test ./internal/hypermind/...
go test ./internal/atenspace/...
go test ./internal/integration/...

# Test all integration packages
go test ./internal/tensorlogic/... ./internal/hypermind/... \
        ./internal/atenspace/... ./internal/integration/...
```

## Integration Points with Boundary

### Variables System
All Boundary variables are now integrated with the Tensor Logic framework, enabling:
- Tensor equation representation of variables
- Einstein summation operations
- Neural, symbolic, and hybrid variable types

### Scope System
Boundary's scope hierarchy is enhanced with Hypermind:
- Distributed P2P scope management
- State propagation across peer networks
- Peer discovery using DHT

### Domain Model
Boundary's domain model defines the ATenSpace:
- Entities, aggregates, and resources as atoms
- Relationships as links in the hypergraph
- Domain boundaries define the Space structure
- Tensor operations on domain objects

## Future Enhancements

Potential future enhancements include:
1. GPU acceleration for tensor operations
2. Advanced peer routing algorithms
3. Distributed tensor storage
4. Real-time synchronization protocols
5. Machine learning model integration
6. Enhanced visualization tools

## References

1. **Tensor Logic**: https://tensor-logic.org/
   - Paper: "Tensor Logic: The Language of AI" (arXiv:2510.12269)
   - PyPI: https://pypi.org/project/tensorlogic/

2. **Hypermind**: https://github.com/o9nn/hypermind
   - Decentralized P2P deployment platform
   - Hyperswarm DHT-based mesh networking

3. **ATenSpace**: https://github.com/o9nn/ATenSpace
   - ATen Tensors + OpenCog AtomSpace
   - Hybrid tensor/hypergraph framework

## License

Copyright (c) HashiCorp, Inc.
SPDX-License-Identifier: BUSL-1.1
