package arm

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Matrix struct {
	r int         // Rows of data
	c int         // Cols of data
	d [][]float64 // Data
}

//
// Matrix MANIPULATION
//

// Make a new Matrix with rows and cols
func NewMtx(rows int, cols int) Matrix {
	d := make([][]float64, rows)
	for i := range d {
		d[i] = make([]float64, cols)
		for j := range d[i] {
			d[i][j] = 0.0
		}
	}
	return Matrix{rows, cols, d}
}

func NewIdentityMtx(rows, cols int) Matrix {
	mtx := NewMtx(rows, cols)
	for i := range mtx.d {
		mtx.d[i][i] = 1
	}
	return mtx
}

// Make a new Matrix from 2D array
func NewMtxArr(data [][]float64) Matrix {
	return Matrix{len(data), len(data[0]), data}
}

// Return Matrix rows
func (m Matrix) R() int {
	return m.r
}

// Return Matrix cols
func (m Matrix) C() int {
	return m.c
}

// Return Matrix data
func (m Matrix) D() [][]float64 {
	return m.d
}

// Copy Matrix n into Matrix m
func (m *Matrix) Cpy(n Matrix) {
	m.r = n.r
	m.c = n.c
	m.d = make([][]float64, m.r)
	for i := range n.d {
		m.d[i] = make([]float64, n.c)
		copy(m.d[i], n.d[i])
	}
}

// Transpose data in Matrix
func (m *Matrix) Trans() {
	t := NewMtx(m.c, m.r)
	for i := range m.d {
		for j := range m.d[i] {
			t.d[j][i] = m.d[i][j]
		}
	}
	m.Cpy(t)
}

// Transpose data in Matrix m
func Trans(m Matrix) Matrix {
	t := NewMtx(m.c, m.r)
	for i := range m.d {
		for j := range m.d[i] {
			t.d[j][i] = m.d[i][j]
		}
	}
	return t
}

// Randomize Matrix data values
func (m *Matrix) Rand() {
	r := rand.New(rand.NewSource(time.Now().UnixMilli() * rand.Int63()))
	for i := range m.d {
		for j := range m.d[i] {
			m.d[i][j] = r.Float64()*2 - 1
		}
	}
}

// Get biggest float from data
func LongestFloat(data [][]float64) int {
	o := 0
	for _, rd := range data {
		for _, cd := range rd {
			if len(fmt.Sprintf("%f", cd)) > o {
				o = len(fmt.Sprintf("%f", cd))
			}
		}
	}
	return o
}

// Print Matrix data
func (m Matrix) Print() {
	l := LongestFloat(m.d)
	fmt.Printf("  ")
	for i := range m.d[0] {
		fmt.Printf("|%d", i)
		for j := 1; j < l; j++ {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	for i, rd := range m.d {
		fmt.Printf(" %d|", i)
		for _, cd := range rd {
			s := fmt.Sprintf("%f", cd)
			for i := len(s); i < l; i++ {
				fmt.Print(" ")
			}
			fmt.Printf("%s ", s)
		}
		fmt.Println()
	}
	fmt.Println()
}

// Format matrix data to printable string
func (m Matrix) Format() string {
	out := ""
	l := LongestFloat(m.d)
	out += "\n"
	for _, rd := range m.d {
		for _, cd := range rd {
			s := fmt.Sprintf("%f", cd)
			for i := len(s); i < l; i++ {
				out += " "
			}
			out += s + " "
		}
		out += "\n"
	}
	out += "\n"
	return out
}

// Apply function fn to data in Matrix
func (m *Matrix) Map(fn func(float64) float64) {
	for _, rd := range m.d {
		for j := range rd {
			rd[j] = fn(rd[j])
		}
	}
}

// Apply function fn to data in Matrix m
func Map(m Matrix, fn func(float64) float64) Matrix {
	o := NewMtx(m.r, m.c)
	for i, rd := range m.d {
		for j := range rd {
			o.d[i][j] = fn(rd[j])
		}
	}
	return o
}

//
// Matrix OPERATIONS
//

// Add i to every num in Matrix
func (m *Matrix) Add(i float64) {
	for _, rd := range m.d {
		for j := range rd {
			rd[j] += i
		}
	}
}

// Add Matrix n to Matrix
func (m *Matrix) AddMtx(n Matrix) {
	if m.r != n.r || m.c != n.c {
		return
	}
	o := NewMtx(m.r, n.c)
	for i, rd := range o.d {
		for j := range rd {
			rd[j] = m.d[i][j] + n.d[i][j]
		}
	}
	m.Cpy(o)
}

// Multiply every num in Matrix by i
func (m *Matrix) Mult(i float64) {
	for _, rd := range m.d {
		for j := range rd {
			rd[j] *= i
		}
	}
}

// Divide every num in Matrix by i
func (m *Matrix) Div(i float64) {
	for _, rd := range m.d {
		for j := range rd {
			rd[j] /= i
		}
	}
}

// Hadamard multiplication of Matrix by Matrix n
func (m *Matrix) MultMtx(n Matrix) {
	for i, rd := range m.d {
		for j := range rd {
			m.d[i][j] *= n.d[i][j]
		}
	}
}

// Dot product of Matrix and Matrix n copied to Matrix m
func (m *Matrix) Dot(n Matrix) {
	if m.c != n.r {
		fmt.Println("DOT: M != N")
		return
	}
	o := NewMtx(m.r, n.c)
	for i := 0; i < m.r; i++ {
		for j := 0; j < n.c; j++ {
			for k := 0; k < m.c; k++ {
				o.d[i][j] += m.d[i][k] * n.d[k][j]
			}
		}
	}
	m.Cpy(o)
}

// Add Matrix n data to Matrix data, returning new Matrix
func AddMtx(m Matrix, n Matrix) Matrix {
	o := NewMtx(m.r, n.c)
	for i, rd := range o.d {
		for j := range rd {
			rd[j] = m.d[i][j] + n.d[i][j]
		}
	}
	return o
}

// Subtract Matrix n data from Matrix m, returning new Matrix
func SubMtx(m Matrix, n Matrix) Matrix {
	o := NewMtx(m.r, n.c)
	for i, rd := range o.d {
		for j := range rd {
			rd[j] = m.d[i][j] - n.d[i][j]
		}
	}
	return o
}

// Hadamard multiplication of Matrix m by Matrix n, returning new Matrix
func MultMtx(m Matrix, n Matrix) Matrix {
	o := NewMtx(m.r, n.c)
	for i, rd := range o.d {
		for j := range rd {
			o.d[i][j] = m.d[i][j] * n.d[i][j]
		}
	}
	return o
}

// Multiply every num in Matrix m by i
func Mult(m Matrix, n float64) Matrix {
	o := NewMtx(m.r, m.c)
	for i := range m.d {
		for j := range m.d[i] {
			o.d[i][j] = m.d[i][j] * n
		}
	}
	return o
}

// Multiply every num in Matrix m by i
func Sqrt(m Matrix) Matrix {
	o := NewMtx(m.r, m.c)
	for i := range m.d {
		for j := range m.d[i] {
			o.d[i][j] = math.Sqrt(m.d[i][j])
		}
	}
	return o
}

// Dot product of Matrix m and Matrix n, returning new Matrix
func Dot(m Matrix, n Matrix) Matrix {
	if m.c != n.r {
		fmt.Println("!!! DOT: M != N !!!")
		return Matrix{}
	}
	o := NewMtx(m.r, n.c)
	for i := 0; i < m.r; i++ {
		for j := 0; j < n.c; j++ {
			for k := 0; k < m.c; k++ {
				o.d[i][j] += m.d[i][k] * n.d[k][j]
			}
		}
	}
	return o
}

// Hadamard multiplication of Matrixes m, returning new Matrix
func MultMtxArr(m []Matrix) Matrix {
	if len(m) < 2 {
		return NewMtx(0, 0)
	}
	o := m[0]
	for i := 1; i < len(m); i++ {
		o = MultMtx(o, m[i])
	}
	return o
}

// Dot product of Matrices m, returning new Matrix
func DotArr(m []Matrix) Matrix {
	if len(m) == 0 {
		return NewMtx(0, 0)
	}
	if len(m) < 2 {
		return m[0]
	}
	o := m[0]
	for i := 0; i < len(m); i++ {
		o = Dot(o, m[i])
	}
	return o
}

// Concat two vectors in top-bottom(1) or right-left(0) direction, returning matrix
func Concat(m Matrix, n Matrix, dir bool) Matrix {
	var o Matrix
	if dir {
		o = NewMtx(m.r+n.r, m.r+n.r)
	} else {
		o = NewMtx(m.c+n.c, m.c+n.c)
	}

	for r := 0; r < o.r; r++ {
		for c := 0; c < o.c; c++ {
			if dir {
				if r < m.r && c < m.c {
					o.d[r][c] = m.d[r][c]
				} else if r >= m.r && c < m.c {
					o.d[r][c] = n.d[r-m.r][c]
				} else if r < m.r && c == m.c {
					o.d[r][c] = 0
				} else {
					o.d[r][c] = 1
				}
			} else {
				if r < m.r && c < m.c {
					o.d[r][c] = m.d[r][c]
				} else if c >= m.c && r < m.r {
					o.d[r][c] = n.d[r][c-m.c]
				} else if r == m.r && c < m.c {
					o.d[r][c] = 0
				} else {
					o.d[r][c] = 1
				}
			}
		}
	}

	return o
}
