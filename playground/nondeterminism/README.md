



1. Run the cadence server via docker-compose up

2. Create samples-domain

```
cadence --env development --domain samples-domain domain register
```

3. Run the worker in this folder

```
go run playground/nondeterminism/*.go
```

4. Trigger a workflow

cadence \
	--env development \
	--domain samples-domain \
	workflow start \
		--tasklist taylan-test-tl \
		--workflow_type NonDeterminismSimulatorWorkflow \
		--execution_timeout 60 \
		--decision_timeout 60

5. While workflow sleeps, shut down and change code