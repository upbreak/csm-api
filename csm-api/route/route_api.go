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
	handlerWeatherSrt := &handler.HandlerWeatherSrtNcst{
		Service: &service.ServiceWeather{
			ApiKey:       apiConfig,
			SafeDB:       safeDB,
			SafeTDB:      safeDB,
			Store:        r,
			SitePosStore: r,
		},
		SitePosService: &service.ServiceSitePos{
			DB:    safeDB,
			Store: r,
		},
	}

	// 기상청 기상특보통보문 조회
	handlerWeatherWrn := &handler.HandlerWeatherWrnMsg{
		Service: &service.ServiceWeather{
			ApiKey:       apiConfig,
			SafeDB:       safeDB,
			SafeTDB:      safeDB,
			Store:        r,
			SitePosStore: r,
		},
	}

	// 공휴일 조회
	restDatehandler := &handler.HandlerRestDate{
		Service: &service.ServiceRestDate{
			ApiKey: apiConfig,
		},
	}

	router.Get("/map/point", roadAddressHandler.AddressPoint) // 지도 좌표
	router.Get("/weather/srt", handlerWeatherSrt.ServeHTTP)   // 기상청 초단기 실황
	router.Get("/weather/wrn", handlerWeatherWrn.ServeHTTP)   // 기상청 기상특보통보문 조회
	router.Get("/rest-date", restDatehandler.ServeHTTP)       // 공휴일 조회
	//router.Get("/weather/") // 날씨

	return router
}
