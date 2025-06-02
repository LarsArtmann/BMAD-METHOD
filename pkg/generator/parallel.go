package generator

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ParallelGenerator handles parallel file generation
type ParallelGenerator struct {
	workerCount   int
	maxRetries    int
	retryDelay    time.Duration
	resultChannel chan GenerationResult
}

// GenerationResult represents the result of generating a single file
type GenerationResult struct {
	Filename     string
	TemplateName string
	Success      bool
	Error        error
	Duration     time.Duration
	Size         int64
}

// GenerationTask represents a file generation task
type GenerationTask struct {
	Filename     string
	TemplateName string
	Context      *GenerationContext
	Generator    *Generator
}

// GenerationSummary provides statistics about the generation process
type GenerationSummary struct {
	TotalFiles    int
	SuccessCount  int
	FailureCount  int
	TotalDuration time.Duration
	TotalSize     int64
	Results       []GenerationResult
}

// NewParallelGenerator creates a new parallel generator
func NewParallelGenerator(workerCount int) *ParallelGenerator {
	if workerCount <= 0 {
		workerCount = 4 // Default worker count
	}

	return &ParallelGenerator{
		workerCount:   workerCount,
		maxRetries:    3,
		retryDelay:    100 * time.Millisecond,
		resultChannel: make(chan GenerationResult, workerCount*2),
	}
}

// GenerateFiles generates multiple files in parallel
func (pg *ParallelGenerator) GenerateFiles(ctx context.Context, tasks []GenerationTask) (*GenerationSummary, error) {
	if len(tasks) == 0 {
		return &GenerationSummary{}, nil
	}

	// Create task channel
	taskChan := make(chan GenerationTask, len(tasks))
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < pg.workerCount; i++ {
		wg.Add(1)
		go pg.worker(ctx, &wg, taskChan)
	}

	// Send tasks to workers
	go func() {
		defer close(taskChan)
		for _, task := range tasks {
			select {
			case taskChan <- task:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Collect results
	results := make([]GenerationResult, 0, len(tasks))
	summary := &GenerationSummary{
		TotalFiles: len(tasks),
	}

	// Start result collector
	resultDone := make(chan struct{})
	go func() {
		defer close(resultDone)
		for i := 0; i < len(tasks); i++ {
			select {
			case result := <-pg.resultChannel:
				results = append(results, result)
				summary.TotalDuration += result.Duration
				summary.TotalSize += result.Size
				
				if result.Success {
					summary.SuccessCount++
				} else {
					summary.FailureCount++
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for workers to complete
	wg.Wait()

	// Wait for result collection to complete
	select {
	case <-resultDone:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	summary.Results = results
	return summary, nil
}

// worker processes generation tasks
func (pg *ParallelGenerator) worker(ctx context.Context, wg *sync.WaitGroup, taskChan <-chan GenerationTask) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			
			result := pg.processTask(ctx, task)
			
			select {
			case pg.resultChannel <- result:
			case <-ctx.Done():
				return
			}
			
		case <-ctx.Done():
			return
		}
	}
}

// processTask processes a single generation task with retries
func (pg *ParallelGenerator) processTask(ctx context.Context, task GenerationTask) GenerationResult {
	start := time.Now()
	
	var lastErr error
	for attempt := 0; attempt < pg.maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return GenerationResult{
				Filename:     task.Filename,
				TemplateName: task.TemplateName,
				Success:      false,
				Error:        ctx.Err(),
				Duration:     time.Since(start),
			}
		default:
		}

		err := task.Generator.generateFile(task.Filename, task.TemplateName, task.Context)
		if err == nil {
			// Success - get file size
			size := pg.getFileSize(task.Generator.config.OutputDir, task.Filename)
			
			return GenerationResult{
				Filename:     task.Filename,
				TemplateName: task.TemplateName,
				Success:      true,
				Duration:     time.Since(start),
				Size:         size,
			}
		}

		lastErr = err
		
		// Wait before retry (except on last attempt)
		if attempt < pg.maxRetries-1 {
			select {
			case <-time.After(pg.retryDelay * time.Duration(attempt+1)):
			case <-ctx.Done():
				return GenerationResult{
					Filename:     task.Filename,
					TemplateName: task.TemplateName,
					Success:      false,
					Error:        ctx.Err(),
					Duration:     time.Since(start),
				}
			}
		}
	}

	return GenerationResult{
		Filename:     task.Filename,
		TemplateName: task.TemplateName,
		Success:      false,
		Error:        fmt.Errorf("failed after %d attempts: %w", pg.maxRetries, lastErr),
		Duration:     time.Since(start),
	}
}

// getFileSize gets the size of a generated file
func (pg *ParallelGenerator) getFileSize(outputDir, filename string) int64 {
	// This is a simplified implementation
	// In practice, you'd use filepath.Join and os.Stat
	return 0
}

// WorkerPool represents a pool of workers for file generation
type WorkerPool struct {
	workerCount int
	taskQueue   chan func()
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	pool := &WorkerPool{
		workerCount: workerCount,
		taskQueue:   make(chan func(), workerCount*2),
		ctx:         ctx,
		cancel:      cancel,
	}

	// Start workers
	for i := 0; i < workerCount; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}

	return pool
}

// Submit submits a task to the worker pool
func (wp *WorkerPool) Submit(task func()) {
	select {
	case wp.taskQueue <- task:
	case <-wp.ctx.Done():
		// Pool is shutting down
	}
}

// worker processes tasks from the queue
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case task := <-wp.taskQueue:
			if task != nil {
				task()
			}
		case <-wp.ctx.Done():
			return
		}
	}
}

// Shutdown gracefully shuts down the worker pool
func (wp *WorkerPool) Shutdown() {
	close(wp.taskQueue)
	wp.cancel()
	wp.wg.Wait()
}

// GetOptimalWorkerCount returns the optimal number of workers based on system resources
func GetOptimalWorkerCount() int {
	// Simple heuristic: number of CPU cores
	// In practice, you might want to consider I/O vs CPU bound operations
	return 4 // Default to 4 workers for now
}