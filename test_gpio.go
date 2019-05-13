package main
import 	"github.com/warthog618/gpio"
import	"os"
import "os/signal"
import "fmt"
import "time"


var pin1 *gpio.Pin // dah  when grounded ie false
var pin2 *gpio.Pin // dit

var entered  = false //down
var entered1 = false
var WAITFORBOUNCE int64 = 22000000
var last time.Time
var WPM int64 = 10
var STD_1200 int64 = 1200
var DIT_TIME_ms int64 = STD_1200/WPM
var STD_3 int64 = 3
var DAH_TIME_ms int64 = STD_3*DIT_TIME_ms

type Mark int
const  (
	Unknown Mark = iota
	DAH
	DIT
	SPACE
	)
var marks [30]Mark
var marIndex = 0

func save_mark (mark Mark) {
	if marIndex+1 < len(marks) {
		marks[marIndex] = mark
		marIndex++  
		fmt.Println(marks)
	}      
} 
    
func set_last() {
	last = time.Now()
	fmt.Print("last: ")
	fmt.Println(last)
}

func getMS(last time.Time) int64 {
	t := time.Now()
	diff := t.Sub(last)
	ns := diff.Nanoseconds()
	ms := ns/1000000
	return ms
}

func waitForDitTimeDown() bool {
	for {
		key := key()
		if key == 1 { // Dit key down
			ms := getMS(last) 
			if ms > DIT_TIME_ms {
				set_last()
				fmt.Print(ms)
				fmt.Print(" ")
				fmt.Print(DIT_TIME_ms)
				fmt.Println(" continue dit")
				return true
			}	
		}
	}
	return false
}

func waitForDahTimeDown() bool {
	for {
		key := key()
		fmt.Println(key)
		if key == 2 { // Dah key down 
			ms := getMS(last) 
			if ms > DAH_TIME_ms {
				set_last()
				fmt.Print(ms)
				fmt.Print(" ")
				fmt.Print(DAH_TIME_ms)
				fmt.Print(ms)
				fmt.Println(" continue dah")
				return true
			}	
		}
		break
	}
	return false
}

func getNano(last time.Time) int64 {
	t := time.Now()
	diff := t.Sub(last)
	ns := diff.Nanoseconds()
	//fmt.Println(last)
	//fmt.Println(t)
	//fmt.Println(ns)
	//fmt.Println("--------")
	return ns
}

func waitForStableUp() bool {
	for {
		key := key()
		if key == 3 {
			ns := getNano(last) 
			if ns > WAITFORBOUNCE {
				fmt.Println("key up")
				//entered = false
				//entered1 = false
				return true
			}
			break;	
		}
		break
	}
	return false
}

func waitForStableDown() bool {
	for {
		key := key()
		fmt.Println(key)
		if key != 3 {
			ns := getNano(last) 
			if ns > WAITFORBOUNCE {
				fmt.Print(ns)
				fmt.Println(" key down")
				return true
			}
		} else {
			fmt.Println(key)
			break
		}
	}
	fmt.Println("waitForStableDown: key up")
	return false
}

func watch_pin_goBoth (pin *gpio.Pin, err error ) {
	
	err = pin.Watch(gpio.EdgeBoth, func(pin *gpio.Pin) {
		Start: 
		if !entered { // possible 3 on start up
			key := key()
			if key != 3 {
				set_last()// real first bounce compare time to this one
				fmt.Println("pressed_started")
				entered = waitForStableDown() // if true we have key down for long time
				if !entered {
					fmt.Println("back to start !entered ")
					fmt.Print("entered: ")
					fmt.Println(entered)
					goto Start
				} else {
					fmt.Print("entered: ")
					fmt.Println(entered)
				}
			} else { //== 3
				//ie both up
				entered = !waitForStableUp() // we have key up?
				return
			}
		}
		// down for "long" time
		// possible bounch
		if !entered1 {
			fmt.Println("entered1 false")
			key := key()
			if key != 3 {
				mark := key_read()
				if mark == DIT {
					for {
						save_mark(mark)
						fmt.Println(" we have a dit")
						if waitForDitTimeDown() { // if timed out return true
							continue // another dit
						} else {
							fmt.Println(" end dit")
							entered  = false
							entered1 = false
							pin1.PullUp()
							pin2.PullUp()
							return
						}
					} 
					return
				} else if mark == DAH {
					for {
						save_mark(mark)
						fmt.Println(" we have a dah")
						if waitForDahTimeDown() { // if timed out return true
							continue // another dah
						} else {
							fmt.Println(" end dah")
							entered  = false
							entered1 = false
							pin1.PullUp()
							pin2.PullUp()
							return
						}
					}
					return
				}
				fmt.Print("not a dit or a Dah")
			} else  { // key == 3 up start over?
				fmt.Println(" start over")
				entered  = false
				entered1 = false
				pin1.PullUp()
				pin2.PullUp()
				goto Start
			}
			// start next mark timing
			fmt.Println(marks)
			// fixme end loop on pin up with time
			
			entered  = false
			entered1 = false
			pin1.PullUp()
			pin2.PullUp()
			goto Start
		}
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	
	err := gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()
	
	pin1 = gpio.NewPin(gpio.GPIO13) // dah
	pin1.Input()
	pin1.PullUp()
	pin2 = gpio.NewPin(gpio.GPIO27) // dit
	pin2.Input()
	pin2.PullUp()
	

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)
	
	defer pin1.Unwatch()
	defer pin2.Unwatch()

	
	watch_pin_goBoth(pin1, err)
	watch_pin_goBoth(pin2, err)
	

	// In a real application the main thread would do something useful here.
	// But we'll just run for a minute then exit.
	fmt.Println("Watching Pin 13...")
	select {
	case <-time.After(time.Minute):
	case <-quit:
	}
}

func key() int8 {
	var k0 = B2I(pin1.Read())
	var k1 = B2I(pin2.Read())
	k0 <<= 1
	var key = (k0|k1)
	return key
}

// key is 2 if dah and a 1 if dit single state key 3 if neither closed
// this returns Mark: DIT DAH or SPACE as 1, 2 or 3
func  key_read() Mark {
	var key = key()
	switch key {
	case 1: 
		return DAH
	case 2: 
		return DIT
	case 3: 
		return SPACE
	default:
		return Unknown
	}
}

func B2I( b gpio.Level) int8 {
	if b == true {
		return 1
	}
	return 0
}
		
