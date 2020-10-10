#######################################################
#################### Generate Binary ##################

MOCK_DIR=mocks
BIN_FILE=./bin/std

bin:
	cd ./cmd && go build -o ../bin/std
clean:
	@if [ -d "$(MOCK_DIR)" ]; then \
		rm -r -f $(MOCK_DIR); \
	fi
	@if [ -d "$(BIN_FILE)" ]; then \
		rm -r -f $(BIN_FILE); \
	fi
	@echo "STD Cleaned"

mocks:
	@echo making mocks to run tests
	mockgen -destination=./mocks/connector.go -package=mocks -source=./repo/connector.go ./repo/connector.go Connector
	mockgen -destination=./mocks/interactor.go -package=mocks -source=./utils/interactor.go ./utils/interactor.go Interactor
	mockgen -destination=./mocks/repo.go -package=mocks -source=./repo/repo.go ./repo/repo.go Repo

test: 
	go test ./...

help:
	@echo 'build - force create a binary and drop in ./bin'
	@echo 'test  - run tests'
	@echo 'mocks - generate interface makes meant for testing'
	@echo 'clean - clean generated code and start from scratch'

vault:
	@echo 'installing hashicorp vault'
	brew tap hashicorp/tap
	brew install hashicorp/tap/vault
