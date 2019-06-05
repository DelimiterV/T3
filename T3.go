// T3.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var myDebug int = 0

type ppl struct {
	id    uint8
	state int
	mtime uint8
	floor uint8
}

var peoples []ppl

type lft struct {
	hold     []uint8
	state    uint8
	curfloor uint8
	aim      uint8
}

var lift lft

type flr struct {
	hold []uint8
}

var floors []flr

type prm struct {
	nPeoples []ppl
	mFloors  uint8
}

var param prm

type ask struct {
	iduser uint8
	floor  uint8
	cursec uint8
}

var asks []ask

type rz struct {
	iduser uint8
	second uint8
}

var street []rz

func showDecission() {
	for k := 0; k < len(street); k++ {
		e := 0
		for m := 0; m < len(street) && e == 0; m++ {
			if street[m].iduser == uint8(k) {
				fmt.Println(street[m].iduser, ",", street[m].second)
				e = 1
			}
		}
	}

}

func showCurState(cursec int, curlift int, curstate int) {
	fmt.Println("########   ", cursec, "sec  #####", curlift, "(", curstate, ")##########")
	for t := len(floors) - 1; t > 0; t-- {
		fmt.Println(t, "-------------------------------")
		if lift.curfloor == uint8(t) {
			fmt.Print("[")
			for u := 0; u < len(lift.hold); u++ {
				fmt.Print(lift.hold[u], ",")
			}
			fmt.Print("]")
			for m := 0; m < len(floors[t].hold); m++ {
				fmt.Print(floors[t].hold[m], ",")
			}
			fmt.Println("")
		} else {
			for m := 0; m < len(floors[t].hold); m++ {
				fmt.Print(floors[t].hold[m], ",")
			}
			fmt.Println("")
		}
		fmt.Println("==================================")
	}
	fmt.Print("STREET:")
	for a := 0; a < len(street); a++ {
		fmt.Print(street[a], "'")
	}
	fmt.Println("")
	fmt.Print("ASKfloor:")
	for o := 0; o < len(asks); o++ {

		fmt.Print(asks[o].floor, asks[o].cursec, "\"")

	}
	fmt.Println("")
}
func live(nfloors int) {
	ex := 0
	var zask ask
	var mstreet rz
	cursecond := 0
	floors = make([]flr, nfloors+1)
	lift.curfloor = 1
	lift.state = 0
	for ex == 0 {
		for u := 0; u < len(peoples); u++ {
			//корректировка людей
			if peoples[u].mtime == uint8(cursecond) {
				peoples[u].state = 1
				// добавляем заявку если нет людей
				if len(floors[peoples[u].floor].hold) == 0 {
					zask.cursec = uint8(cursecond)
					zask.floor = peoples[u].floor
					asks = append(asks, zask)
				}
				// добавляем на этаж
				floors[peoples[u].floor].hold = append(floors[peoples[u].floor].hold, peoples[u].id)
			}
		}
		//корректировка лифта
		switch lift.state {
		case 0: // 0-ожидание вызова

			if len(asks) > 0 {
				lift.aim = asks[0].floor
				// удаление заявки
				t := 0
				asks = append(asks[:t], asks[t+1:]...)
				lift.state = 1
			}

		case 1: // 1-движение к вызову
			if lift.aim > lift.curfloor {
				lift.curfloor++
			}
			if lift.aim == lift.curfloor {
				//достигнута цель
				lift.hold = append(lift.hold, floors[lift.curfloor].hold...)
				floors[lift.curfloor].hold = make([]uint8, 0)
				lift.state = 2
				//cursecond--
			}

		case 2: //2-движение на 1-й этаж
			if lift.curfloor > 1 {
				lift.hold = append(lift.hold, floors[lift.curfloor].hold...)
				floors[lift.curfloor].hold = make([]uint8, 0)
				for op := 0; op < len(asks); op++ {
					if lift.curfloor == asks[op].floor {
						asks = append(asks[:op], asks[op+1:]...)
						op--
					}
				}
				lift.curfloor--

			}
			if lift.curfloor == 1 {
				//cursecond--
				for p := 0; p < len(lift.hold); p++ {
					mstreet.iduser = lift.hold[p]
					mstreet.second = uint8(cursecond)
					street = append(street, mstreet)
				}
				lift.hold = make([]uint8, 0)
				lift.state = 0
				if len(asks) > 0 {
					lift.aim = asks[0].floor
					// удаление заявки
					t := 0
					asks = append(asks[:t], asks[t+1:]...)
					lift.state = 1
				}
			}

		}
		if myDebug == 1 {
			showCurState(cursecond, int(lift.curfloor), int(lift.state))
			time.Sleep(1 * time.Second)
		}

		cursecond++
		if len(street) == len(peoples) {
			ex = 1
		}
	}
	showDecission()
}
func readParams(filename string) bool {
	rez := true
	rstr := ""
	cnt := 0
	var uppl ppl
	dat, e := ioutil.ReadFile(filename)
	if e == nil {
		l := 0
		k := 0
		for z := 0; z < len(dat); z++ {
			switch dat[z] {
			case 0x0d:
				k = z
			case 0x0a:
				rstr = string(dat[l:k])
				if l == 0 {
					fmt.Println("Param", rstr)
					rstr = strings.Trim(rstr, " ")
					uh := strings.Split(rstr, " ")
					nf, _ := strconv.Atoi(uh[1])
					np, _ := strconv.Atoi(uh[0])
					param.mFloors = uint8(nf)
					param.nPeoples = make([]ppl, np)
				} else {
					fmt.Println("People", rstr)
					ut := strings.Split(rstr, " ")
					vr, _ := strconv.Atoi(ut[0])
					et, _ := strconv.Atoi(ut[1])
					uppl.floor = uint8(et)
					uppl.id = uint8(cnt)
					uppl.mtime = uint8(vr)
					uppl.state = 0
					peoples = append(peoples, uppl)
					cnt++
				}
				//fmt.Println(rstr)
				l = z + 1
			}
		}
	} else {
		rez = false
	}
	return rez
}
func main() {
	fmt.Println("Usage: T3 -debugMode=0 или 1 \r\n")

	debugPtr := flag.Int("debugMode", 0, "an int")

	flag.Parse()

	if *debugPtr == 1 {
		myDebug = 1
	} else {
		myDebug = 0
	}

	fmt.Println("Лифт! Нумерация(ID) людей с нуля")
	if readParams("samples.txt") {
		live(int(param.mFloors))
	}

}
