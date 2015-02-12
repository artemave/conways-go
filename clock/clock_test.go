package clock_test

import (
	"time"

	. "github.com/artemave/conways-go/clock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var clock *Clock
var delayArg chan time.Duration
var tickerC chan time.Time
var stopCallCount chan struct{}

type TestTicker struct {
	Ch chan time.Time
}

func (t TestTicker) C() <-chan time.Time {
	return t.Ch
}

func (t TestTicker) Stop() {
	stopCallCount <- struct{}{}
}

func tickerFactory(delay time.Duration) Ticker {
	stopCallCount = make(chan struct{})
	tickerC = make(chan time.Time)
	delayArg = make(chan time.Duration, 1)

	go func() {
		defer func() { recover() }()
		delayArg <- delay
	}()
	return TestTicker{Ch: tickerC}
}

var _ = Describe("Clock", func() {
	Context("When clock is on", func() {
		BeforeEach(func() {
			clock = NewClock(time.Duration(12), tickerFactory)
			clock.StartClock()
		})
		AfterEach(func() {
			clock.Cleanup()
		})

		It("Emits tick every _delay_", func() {
			tickerC <- time.Time{}
			Eventually(func() bool {
				<-clock.NextTick()
				return true
			}).Should(BeTrue())

			Eventually(func() time.Duration {
				return <-delayArg
			}).Should(Equal(time.Duration(12) * time.Millisecond))
		})

		Describe("StopClock()", func() {
			It("pauses ticker", func() {
				clock.StopClock()
				Eventually(stopCallCount).Should(Receive())
			})
		})
	})
})
