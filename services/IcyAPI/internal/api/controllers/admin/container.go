package admin

import (
	"IcyAPI/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// IsRunningInContainer checks if the program is running inside a container
func IsRunningInContainer() bool {
	_, err1 := os.Stat("/.dockerenv")
	_, err2 := os.Stat("/run/.containerenv")
	return err1 == nil || err2 == nil
}

// GetContainerID returns the container ID from the hostname (works inside a container)
func GetContainerID() (string, error) {
	return os.Hostname()
}

// RestartContainer restarts the current container
func RestartContainer() error {
	if !IsRunningInContainer() {
		return fmt.Errorf("not running inside a container")
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	containerID, err := GetContainerID()
	if err != nil {
		return fmt.Errorf("failed to get container ID: %w", err)
	}

	fmt.Printf("Restarting container: %s\n", containerID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := cli.ContainerRestart(ctx, containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	time.Sleep(2 * time.Second)

	containerInfo, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to inspect container: %w", err)
	}

	if !containerInfo.State.Running {
		return fmt.Errorf("container is not running after restart")
	}

	fmt.Println("Container restarted successfully.")
	return nil
}

// ExecuteCommand runs a shell command and returns its output
func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// StatusHandler returns the container's status
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create Docker client"})
		return
	}
	defer cli.Close()

	containerID, err := GetContainerID()
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to get container ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	containerInfo, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error retrieving container status"})
		return
	}

	response := map[string]interface{}{
		"id":      containerID,
		"running": containerInfo.State.Running,
		"status":  containerInfo.State.Status,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// RestartHandler handles restarting the container
func RestartHandler(w http.ResponseWriter, r *http.Request) {
	if err := RestartContainer(); err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Container restarted successfully"})
}

// ExecHandler runs a command inside the container
func ExecHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	var request struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if request.Command == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Command cannot be empty"})
		return
	}

	output, err := ExecuteCommand(request.Command, request.Args...)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error":  fmt.Sprintf("Error executing command: %v", err),
			"output": output,
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"output": output})
}
