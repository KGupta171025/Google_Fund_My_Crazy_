"""
Pure Python Deep Learning Optimizers: SGD and Adam with Bias Correction.
"""

import math
from typing import List, Tuple
from .matrix_ops import Matrix, zeros, shape


class Optimizer:
    def __init__(self, lr: float = 0.001):
        self.lr = lr

    def step(self, param_groups: List[Tuple[Matrix, Matrix]]):
        raise NotImplementedError


class SGD(Optimizer):
    """Stochastic Gradient Descent with Momentum"""
    def __init__(self, lr: float = 0.01, momentum: float = 0.9, weight_decay: float = 1e-4):
        super().__init__(lr)
        self.momentum = momentum
        self.weight_decay = weight_decay
        self.velocities: List[Matrix] = []

    def step(self, param_groups: List[Tuple[Matrix, Matrix]]):
        if not self.velocities:
            for params, _ in param_groups:
                r, c = shape(params)
                self.velocities.append(zeros(r, c))

        for idx, (params, grads) in enumerate(param_groups):
            v = self.velocities[idx]
            r, c = shape(params)
            for i in range(r):
                for j in range(c):
                    g = grads[i][j] + self.weight_decay * params[i][j]
                    v[i][j] = self.momentum * v[i][j] + self.lr * g
                    params[i][j] -= v[i][j]


class Adam(Optimizer):
    """Adam Optimizer with first/second moment tracking and bias correction"""
    def __init__(self, lr: float = 0.001, beta1: float = 0.9, beta2: float = 0.999, eps: float = 1e-8, weight_decay: float = 0.0):
        super().__init__(lr)
        self.beta1 = beta1
        self.beta2 = beta2
        self.eps = eps
        self.weight_decay = weight_decay
        self.t = 0
        self.m: List[Matrix] = []
        self.v: List[Matrix] = []

    def step(self, param_groups: List[Tuple[Matrix, Matrix]]):
        self.t += 1
        if not self.m:
            for params, _ in param_groups:
                r, c = shape(params)
                self.m.append(zeros(r, c))
                self.v.append(zeros(r, c))

        # Bias correction coefficients
        bias_c1 = 1.0 - math.pow(self.beta1, self.t)
        bias_c2 = 1.0 - math.pow(self.beta2, self.t)

        for idx, (params, grads) in enumerate(param_groups):
            m_mat = self.m[idx]
            v_mat = self.v[idx]
            r, c = shape(params)

            for i in range(r):
                for j in range(c):
                    g = grads[i][j]
                    if self.weight_decay > 0.0:
                        g += self.weight_decay * params[i][j]

                    # Update biased first & second moment estimates
                    m_mat[i][j] = self.beta1 * m_mat[i][j] + (1.0 - self.beta1) * g
                    v_mat[i][j] = self.beta2 * v_mat[i][j] + (1.0 - self.beta2) * (g * g)

                    # Compute bias-corrected moments
                    m_hat = m_mat[i][j] / max(bias_c1, 1e-12)
                    v_hat = v_mat[i][j] / max(bias_c2, 1e-12)

                    # Update parameter
                    params[i][j] -= self.lr * m_hat / (math.sqrt(v_hat) + self.eps)
