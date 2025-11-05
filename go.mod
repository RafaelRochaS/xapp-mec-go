module github.com/RafaelRochaS/xapp-mec-go

go 1.15

require (
	gerrit.o-ran-sc.org/r/ric-plt/alarm-go.git/alarm v0.5.0
	gerrit.o-ran-sc.org/r/ric-plt/xapp-frame v0.9.3
	github.com/google/uuid v1.6.0 // indirect
)

replace gerrit.o-ran-sc.org/r/ric-plt/xapp-frame => gerrit.o-ran-sc.org/r/ric-plt/xapp-frame.git v0.9.3

replace gerrit.o-ran-sc.org/r/ric-plt/sdlgo => gerrit.o-ran-sc.org/r/ric-plt/sdlgo.git v0.8.0

replace gerrit.o-ran-sc.org/r/com/golog => gerrit.o-ran-sc.org/r/com/golog.git v0.0.2
