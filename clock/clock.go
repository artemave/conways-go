package clock

import "time"

type Clock struct {
	toggleClock chan bool
	nextTick    chan Tick
}

type Tick struct{}

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

func NewClock(delay time.Duration, tickerFactory func(time.Duration) Ticker) *Clock {
	clock := &Clock{
		nextTick:    make(chan Tick),
		toggleClock: make(chan bool),
	}

	go func() {
		var ticker Ticker
		stopTicker := make(chan struct{})
		clockIsOn := false

		for {
			startClock, ok := <-clock.toggleClock

			if !ok {
				clock.toggleClock = nil
				close(stopTicker)
				return
			}

			if startClock {
				if !clockIsOn {
					clockIsOn = true
					ticker = tickerFactory(delay * time.Millisecond)

					go func() {
						for {
							select {
							case <-ticker.C():
								go func() {
									defer func() { recover() }()
									clock.nextTick <- Tick{}
								}()
							case <-stopTicker:
								return
							}
						}
					}()
				}
			} else {
				if clockIsOn {
					clockIsOn = false
					ticker.Stop()
					stopTicker <- struct{}{}
				}
			}
		}
	}()

	return clock
}

func (c *Clock) StartClock() {
	if c.toggleClock != nil {
		c.toggleClock <- true
	}
}

func (c *Clock) StopClock() {
	if c.toggleClock != nil {
		c.toggleClock <- false
	}
}

func (c *Clock) NextTick() <-chan Tick {
	return c.nextTick
}

func (c *Clock) Cleanup() {
	c.StopClock()
	close(c.nextTick)
	close(c.toggleClock)
}
