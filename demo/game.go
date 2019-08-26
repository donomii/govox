// game.go
package main

//"math/rand"
//"time"

var monsters []Vec3

func InitGame(size int, maze [][]int) {

	/*
		for monsters := []Vec3{}; len(monsters) < 5; monsters = monsters {
			pos := Vec3{rand.Intn(50), 0, rand.Intn(70)}

			if moveOk(pos, maze) {
			monsters = append(monsters, pos)
			}
		}
	*/

	monsters = append(monsters, Vec3{21, 0, 14})
	monsters = append(monsters, Vec3{37, 0, 19})
	monsters = append(monsters, Vec3{20, 0, 26})
	monsters = append(monsters, Vec3{35, 0, 22})
	monsters = append(monsters, Vec3{27, 0, 16})

}

func handleCollision(pl, other Vec3) []Vec3 {
	out := []Vec3{}
	for _, v := range monsters {
		if v[0] == other[0] && v[1] == other[1] && v[2] == other[2] {

		} else {
			out = append(out, v)
		}
	}
	return out
}

func InView(ppos, tpos Vec3) bool {

	scrPos := Vec3{0, 0, 0}
	scrPos[0] = (tpos[0] - ppos[0]) * (tpos[0] - ppos[0])
	scrPos[1] = (tpos[1] - ppos[1]) * (tpos[1] - ppos[1])
	scrPos[2] = (tpos[2] - ppos[2]) * (tpos[2] - ppos[2])

	if scrPos[0] < tiles*tiles {
		if scrPos[1] < tiles*tiles {
			if scrPos[2] < tiles*tiles {
				return true
			}
		}
	}

	return false
	//return !blocks[pos[0]][pos[1]][pos[2]].Active
}

func moveOk(pos Vec3, maze [][]int) bool {
	size := 300 //FIXME
	if maze[pos[0]][pos[2]] != 1 {
		if pos[0] < size-tileRadius {
			if pos[0] > tileRadius {
				//if pos[1] < size-tileRadius {
				//	if pos[1] > tileRadius {
				if pos[2] < size-tileRadius {
					if pos[2] > tileRadius {
						return true
					}
				}
			}

		}
	}
	return false
	//return !blocks[pos[0]][pos[1]][pos[2]].Active
}

func GenerateMaze(x, y int) [][]int {
	grid := make([][]int, y)
	for i := 0; i < y; i++ {
		var row []int
		for j := 0; j < x; j++ {
			row = append(row, 1)
		}
		grid[i] = row
	}
	return grid
}
