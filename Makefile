ifeq ($(OS),Windows_NT)
    BINARY_NAME := fun.exe
    KILL_CMD    := powershell -Command "Get-Process -Name 'fun' -ErrorAction SilentlyContinue | Stop-Process -Force; exit 0"
    RM_CMD      := powershell -Command "Remove-Item -ErrorAction SilentlyContinue $(BINARY_NAME)"
else
    BINARY_NAME := fun
    KILL_CMD    := pkill -x $(BINARY_NAME) 2>/dev/null || true
    RM_CMD      := rm -f $(BINARY_NAME)
endif

.PHONY: all build kill clean

all: build

kill:
	-@$(KILL_CMD)

build: kill
	go build -o $(BINARY_NAME) .

clean:
	-@$(RM_CMD)