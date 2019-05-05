package main
import 	"github.com/stianeikeland/go-rpio"
//import	"reflect"
//import "fmt"


var pin1 rpio.Pin // dah  when grounded ie =0
var pin2 rpio.Pin // dit
const (
	Unknown = iota
	DIT
	DAH
	SPACE
	)

func main() {
	rpio.Open()
	defer rpio.Close()
	
	pin1 = rpio.Pin(13)
	//fmt.Println(reflect.TypeOf(pin1))
	pin2 = rpio.Pin(27)
	pin1.Input()
	pin2.Input()
	
	var pin1_state = rpio.ReadPin(pin1)
	var pin2_state = rpio.ReadPin(pin2)
	
	println(pin1_state)
	println(pin2_state)
	
	var key = key_read()
	println(key)
}

// key is 2 if dah and a 1 if dit single state key 3 if neither closed
// this returns DIT DAH or SPACE as 1, 2 or 3
func  key_read() rpio.State {
	var k0 = rpio.ReadPin(pin1)
	var k1 = rpio.ReadPin(pin2)
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
		

func key_loop(mark long) {
	
	var state uint8 = 3
	var last uint8
	var ultimatic uint8
	var staged uint8 0
	var mcode=0x80
	var ret uint8 
	
	var key uint8 = key_read()
	switch state {
		case 1: // waiting until read for read
			
			break
		case 2: // waiting and reading
			
			fallthrough
		case 3: // idle, spacing
			
			
			
			
			break
		case 4: 
			
			break
		case 5:
			fallthrough
		case 6:
			
			
			
		
		return ret
	}
	
}
