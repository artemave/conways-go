package clock_test

import (
	"time"

	. "github.com/artemave/conways-go/clock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var clock *Clock
var sleepArgument chan (time.Duration)

func sleep(delay time.Duration) {
	go func() {
		defer func() { recover() }()
		sleepArgument <- delay
	}()
}

var _ = Describe("Clock", func() {
	Context("When clock is on", func() {
		BeforeEach(func() {
			sleepArgument = make(chan time.Duration, 1)
			clock = NewClock(time.Duration(123), sleep)
			clock.StartClock()
		})
		AfterEach(func() {
			close(sleepArgument)
		})

		It("Emits tick every _delay_", func() {
			Eventually(func() bool {
				<-clock.NextTick()
				return true
			}).Should(BeTrue())

			Eventually(func() time.Duration {
				return <-sleepArgument
			}).Should(Equal(time.Duration(123) * time.Millisecond))
		})
	})
})
