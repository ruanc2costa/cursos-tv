package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tvtec/controller"
	"tvtec/middleware"
	"tvtec/models"
	"tvtec/repository"
	"tvtec/service"

	"github.com/joho/godotenv"
)

// Middleware para evitar redirecionamentos 307
func noRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Garante que a URL não tenha barras duplas
		path := c.Request.URL.Path
		if strings.Contains(path, "//") {
			cleanPath := strings.ReplaceAll(path, "//", "/")
			log.Printf("URL com barras duplas detectada. Original: %s, Limpa: %s", path, cleanPath)
			c.Request.URL.Path = cleanPath
		}

		// Remove a barra final da URL se existir e não for a raiz
		if path != "/" && strings.HasSuffix(path, "/") {
			cleanPath := strings.TrimSuffix(path, "/")
			log.Printf("URL com barra final detectada. Original: %s, Limpa: %s", path, cleanPath)
			c.Request.URL.Path = cleanPath
		}

		c.Next()
	}
}

// Middleware para log de todas as requisições
func requestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Log da requisição
		log.Printf("Requisição recebida: %s %s", c.Request.Method, c.Request.URL.Path)
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			bodyBytes, _ := c.GetRawData()
			if len(bodyBytes) > 0 {
				log.Printf("Corpo da requisição: %s", string(bodyBytes))
				// Restaura o corpo para que possa ser lido novamente pelo handler
				c.Request.Body = &bytesBodyReader{bytes.NewReader(bodyBytes)}
			}
		}

		// Processa a requisição
		c.Next()

		// Log da resposta
		latency := time.Since(start)
		log.Printf("Resposta: %d, Tempo: %v", c.Writer.Status(), latency)
	}
}

func main() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

	// Inicializar configurações de autenticação
	middleware.InitAuthConfig()

	// Lê a connection string do banco de dados a partir da variável de ambiente DATABASE_URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Variável de ambiente DATABASE_URL não definida")
	}

	// Abre a conexão com o PostgreSQL usando o driver do GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}

	// Configuração de log para debugging
	log.Println("Conectado ao banco de dados PostgreSQL")

	// Executa o AutoMigrate para criar/atualizar as tabelas no banco de dados
	if err := db.AutoMigrate(&models.Aluno{}, &models.Curso{}, &models.Inscricao{}); err != nil {
		log.Fatalf("Erro ao migrar o banco de dados: %v", err)
	}
	log.Println("Migração de banco de dados concluída com sucesso")

	// Instancia os repositórios
	alunoRepo := repository.NewAlunoRepository(db)
	cursoRepo := repository.NewCursoRepository(db)
	inscricaoRepo := repository.NewInscricaoRepository(db)

	// Instancia os serviços, injetando os repositórios necessários
	cursoService := service.NewCursoService(cursoRepo, inscricaoRepo)
	alunoService := service.NewAlunoService(alunoRepo, cursoRepo, inscricaoRepo)
	inscricaoService := service.NewInscricaoService(inscricaoRepo)
	//powerBIService := service.NewPowerBIService() // Novo serviço Power BI

	// Instancia os controllers
	alunoController := controller.NewAlunoController(alunoService)
	cursoController := controller.NewCursoController(cursoService)
	authController := controller.NewAuthController()
	inscricaoController := controller.NewInscricaoController(inscricaoService)
	//powerBIController := controller.NewPowerBIController(powerBIService) // Novo controller Power BI

	// Inicializa o roteador Gin (modo baseado em variável de ambiente)
	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.New()

	// Middleware personalizado para evitar redirecionamentos
	router.Use(noRedirectMiddleware())

	// Middleware de log de requisições (para debugging)
	router.Use(requestLoggerMiddleware())

	// Configuração CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Adiciona middleware de recuperação
	router.Use(gin.Recovery())

	log.Println("Configuração CORS aplicada. Todos os origens permitidas.")

	// Rota de verificação de saúde da API
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Rotas de autenticação
	auth := router.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.GET("/validate", middleware.AuthMiddleware(), authController.ValidateToken)
	}

	// Rotas públicas (sem autenticação)
	// Rotas para Curso (apenas visualização)
	router.GET("/curso", cursoController.ListarCursos)
	router.GET("/curso/:id", cursoController.ObterCursoPorID)
	router.GET("/curso/:id/vagas", cursoController.VerificarDisponibilidadeVagas)

	// Rotas para inscrição de alunos (acessível sem autenticação)
	router.POST("/aluno/inscricao", alunoController.CadastrarAlunoEInscrever)

	//// Rotas para o Power BI (requerem autenticação)
	//powerbi := router.Group("/powerbi")
	//powerbi.Use(middleware.AuthMiddleware())
	//{
	//	powerbi.GET("/token", powerBIController.GetEmbedToken)
	//
	//	// Rota de teste para verificar se o endpoint está funcionando
	//	powerbi.GET("/teste", func(c *gin.Context) {
	//		c.JSON(http.StatusOK, gin.H{
	//			"status":    "ok",
	//			"message":   "Endpoint do PowerBI funcionando",
	//			"timestamp": time.Now().Format(time.RFC3339),
	//		})
	//	})
	//}

	// Rotas protegidas (requerem autenticação de admin)
	admin := router.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		// Administração de Cursos
		admin.POST("/curso", cursoController.CriarCurso)
		admin.PUT("/curso/:id", cursoController.AtualizarCurso)
		admin.DELETE("/curso/:id", cursoController.RemoverCurso)
		admin.GET("/curso/:id/inscricoes", cursoController.ListarInscricoesCurso)

		// Administração de Alunos
		admin.GET("/aluno", alunoController.ListarAlunos)
		admin.GET("/aluno/:id", alunoController.ObterAlunoPorID)
		admin.PUT("/aluno/:id", alunoController.AtualizarAluno)
		admin.DELETE("/aluno/:id", alunoController.RemoverAluno)
		admin.POST("/aluno/:id/curso/:cursoId", alunoController.AdicionarAlunoCurso)
		admin.GET("/aluno/:id/inscricoes", alunoController.ListarInscricoesAluno)

		// NOVAS ROTAS: Administração de Inscrições
		admin.GET("/inscricoes", inscricaoController.ListarInscricoes)
		admin.GET("/inscricoes/:id", inscricaoController.ObterInscricaoPorID)
		admin.POST("/relatorio", inscricaoController.GerarRelatorio)
		admin.DELETE("inscricoes/:id", inscricaoController.CancelarInscricao)
	}

	// Define a porta a partir da variável de ambiente PORT ou utiliza 8080 como padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API iniciada na porta %s", port)
	router.Run(":" + port)
}

// bytesBodyReader é um helper para restaurar o corpo da requisição após leitura
type bytesBodyReader struct {
	*bytes.Reader
}

func (r *bytesBodyReader) Close() error {
	return nil
}
