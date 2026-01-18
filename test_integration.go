package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/azahir21/go-backend-boilerplate/pkg/logger"
	proto "github.com/azahir21/go-backend-boilerplate/internal/user/delivery/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log := logger.NewLogger()

	// Build the main application
	buildCmd := exec.Command("go", "build", "-o", "main", "./cmd")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	log.Info("Building main application...")
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("failed to build main application: %v", err)
	}
	log.Info("Main application built.")

	// Start the main application as a separate process
	cmd := exec.Command("./main")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to start main application: %v", err)
	}
	log.Infof("Main application started with PID: %d", cmd.Process.Pid)

	// Ensure the application process is killed when the test finishes
	defer func() {
		if err := cmd.Process.Signal(syscall.SIGINT); err != nil {
			log.Errorf("failed to send SIGINT to main application: %v", err)
		}
		// Give some time for graceful shutdown
		time.Sleep(2 * time.Second)
		if err := cmd.Process.Kill(); err != nil {
			log.Errorf("failed to kill main application process: %v", err)
		}
		log.Info("Main application process terminated.")
		os.Remove("./main") // Clean up the executable
	}()

	// Wait for servers to start
	log.Info("Waiting for servers to become ready...")
	if err := waitForServerReady("http://localhost:8080/api/v1/ping", 30*time.Second); err != nil {
		log.Fatalf("servers not ready: %v", err)
	}
	log.Info("Servers are ready.")

	// --- Test REST API ---
	log.Info("Testing REST API...")
	restClient := &http.Client{Timeout: 5 * time.Second}

	// Test /ping endpoint
	resp, err := restClient.Get("http://localhost:8080/api/v1/ping")
	if err != nil {
		log.Errorf("REST /ping failed: %v", err)
	} else {
		body, _ := io.ReadAll(resp.Body)
		log.Infof("REST /ping response: %s", string(body))
		resp.Body.Close()
	}

	// --- Test gRPC API ---
	log.Info("Testing gRPC API...")
	grpcConn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("gRPC dial failed: %v", err)
	} else {
		defer grpcConn.Close()
		grpcClient := proto.NewUserServiceClient(grpcConn)

		// Test Register gRPC
		grpcRegisterResp, err := grpcClient.Register(context.Background(), &proto.RegisterRequest{
			Username: "grpcuser",
			Email:    "grpc@example.com",
			Password: "password",
		})
		if err != nil {
			log.Errorf("gRPC Register failed: %v", err)
		} else {
			log.Infof("gRPC Register response: %v", grpcRegisterResp)
		}

		// Test Login gRPC
		grpcLoginResp, err := grpcClient.Login(context.Background(), &proto.LoginRequest{
			Username: "grpcuser",
			Password: "password",
		})
		if err != nil {
			log.Errorf("gRPC Login failed: %v", err)
		} else {
			log.Infof("gRPC Login response: %v", grpcLoginResp)
		}
	}

	// --- Test GraphQL API ---
	log.Info("Testing GraphQL API...")
	graphqlClient := &http.Client{Timeout: 5 * time.Second}

	// Test GraphQL query
	graphqlQuery := `{"query":"query { user(id: 1) { id username email } }"}`
	graphqlResp, err := graphqlClient.Post("http://localhost:8081/graphql", "application/json", bytes.NewBufferString(graphqlQuery))
	if err != nil {
		log.Errorf("GraphQL query failed: %v", err)
	} else {
		body, _ := io.ReadAll(graphqlResp.Body)
		log.Infof("GraphQL query response: %s", string(body))
		graphqlResp.Body.Close()
	}

	// Test GraphQL mutation (Register)
	graphqlMutation := `{"query":"mutation { register(username: \"graphqluser\", email: \"graphql@example.com\", password: \"password\") { id username email } }"}`
	graphqlResp, err = graphqlClient.Post("http://localhost:8081/graphql", "application/json", bytes.NewBufferString(graphqlMutation))
	if err != nil {
		log.Errorf("GraphQL register mutation failed: %v", err)
	} else {
		body, _ := io.ReadAll(graphqlResp.Body)
		log.Infof("GraphQL register mutation response: %s", string(body))
		graphqlResp.Body.Close()
	}

	log.Info("Integration tests completed.")
}

// waitForServerReady polls the given URL until the server responds with a 200 OK.
func waitForServerReady(url string, timeout time.Duration) error {
	client := &http.Client{Timeout: 1 * time.Second}
	start := time.Now()
	for time.Since(start) < timeout {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("server at %s did not become ready within %s", url, timeout)
}
