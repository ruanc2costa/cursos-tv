package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tvtec/controller"
	"tvtec/models"
	"tvtec/repository"
	"tvtec/service"
)

// CORSMiddleware é um middleware personalizado que lida com solicitações CORS
// Permite qualquer origem enquanto mantém suporte a credenciais
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtém a origem da solicitação
		origin := c.Request.Header.Get("Origin")

		// Se a origem estiver presente, configura os cabeçalhos CORS
		// usando a origem exata em vez de "*" para permitir credenciais
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With, x-usuario")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 horas
		}

		// Tratamento especial para solicitações preflight OPTIONS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 204 No Content
			return
		}

		// Prossegue para o próximo middleware ou handler
		c.Next()
	}
}

func main() {
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
	cursoService := service.NewCursoService(cursoRepo)
	alunoService := service.NewAlunoService(alunoRepo, cursoRepo, inscricaoRepo)

	// Instancia os controllers
	alunoController := controller.NewAlunoController(alunoService)
	cursoController := controller.NewCursoController(cursoService)

	// Inicializa o roteador Gin
	router := gin.New()

	// Adiciona middleware de recuperação e logging
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Adiciona middleware personalizado de CORS antes de qualquer outro middleware
	router.Use(CORSMiddleware())

	// Middleware para log de todas as requisições (para debugging)
	router.Use(func(c *gin.Context) {
		log.Printf("[%s] Requisição recebida: %s %s", time.Now().Format(time.RFC3339), c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

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
		aluno.GET("/", func(c *gin.Context) {
			log.Println("Handler ListarAlunos chamado")
			alunoController.ListarAlunos(c)
		})
		aluno.GET("/:id", alunoController.ObterAlunoPorID)
		aluno.POST("/", alunoController.CriarAluno)
		aluno.PUT("/:id", alunoController.AtualizarAluno)
		aluno.DELETE("/:id", alunoController.RemoverAluno)
		aluno.POST("/curso/:cursoId", alunoController.AdicionarAlunoCurso)
	}

	// Rotas para Curso (no singular)
	curso := router.Group("/curso")
	{
		curso.GET("/", cursoController.ListarCursos)
		curso.GET("/:id", cursoController.ObterCursoPorID)
		curso.POST("/", func(c *gin.Context) {
			log.Println("Handler CriarCurso chamado")
			// Log do corpo da requisição para debugging
			bodyBytes, _ := c.GetRawData()
			if len(bodyBytes) > 0 {
				log.Printf("Corpo da requisição: %s", string(bodyBytes))
				// Restaura o corpo para que possa ser lido novamente pelo handler
				c.Request.Body = &bytesBodyReader{bytes.NewReader(bodyBytes)}
			}
			cursoController.CriarCurso(c)
		})
		curso.PUT("/:id", cursoController.AtualizarCurso)
		curso.DELETE("/:id", cursoController.RemoverCurso)
		curso.GET("/:id/vagas", cursoController.VerificarDisponibilidadeVagas)
	}

	// Adiciona handler OPTIONS global para lidar com preflight CORS para todas as rotas
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

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
