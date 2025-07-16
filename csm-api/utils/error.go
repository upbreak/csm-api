package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func CustomErrorf(err error) error {
	if err == nil {
		return nil
	}

	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("unknown error: %w", err)
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")
	shortFunc := parts[len(parts)-1]
	shortFile := filepath.Base(file)

	return fmt.Errorf("%s/%s err: %w", shortFile, shortFunc, err)
}

func CustomMessageErrorf(message string, err error) error {
	if err == nil {
		return nil
	}

	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("unknown error: %w", err)
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")
	shortFunc := parts[len(parts)-1]
	shortFile := filepath.Base(file)

	return fmt.Errorf("%s/%s %s err: %w", shortFile, shortFunc, message, err)
}

func CustomErrorfDepth(depth int, err error) error {
	if err == nil {
		return nil
	}

	pc, file, _, ok := runtime.Caller(depth)
	if !ok {
		return fmt.Errorf("unknown error: %w", err)
	}

	funcName := runtime.FuncForPC(pc).Name()
	funcParts := strings.Split(funcName, ".")
	shortFunc := funcParts[len(funcParts)-1]
	baseFile := filepath.Base(file)

	return fmt.Errorf("%s/%s err: %w", baseFile, shortFunc, err)
}

func CustomMessageErrorfDepth(depth int, message string, err error) error {
	if err == nil {
		return nil
	}

	pc, file, _, ok := runtime.Caller(depth)
	if !ok {
		return fmt.Errorf("unknown error: %w", err)
	}

	funcName := runtime.FuncForPC(pc).Name()
	funcParts := strings.Split(funcName, ".")
	shortFunc := funcParts[len(funcParts)-1]
	baseFile := filepath.Base(file)

	return fmt.Errorf("%s/%s %s err: %w", baseFile, shortFunc, message, err)
}
