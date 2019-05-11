package main
import 	"github.com/warthog618/gpio"
import	"os"
import "os/signal"
import "fmt"
import "time"


var pin1 *gpio.Pin // dah  when grounded ie false
var pin2 *gpio.Pin // dit

var entered  = false
var entered1 = false
var WAITFORBOUNCE int64 = 5000000

type Mark int
const  (
	Unknown Mark = iota
	DIT
	DAH
	SPACE
	)
var marks [30]Mark
var marIndex = 0

func save_mark (mark Mark) {
	marks[marIndex] = mark
	marIndex++        
    } 


func getNano(last time.Time) int64 {
	t := time.Now()
	diff := t.Sub(last)
	ns := diff.Nanoseconds()
	return ns
}

func watch_pin_goUp (pin *gpio.Pin) {

}

func watch_pin_goDown (pin *gpio.Pin, err error ) {
	// init times
	last := time.Now()
	
	err = pin.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
		if !entered { // ignore first one seems to do a false 3 on start up
			//fmt.Println(" 0 key is %v", key())
			key := key()
			if key != 3 {
				last = time.Now()// real first bounce compare time to this one
				fmt.Println("pressed_started")
				//fmt.Println(" 1 key is %v", key())
				entered = true;
				return
			}
		}
		if !entered1 {
			key := key()
			if key == 3 {
				entered  = false
				entered1 = false
				pin1.PullUp()
				pin2.PullUp()
				return
			}
			for {
				ns := getNano(last) 
				if ns > WAITFORBOUNCE {
					break;
				}	
			}
			// fixme loop on dah/dit time
			fmt.Println(last)
			t := time.Now()
			fmt.Println(t)
			//
			mark := key_read()
			if mark == DIT {
				fmt.Println(" we have a dit")
			}
			if mark == DAH {
				fmt.Println(" we have a dah")
			}
			save_mark(mark)
			fmt.Println(marks)
			// fixme end loop on pin up
			entered  = false
			entered1 = false
			pin1.PullUp()
			pin2.PullUp()
			
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
	
	//entered  := false
	//entered1 := false
	//entered2 := false

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)
	
	defer pin1.Unwatch()
	defer pin2.Unwatch()

	
	watch_pin_goDown(pin1, err)
	watch_pin_goDown(pin2, err)
	

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
		
