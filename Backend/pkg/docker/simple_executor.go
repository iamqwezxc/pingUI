// pkg/docker/simple_executor.go
package docker

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type SimpleExecutor struct{}

func NewDockerExecutor() (*SimpleExecutor, error) {
	// Проверяем что Docker доступен
	if err := exec.Command("docker", "version").Run(); err != nil {
		return nil, fmt.Errorf("Docker is not available: %v", err)
	}
	return &SimpleExecutor{}, nil
}

func (d *SimpleExecutor) BuildImage() error {
	// Ничего не делаем, используем готовые образы
	return nil
}

func (d *SimpleExecutor) ExecuteCommand(command string, timeout int) (string, error) {
	// Валидация команды
	if strings.TrimSpace(command) == "" {
		return "", fmt.Errorf("command cannot be empty")
	}

	// Запрещаем опасные команды
	forbiddenPatterns := []string{
		"rm -rf", "mkfs", "dd", "shutdown", "reboot", "halt",
		"poweroff", "> /dev/", "| sudo", "chmod 777", "passwd",
	}

	for _, pattern := range forbiddenPatterns {
		if strings.Contains(command, pattern) {
			return "", fmt.Errorf("dangerous command not allowed")
		}
	}

	// Используем готовый alpine образ БЕЗ run_command.sh
	dockerCmd := fmt.Sprintf(
		"docker run --rm alpine sh -c \"%s\"",
		strings.ReplaceAll(command, "\"", "\\\""),
	)

	cmd := exec.Command("sh", "-c", dockerCmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Устанавливаем таймаут
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		cmd.Process.Kill()
		return "", fmt.Errorf("command timed out after %d seconds", timeout)
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("execution failed: %v, stderr: %s", err, stderr.String())
		}
	}

	return stdout.String(), nil
}

func (d *SimpleExecutor) Cleanup() error {
	// Очищаем остановленные контейнеры
	cmd := exec.Command("sh", "-c", "docker ps -a -f status=exited -q | xargs -r docker rm 2>/dev/null || true")
	return cmd.Run()
}
