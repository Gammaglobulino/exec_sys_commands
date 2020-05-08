package cockroachlab

import (
	"../sending_receiving_command_loop"
	"strconv"
	"strings"
	"time"
)

func CreateNewServerClusterPipe(numberOfNodes int, ipAdsr string, portNumber int) sending_receiving_command_loop.CommandPipe {
	pipe := sending_receiving_command_loop.CommandPipe{Name: "Remote start a local cockroach cluster"}
	tcpPort := 8080
	joinParam := createJoinParam(numberOfNodes, ipAdsr, portNumber)
	for i := 0; i < numberOfNodes; i++ {

		cockroachNode := sending_receiving_command_loop.NewCommandString("start cockroach node"+strconv.Itoa(i), "cockroach start", "windows")
		cockroachNode.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
		cockroachNode.(*sending_receiving_command_loop.RemoteCommand).AddParam("--store=node-" + strconv.Itoa(i))
		cockroachNode.(*sending_receiving_command_loop.RemoteCommand).AddParam("--listen-addr=" + ipAdsr + ":" + strconv.Itoa(portNumber+i))
		cockroachNode.(*sending_receiving_command_loop.RemoteCommand).AddParam("--http-addr=" + ipAdsr + ":" + strconv.Itoa(tcpPort+i))
		cockroachNode.(*sending_receiving_command_loop.RemoteCommand).AddParam(joinParam)
		pipe.AddCommand(cockroachNode)
	}
	initCluster := sending_receiving_command_loop.NewCommandString("init cockroach cluster", "cockroach init", "windows")
	initCluster.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	initCluster.(*sending_receiving_command_loop.RemoteCommand).AddParam("--host=" + ipAdsr + ":" + strconv.Itoa(portNumber))
	pipe.AddCommand(initCluster)
	return pipe
}

func createJoinParam(numberOfNodes int, ipAddr string, port int) string {
	joinParam := "--join="
	for i := 0; i < numberOfNodes; i++ {
		joinParam = joinParam + ipAddr + ":" + strconv.Itoa(port+i) + ","
	}
	joinParam = strings.TrimSuffix(joinParam, ",")

	return joinParam
}
func VerifyNodeStatusOn(host string) (string, error) {
	cockroachNodeStatus := sending_receiving_command_loop.NewCommandString("cockroach node status", "cockroach node status", "windows")
	cockroachNodeStatus.(*sending_receiving_command_loop.RemoteCommand).AddParam("--host=" + host)
	cockroachNodeStatus.(*sending_receiving_command_loop.RemoteCommand).AddParam("--insecure")
	cockroachNodeStatus.(*sending_receiving_command_loop.RemoteCommand).ExecutionTimeout = 10 * time.Second
	msg, err := cockroachNodeStatus.Execute()
	if err != nil {
		return msg, err
	}
	return cockroachNodeStatus.(*sending_receiving_command_loop.RemoteCommand).CommandOutput, nil
}
