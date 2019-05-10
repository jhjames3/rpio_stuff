package main
import 	"github.com/warthog618/gpio"
import	"os"
import "os/signal"
import "fmt"
import "time"


var pin1 *gpio.Pin // dah  when grounded ie false
var pin2 *gpio.Pin // dit
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
	
	//pin1 = gpio.NewPin(gpio.GPIO13) //pin 13 
	//pin1.Input()
	
	//pin2 = gpio.NewPin(gpio.GPIO27) // pin 27 
	//pin2.Input()
	
	//var pin1_state = pin1.Read()
	//var pin2_state = pin2.Read()
	
	//println(pin1_state)
	//println(pin2_state)
	
	//var key = key_read()
	//println(key)
	
	pin := gpio.NewPin(gpio.GPIO13)
	pin.Input()
	pin.PullUp()
	
	entered  := false
	entered1 := false

	// capture exit signals to ensure resources are released on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)
	last := time.Now()
	t := last
	
	err = pin.Watch(gpio.EdgeFalling, func(pin *gpio.Pin) {
		if !entered {
			last = time.Now() # first bounce
			entered = true;
			fmt.Println("pressed_started")
			fmt.Println(" Pin 13 is %v", pin.Read())
			return
		} 
		if !entered {
			t = time.Now()
			//sec := t.UnixNano()
			diff := t.Sub(last)
			ns := diff.Nanoseconds()
			fmt.Print("another bounce")
			fmt.Println(int64(ns))
			return

		}
		
		
		
		
		//fmt.Println(" Pin 13 is %v", pin.Read())
		fmt.Print(t)
		fmt.Print(" ")
		
		
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
		

//func key_loop(mark long) {
	
	//var state uint8 = 3
	//var last uint8
	//var ultimatic uint8
	//var staged uint8 0
	//var mcode=0x80
	//var ret uint8 
	
	//var key uint8 = key_read()
	//switch state {
		//case 1: // waiting until read for read
			
			//break
		//case 2: // waiting and reading
			
			//fallthrough
		//case 3: // idle, spacing
			
			
			
			
			//break
		//case 4: 
			
			//break
		//case 5:
			//fallthrough
		//case 6:
			
			
			
		
		//return ret
	//}
	
//}
