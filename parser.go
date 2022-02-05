package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"image"
	"image/color"
	"image/png"
	"runtime"
)

var height = 0
var width = 0
var RGBA = false

var numbers = map[string]int {
	"zero": 0,
	"jeden": 1,
	"dwa": 2,
	"trzy": 3,
	"cztery": 4,
	"pięć": 5,
	"sześć": 6,
	"siedem": 7,
	"osiem": 8,
	"dziewięć": 9,
	"dziesięć": 10,
	"jedenaście": 11,
	"dwanaście": 12,
	"trzynaście": 13,
	"czternaście": 14,
	"piętnaście": 15,
	"szesnaście": 16,
	"siedemnaście": 17,
	"osiemnaście": 18,
	"dziewiętnaście": 19,
	"dwadzieścia": 20,
	"trzydzieści": 30,
	"czterdzieści": 40,
	"pięćdziesiąt": 50,
	"sześćdziesiąt": 60,
	"siedemdziesiąt": 70,
	"osiemdziesiąt": 80,
	"dziewięćdziesiąt": 90,
	"sto": 100,
	"dwieście": 200,
	"trzysta": 300,
	"czterysta": 400,
	"pięćset": 500,
	"sześćset": 600,
	"siedemset": 700,
	"osiemset": 800,
	"dziewięćset": 900,
	"tysiąc": 1000,
	"tysiące": 1000,
	"tysięcy": 1000,
	
}

func changeSTRtoINT(line string) (int){
	var str []byte
	num := 0
	line += " "
	for i := 0; i < len(line); i++{
		if line[i] != 32 {
			str = append(str, line[i])		
		}else{
			if len(str) != 0 {
				key, err := numbers[string(str)]
				if err == false {
					fmt.Println("Error: Liczba")
					os.Exit(1)
				}
				// TODO sprawdzania czy liczba jest poprawnie gramatycznie
				if string(str) == "tysięcy" || string(str) == "tysiące" {
					num *= key
				}else{
					num += key
				}
			}
			str = nil
		}
	}
	return num
}

func parse_flags(line string){
	if string(line)[0:4] != "CIF:" {
		fmt.Println("Error: Flaga")
		os.Exit(1)
	}
	fmt.Printf("Flagi:%s\n", strings.Trim(fmt.Sprint(strings.Fields(line)[1:]), "[]"))
}

func parse_version(line string){
	if string(line)[0:6] != "WERSJA" {
		fmt.Println("Error: Wersja")
		os.Exit(1)
	}
	fmt.Printf("Wersja: %d\n",changeSTRtoINT(string(line)[6:]))
}

func parse_data(line string){
	var str []byte
	counter := 0
	for i := 0; i < len(line); i++ {
		if string(line[i]) != "," {
			str = append(str, line[i])
		}else {
			if len(str) != 0 {
				if string(line[i-1]) == " " {
					fmt.Println("Error: Space before \",\"")
					os.Exit(1)
				} 
				counter++
				strSPACE := strings.Fields(string(str))
				if counter == 1 {
					if string(str)[0:7] != "ROZMIAR" || strSPACE[1] != "szerokość:" {
						fmt.Println("Error: Rozmiar")
						os.Exit(1)
					}
					var widthSTR string
					for i := 2; i < len(strSPACE); i++ {
						widthSTR += strSPACE[i] + " "
					}
					width = changeSTRtoINT(widthSTR)
				}else if counter == 2 {
					if strSPACE[0] != "wysokość:" {
						fmt.Println("Error: Rozmiar")
						os.Exit(1)
					}
					var heightSTR string
					for i := 1; i < len(strSPACE); i++ {
						heightSTR += strSPACE[i] + " "
					}
					height = changeSTRtoINT(heightSTR)
				}
			}
			str = nil
		}
	}
	strSPACE := strings.Fields(string(str))
	if strSPACE[0] != "bitów_na_piksel:" {
		fmt.Println("Error: bitów_na_piksel")
		os.Exit(1)
	}
	var rgbaSTR string
	for i := 1; i < len(strSPACE); i++ {
		rgbaSTR += strSPACE[i] + " "
	}
	if changeSTRtoINT(rgbaSTR) == 24 {
		RGBA = false
	}else if changeSTRtoINT(rgbaSTR) == 32 {
		RGBA = true
	}else{
		fmt.Println("Error: bitów_na_piksel")
		os.Exit(1)
	}

}

func ImageSet(line []string, start int, end int, img *image.NRGBA, row bool, channel chan bool) {
	x := 0
	y := 0
	tmp := 0
	var tmp2 int
	if row == true {
		x = 0
		y = start
		tmp = height 
		tmp2 = height * end
	}else{
		x = start
		y = 0
		tmp = 1
		tmp2 = height*width
	}
	for k := start * tmp; k < tmp2; k++ {
		checkSTR := strings.Fields(line[k])
		if len(checkSTR) > 1 {
			for j := 0; j < len(checkSTR); j++ {
				if checkSTR[j] == ";" {
					fmt.Println("Error: Space before \";\"")
					os.Exit(1)
				}
			}
		}
		str := strings.Split(line[k], ";") 
		if len(str) <= 1 {
			continue
		}
		if x + 1 > width && row == true{
			x = 0
			y++
		}
		if x + 1 > end && row == false {
			x = start
			y++
			k += width - end + start
			if x == start { x++; continue }
		}
		if RGBA == false {
			if !(len(str) >= 3) {
				fmt.Println("Error: nil value detected")
				os.Exit(1)
			}
			img.Set(x, y, color.NRGBA{uint8(changeSTRtoINT(str[0])), uint8(changeSTRtoINT(str[1])), uint8(changeSTRtoINT(str[2])), 0xff})
		}else{
			if !(len(str) >= 4) {
				fmt.Println("Error: nil value detected")
				os.Exit(1)
			}
			img.Set(x, y, color.NRGBA{uint8(changeSTRtoINT(str[0])), uint8(changeSTRtoINT(str[1])), 
			uint8(changeSTRtoINT(str[2])), uint8(changeSTRtoINT(str[3]))})
		}
		x++;
	}
	channel <- true
}
func makeImage(line []string) {
	up := image.Point{0,0}
	down := image.Point{width,height}
	img := image.NewNRGBA(image.Rectangle{up, down})
	cpu := runtime.NumCPU() + 1
	var n int
	var row bool
	if width > height {
		n = width / cpu
		row = false
	}else {
		n = height / cpu
		row = true
	}
	
	var channels []chan bool
	for i := 0; i < cpu-1; i++ {
		channels = append(channels, make(chan bool))
		go ImageSet(line, (i*n), (i+1)*n+1, img, row, channels[i])
	}
	last := n * cpu - 1
	if n * cpu != height || n * cpu != width {
		if row == false {
			last = width
		}else {
			last = height
		}
	}
	channels = append(channels, make(chan bool))
	go ImageSet(line, (cpu-1)*n, last+1, img, row, channels[cpu-1])

	for j := 0; j < cpu; j++ {
		<-channels[j]	
	}
	file, _ := os.Create(os.Args[2])
	png.Encode(file, img) 
}

func parser(line string, line_number int){
	if line_number == 1 {
		parse_flags(line)
	}else if line_number == 2 {
		parse_version(line)
	}else if line_number == 3 {
		parse_data(line)
	}
}

func parse_file(filename string){
	file_open, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	count := 0
	newline_counter := 0
	file_read := strings.Split(string(file_open), "\n")
	for i := 0; i < len(file_read); i++ {
		if len(file_read[i]) > 0{
			count++
			parser(file_read[i], count)
		}else{
			newline_counter++
		}
		if count == 3 {
			break
		}
	}
	checkMETA := false
	checkDATA := false
	if newline_counter != 0 {
		count += newline_counter
	}
	for j := count; j < len(file_read); j++ {
		// file_read[j] += " "
		str := strings.Fields(file_read[j]) 
		if len(str) > 0 {
			if str[0] == "METADANE" {
				checkMETA = true
				if len(str[1]) > 1 && len(str[2:]) < 1 {
					fmt.Printf("Klucz: %s Wartość: %s\n", str[1], file_read[j+1])
					j++ 
				}else{
					fmt.Printf("Klucz: %s Wartość: %s\n", str[1], strings.Trim(fmt.Sprint(str[2:]), "[]")) 
				}
			}else {
				makeImage(file_read[j:])
				checkDATA = true
				break;
			}
		}
	}
	if checkMETA != true || checkDATA != true {
		fmt.Println("Error: No data or metadata")
		os.Exit(1)
	}
}
