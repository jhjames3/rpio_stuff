package main
import 	"github.com/warthog618/gpio"
import	"os"
import "os/signal"
import "fmt"
import "time"
import "strconv"
import "log"


//var pin1 *gpio.Pin // dah  when grounded ie false
//var pin2 *gpio.Pin // dit
var pin *gpio.Pin
var dah_history uint8 = 0
var dah_temp uint8 = 0
var test uint8 = 0
var testup uint8 = 0

type Mark int
const  (
	Unknown Mark = iota
	DIT
	DAH
	SPACE
	)

func main() {
	err := gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()

	pin = gpio.NewPin(gpio.GPIO13)
	pin.Input()
	pin.PullUp()
	//entered  := false
	// initialize test to stop
	testtemp, err := strconv.ParseInt("11111000",2,64)
	if err != nil {
		log.Fatal(err)
	}
	test = uint8(testtemp) // set up test looking for pulldown
	// init test for up
	testupTemp, err := strconv.ParseInt("00000111",2,64)
	if err != nil {
		log.Fatal(err)
	}
	testup = uint8(testupTemp) // set up test looking for pulldown
	// initialize pattern in history to be up
	dah_temp1, err := strconv.ParseInt("11111111",2,64)
	if err != nil {
		log.Fatal(err)
	}
	dah_temp =  uint8(dah_temp1)
	dah_history = uint8(dah_temp1)
	dah_history_view := int64(dah_history)
	fmt.Println(strconv.FormatInt(dah_history_view,2))
	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)
	
	// start Watch for Dah key pushed
	err = pin.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
		//if !entered { // fix me; ignore first drop
			//entered = true;
			//return
		//}
		start := time.Now()
		fmt.Println(start)
		fmt.Println("pressed_started")
		t := time.Now()
		// look for pressed pattern
		for {
			pressed := test_for_press_only()
			
			if pressed == 1 {
				fmt.Println("pressed stableized")
				break;
			}
		}
		fmt.Println("left 1st loop")
		// look for unpressed pattern
		for {
			pressed := test_for_up_only()
			//fmt.Println(pressed)
			if pressed == 0 {
				fmt.Println("pressed Up")
				t = time.Now()
				break;
			}
		}
		fmt.Println("left 2nd loop")
		diff := t.Sub(start)
		ns := diff.Nanoseconds()
		fmt.Println(start)
		fmt.Println(" Pin 13 is %v", pin.Read())
		fmt.Print(t)
		fmt.Print(" ")
		fmt.Print(int64(ns))
	})
	if err != nil {
		panic(err)
	}
	defer pin.Unwatch()

	// In a real application the main thread would do something useful here.
	// But we'll just run for a minute then exit.
	fmt.Println("Watching Pin 13...")
	select {
	case <-time.After(time.Minute):
	case <-quit:
	}

}

	// key is 2 if dah and a 1 if dit single state key 3 if neither closed
	// this returns Mark: DIT DAH or SPACE as 1, 2 or 3
//func  key_read() Mark {
	//var k0 = B2I(pin1.Read())
	//var k1 = B2I(pin2.Read())
	//k0 <<= 1
	//var key = (k0|k1)
	//switch key {
	//case 1: 
		//return DIT
	//case 2: 
		//return DAH
	//case 3: 
		//return SPACE
	//default:
		//return Unknown
	//}
//}

func B2I( b gpio.Level) uint8 {
	if b == true {
		return 1
	}
	return 0
}

func test_for_press_only() uint8{
	//dah_history_view := int64(dah_history)
	//fmt.Println(strconv.FormatInt(dah_history_view,2))
	var pressed uint8 = 0
	k0_temp := pin.Read()
	var k0 uint8 = B2I(k0_temp)
	dah_history <<= 1
	dah_history |= k0
	//dah_history_view := int64(dah_history)
	//fmt.Println(strconv.FormatInt(dah_history_view,2))
	if (dah_history == test) { // test = 11111000
		pressed = 1
		//dah_hist, err := strconv.ParseInt("00000000",2,64)
		//if err != nil {
			//log.Fatal(err)
		//}
		dah_history = uint8(0)
	}
	return pressed
}	


func test_for_up_only() uint8{
	//dah_history_view := int64(dah_history)
	//fmt.Println(strconv.FormatInt(dah_history_view,2))
	var pressed uint8 = 1
	k0_temp := pin.Read()
	var k0 uint8 = B2I(k0_temp)
	dah_history <<= 1
	dah_history |= k0
	//dah_history_view := int64(dah_history)
	//fmt.Println(strconv.FormatInt(dah_history_view,2))
	if (dah_history == dah_temp) { // test = 00000111
		pressed = 0
		dah_history = dah_temp
	}
	return pressed
}	
	


