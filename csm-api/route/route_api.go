package route

import (
	"csm-api/config"
	"csm-api/handler"
	"csm-api/service"
	"csm-api/store"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func ApiRoute(apiConfig *config.ApiConfig, safeDB *sqlx.DB, r *store.Repository) chi.Router {
	router := chi.NewRouter()

	// 지도 좌표
	roadAddressHandler := &handler.HandlerRoadAddress{
		Service: &service.ServiceAddressSearch{
			ApiKey: apiConfig,
		},
	}

	// 기상청 초단기 실황
	handlerWhetherSrt := &handler.HandlerWhetherSrtNcst{
		Service: &service.ServiceWhether{
			ApiKey: apiConfig,
		},
		SitePosService: &service.ServiceSitePos{
			DB:    safeDB,
			Store: r,
		},
	}

	// 기상청 기상특보통보문 조회
	handlerWhetherWrn := &handler.HandlerWhetherWrnMsg{
		Service: &service.ServiceWhether{
			ApiKey: apiConfig,
		},
	}

	// 공휴일 조회
	restDatehandler := &handler.HandlerRestDate{
		Service: &service.ServiceRestDate{
			ApiKey: apiConfig,
		},
	}

	router.Get("/map/point", roadAddressHandler.AddressPoint) // 지도 좌표
	router.Get("/whether/srt", handlerWhetherSrt.ServeHTTP)   // 기상청 초단기 실황
	router.Get("/whether/wrn", handlerWhetherWrn.ServeHTTP)   // 기상청 기상특보통보문 조회
	router.Get("/rest-date", restDatehandler.ServeHTTP)       // 공휴일 조회

	return router
}
