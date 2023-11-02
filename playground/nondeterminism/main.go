package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/uber-go/tally"
	apiv1 "github.com/uber/cadence-idl/go/proto/api/v1"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/compatibility"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/grpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	HostPort       = "127.0.0.1:7833"
	Domain         = "samples-domain"
	TaskListName   = "taylan-test-tl"
	ClientName     = "taylan-nondeterminism-worker"
	CadenceService = "cadence-frontend"
)

func main() {
	logger, cadenceClient := BuildLogger(), BuildCadenceClient()
	workerOptions := worker.Options{
		Logger:                         logger,
		MetricsScope:                   tally.NewTestScope(TaskListName, nil),
		NonDeterministicWorkflowPolicy: worker.NonDeterministicWorkflowPolicyFailWorkflow,
	}

	w := worker.New(
		cadenceClient,
		Domain,
		TaskListName,
		workerOptions)

	w.RegisterWorkflowWithOptions(NonDeterminismSimulatorWorkflow, workflow.RegisterOptions{Name: "NonDeterminismSimulatorWorkflow"})
	w.RegisterActivityWithOptions(ActivityA, activity.RegisterOptions{Name: "ActivityA"})
	w.RegisterActivityWithOptions(ActivityB, activity.RegisterOptions{Name: "ActivityB"})

	err := w.Start()
	if err != nil {
		panic("Failed to start worker: " + err.Error())
	}
	logger.Info("Started Worker.", zap.String("worker", TaskListName))

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT)
	fmt.Println("Cadence worker started, press ctrl+c to terminate...")
	<-done
}

func BuildCadenceClient() workflowserviceclient.Interface {
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			CadenceService: {Unary: grpc.NewTransport().NewSingleOutbound(HostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		panic("Failed to start dispatcher: " + err.Error())
	}

	clientConfig := dispatcher.ClientConfig(CadenceService)

	return compatibility.NewThrift2ProtoAdapter(
		apiv1.NewDomainAPIYARPCClient(clientConfig),
		apiv1.NewWorkflowAPIYARPCClient(clientConfig),
		apiv1.NewWorkerAPIYARPCClient(clientConfig),
		apiv1.NewVisibilityAPIYARPCClient(clientConfig),
	)
}

func BuildLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	var err error
	logger, err := config.Build()
	if err != nil {
		panic("Failed to setup logger: " + err.Error())
	}

	return logger
}
