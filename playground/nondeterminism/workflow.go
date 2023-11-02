package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"go.uber.org/cadence/workflow"
)

var (
	callCount = int32(1)
)

func NonDeterminismSimulatorWorkflow(ctx workflow.Context) error {
	workflowInfo := workflow.GetInfo(ctx)

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: 1 * time.Minute,
		StartToCloseTimeout:    1 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx).Sugar()
	logger.Info("### NonDeterminismSimulatorWorkflow started")
	logger.Info("### workflowInfo: %+v", workflowInfo)

	activityName := "ActivityA"
	if atomic.LoadInt32(&callCount) == 0 {
		fmt.Printf("### Activity %s is going to be called\n", activityName)
		logger.Infof("### Activity %s is going to be called", activityName)
		err := workflow.ExecuteActivity(ctx, activityName).Get(ctx, nil)
		if err != nil {
			logger.Errorf("%s failed", activityName, err)
			return err
		}
	}

	selector := workflow.NewSelector(ctx)
	timer := workflow.NewTimer(ctx, 15*time.Second)
	selector.AddFuture(timer, func(workflow.Future) {
		logger.Infof("Timer future is called")
	})

	logger.Info("Workflow will wait on timer")
	selector.Select(ctx)

	logger.Info("Timer returned. calling another activity")

	err := workflow.ExecuteActivity(ctx, "ActivityB").Get(ctx, nil)
	if err != nil {
		logger.Errorf("%s failed", activityName, err)
		return err
	}

	logger.Info("Workflow finished")
	return nil
}
