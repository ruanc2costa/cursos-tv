//package controller
//
//import (
//	"log"
//	"net/http"
//	"tvtec/service"
//
//	"github.com/gin-gonic/gin"
//)
//
//type PowerBIController struct {
//	powerBIService service.PowerBIService
//}
//
//func NewPowerBIController(powerBIService service.PowerBIService) *PowerBIController {
//	return &PowerBIController{
//		powerBIService: powerBIService,
//	}
//}
//
//// GetEmbedToken gera um token de incorporação para um dashboard do Power BI
//func (c *PowerBIController) GetEmbedToken(ctx *gin.Context) {
//	dashboardID := ctx.Query("dashboardId")
//	groupID := ctx.Query("groupId")
//	datasetID := ctx.Query("datasetId")
//
//	if dashboardID == "" || groupID == "" || datasetID == "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"error": "dashboardId, groupId e datasetId são obrigatórios",
//		})
//		return
//	}
//
//	log.Printf("Gerando token para dashboard: %s, grupo: %s, dataset: %s", dashboardID, groupID, datasetID)
//
//	embedConfig, err := c.powerBIService.GetDashboardEmbedToken(dashboardID, groupID, datasetID)
//	if err != nil {
//		log.Printf("Erro ao gerar token de incorporação: %v", err)
//		ctx.JSON(http.StatusInternalServerError, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	log.Printf("Token de incorporação gerado com sucesso: %s", embedConfig.EmbedToken.TokenID)
//	ctx.JSON(http.StatusOK, embedConfig)
//}
