package clock

import "time"

type Clock struct {
	Delay       time.Duration
	toggleClock chan bool
	nextTick    chan Tick
}

type Tick struct{}

func NewClock(delay time.Duration) *Clock {
	clock := &Clock{
		Delay:       delay,
		nextTick:    make(chan Tick),
		toggleClock: make(chan bool),
	}

	go func() {
		clockIsOn := false
		for {
			select {
			case clockIsOn = <-clock.toggleClock:
			default:
				if clockIsOn {
					select {
					case clock.nextTick <- Tick{}:
					default:
						// client didn't read previous tick yet?
						// do nothing
					}
					time.Sleep(clock.Delay * time.Millisecond)
				} else {
					// throttle `default:` to run only 5x per second
					time.Sleep(200 * time.Millisecond)
				}
			}
		}
	}()

	return clock
}

func (c *Clock) StartClock() {
	go func() {
		c.toggleClock <- true
	}()
}

func (c *Clock) StopClock() {
	go func() {
		c.toggleClock <- false
	}()
}

func (c *Clock) NextTick() chan Tick {
	return c.nextTick
}
