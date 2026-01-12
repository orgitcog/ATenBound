// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package tensorlogic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFramework(t *testing.T) {
	ctx := context.Background()

	t.Run("creates framework successfully", func(t *testing.T) {
		f, err := NewFramework(ctx)
		require.NoError(t, err)
		require.NotNil(t, f)
		assert.NotNil(t, f.Variables)
		assert.NotNil(t, f.Equations)
		assert.Equal(t, 0, len(f.Variables))
		assert.Equal(t, 0, len(f.Equations))
	})
}

func TestFramework_RegisterVariable(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*Framework, *Variable)
		wantErr bool
		errMsg  string
	}{
		{
			name: "register symbolic variable",
			setup: func() (*Framework, *Variable) {
				f, _ := NewFramework(ctx)
				v := &Variable{
					Name:    "x",
					Indices: []string{"i", "j"},
					Shape:   []int{3, 3},
					Data:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
					Type:    SymbolicType,
				}
				return f, v
			},
			wantErr: false,
		},
		{
			name: "register neural variable",
			setup: func() (*Framework, *Variable) {
				f, _ := NewFramework(ctx)
				v := &Variable{
					Name:    "weight",
					Indices: []string{"in", "out"},
					Shape:   []int{128, 64},
					Type:    NeuralType,
				}
				return f, v
			},
			wantErr: false,
		},
		{
			name: "register probabilistic variable",
			setup: func() (*Framework, *Variable) {
				f, _ := NewFramework(ctx)
				v := &Variable{
					Name:    "prob",
					Indices: []string{"state"},
					Shape:   []int{10},
					Type:    ProbabilisticType,
				}
				return f, v
			},
			wantErr: false,
		},
		{
			name: "error on nil variable",
			setup: func() (*Framework, *Variable) {
				f, _ := NewFramework(ctx)
				return f, nil
			},
			wantErr: true,
			errMsg:  "variable is nil",
		},
		{
			name: "error on empty name",
			setup: func() (*Framework, *Variable) {
				f, _ := NewFramework(ctx)
				v := &Variable{
					Name: "",
					Type: SymbolicType,
				}
				return f, v
			},
			wantErr: true,
			errMsg:  "variable name is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, v := tt.setup()
			err := f.RegisterVariable(ctx, v)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Contains(t, f.Variables, v.Name)
				assert.Equal(t, v, f.Variables[v.Name])
			}
		})
	}
}

func TestFramework_DefineEquation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*Framework, *TensorEquation)
		wantErr bool
		errMsg  string
	}{
		{
			name: "define simple equation",
			setup: func() (*Framework, *TensorEquation) {
				f, _ := NewFramework(ctx)
				eq := &TensorEquation{
					Left: Variable{
						Name:    "result",
						Indices: []string{"i", "k"},
					},
					Right:     "A_ij * B_jk",
					Operation: "join",
				}
				return f, eq
			},
			wantErr: false,
		},
		{
			name: "error on nil equation",
			setup: func() (*Framework, *TensorEquation) {
				f, _ := NewFramework(ctx)
				return f, nil
			},
			wantErr: true,
			errMsg:  "equation is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, eq := tt.setup()
			err := f.DefineEquation(ctx, eq)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.Contains(t, f.Equations, eq)
			}
		})
	}
}

func TestFramework_Evaluate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*Framework, string)
		wantErr bool
		errMsg  string
	}{
		{
			name: "evaluate existing variable",
			setup: func() (*Framework, string) {
				f, _ := NewFramework(ctx)
				v := &Variable{
					Name:    "test",
					Indices: []string{"i"},
					Shape:   []int{5},
					Data:    []float64{1, 2, 3, 4, 5},
					Type:    SymbolicType,
				}
				_ = f.RegisterVariable(ctx, v)
				return f, "test"
			},
			wantErr: false,
		},
		{
			name: "error on non-existent variable",
			setup: func() (*Framework, string) {
				f, _ := NewFramework(ctx)
				return f, "nonexistent"
			},
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, varName := tt.setup()
			result, err := f.Evaluate(ctx, varName)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, varName, result.Name)
				assert.NotNil(t, result.Data)
			}
		})
	}
}

func TestFramework_Project(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*Framework, *Variable, []string)
		wantErr bool
		errMsg  string
	}{
		{
			name: "project variable successfully",
			setup: func() (*Framework, *Variable, []string) {
				f, _ := NewFramework(ctx)
				v := &Variable{
					Name:    "matrix",
					Indices: []string{"i", "j"},
					Shape:   []int{3, 3},
					Type:    SymbolicType,
				}
				return f, v, []string{"i"}
			},
			wantErr: false,
		},
		{
			name: "error on nil variable",
			setup: func() (*Framework, *Variable, []string) {
				f, _ := NewFramework(ctx)
				return f, nil, []string{"i"}
			},
			wantErr: true,
			errMsg:  "variable is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, v, indices := tt.setup()
			result, err := f.Project(ctx, v, indices)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, indices, result.Indices)
			}
		})
	}
}

func TestFramework_Join(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func() (*Framework, *Variable, *Variable)
		wantErr bool
		errMsg  string
	}{
		{
			name: "join two variables successfully",
			setup: func() (*Framework, *Variable, *Variable) {
				f, _ := NewFramework(ctx)
				v1 := &Variable{
					Name:    "A",
					Indices: []string{"i", "j"},
					Type:    SymbolicType,
				}
				v2 := &Variable{
					Name:    "B",
					Indices: []string{"j", "k"},
					Type:    SymbolicType,
				}
				return f, v1, v2
			},
			wantErr: false,
		},
		{
			name: "error on nil first variable",
			setup: func() (*Framework, *Variable, *Variable) {
				f, _ := NewFramework(ctx)
				v2 := &Variable{Name: "B"}
				return f, nil, v2
			},
			wantErr: true,
			errMsg:  "one or both variables are nil",
		},
		{
			name: "error on nil second variable",
			setup: func() (*Framework, *Variable, *Variable) {
				f, _ := NewFramework(ctx)
				v1 := &Variable{Name: "A"}
				return f, v1, nil
			},
			wantErr: true,
			errMsg:  "one or both variables are nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, v1, v2 := tt.setup()
			result, err := f.Join(ctx, v1, v2)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, HybridType, result.Type)
			}
		})
	}
}

func TestFramework_IntegrateWithBoundary(t *testing.T) {
	ctx := context.Background()

	t.Run("integration succeeds", func(t *testing.T) {
		f, err := NewFramework(ctx)
		require.NoError(t, err)

		err = f.IntegrateWithBoundary(ctx)
		assert.NoError(t, err)
	})
}

func TestVariableTypes(t *testing.T) {
	tests := []struct {
		name     string
		varType  VariableType
		expected string
	}{
		{"symbolic type", SymbolicType, "symbolic"},
		{"neural type", NeuralType, "neural"},
		{"probabilistic type", ProbabilisticType, "probabilistic"},
		{"hybrid type", HybridType, "hybrid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.varType))
		})
	}
}

func TestVariable_Creation(t *testing.T) {
	v := &Variable{
		Name:    "test_var",
		Indices: []string{"i", "j", "k"},
		Shape:   []int{2, 3, 4},
		Data:    make([]float64, 24),
		Type:    NeuralType,
	}

	assert.Equal(t, "test_var", v.Name)
	assert.Equal(t, 3, len(v.Indices))
	assert.Equal(t, 3, len(v.Shape))
	assert.Equal(t, 24, len(v.Data))
	assert.Equal(t, NeuralType, v.Type)
}

func TestTensorEquation_Creation(t *testing.T) {
	eq := &TensorEquation{
		Left: Variable{
			Name:    "C",
			Indices: []string{"i", "k"},
		},
		Right:     "A_ij * B_jk",
		Operation: "join",
	}

	assert.Equal(t, "C", eq.Left.Name)
	assert.Equal(t, "A_ij * B_jk", eq.Right)
	assert.Equal(t, "join", eq.Operation)
}

func TestFramework_MultipleVariables(t *testing.T) {
	ctx := context.Background()
	f, err := NewFramework(ctx)
	require.NoError(t, err)

	// Register multiple variables
	vars := []*Variable{
		{Name: "v1", Type: SymbolicType, Indices: []string{"i"}},
		{Name: "v2", Type: NeuralType, Indices: []string{"i", "j"}},
		{Name: "v3", Type: ProbabilisticType, Indices: []string{"k"}},
	}

	for _, v := range vars {
		err := f.RegisterVariable(ctx, v)
		require.NoError(t, err)
	}

	assert.Equal(t, 3, len(f.Variables))
	for _, v := range vars {
		assert.Contains(t, f.Variables, v.Name)
	}
}

func TestFramework_MultipleEquations(t *testing.T) {
	ctx := context.Background()
	f, err := NewFramework(ctx)
	require.NoError(t, err)

	equations := []*TensorEquation{
		{
			Left:      Variable{Name: "e1"},
			Right:     "expr1",
			Operation: "join",
		},
		{
			Left:      Variable{Name: "e2"},
			Right:     "expr2",
			Operation: "project",
		},
	}

	for _, eq := range equations {
		err := f.DefineEquation(ctx, eq)
		require.NoError(t, err)
	}

	assert.Equal(t, 2, len(f.Equations))
}
