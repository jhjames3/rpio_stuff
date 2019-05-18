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
var last time.Time
const WAITFORBOUNCE int64 = 1700000

const WPM int64 = 10
const STD_1200 int64 = 1200
const STD_1300 int64 = 1400
const DIT_TIME_ms int64 = STD_1200/WPM
const DIT1_TIME_ms int64 = STD_1300/WPM
const STD_3 int64 = 3
const DAH_TIME_ms int64 = STD_3*DIT_TIME_ms
const LETTERTIME int64 = DAH_TIME_ms
const STD_7 int64 = 7
const WORDTIME int64 = STD_7*DIT_TIME_ms
const SIZEMARKS = 50

type Mark int
const  (
	Unknown Mark = iota
	DAH
	DIT
	SPACE // end word
	LETTER // END LETTER
	)
var marks [SIZEMARKS]Mark
var marIndex = 0
var mark Mark = Unknown

func clearMarks() {
	for i,_ := range marks {
		marks[i] = Unknown
	}
}

func save_mark (mark Mark) {
	if marIndex+1 < len(marks) {
		marks[marIndex] = mark
		marIndex++  
		fmt.Println(marks)
	} else {
		fmt.Print("marks index too big")
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
		if key == 2 { // Dit key down
			ms := getMS(last) 
			if ms > DIT1_TIME_ms {
				//fmt.Print(ms)
				//fmt.Print(" ")
				//fmt.Print(DIT_TIME_ms)
				//fmt.Println(" continue dit ms")
				return true
			} else {
				continue
			}	
		} else {
			break
		}
	}
	return false
}

func waitForDahTimeDown() bool {
	for {
		key := key()
		ms := getMS(last) 
		if key == 1 { // Dah key down 
			if ms > DAH_TIME_ms {
				// debug
				//fmt.Print(ms)
				//fmt.Print(" ")
				//fmt.Print(DAH_TIME_ms)
				//fmt.Println(" ms")
				//fmt.Println(" continue dah ms")
				return true
			} else {
				continue
			}	
		} else {
			break
		}
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

func waitForLetterUp() bool {
	for {
		key := key()
		if key == 3 {
			ns := getNano(last) 
			if ns > LETTERTIME {
				fmt.Println("key up")
				//entered = false
				//entered1 = false
				return true
			} else {
				continue
			}	
		} else {
			break 
		}
	}
	return false
}

func waitForWordUP() bool {
	for {
		key := key()
		if key == 3 {
			ns := getNano(last) 
			if ns > WORDTIME {
				fmt.Println("key up")
				//entered = false
				//entered1 = false
				return true
			} else {
				continue
			}	
		} else {
			break 
		}
	}
	return false
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
			} else {
				continue
			}	
		} else {
			break 
		}
	}
	return false
}

func waitForStableDown() bool {
	for {
		key := key()
		//fmt.Println(key)
		if key != 3 {
			ns := getNano(last) 
			if ns > WAITFORBOUNCE {
				fmt.Print(ns)
				fmt.Println(" key down")
				return true
			}
		} else {
			//fmt.Println(key)
			continue
		}
	}
	fmt.Println("waitForStableDown: key up")
	return false
}

func watch_pin_goBoth (pin *gpio.Pin, err error ) {
	var isWord = false
	err = pin.Watch(gpio.EdgeBoth, func(pin *gpio.Pin) {
		entered  = false
		entered1 = false
		//Start: 
		if !entered { // possible 3 on start up
			key := key()
			mark = key_read()
			if key == 1 || key == 2 { // dit or dah down
				fmt.Println("gitgit")
				set_last()// real first bounce compare time to this one (up or down)
				ns := getNano(last) 
				fmt.Print(ns)
				fmt.Println(" key down")
				entered = waitForStableDown() // if true we have key down for long time
				
				if !entered {
					//fmt.Println("back to start !entered ")
					//fmt.Print("entered: ")
					//fmt.Println(entered)
					//goto Start
					return
				} else {// true
					
					save_mark(mark)
					//fmt.Print("first mark: ")
					//fmt.Println(entered)
					goto CheckNext
				}
			} else { //== 3
				//ie both up
					fmt.Println("both up ")
				 // we have key up wait for next down 
					return
			}
		}
		// down for "long" time
		// possible bounch
		CheckNext:
		if !entered1 {
			fmt.Println("entered1 false")
			key := key()
			if key != 3 {
				mark = key_read()
				fmt.Print(mark)
				fmt.Println(" mark")
				if mark == DIT {
					// start For for number of marks
					//------------------------------
					for {
						//fmt.Println(" we have a start dit")
						if waitForDitTimeDown() { // if timed out return true
							//fmt.Println(" we continue dit")
							save_mark(mark)
							set_last()
							continue // another dit in this time down
						} else {
							isPossibleLetter := waitForLetterUp()
							if (isPossibleLetter) {
								isWord := waitForWordUP() 
									if isWord {
										save_mark(SPACE)
										message := createMessageForWord()
										fmt.Println(message)
										clearMarks()
										fmt.Println(marks)
										//sendTcp(message)
										return
									} else {
										save_mark(LETTER)
										return
									}
								}
							}
							fmt.Println(" end dit")
							// entered  = false
							// entered1 = false
							// pin1.PullUp()
							// pin2.PullUp()
							return
						}
					} // end for loop
				} else if mark == DAH {
					for {
						//fmt.Println(" we have a start dah")
						if waitForDahTimeDown() { // if timed out return true
							
							//fmt.Println(" we continue dit")
							save_mark(mark)
							set_last()
							continue // another dah
						} else {
							fmt.Println("Dah Up check letter ")
							isPossibleLetter := waitForLetterUp()
							if (isPossibleLetter) {
								fmt.Println("Dah Up check word ")
								isWord := waitForWordUP()
								if isWord {
									fmt.Println("Is word ")
									save_mark(SPACE)
									message := createMessageForWord()
									fmt.Println(message)
									//sendTcp(message)
									return
								} else {
									fmt.Println("Is letter ")
									save_mark(LETTER)
									return
								}
							}
							if isWord {
								save_mark(SPACE)
								message := createMessageForWord()
								fmt.Println(message)
								//sendTcp(message)
								return
							} else {
								save_mark(LETTER)
								return
							}
						}
						if isWord {
							save_mark(SPACE)
							message := createMessageForWord()
							fmt.Println(message)
							clearMarks()
							fmt.Println(marks)
							//sendTcp(message)
							return
						} else {
							save_mark(LETTER)
							return
						}
						fmt.Println(" end dah")
						return
					} // end for loop
				}
				fmt.Print("not a dit or a Dah")
			} else  { // key == 3 up start over?
				fmt.Println(" start over")
				entered  = false
				entered1 = false
				pin1.PullUp()
				pin2.PullUp()
				return // wait for next down
			}
			// start next mark timing
			fmt.Println(marks)
			// fixme end loop on pin up with time
			fmt.Print("should not get here")
			entered  = false
			entered1 = false
			pin1.PullUp()
			pin2.PullUp()
			return
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

	// setup tcp
	//openConnection("localhost:6666") 
	

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
	fmt.Print("Watching Pins 13 and 27...wpm = ")
	fmt.Println(WPM)
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

func createMessageForWord() string {
	var s = "[["
	var sz = len(s)
	for _, mark := range marks {
		if mark == DIT {
			s += "\"DIT\","
		} else if mark == DAH {
			s += "\"DAH\","
		} else if mark == LETTER {
			sz = len(s)
			if sz > 0 && s[sz-1] == ',' {
			    s = s[:sz-1]
			}
			s += "]["
		} else if mark == SPACE {
			sz = len(s)
			if sz > 0 && s[sz-1] == ',' {
			    s = s[:sz-1]
			}
			s += "]]"
			return s
		}
	} // end loop
	return s;
}
		
