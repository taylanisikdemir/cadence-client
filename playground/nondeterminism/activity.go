package main

import (
	"context"

	"go.uber.org/cadence/activity"
)

func ActivityA(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("ActivityA is called")
	return nil
}

func ActivityB(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	logger.Info("ActivityB is called")
	return nil
}
