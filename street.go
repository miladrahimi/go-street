package main

import (
	"time"
	"github.com/nsf/termbox-go"
	"os"
	"math/rand"
	"strconv"
)

const (
	DefaultStreetWidth  = 80
	DefaultStreetHeight = 12
)

type Street struct {
	width    int
	height   int
	pixels   [][]bool
	interval int
	ticker   *time.Ticker
}

func NewStreet(height, width int) *Street {
	pixels := make([][]bool, height)

	for i := int(0); i < height; i++ {
		pixels[i] = make([]bool, width)

		for j := int(0); j < width; j++ {
			pixels[i][j] = false
		}
	}

	ticker := time.NewTicker(time.Second)

	return &Street{
		width:    width,
		height:   height,
		pixels:   pixels,
		interval: 1000,
		ticker:   ticker,
	}
}

func (street *Street) NextState() {
	for i := int(1); i < street.height-1; i++ {
		for j := int(street.width - 1); j > 0; j-- {
			street.pixels[i][j] = street.pixels[i][j-1]
		}

		p := street.pixels

		switch {
		case p[i][1] == true && p[i][2] == true && p[i][3] == true && p[i][4] == true && p[i][5] == true:
			p[i][0] = generatePixel(40)
		case p[i][1] == true && p[i][2] == true && p[i][3] == true && p[i][4] == true && p[i][5] == false:
			p[i][0] = generatePixel(70)
		case p[i][1] == true && p[i][2] == true && p[i][3] == true && p[i][4] == false:
			p[i][0] = generatePixel(90)
		case p[i][1] == true && p[i][2] == true && p[i][3] == false:
			p[i][0] = true
		case p[i][1] == true && p[i][2] == false:
			p[i][0] = true
		case p[i][1] == false && p[i][2] == false && p[i][3] == false && p[i][4] == false && p[i][5] == false:
			p[i][0] = generatePixel(20)
		case p[i][1] == false && p[i][2] == false && p[i][3] == false && p[i][4] == false && p[i][5] == true:
			p[i][0] = generatePixel(15)
		case p[i][1] == false && p[i][2] == false && p[i][3] == false && p[i][4] == true:
			p[i][0] = generatePixel(5)
		case p[i][1] == false && p[i][2] == false && p[i][3] == true:
			p[i][0] = false
		case p[i][1] == false && p[i][2] == true:
			p[i][0] = false
		default:
			p[i][0] = generatePixel(35)
		}
	}
}

func (street *Street) Draw() {
	for i := 0; i < street.height; i++ {
		for j := 0; j < street.width; j++ {
			if street.pixels[i][j] {
				termbox.SetCell(j+2, i+2, 0x2588, termbox.ColorWhite, termbox.ColorBlack)
			} else {
				termbox.SetCell(j+2, i+2, 0x2000, termbox.ColorBlack, termbox.ColorBlack)
			}
		}
	}
}

func (street *Street) DrawWalls(c termbox.Attribute) {
	for i := 1; i < street.width+3; i++ {
		termbox.SetCell(i, 1, 0x2000, c, c)
		termbox.SetCell(i, street.height+2, 0x2000, c, c)
	}

	for i := 0; i < street.height; i++ {
		termbox.SetCell(1, i+2, 0x2000, c, c)
		termbox.SetCell(street.width+2, i+2, 0x2000, c, c)
	}
}

type Level struct {
	number int
	street *Street
}

func NewLevel(street *Street) *Level {
	return &Level{number: 1, street: street}
}

func (levelBoard *Level) Draw() {
	text := " LEVEL: " + strconv.Itoa(levelBoard.number) + " "
	i := 3

	for _, c := range text {
		termbox.SetCell(i, levelBoard.street.height+2, c, termbox.ColorWhite, termbox.ColorBlack)
		i++
	}
}

type Player struct {
	x      int
	y      int
	street *Street
	level  *Level
}

func NewPlayer(street *Street, level *Level) *Player {
	return &Player{x: street.width / 2, y: street.height - 1, street: street, level: level}
}

func (player *Player) Init() {
	go player.street.DrawWalls(termbox.ColorWhite)
	go player.level.Draw()

	for i := 0; i < player.street.width*3; i++ {
		player.street.NextState()
	}
}

func (player *Player) Draw() {
	termbox.SetCell(player.x+2, player.y+2, 0x2588, termbox.ColorMagenta, termbox.ColorBlack)
}

func (player *Player) checkWinning() {
	if player.y == 0 {
		player.x = player.street.width / 2
		player.y = player.street.height - 1

		if player.street.interval > 200 {
			player.street.interval -= 200
		} else {
			player.street.interval /= 2
		}

		player.street.ticker = time.NewTicker(time.Duration(player.street.interval) * time.Millisecond)
		player.level.number++

		go func() {
			player.street.DrawWalls(termbox.ColorGreen)
			time.Sleep(time.Second / 7)
			player.street.DrawWalls(termbox.ColorWhite)
			time.Sleep(time.Second / 10)
			player.street.DrawWalls(termbox.ColorGreen)
			time.Sleep(time.Second / 7)
			player.street.DrawWalls(termbox.ColorWhite)
		}()
	}
}

func (player *Player) checkLoosing() {
	if player.street.pixels[player.y][player.x] == true {
		time.Sleep(time.Second / 20)

		player.x = player.street.width / 2
		player.y = player.street.height - 1

		player.street.interval = 1000
		player.street.ticker = time.NewTicker(time.Duration(player.street.interval) * time.Millisecond)
		player.level.number = 1

		go func() {
			player.street.DrawWalls(termbox.ColorRed)
			time.Sleep(time.Second / 7)
			player.street.DrawWalls(termbox.ColorWhite)
			time.Sleep(time.Second / 10)
			player.street.DrawWalls(termbox.ColorRed)
			time.Sleep(time.Second / 7)
			player.street.DrawWalls(termbox.ColorWhite)
		}()
	}
}

func (player *Player) listenToKeyboard() {
	for {
		event := termbox.PollEvent()

		go func() {
			if event.Type == termbox.EventKey {
				switch event.Key {
				case termbox.KeyCtrlC:
					os.Exit(1)
				case termbox.KeyArrowUp:
					if player.y > 0 {
						player.y--
					}

					player.checkWinning()
				case termbox.KeyArrowDown:
					if player.y < player.street.height-1 {
						player.y++
					}
				case termbox.KeyArrowLeft:
					if player.x > 0 {
						player.x--
					}
				case termbox.KeyArrowRight:
					if player.x < player.street.width-1 {
						player.x++
					}
				}
			}
		}()
	}
}

func generatePixel(percent int) bool {
	if rand.Intn(100) > (100 - percent) {
		return true
	} else {
		return false
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	termbox.SetInputMode(termbox.InputEsc)

	street := NewStreet(DefaultStreetHeight, DefaultStreetWidth)
	level := NewLevel(street)
	player := NewPlayer(street, level)

	defer street.ticker.Stop()
	defer termbox.Close()

	go player.Init()
	go player.listenToKeyboard()

	go func() {
		for {
			<-street.ticker.C
			street.NextState()
		}
	}()

	go func() {
		for {
			street.Draw()
			player.Draw()
			player.checkLoosing()
			level.Draw()
			termbox.SetCursor(1, 1)
			termbox.HideCursor()
			termbox.Flush()
		}
	}()

	exit := make(chan bool)
	<-exit
}
