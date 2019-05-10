package main
import 	"github.com/warthog618/gpio"
import	"os"
import "os/signal"
import "fmt"
import "time"


var pin1 *gpio.Pin // dah  when grounded ie false
var pin2 *gpio.Pin // dit
var marks Mark[]
var entered  = false
var entered1 = false
vaf entered2 = false
type Mark int
const  (
	Unknown Mark = iota
	DIT
	DAH
	SPACE
	)



func getNano(last time.Time) int64 {
	t := time.Now()
	diff := t.Sub(last)
	ns := diff.Nanoseconds()
	return ns
}

func watch_pin_goUp (pin gpio.Pin) {

}

func watch_pin_goDown (pin gpio.Pin) {
	// init times
	last := time.Now()
	t := last
	
	err = pin.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
		if !entered {
			entered = true;
			return
		}
		if !entered1 {
			last = time.Now()// real first bounce compare time to this one
			fmt.Println("pressed_started")
			fmt.Println(" Pin 13 is %v", pin.Read())
			entered1 = true;
		}
		if !entered2 {
			ns := getNano(last) 
			if ns > 2000000 {
				entered2 = true
				return;
			}
			for {
				ns = getNano(last) 
				if ns > 2000000 {
					break;
				}	
			}
			// fixme loop on dah/dit time
			fmt.Println(" Pin 13 is %v", pin.Read())
			fmt.Print(t)
			fmt.Print(" we have a dah")	
			mark := key_read()
			save_mark(mark)
			// fixme end loop on pin up
		}
		
	})
	if err != nil {
		panic(err)
	}
}

func save_mark (mark Mark) {
	marks.append(mark)
	fmt.Print(" marks: ")
	for _,element := range marks{
        //fmt.Println(index)
        fmt.Println(element)        
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
	
	entered  := false
	entered1 := false
	entered2 := false

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)
	
	defer pin.Unwatch()

	for {
		watch_pin_goDown(pin1)
	}

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
func  key_read() Mark {
	var k0 = B2I(pin1.Read())
	var k1 = B2I(pin2.Read())
	k0 <<= 1
	var key = (k0|k1)
	switch key {
	case 1: 
		return DIT
	case 2: 
		return DAH
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
		