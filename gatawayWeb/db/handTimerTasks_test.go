package db

import (
	"testing"
)

//HandleDayTasks()
func TestHandleDayTasks(t *testing.T) {
	Newdb()
	HandleDayTasks()
	HandleHourTasks()

}
