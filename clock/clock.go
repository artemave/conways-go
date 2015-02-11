package clock

import "time"

type Clock struct {
	delay       time.Duration
	toggleClock chan bool
	nextTick    chan Tick
}

type Tick struct{}

func NewClock(delay time.Duration, sleep func(time.Duration)) *Clock {
	clock := &Clock{
		delay:       delay,
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
					clock.nonBlockingTick()
					sleep(clock.delay * time.Millisecond)
				} else {
					// throttle `default:` to run only 5x per second
					sleep(200 * time.Millisecond)
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

func (c *Clock) nonBlockingTick() {
	select {
	case c.nextTick <- Tick{}:
	default:
	}
}
