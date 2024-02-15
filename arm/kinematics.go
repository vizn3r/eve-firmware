package arm

import (
	"eve-firmware/cmds"
	"eve-firmware/util"
	"fmt"
	"math"
	"strconv"
)

type Position struct {
	Rotation       Matrix
	JointTheta     []float64
	JointRotations []Matrix

	JointDisplacemets []Matrix

	// Homogeneous Transformation Matricesj
	HTMatrix   Matrix
	HTMatrices []Matrix

	// Parameter Table
	ParamTable  [][]float64
	LinkLengths []float64
}

type Joint struct {
	Rotation [][]float64
	Position []float64
}

type RobotConfig struct {
	ParameterTable [][]float64
	LinkLengths    []float64
	Joints         []Joint
}

var POSITION = new(Position)

func sin(i float64) float64 {
	return math.Sin(i)
}

func cos(i float64) float64 {
	return math.Cos(i)
}

func degToRad(deg float64) float64 {
	return (deg / 180) * math.Pi
}

func getRotation(t Matrix) Matrix {
	return NewMtxArr([][]float64{
		{t.d[0][0], t.d[0][1], t.d[0][2]},
		{t.d[1][0], t.d[1][1], t.d[1][2]},
		{t.d[2][0], t.d[2][1], t.d[2][2]},
	})
}

func getPosition(t Matrix) Matrix {
	return NewMtxArr([][]float64{
		{t.d[0][3]},
		{t.d[1][3]},
		{t.d[2][3]},
	})
}

func HTMFromTo(from int, to int) Matrix {
	return DotArr(POSITION.HTMatrices[from:to])
}

func HTMFromArr(p []float64, deg bool) Matrix {
	// 0 - Theta
	t := p[0]
	// 1 - Alpha
	a := p[1]
	// 2 - r
	r := p[2]
	// 3 - d
	d := p[3]
	if deg {
		t = degToRad(t)
		a = degToRad(a)
		return NewMtxArr([][]float64{
			{cos(t), -sin(t) * cos(a), sin(t) * sin(a), r * cos(t)},
			{sin(t), cos(t) * cos(a), -cos(t) * sin(a), r * sin(t)},
			{0, sin(a), cos(a), d},
			{0, 0, 0, 1},
		})
	}

	return NewMtxArr([][]float64{
		{cos(t), -sin(t) * cos(a), sin(t) * sin(a), r * cos(t)},
		{sin(t), cos(t) * cos(a), -cos(t) * sin(a), r * sin(t)},
		{0, sin(a), cos(a), d},
		{0, 0, 0, 1},
	})
}

func rotateY(r Matrix, angle float64) Matrix {
	return Dot(NewMtxArr([][]float64{
		{cos(angle), -sin(angle), 0},
		{sin(angle), cos(angle), 0},
		{0, 0, 1},
	}), r)
}

func rotateX(r Matrix, angle float64) Matrix {
	return Dot(NewMtxArr([][]float64{
		{1, 0, 0},
		{0, cos(angle), -sin(angle)},
		{0, sin(angle), cos(angle)},
	}), r)
}

func posZ(d Matrix) Matrix {
	return NewMtxArr([][]float64{
		{0},
		{0},
		{d.d[2][0]},
	})
}

func posX(d Matrix) Matrix {
	return NewMtxArr([][]float64{
		{cos(d.d[0][0])},
		{sin(d.d[0][1])},
		{0},
	})
}

func HTMCalculate() {
	T_rz := []Matrix{}
	for _, r := range POSITION.JointRotations {
		T_rz = append(T_rz, Concat(rotateY(r, 0), NewMtx(3, 1), false))
	}

	T_z := []Matrix{}
	for _, d := range POSITION.JointDisplacemets {
		T_z = append(T_z, Concat(NewIdentityMtx(3, 3), posZ(d), false))
	}

	T_rx := []Matrix{}
	for _, r := range POSITION.JointRotations {
		T_rx = append(T_rx, Concat(rotateX(r, 0), NewMtx(3, 1), false))
	}

	T_x := []Matrix{}
	for _, d := range POSITION.JointDisplacemets {
		T_x = append(T_x, Concat(NewIdentityMtx(4, 4), posX(d), false))
	}

	for i := range POSITION.HTMatrices {
		POSITION.HTMatrices[i] = DotArr([]Matrix{T_rz[i], T_z[i], T_x[i], T_rx[i]})
		POSITION.HTMatrices[i].Print()
	}
}

func InitPosition() {
	cmds.COMMANDS = append(cmds.COMMANDS, cmds.Command{
		Call: 'K',
		Type: cmds.USER,
		Funcs: []cmds.CommandFunc{
			{
				NumArgs: 0,
				Desc:    "Print Position matrix",
				Func: func(c cmds.CommandCtx) string {
					return POSITION.Rotation.Format()
				},
			},
			{
				NumArgs: 0,
				Desc:    "Print Homogeneous Transformation matrices",
				Func: func(cmds.CommandCtx) string {
					out := ""
					for i, m := range POSITION.HTMatrices {
						out += "H" + strconv.Itoa(i) + "_" + strconv.Itoa(i+1) + "\n"
						out += m.Format() + "\n"
					}
					return out
				},
			},
			{
				NumArgs: 2,
				Desc:    "Print Homogeneous Transformation matrix from n to m",
				Args:    "<n, m>",
				Func: func(c cmds.CommandCtx) string {
					m, n := c.IntArgs[0], c.IntArgs[1]
					return "H" + strconv.Itoa(m) + "_" + strconv.Itoa(n) + "\n" + HTMFromTo(m, n).Format()
				},
			},
			{
				NumArgs: 0,
				Desc:    "Test the inverse kinematics output",
				Func: func(c cmds.CommandCtx) string {
					Inverse(HTMFromTo(0, 6))
					return ""
				},
			},
		},
	})
	CalculatePosition()
}

func CalculatePosition() {
	var conf RobotConfig
	util.ParseJSON("./conf/robot.json", &conf)
	POSITION.ParamTable = conf.ParameterTable
	POSITION.LinkLengths = conf.LinkLengths

	for _, j := range conf.Joints {
		POSITION.JointRotations = append(POSITION.JointRotations, NewMtxArr(j.Rotation))
	}

	HTMCalculate()
}

func pow2(n float64) float64 {
	return n * n
}

func Inverse(t Matrix) {
	CalculatePosition()

	Pt := getRotation(t)
	Pp := getPosition(t)

	fmt.Println("Rotation")
	Pt.Print()
	fmt.Println("Position")
	Pp.Print()

	d := POSITION.LinkLengths
	fmt.Println("Lenghts", d)

	// Calculation of W (T0_3)
	W := NewMtxArr([][]float64{
		{Pp.d[0][0] - d[5]*Pt.d[2][0]},
		{Pp.d[1][0] - d[5]*Pt.d[2][1]},
		{Pp.d[2][0] - d[5]*Pt.d[2][2]},
	})

	fmt.Println("W - calculated")
	W.Print()

	r := math.Sqrt(pow2(W.d[0][0]) + pow2(W.d[0][0]))

	// Calculation of the first part (fist 3 angles set the end effector position)
	theta1 := math.Atan(W.d[0][0] / W.d[2][0])
	theta2 := math.Atan((pow2(d[1])-pow2(d[2])-pow2(r))/(-2*d[1]*r)) - math.Atan(W.d[2][0]/W.d[0][0])
	theta3 := math.Acos(pow2(r) - pow2(d[1]) - pow2(d[2])/(-2*d[1]*d[2]))
	fmt.Println(theta1, theta2, theta3)
}

// "ParameterTable": [
//   [ 0, 90, 0, 55 ],
//   [ 0, 0, 168, 0 ],
//   [ 90, 90, 48, 0 ],
//   [ 180, 90, 0, 115 ],
//   [ 0, -90, 91, 0 ],
//   [ 0, 0, 0, 0 ]
//  ],
