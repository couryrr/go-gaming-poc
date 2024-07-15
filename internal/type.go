package internal

import rl "github.com/gen2brain/raylib-go/raylib"

type Material struct {
    Position rl.Vector2
    Velocity rl.Vector2
    Acceleration rl.Vector2
}

type Cell struct {
    Position rl.Vector2
    CellType CellType
    IsConveyor bool
    IsGoal bool
    Symbol string
    Acceleration rl.Vector2
}

type CellType interface{
    GetColor() rl.Color
    GetType() string
}

type Ground struct {}
type Resource struct {
    ObjectTimer    float32
    ObjectInterval float32
}

type Obstacle struct {}
type Conveyor struct {}
type Lift struct {}

func(g *Ground) GetType() string {
    return "Ground"
}

func(g *Ground) GetColor() rl.Color {
    return rl.DarkGreen
}

func(r *Resource) GetType() string {
    return "Resource"
}

func(r *Resource) GetColor() rl.Color {
    return rl.Beige
}

func(o *Obstacle) GetType() string {
    return "Obstacle"
}

func(o *Obstacle) GetColor() rl.Color {
    return rl.Brown
}

func(c *Conveyor) GetType() string {
    return "Conveyor"
}

func(c *Conveyor) GetColor() rl.Color {
    return rl.Black
}

func(l *Lift) GetType() string {
    return "Lift"
}

func(l *Lift) GetColor() rl.Color {
    return rl.Blue
}

