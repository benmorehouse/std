#######################################################
#################### Generate Binary ##################

bin:
	cd ./cmd && go build -o ../bin/std
mocks:
	@echo making mocks to run tests
	mockgen -destination=./mocks/connector.go -source=./repo/connector.go ./repo/connector.go Connector
	mockgen -destination=./mocks/interactor.go -source=./utils/interactor.go ./utils/interactor.go Interactor
	mockgen -destination=./mocks/repo.go -source=./repo/repo.go ./repo/repo.go Repo
help:
	@echo 'build - force create a binary and drop in ./bin'
