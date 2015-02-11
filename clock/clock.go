package clock

import "time"

type Clock struct {
	toggleClock chan bool
	nextTick    chan Tick
}

type Tick struct{}

func NewClock(delay time.Duration, tickerFactory func(time.Duration) *time.Ticker) *Clock {
	clock := &Clock{
		nextTick:    make(chan Tick),
		toggleClock: make(chan bool),
	}

	go func() {
		var ticker *time.Ticker
		var stopTicker chan struct{}

		for {
			startClock, ok := <-clock.toggleClock

			clockIsOn := stopTicker != nil

			if !ok {
				clock.toggleClock = nil
				go func() {
					defer func() { recover() }()
					close(stopTicker)
				}()
				return
			}

			if startClock && !clockIsOn {
				ticker = tickerFactory(delay * time.Millisecond)
				stopTicker = make(chan struct{})

				go func() {
					for {
						select {
						case <-ticker.C:
							go func() {
								defer func() { recover() }()
								clock.nextTick <- Tick{}
							}()
						case <-stopTicker:
							stopTicker = nil
							return
						}
					}
				}()
			}

			if !startClock && clockIsOn {
				ticker.Stop()
				close(stopTicker)
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

func (c *Clock) NextTick() chan Tick {
	return c.nextTick
}

func (c *Clock) Cleanup() {
	c.StopClock()
	close(c.nextTick)
	close(c.toggleClock)
}
