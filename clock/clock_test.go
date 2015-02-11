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

func tickerFactory(delay time.Duration) *time.Ticker {
	go func() {
		defer func() { recover() }()
		delayArg <- delay
	}()
	return &time.Ticker{C: tickerC}
}

var _ = Describe("Clock", func() {
	Context("When clock is on", func() {
		BeforeEach(func() {
			tickerC = make(chan time.Time)
			delayArg = make(chan time.Duration, 1)
			clock = NewClock(time.Duration(12), tickerFactory)
			clock.StartClock()
		})
		AfterEach(func() {
			clock.Cleanup()
			close(delayArg)
			close(tickerC)
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
	})
})
