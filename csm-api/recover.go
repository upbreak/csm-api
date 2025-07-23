package main

import (
	"context"
	"csm-api/entity"
	"csm-api/utils"
	"fmt"
	"log"
)

func Recover(message string) {
	if r := recover(); r != nil {
		log.Printf("[panic][Scheduler-ModifyWorkerDeadlineSchedule]: %v", r)
		_ = entity.WriteErrorLog(context.Background(), utils.CustomMessageErrorf(fmt.Sprintf("panic %s", message), fmt.Errorf("%v", r)))
	}
}
