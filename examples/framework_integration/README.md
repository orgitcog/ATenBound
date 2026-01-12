# Framework Integration Example

This example demonstrates the integration of Tensor Logic, Hypermind, and ATenSpace frameworks in Boundary.

## Running the Example

```bash
cd examples/framework_integration
go run main.go
```

## What This Example Demonstrates

### 1. Individual Framework Usage

**Tensor Logic:**
- Creating tensor variables with named indices
- Defining tensor equations
- Evaluating tensor expressions
- Supporting symbolic, neural, and hybrid variable types

**Hypermind:**
- Creating distributed scope hierarchies
- Connecting P2P peers to scopes
- Propagating state across the network
- Peer discovery using DHT

**ATenSpace:**
- Creating atoms (nodes) in the hypergraph
- Establishing links (edges) between atoms
- Attaching tensors to atoms
- Defining domain boundaries where Space is defined by Boundary

### 2. Unified Integration

The unified framework demonstrates:
- Creating scopes that exist across all three frameworks simultaneously
- Querying scope information from all frameworks at once
- State propagation that updates all frameworks
- Domain boundary definitions that span the integrated system

## Expected Output

When you run the example, you'll see:

```
=== Example 1: Individual Framework Usage ===

--- Tensor Logic Framework ---
Evaluated tensor variable: organization_scope (type: hybrid)
Variables registered: 2
Equations defined: 1

--- Hypermind Multi-Scope Architecture ---
Distributed scopes: 3 (global, org, project)
Connected peers: 2
Peers for org-acme: 2

--- ATenSpace (Space defined by Boundary) ---
Atoms in space: 3 (global, org, user)
Links in space: 2 (scope, membership)
Domain boundaries: 1
Atoms in org boundary: 2
Links for org-acme: 2

=== Example 2: Unified Integration ===

--- Unified Framework Integration ---
✓ Integrated Tensor Logic, Hypermind, and ATenSpace with Boundary
✓ Created scope 'global' across all frameworks
✓ Created scope 'org-acme' across all frameworks
✓ Created scope 'project-alpha' across all frameworks
✓ Defined domain boundary
✓ Propagated state across P2P network

--- Scope Query Results ---
Scope ID: org-acme
✓ Tensor Logic: Variable 'org-acme' (type: hybrid)
✓ Hypermind: Distributed scope (type: org, state entries: 3)
✓ ATenSpace: Atom (type: aggregate, has tensor: true)

✓ All frameworks successfully integrated and working together!
```

## Key Concepts

### Tensor Logic Integration
All Boundary variables can now be represented as tensor equations, enabling:
- Mathematical operations on scope hierarchies
- Neural network integration for AI-powered access control
- Symbolic reasoning about permissions and relationships

### Hypermind Multi-Scope Architecture
Boundary scopes become distributed P2P entities:
- Scopes can be replicated across multiple nodes
- State changes propagate through the P2P network
- Peers can discover each other using DHT

### ATenSpace Domain Model
The Boundary domain model defines the Space:
- Entities, aggregates, and resources become atoms in a hypergraph
- Relationships become links with strengths
- Tensors can be attached to any domain object
- Domain boundaries define the structure of the Space

## Further Reading

See [FRAMEWORK_INTEGRATION.md](../../FRAMEWORK_INTEGRATION.md) for complete documentation.
