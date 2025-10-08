package aiassist

import (
	"math"
	"sync"
)

// Простая реализация LinUCB для маленькой размерности фич
type LinUCB struct {
	mu    sync.Mutex
	d     int
	alpha float64
	A     [][]float64 // d x d
	b     []float64   // d
}

func NewLinUCB(dim int, alpha float64) *LinUCB {
	A := make([][]float64, dim)
	for i := range A {
		A[i] = make([]float64, dim)
		A[i][i] = 1.0
	}
	return &LinUCB{
		d: dim, alpha: alpha, A: A, b: make([]float64, dim),
	}
}

// Выбор лучшего «рукава» среди arms (каждый arm — вектор длины d)
func (l *LinUCB) Choose(arms [][]float64) (int, float64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	Ainv := inv(l.A)
	theta := matVec(Ainv, l.b)

	bestIdx := 0
	bestU := math.Inf(-1)
	for i, x := range arms {
		est := dot(theta, x)
		varU := quadForm(Ainv, x)
		u := est + l.alpha*math.Sqrt(varU)
		if u > bestU {
			bestU, bestIdx = u, i
		}
	}
	return bestIdx, bestU
}

// Апдейт по выбранному арм/фичам и наблюдаемой награде
func (l *LinUCB) Update(x []float64, reward float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	outerAdd(l.A, x)
	axpy(l.b, x, reward)
}

// ===== линал утилиты (наивно, без зависимостей) =====

func dot(a, b []float64) float64 {
	s := 0.0
	for i := range a {
		s += a[i] * b[i]
	}
	return s
}

func axpy(y, x []float64, a float64) {
	for i := range y {
		y[i] += a * x[i]
	}
}

func outerAdd(A [][]float64, x []float64) {
	for i := range x {
		for j := range x {
			A[i][j] += x[i] * x[j]
		}
	}
}

func matVec(A [][]float64, x []float64) []float64 {
	y := make([]float64, len(x))
	for i := range A {
		s := 0.0
		for j := range x {
			s += A[i][j] * x[j]
		}
		y[i] = s
	}
	return y
}

func quadForm(A [][]float64, x []float64) float64 {
	Ax := matVec(A, x)
	return dot(x, Ax)
}

// Наивная инверсия матрицы методом Гаусса (достаточно для d~8–16)
func inv(M [][]float64) [][]float64 {
	n := len(M)
	A := make([][]float64, n)
	for i := range A {
		A[i] = make([]float64, 2*n)
		for j := 0; j < n; j++ {
			A[i][j] = M[i][j]
		}
		A[i][n+i] = 1
	}
	for i := 0; i < n; i++ {
		// частичный выбор главного элемента
		p := i
		for r := i; r < n; r++ {
			if math.Abs(A[r][i]) > math.Abs(A[p][i]) {
				p = r
			}
		}
		A[i], A[p] = A[p], A[i]
		piv := A[i][i]
		if math.Abs(piv) < 1e-12 {
			continue
		}
		for j := 0; j < 2*n; j++ {
			A[i][j] /= piv
		}
		for r := 0; r < n; r++ {
			if r == i {
				continue
			}
			f := A[r][i]
			for j := 0; j < 2*n; j++ {
				A[r][j] -= f * A[i][j]
			}
		}
	}
	Inv := make([][]float64, n)
	for i := range Inv {
		Inv[i] = make([]float64, n)
		copy(Inv[i], A[i][n:])
	}
	return Inv
}
