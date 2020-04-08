package main

//  ________________       ______________      ________________    _________           __           ________
// |                |     |              |    |                |  |   ___   \         |  |          \        \
//  ¯¯¯¯¯|   |¯¯¯¯¯¯      |  |¯¯¯¯¯¯¯¯¯¯¯      ¯¯¯¯¯|   |¯¯¯¯¯¯   |  |   |   |         ¯¯            \    \¯¯¯
//		 |   |            |  |___________           |   |         |   ¯¯¯   /          __             \    \
//       |   |            |              |          |   |         |        \          |  |             \    \
//       |   |            |  |¯¯¯¯¯¯¯¯¯¯¯           |   |         |   | \   \         |  |              \    \
//       |   |            |  |___________           |   |         |   |  \   \        |  |            ___\    \
//	     |   |			  | 	         |          |   |         |   |   \   \       |  |            \        \
//        ¯¯¯              ¯¯¯¯¯¯¯¯¯¯¯¯¯¯            ¯¯¯           ¯¯¯      ¯¯¯        ¯¯              ¯¯¯¯¯¯¯¯¯
import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	heigth = 21
	width  = 12
)

var (
	//Bord you dont mess with it
	Bord [heigth][width]int // x,y
	//T = time
	T int64
	// Fall for how many time he need to fall
	Fall int
	//Delay the player
	Delay int
	// Line the obj found
	Line int
)

type dir int

type obj struct {
	place     [8]int
	dir       dir
	lastplace [8]int
	d         bool
	dead      []int
	touch     bool
}

func init() {
	edge()
	Delay = 300
}

func main() {
	fmt.Println("ready")
	var r string
	x := [8]int{5, 0, 6, 0, 5, 1, 6, 1}
	bob := newobj(x)
	T = time.Now().UnixNano()
	for {
		T = time.Now().UnixNano()
		add(bob)
		view(bob)
		_, err := fmt.Scanln(&r)
		if err != nil {
			panic(err)
		}
		bob.lastplace = bob.place
		bob.dir = stod(r)
		move(bob)
		times(bob)
		view(bob)
		fmt.Printf("%v \n", Fall)
		fmt.Println(r)

	}
}

func view(b *obj) {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()

	for y := 0; y < heigth; y++ {
		for x := 0; x < width; x++ {
			if x < width-1 {
				if Bord[y][x] == 0 {
					fmt.Print(" ")
				} else {
					fmt.Print(Bord[y][x])
				}

			} else {
				fmt.Println(Bord[y][x])
			}
		}

	}
	if len(b.dead) != 0 {
		fmt.Printf("pl: %v \nlp: %v \nded:%v \n", b.place, b.lastplace, b.dead[len(b.dead)-8:])
	}
}

func edge() {
	for i := 0; i < width-1; i++ {
		Bord[heigth-1][i] = 1
	}
	for x := 0; x < heigth; x++ {
		Bord[x][width-1] = 1
		Bord[x][0] = 1
	}
}

func newobj(p [8]int) *obj {
	//Iblock = [8]int{5, 0, 6, 0, 5, 1, 6, 1}
	//jblock = [8]int{5, 0, 6, 0, 5, 1, 5, 2}
	//lblock = [8]int{5, 0, 4, 0, 5, 1, 5, 2}
	//sblock
	//Tblock = [8]int{5, 0, 4, 0, 6, 0, 5, 1}
	//zblock
	name := obj{p, dir(0), p, false, nil, false}

	return &name
}

func add(bob *obj) {
	b := bob
	for out := 0; out < 8; out++ {
		if b.place[out] > 21 { // can walk tour the edge
			b.place = b.lastplace
		}
	}

	if len(b.dead) != 0 {
		for rip := 0; rip < len(bob.dead); rip += 2 { // add the dead
			fmt.Println(b.dead)
			Bord[b.dead[rip+1]][b.dead[rip]] = 3
		}
		//panic(bob.dead)
	}

	for i := 0; i < 8; i += 2 {
		if Bord[b.place[i+1]][b.place[i]] != 1 {
			Bord[b.place[i+1]][b.place[i]] = 2
		} else {
			if b.place[i+1] == 20 {
				b.d = true
			}
			if b.d {
				for _, val := range b.lastplace { //add to dead
					bob.dead = append(bob.dead, val)
					bob.newlook()
					b.d = false
				}
			} else {
				bob.place = bob.lastplace
			}

		}
	}

}

func move(bob *obj) {

	switch int(bob.dir) {
	case 0:
		del(bob)
		if Line < 21 {
			for i := 1; i < 8; i += 2 {
				bob.place[i]++

			}
		} else {
			bob.newlook()
		}
	case 1: //left
		del(bob)
		for i := 0; i < 8; i += 2 {
			bob.place[i]++
		}
	case 2: //right
		del(bob)
		for i := 0; i < 8; i += 2 {
			bob.place[i]--
		}
	default:
		del(bob)
	}
	add(bob)

}

func del(bob *obj) {
	for i := 0; i < 8; i += 2 {
		Bord[bob.place[i+1]][bob.place[i]] = 0
	}
}
func stod(s string) dir {
	var d dir
	switch s {
	case "d":
		d = dir(1)
	case "a":
		d = dir(2)
	case "s":
		d = dir(0)
	}
	return d
}
func times(bob *obj) {
	now := time.Now().UnixNano()
	down := (int(now-T) / 1000000)
	Fall = down / Delay
	Line += Fall
	if Line >= 20 {
		Fall = 0
		Line = 0
		T = time.Now().UnixNano()
	}
	del(bob)

	for x := 0; x < Fall; x++ {
		for i := 1; i < 8; i += 2 {
			if bob.place[i] < 20 {
				bob.place[i]++
			} else {
				bob.place[i] = 20
			}
		}
	}
}

func check() {
	a := Bord[20][1:11]
	for line := 0; line < heigth-1; line++ {
		for num := 0; num < len(a); num++ {
			if a[0] == a[num] && a[num] != 0 {
				linedown(line)
			}
		}
	}
}
func linedown(line int) {
	for x := line; x >= 0; x-- {
		Bord[line] = Bord[x]
	}
}
func (bop *obj) newlook() {
	bop.place = [8]int{5, 0, 6, 0, 5, 1, 6, 1}

}
