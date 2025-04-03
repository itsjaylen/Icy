// Package admin provides functionality for managing Docker containers inside the stack.
package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/itsjaylen/IcyAPI/internal/utils"
	logger "itsjaylen/IcyLogger"
)

var (
	// ErrNotRunningInContainer is returned when the application is not running inside a container.
	ErrNotRunningInContainer = errors.New("not running inside a container")

	// ErrContainerNotRunning is returned when the container does not start after a restart.
	ErrContainerNotRunning = errors.New("container is not running after restart")
)

// IsRunningInContainer checks if the program is running inside a container.
func IsRunningInContainer() bool {
	_, err1 := os.Stat("/.dockerenv")
	_, err2 := os.Stat("/run/.containerenv")

	return err1 == nil || err2 == nil
}

// GetContainerID returns the container ID from the hostname (works inside a container).
func GetContainerID() (string, error) {
	return os.Hostname()
}

// RestartContainer restarts the current container.
func RestartContainer(ctx context.Context) error {
	if !IsRunningInContainer() {
		return ErrNotRunningInContainer
	}

	cli, cliErr := client.NewClientWithOpts(client.FromEnv)
	if cliErr != nil {
		return fmt.Errorf("failed to create Docker client: %w", cliErr)
	}
	defer func() {
		if closeErr := cli.Close(); closeErr != nil {
			logger.Error.Println("Failed to close Docker client:", closeErr)
		}
	}()

	containerID, idErr := GetContainerID()
	if idErr != nil {
		return fmt.Errorf("failed to get container ID: %w", idErr)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if restartErr := cli.ContainerRestart(ctx, containerID, container.StopOptions{}); restartErr != nil {
		return fmt.Errorf("failed to restart container: %w", restartErr)
	}

	time.Sleep(2 * time.Second)

	containerInfo, inspectErr := cli.ContainerInspect(ctx, containerID)
	if inspectErr != nil {
		return fmt.Errorf("failed to inspect container: %w", inspectErr)
	}

	if !containerInfo.State.Running {
		return ErrContainerNotRunning
	}

	logger.Info.Printf("Container %s restarted successfully.", containerID)

	return nil
}

// ExecuteCommand runs a shell command and returns its output.
func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

// StatusHandler returns the container's status.
func StatusHandler(writer http.ResponseWriter, request *http.Request) {
	if !IsRunningInContainer() {
		utils.WriteJSONResponse(writer, http.StatusInternalServerError,
			map[string]string{"error": "Not running in a container"})

		return
	}

	cli, cliErr := client.NewClientWithOpts(client.FromEnv)
	if cliErr != nil {
		utils.WriteJSONResponse(writer, http.StatusInternalServerError,
			map[string]string{"error": "Failed to create Docker client"})

		return
	}
	defer func() {
		if closeErr := cli.Close(); closeErr != nil {
			logger.Error.Println("Failed to close Docker client:", closeErr)
		}
	}()

	containerID, idErr := GetContainerID()
	if idErr != nil {
		utils.WriteJSONResponse(writer, http.StatusInternalServerError,
			map[string]string{"error": "Failed to get container ID"})

		return
	}

	containerInfo, inspectErr := cli.ContainerInspect(request.Context(), containerID)
	if inspectErr != nil {
		utils.WriteJSONResponse(writer, http.StatusInternalServerError,
			map[string]string{"error": "Error retrieving container status"})

		return
	}

	response := map[string]interface{}{
		"id":      containerID,
		"running": containerInfo.State.Running,
		"status":  containerInfo.State.Status,
	}

	utils.WriteJSONResponse(writer, http.StatusOK, response)
}

// RestartHandler handles restarting the container.
func RestartHandler(writer http.ResponseWriter, request *http.Request) {
	if err := RestartContainer(request.Context()); err != nil {
		utils.WriteJSONResponse(writer, http.StatusInternalServerError, map[string]string{"error": err.Error()})

		return
	}

	utils.WriteJSONResponse(writer, http.StatusOK, map[string]string{"message": "Container restarted successfully"})
}

// ExecHandler runs a command inside the container.
func ExecHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		utils.WriteJSONResponse(writer, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})

		return
	}

	var req struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}

	if decodeErr := json.NewDecoder(request.Body).Decode(&req); decodeErr != nil {
		utils.WriteJSONResponse(writer, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})

		return
	}

	if req.Command == "" {
		utils.WriteJSONResponse(writer, http.StatusBadRequest, map[string]string{"error": "Command cannot be empty"})

		return
	}

	output, execErr := ExecuteCommand(req.Command, req.Args...)
	if execErr != nil {
		utils.WriteJSONResponse(writer, http.StatusInternalServerError, map[string]string{
			"error":  fmt.Sprintf("Error executing command: %v", execErr),
			"output": output,
		})

		return
	}

	utils.WriteJSONResponse(writer, http.StatusOK, map[string]string{"output": output})
}
