package adminapi

import (
	"fmt"
	"net/http"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/router"
	"github.com/ServiceComb/go-chassis/metrics"
	// "github.com/ServiceComb/go-chassis/server/restful"
	"github.com/emicklei/go-restful"
	"github.com/go-chassis/mesher/adminapi/health"
	"github.com/go-chassis/mesher/adminapi/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//GetWebService creates route and returns all admin api's
func GetWebService() restful.WebService {
	restfulWebService := new(restful.WebService)
	restfulWebService.Route(restfulWebService.GET("/v1/mesher/version").To(GetVersion))
	restfulWebService.Route(restfulWebService.GET("/v1/mesher/metrics").To(GetMetrics))
	restfulWebService.Route(restfulWebService.GET("/v1/mesher/routeRule").To(RouteRule))
	restfulWebService.Route(restfulWebService.GET("/v1/mesher/routeRule/:serviceName").To(RouteRuleByService))
	restfulWebService.Route(restfulWebService.GET("/v1/mesher/health").To(MesherHealth))
	return *restfulWebService
}

//GetVersion writes version in response header
func GetVersion(req *restful.Request, resp *restful.Response) {
	version := version.Ver()
	resp.WriteHeaderAndJson(http.StatusOK, version, common.JSON)
}

//GetMetrics returns metrics data
func GetMetrics(req *restful.Request, resp *restful.Response) {
	promhttp.HandlerFor(metrics.GetSystemPrometheusRegistry(), promhttp.HandlerOpts{}).ServeHTTP(resp.ResponseWriter, req.Request)
}

//RouteRule returns all router configs
func RouteRule(req *restful.Request, resp *restful.Response) {
	routerConfig := &model.RouterConfig{
		Destinations: router.DefaultRouter.FetchRouteRule(),
	}
	resp.WriteHeaderAndJson(http.StatusOK, routerConfig, "text/vnd.yaml")
}

//RouteRuleByService returns route config for particular service
func RouteRuleByService(req *restful.Request, resp *restful.Response) {
	serviceName := req.Request.URL.Query().Get(":serviceName")
	routeRule := router.DefaultRouter.FetchRouteRuleByServiceName(serviceName)
	if routeRule == nil {
		resp.WriteHeaderAndJson(http.StatusNotFound, fmt.Sprintf("%s routeRule not found", serviceName), common.JSON)
		return
	}
	resp.WriteHeaderAndJson(http.StatusOK, routeRule, "text/vnd.yaml")
}

//MesherHealth returns mesher health
func MesherHealth(req *restful.Request, resp *restful.Response) {
	healthResp := health.GetMesherHealth()
	if healthResp.Status == health.Red {
		resp.WriteHeaderAndJson(http.StatusInternalServerError, healthResp, common.JSON)
		return
	}
	resp.WriteHeaderAndJson(http.StatusOK, healthResp, common.JSON)
}
