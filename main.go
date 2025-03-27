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
	// Lê a connection string do banco de dados a partir da variável de ambiente DATABASE_URL
	godotenv.Load()
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
	cursoService := service.NewCursoService(cursoRepo)
	alunoService := service.NewAlunoService(alunoRepo, cursoRepo, inscricaoRepo)

	// Instancia os controllers
	alunoController := controller.NewAlunoController(alunoService)
	cursoController := controller.NewCursoController(cursoService)

	// Inicializa o roteador Gin (mantendo o modo de debug para logs detalhados)
	gin.SetMode(gin.DebugMode)
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

	// Rotas para Aluno (no singular)
	aluno := router.Group("/aluno")
	{
		aluno.GET("", alunoController.ListarAlunos) // Sem barra no final
		aluno.GET("/:id", alunoController.ObterAlunoPorID)
		aluno.POST("", alunoController.CriarAluno) // Sem barra no final
		aluno.PUT("/:id", alunoController.AtualizarAluno)
		aluno.DELETE("/:id", alunoController.RemoverAluno)
		aluno.POST("/curso/:cursoId", alunoController.AdicionarAlunoCurso)
	}

	// Rotas para Curso (no singular)
	curso := router.Group("/curso")
	{
		// Importante: Rotas sem barra no final
		curso.GET("", cursoController.ListarCursos)
		curso.GET("/:id", cursoController.ObterCursoPorID)
		curso.POST("", func(c *gin.Context) {
			log.Println("Handler CriarCurso chamado")
			cursoController.CriarCurso(c)
		})
		curso.PUT("/:id", cursoController.AtualizarCurso)
		curso.DELETE("/:id", cursoController.RemoverCurso)
		curso.GET("/:id/vagas", cursoController.VerificarDisponibilidadeVagas)
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
