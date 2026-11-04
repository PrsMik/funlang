ifeq ($(OS),Windows_NT)
    BINARY_NAME := fun.exe
    KILL_CMD := powershell -Command "Get-Process -Name 'fun' -ErrorAction SilentlyContinue | Stop-Process -Force; Start-Sleep -Milliseconds 1000"
    RM_CMD := rm -f $(BINARY_NAME) || del /Q /F $(BINARY_NAME) 2>nul || exit 0
else
    BINARY_NAME := fun
    KILL_CMD := pkill -x $(BINARY_NAME) 2>/dev/null || true
    RM_CMD := rm -f $(BINARY_NAME)
endif

.PHONY: all build kill clean

all: build

build: clean
	go build -o $(BINARY_NAME) .

kill:
	-@$(KILL_CMD)

clean: kill
	-@$(RM_CMD)
