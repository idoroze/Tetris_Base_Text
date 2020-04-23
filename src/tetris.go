package main

// https://github.com/idoroze/Tetris_Base_Text
//  ________________       ______________      ________________    _________           __        ________
// |                |     |              |    |                |  |   ___   \         |  |      /        \
//  ¯¯¯¯¯|   |¯¯¯¯¯¯      |  |¯¯¯¯¯¯¯¯¯¯¯      ¯¯¯¯¯|   |¯¯¯¯¯¯   |  |   |   |         ¯¯      /  \¯¯¯¯\  \
//		 |   |            |  |___________           |   |         |   ¯¯¯   /          __       \   \   ¯¯¯
//       |   |            |              |          |   |         |        \          |  |        \   \
//       |   |            |  |¯¯¯¯¯¯¯¯¯¯¯           |   |         |   | \   \         |  |       __ \   \
//       |   |            |  |___________           |   |         |   |  \   \        |  |     /  /   \   \
//	     |   |			  | 	         |          |   |         |   |   \   \       |  |     \   ¯¯¯¯   /
//        ¯¯¯              ¯¯¯¯¯¯¯¯¯¯¯¯¯¯            ¯¯¯           ¯¯¯      ¯¯¯        ¯¯       ¯¯¯¯¯¯¯¯¯¯
import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	heigth = 21
	width  = 12
)

var (
	//Bord you dont mess with it
	Bord [heigth][width]int // x,y
	// Fall for how many time he need to fall
	Fall int
	//Delay the player
	Delay time.Duration
	// Line the obj found
	Line int
	//Histroy show the histroy of the game
	Histroy deBug
	//Hide the place,lastpalce,dead
	Hide bool
)

type dir int

type obj struct {
	place     [8]int
	dir       dir
	lastplace [8]int
	d         bool
	dead      []int
}
type deBug struct {
	place     [200][]int
	lastplace [200][]int
	dir       []dir
	dead      []int
	counter   int
	activate  bool
	line      []int
}

func init() {
	edge()
	Delay = 3000 * time.Millisecond

}

func main() {
	fmt.Println("ready")

	T := time.Now()
	x := [8]int{5, 0, 6, 0, 5, 1, 6, 1}
	bob := newobj(x)
	view(bob)
	// key event
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		debug(bob)
		r, key, err := keyboard.GetKey()
		if err != nil {

			panic(err)
		} else if key == keyboard.KeyEsc {
			break
		}

		bob.lastplace = bob.place
		bob.dir = dir(r)
		view(bob)
		fmt.Printf("%v , %v\n", time.Since(T), Delay)
		if (time.Since(T)) >= Delay {
			go down(bob)
			T = time.Now()
		}

		if Histroy.activate {
			Histroy.Print()
			break
		}
	}
}

func view(b *obj) {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()

	add(b)
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
	if Hide {
		if len(b.dead) != 0 {
			fmt.Printf("pl: %v \nlp: %v \nded:%v \n", b.place, b.lastplace, b.dead[len(b.dead)-8:])
		} else {
			fmt.Printf("pl: %v \nlp: %v \n", b.place, b.lastplace)
		}
	}
	del(b)
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
	name := obj{p, dir(0), p, false, nil}

	return &name
}

func add(b *obj) {
	move(b)
	for out := 0; out < 8; out++ {
		if b.place[out] > 21 { // can't walk tour the edge
			b.place = b.lastplace
		}
	}

	if len(b.dead) != 0 {
		for rip := 0; rip < len(b.dead); rip += 2 { // add the dead
			Bord[b.dead[rip+1]][b.dead[rip]] = 3
		}
	}

	for i := 0; i < 8; i += 2 {

		if Bord[b.place[i+1]][b.place[i]] == 0  {
			Bord[b.place[i+1]][b.place[i]] = 2
		} else {

			if b.place[i+1] >= 20 || Bord[b.place[i+1]][b.place[i]] == 3 {
				b.d = true
			}

			if b.d {
				for _, val := range b.lastplace { //add to dead
					b.dead = append(b.dead, val)
					b.newlook()
					b.d = false
				}
			} else {
				b.place = b.lastplace
			}

		}
	}

}

func del(bob *obj) {
	for i := 0; i < 8; i += 2 {
		Bord[bob.place[i+1]][bob.place[i]] = 0
		Bord[bob.lastplace[i+1]][bob.lastplace[i]] = 0
	}
	if len(bob.dead) != 0 {
		for i := 0; i < len(bob.dead); i += 2 {
			Bord[bob.dead[i+1]][bob.dead[i]] = 0
		}
	}

}

func move(bob *obj) {
	switch int(bob.dir) {
	case 119, 87: //up change pos
		change(bob)
	case 83, 115: //down

		if bob.place[1] < 20 && bob.place[3] < 20 && bob.place[5] < 20 && bob.place[7] < 20 {

			for i := 1; i < 8; i += 2 {
				bob.place[i]++

			}
		} else {
			bob.place = bob.lastplace
		}

	case 100, 68: //left
		for i := 0; i < 8; i += 2 {
			bob.place[i]++
		}
	case 65, 97: //right
		for i := 0; i < 8; i += 2 {
			bob.place[i]--
		}
	case 72: //debug
		Histroy.activate = !Histroy.activate
	case 104:
		Hide = !Hide
	default:
		bob.dir = -1
	}

}

func (bop *obj) newlook() {
	Shape := map[int][8]int{
		0: [8]int{5, 0, 6, 0, 5, 1, 6, 1}, //shape cube
		1: [8]int{5, 0, 6, 0, 5, 1, 5, 2}, //shape r
		2: [8]int{5, 0, 4, 0, 5, 1, 5, 2}, //shape j
		3: [8]int{5, 0, 6, 0, 1, 5, 1, 4}, //shape S
		4: [8]int{5, 0, 4, 0, 6, 0, 5, 1}, //shape L
		5: [8]int{5, 0, 4, 0, 1, 5, 1, 6}, //shape Z
	}

	bop.place = Shape[rand.Intn(6)]
}

//debug
func debug(d *obj) {
	Histroy.dead = d.dead
	pl := []int{}
	lp := []int{}
	for _, val := range d.lastplace {
		lp = append(lp, val)
	}
	for _, val := range d.place {
		pl = append(pl, val)
	}
	Histroy.place[Histroy.counter] = append(Histroy.place[Histroy.counter], pl...)
	Histroy.lastplace[Histroy.counter] = append(Histroy.lastplace[Histroy.counter], lp...)
	Histroy.dir = append(Histroy.dir, d.dir)
	Histroy.line = append(Histroy.line, Line)
	Histroy.counter++

}
func (*deBug) Print() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Print("place : lastplace \n")
	for num := 0; num < len(Histroy.place); num++ {
		if Histroy.place[num] != nil {
			fmt.Printf("\n\t%v:%v\n", Histroy.place[num], Histroy.lastplace[num])
		} else {
			fmt.Print("\n")
			break
		}
	}
	fmt.Print("dead:\n")
	if len(Histroy.dead) != 0 {

		for id := range Histroy.dead {
			if (id+1)%8 == 0 {
				fmt.Print("\n")
			} else {
				if id%8 == 0 {
					fmt.Printf("\t%v\n", Histroy.dead[id:id+8])
				}

			}

		}

	}

	fmt.Printf("dir: \n\t%v\nLine:\n\t%v\n", Histroy.dir, Histroy.line)

}

func down(bob *obj) {
	if bob.place[1] < 20 && bob.place[3] < 20 && bob.place[5] < 20 && bob.place[7] < 20 {

		for i := 1; i < 8; i += 2 {
			bob.place[i]++

		}
	} else {
		bob.place = bob.lastplace
	}

}
func change(bob *obj) {

}
