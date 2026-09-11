"""
Pure-Python Matrix & Tensor Math Library from First Principles.
Zero third-party dependencies required.
Implements matrix multiplication, broadcasting, activation functions, and analytical gradients.
"""

import math
import random
from typing import List, Tuple, Union

Matrix = List[List[float]]
Vector = List[float]


def zeros(rows: int, cols: int) -> Matrix:
    """Create a rows x cols matrix initialized to 0.0"""
    return [[0.0 for _ in range(cols)] for _ in range(rows)]


def ones(rows: int, cols: int) -> Matrix:
    """Create a rows x cols matrix initialized to 1.0"""
    return [[1.0 for _ in range(cols)] for _ in range(rows)]


def randn(rows: int, cols: int, mean: float = 0.0, std: float = 1.0) -> Matrix:
    """Create a rows x cols matrix with Box-Muller Gaussian normal distribution"""
    mat = zeros(rows, cols)
    for i in range(rows):
        for j in range(cols):
            # Box-Muller transform
            u1 = max(1e-15, random.random())
            u2 = random.random()
            z = math.sqrt(-2.0 * math.log(u1)) * math.cos(2.0 * math.pi * u2)
            mat[i][j] = mean + z * std
    return mat


def he_uniform(rows: int, cols: int) -> Matrix:
    """He / Kaiming initialization for deep ReLU networks"""
    limit = math.sqrt(6.0 / rows)
    return [[random.uniform(-limit, limit) for _ in range(cols)] for _ in range(rows)]


def xavier_uniform(rows: int, cols: int) -> Matrix:
    """Glorot / Xavier initialization for tanh/sigmoid networks"""
    limit = math.sqrt(6.0 / (rows + cols))
    return [[random.uniform(-limit, limit) for _ in range(cols)] for _ in range(rows)]


def shape(mat: Matrix) -> Tuple[int, int]:
    """Return (rows, cols) of a matrix"""
    if not mat or not mat[0]:
        return 0, 0
    return len(mat), len(mat[0])


def transpose(mat: Matrix) -> Matrix:
    """Transpose matrix A^T"""
    r, c = shape(mat)
    trans = zeros(c, r)
    for i in range(r):
        for j in range(c):
            trans[j][i] = mat[i][j]
    return trans


def matmul(A: Matrix, B: Matrix) -> Matrix:
    """Matrix multiplication C = A * B"""
    rA, cA = shape(A)
    rB, cB = shape(B)
    if cA != rB:
        raise ValueError(f"Incompatible shapes for matmul: ({rA}, {cA}) x ({rB}, {cB})")

    C = zeros(rA, cB)
    # Optimized cache-friendly traversal
    for i in range(rA):
        for k in range(cA):
            a_ik = A[i][k]
            if a_ik == 0.0:
                continue
            for j in range(cB):
                C[i][j] += a_ik * B[k][j]
    return C


def add(A: Matrix, B: Union[Matrix, Vector, float]) -> Matrix:
    """Element-wise addition with broadcasting support for row vector bias or scalar"""
    r, c = shape(A)
    res = zeros(r, c)

    if isinstance(B, (int, float)):
        for i in range(r):
            for j in range(c):
                res[i][j] = A[i][j] + float(B)
    elif isinstance(B, list) and len(B) > 0 and isinstance(B[0], (int, float)):
        # 1D Bias vector of length c broadcast across all rows
        if len(B) != c:
            raise ValueError(f"Bias length {len(B)} does not match matrix columns {c}")
        for i in range(r):
            for j in range(c):
                res[i][j] = A[i][j] + float(B[j])
    else:
        # 2D Matrix
        rB, cB = shape(B)  # type: ignore
        if r != rB or c != cB:
            raise ValueError(f"Cannot add shapes ({r},{c}) and ({rB},{cB})")
        for i in range(r):
            for j in range(c):
                res[i][j] = A[i][j] + B[i][j]  # type: ignore
    return res


def subtract(A: Matrix, B: Matrix) -> Matrix:
    """Element-wise subtraction A - B"""
    r, c = shape(A)
    res = zeros(r, c)
    for i in range(r):
        for j in range(c):
            res[i][j] = A[i][j] - B[i][j]
    return res


def multiply(A: Matrix, B: Union[Matrix, float]) -> Matrix:
    """Hadamard element-wise multiplication A .* B or scalar scaling"""
    r, c = shape(A)
    res = zeros(r, c)
    if isinstance(B, (int, float)):
        for i in range(r):
            for j in range(c):
                res[i][j] = A[i][j] * float(B)
    else:
        for i in range(r):
            for j in range(c):
                res[i][j] = A[i][j] * B[i][j]
    return res


# Activation Functions & Analytical Derivatives

def relu(x: float) -> float:
    return max(0.0, x)


def relu_derivative(x: float) -> float:
    return 1.0 if x > 0 else 0.0


def sigmoid(x: float) -> float:
    if x >= 0:
        return 1.0 / (1.0 + math.exp(-x))
    else:
        z = math.exp(x)
        return z / (1.0 + z)


def sigmoid_derivative(x: float) -> float:
    s = sigmoid(x)
    return s * (1.0 - s)


def gelu(x: float) -> float:
    """Gaussian Error Linear Unit approximation"""
    return 0.5 * x * (1.0 + math.tanh(math.sqrt(2.0 / math.pi) * (x + 0.044715 * math.pow(x, 3))))


def apply_elementwise(mat: Matrix, fn) -> Matrix:
    """Apply an arbitrary scalar function elementwise"""
    r, c = shape(mat)
    res = zeros(r, c)
    for i in range(r):
        for j in range(c):
            res[i][j] = fn(mat[i][j])
    return res


def softmax(mat: Matrix) -> Matrix:
    """Numerically stable row-wise softmax"""
    r, c = shape(mat)
    res = zeros(r, c)
    for i in range(r):
        max_val = max(mat[i])
        exp_sum = sum(math.exp(val - max_val) for val in mat[i])
        for j in range(c):
            res[i][j] = math.exp(mat[i][j] - max_val) / max(exp_sum, 1e-12)
    return res


def mse_loss(y_pred: Matrix, y_true: Matrix) -> Tuple[float, Matrix]:
    """Mean Squared Error loss and its analytical gradient dL/dy_pred"""
    r, c = shape(y_pred)
    total_err = 0.0
    grad = zeros(r, c)
    n = r * c

    for i in range(r):
        for j in range(c):
            diff = y_pred[i][j] - y_true[i][j]
            total_err += diff * diff
            grad[i][j] = (2.0 / n) * diff

    return total_err / n, grad
