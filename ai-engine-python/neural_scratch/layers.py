"""
Modular Deep Neural Network Layers implemented from scratch in pure Python.
"""

import math
import random
from typing import List, Optional, Tuple
from .matrix_ops import (
    Matrix, Vector, zeros, he_uniform, matmul, transpose, add, multiply,
    relu, relu_derivative, sigmoid, sigmoid_derivative, gelu, apply_elementwise, shape
)


class Layer:
    """Base neural layer interface"""
    def forward(self, x: Matrix, training: bool = True) -> Matrix:
        raise NotImplementedError

    def backward(self, grad_output: Matrix) -> Matrix:
        raise NotImplementedError

    def get_params_and_grads(self) -> List[Tuple[Matrix, Matrix]]:
        """Return list of (parameter_tensor, gradient_tensor) for optimizer"""
        return []


class Dense(Layer):
    """Fully Connected / Linear Layer: y = xW + b"""
    def __init__(self, in_features: int, out_features: int, init_scheme: str = "he"):
        self.in_features = in_features
        self.out_features = out_features
        
        # Initialize weights
        if init_scheme == "he":
            self.weights = he_uniform(in_features, out_features)
        else:
            limit = math.sqrt(1.0 / in_features)
            self.weights = [[random.uniform(-limit, limit) for _ in range(out_features)] for _ in range(in_features)]
            
        self.bias = [0.0 for _ in range(out_features)]
        
        # Gradients
        self.grad_weights = zeros(in_features, out_features)
        self.grad_bias = [0.0 for _ in range(out_features)]
        
        self.last_input: Optional[Matrix] = None

    def forward(self, x: Matrix, training: bool = True) -> Matrix:
        self.last_input = x
        out = matmul(x, self.weights)
        return add(out, self.bias)

    def backward(self, grad_output: Matrix) -> Matrix:
        if self.last_input is None:
            raise RuntimeError("Must call forward() before backward()")

        # grad_weights = X^T * grad_output
        input_t = transpose(self.last_input)
        self.grad_weights = matmul(input_t, grad_output)

        # grad_bias = column sum of grad_output
        r, c = shape(grad_output)
        self.grad_bias = [0.0 for _ in range(c)]
        for j in range(c):
            for i in range(r):
                self.grad_bias[j] += grad_output[i][j]

        # grad_input = grad_output * W^T
        weights_t = transpose(self.weights)
        return matmul(grad_output, weights_t)

    def get_params_and_grads(self) -> List[Tuple[Matrix, Matrix]]:
        # Wrap bias as 1xN matrix for uniform optimizer interface
        return [
            (self.weights, self.grad_weights),
            ([self.bias], [self.grad_bias])
        ]


class ReLU(Layer):
    """Rectified Linear Unit activation layer"""
    def __init__(self):
        self.last_input: Optional[Matrix] = None

    def forward(self, x: Matrix, training: bool = True) -> Matrix:
        self.last_input = x
        return apply_elementwise(x, relu)

    def backward(self, grad_output: Matrix) -> Matrix:
        if self.last_input is None:
            raise RuntimeError("Must call forward() before backward()")
        
        r, c = shape(grad_output)
        grad_input = zeros(r, c)
        for i in range(r):
            for j in range(c):
                grad_input[i][j] = grad_output[i][j] * (1.0 if self.last_input[i][j] > 0 else 0.0)
        return grad_input


class Sigmoid(Layer):
    """Sigmoid activation layer"""
    def __init__(self):
        self.last_output: Optional[Matrix] = None

    def forward(self, x: Matrix, training: bool = True) -> Matrix:
        out = apply_elementwise(x, sigmoid)
        self.last_output = out
        return out

    def backward(self, grad_output: Matrix) -> Matrix:
        if self.last_output is None:
            raise RuntimeError("Must call forward() before backward()")
        
        r, c = shape(grad_output)
        grad_input = zeros(r, c)
        for i in range(r):
            for j in range(c):
                s = self.last_output[i][j]
                grad_input[i][j] = grad_output[i][j] * s * (1.0 - s)
        return grad_input


class Dropout(Layer):
    """Inverted Dropout regularization"""
    def __init__(self, p: float = 0.2):
        self.p = p
        self.mask: Optional[Matrix] = None

    def forward(self, x: Matrix, training: bool = True) -> Matrix:
        if not training or self.p <= 0.0:
            return x
        
        scale = 1.0 / (1.0 - self.p)
        r, c = shape(x)
        self.mask = zeros(r, c)
        out = zeros(r, c)
        for i in range(r):
            for j in range(c):
                keep = 1.0 if random.random() > self.p else 0.0
                self.mask[i][j] = keep * scale
                out[i][j] = x[i][j] * self.mask[i][j]
        return out

    def backward(self, grad_output: Matrix) -> Matrix:
        if self.mask is None:
            return grad_output
        return multiply(grad_output, self.mask)
