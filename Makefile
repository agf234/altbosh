all: compile run

compile:
	go build -ldflags "-w" altbosh.go 

run:
	./altbosh

debug:
	go build -gcflags "-N -l" -o gdb_sandbox altbosh.go;
	gdb gdb_sandbox list;
