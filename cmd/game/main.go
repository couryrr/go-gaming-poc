package main

import (
	"fmt"

	"github.com/couryrr/ai-world-eater/internal"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
    width    = 2000
    height   = 1100
    gridSize = 50
)

var (
    grid [][] internal.Cell
    conveyorAccSymbol = "U"
    conveyorAcc = rl.Vector2{X: 0, Y: -1}
    materials []internal.Material
    showGrid bool = false
    score int = 0
    resources []internal.Resource
)

func init(){
    grid = make([][]internal.Cell, (height-100)/gridSize)
    for r := range (height-100)/gridSize {
        grid[r] = make([]internal.Cell, width/gridSize)
        for c := range width/gridSize {
            grid[r][c] = internal.Cell {
                Position: rl.Vector2 {
                    X: float32(c*gridSize),
                    Y: float32(r*gridSize),
                },
                CellType: &internal.Ground{},
                Acceleration: rl.Vector2Zero(),
            }
        }
    }
    
    cell := grid[9][20]
    cell.CellType = &internal.Resource{ObjectInterval: 5}
    grid[9][20] = cell
    
    cell = grid[2][10]
    cell.CellType = &internal.Lift{}
    grid[2][10] = cell
}

func main(){
    rl.InitWindow(width, height, "ai-world-eater")
    defer rl.CloseWindow()
    rl.SetTargetFPS(60)

    for !rl.WindowShouldClose(){
        rl.BeginDrawing()
        rl.HideCursor()
        rl.ClearBackground(rl.White)
        updateConveyorHightlight()
        updateGrid()
        updateMaterial()
        drawGrid()
        drawConveyorHighlight()
        drawMaterial()
        rl.DrawText(fmt.Sprintf("Score: %d", score), 1890, 1070, 25, rl.Black)
        rl.EndDrawing()
    }
}

func drawMaterial(){
    for _, material := range materials{
        rl.DrawCircleV(material.Position, 10, rl.Red)
    }
}

func updateMaterial (){
    for idx := range materials {
        material := materials[idx]
        mX := int32(material.Position.X / gridSize)
        mY := int32(material.Position.Y / gridSize)

        cell := grid[mY][mX]
        cellCenter := rl.Vector2AddValue(cell.Position, gridSize/2)

        if cell.CellType.GetType() == "Lift" {
            score += 1
            material.Velocity = rl.Vector2Zero()
        } else if cell.CellType.GetType() == "Resource" {
            material.Velocity = rl.Vector2Scale(rl.Vector2Add(material.Acceleration, 
                rl.Vector2{X:0, Y: -1}), .75)
        } else if cell.CellType.GetType() == "Conveyor" {
            if rl.Vector2Distance(cellCenter, material.Position) <= 0 {
                material.Velocity = rl.Vector2Add(material.Acceleration, 
                    cell.Acceleration) 
            }
            
            //FIXME: Move material toward center after stopping
            if rl.Vector2Equals(material.Velocity, rl.Vector2Zero()){
                material.Velocity = rl.Vector2Normalize(rl.Vector2Subtract(cellCenter, material.Position))
            }

        } else {
           material.Velocity = rl.Vector2Zero() 
        }

        material.Position = rl.Vector2Add(material.Position, material.Velocity)
        materials[idx] = material
    }
}

func mouseToGrid() (int32, int32){
    mouseX := (rl.GetMouseX() / gridSize) * gridSize
    mouseY := (rl.GetMouseY() / gridSize) * gridSize
    return mouseX, mouseY
}

func updateGrid(){
    if rl.IsKeyPressed(rl.KeyT) {
        showGrid = !showGrid
    }

    mouseX, mouseY := mouseToGrid()
    cell := grid[mouseY/gridSize][mouseX/gridSize]
    if cell.CellType.GetType() != "Lift"  && cell.CellType.GetType() != "Obstacle"{
        if rl.IsMouseButtonPressed(rl.MouseLeftButton){
            mouseX, mouseY := mouseToGrid()
            cell.CellType = &internal.Conveyor{}
            cell.Acceleration = conveyorAcc
            cell.Symbol= conveyorAccSymbol
            grid[mouseY/gridSize][mouseX/gridSize] = cell
        }
    }
}

func drawGrid(){
    dt := rl.GetFrameTime()
    for _, r := range grid {
        for _, c := range r {
            cellCenter := rl.Vector2AddValue(c.Position, gridSize/2)
            rl.DrawRectangle(
                int32(c.Position.X),
                int32(c.Position.Y), 
                gridSize,
                gridSize,
                c.CellType.GetColor())

            if c.CellType.GetType() == "Conveyor" {
                rotation := rl.Vector2Angle(c.Acceleration, rl.Vector2{X: 1, Y: 0})
                rl.DrawPoly(cellCenter, 3, 10, -(rotation*rl.Rad2deg), rl.Yellow) 
            }
            
            if c.CellType.GetType() == "Resource" {
                r := c.CellType.(*internal.Resource)
                if r.ObjectTimer >= r.ObjectInterval {
                    materials = append(materials, internal.Material{Position: cellCenter})
                    r.ObjectTimer = 0
                }
                r.ObjectTimer += dt
            }

            if showGrid {
                rl.DrawCircleV(cellCenter, 5, rl.Gray)
                rl.DrawRectangleLines(
                    int32(c.Position.X),
                    int32(c.Position.Y), 
                    gridSize,
                    gridSize,
                    rl.Gray)
            }
        }
    }
}

func updateConveyorHightlight(){
    if rl.IsKeyPressed(rl.KeyR){
        switch conveyorAccSymbol{
        case "U":
            conveyorAccSymbol = "R"
            conveyorAcc = rl.Vector2{X: 1, Y: 0}
        case "R":
            conveyorAccSymbol = "D"
            conveyorAcc = rl.Vector2{X: 0, Y: 1}
        case "D":
            conveyorAccSymbol = "L"
            conveyorAcc = rl.Vector2{X: -1, Y: 0}
        case "L":
            conveyorAccSymbol = "U"
            conveyorAcc = rl.Vector2{X:0, Y: -1}
        }
    }
}

func drawConveyorHighlight() {
    mouseX, mouseY := mouseToGrid()  
    cellCenter := rl.Vector2{X: float32(mouseX+(gridSize/2)), Y: float32(mouseY+(gridSize/2))}
    rotation := rl.Vector2Angle(conveyorAcc, rl.Vector2{X: 1, Y: 0})
    rl.DrawRectangle(mouseX, mouseY, gridSize, gridSize, rl.Black)
    rl.DrawPoly(cellCenter, 3, 10, -(rotation*rl.Rad2deg), rl.Yellow) 
}

