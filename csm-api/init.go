package main

import (
	"context"
	"csm-api/config"
	"csm-api/service"
	"csm-api/store"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

type Init struct {
	Service service.RestDateApiService
}

func (i *Init) New() (*Init, error) {
	apiCfg, err := config.GetApiConfig()
	if err != nil {
		return nil, err
	}

	return &Init{
		Service: &service.ServiceRestDate{
			ApiKey: apiCfg,
		},
	}, err
}

func (i *Init) RunInitializations(ctx context.Context, cfg *config.DBConfigs) (err error) {
	eg, ctx := errgroup.WithContext(ctx)

	_, safeCleanup, safeErr := store.New(ctx, cfg.Safe)
	if safeErr != nil {
		safeCleanup()
		return fmt.Errorf("store.New: %w", safeErr)
	}

	defer func() {
		if err != nil {
			safeCleanup()
		}
	}()

	eg.Go(func() error {
		fmt.Println("[init] restDelInfo API call & db save")

		time.Sleep(2 * time.Second) // 2초 대기

		return nil
	})

	if err = eg.Wait(); err != nil {
		return err
	}
	return
}
