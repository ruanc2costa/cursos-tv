package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tvtec/controller"
	"tvtec/models"
	"tvtec/repository"
	"tvtec/service"
)

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

	// Executa o AutoMigrate para criar/atualizar as tabelas no banco de dados
	if err := db.AutoMigrate(&models.Aluno{}, &models.Curso{}, &models.Inscricao{}); err != nil {
		log.Fatalf("Erro ao migrar o banco de dados: %v", err)
	}

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

	// Inicializa o roteador Gin com o middleware de Recovery
	router := gin.Default() // Usa Default() que já inclui Logger e Recovery

	// Configuração do middleware CORS para permitir requisições de outros domínios
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ou especifique as origens permitidas
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-usuario"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rotas para Aluno (no singular e sem prefixo api/v1)
	aluno := router.Group("/aluno")
	{
		aluno.GET("/", alunoController.ListarAlunos)
		aluno.GET("/:id", alunoController.ObterAlunoPorID)
		aluno.POST("/", alunoController.CriarAluno)
		aluno.PUT("/:id", alunoController.AtualizarAluno)
		aluno.DELETE("/:id", alunoController.RemoverAluno)
		aluno.POST("/curso/:cursoId", alunoController.AdicionarAlunoCurso)
	}

	// Rotas para Curso (no singular e sem prefixo api/v1)
	curso := router.Group("/curso")
	{
		curso.GET("/", cursoController.ListarCursos)
		curso.GET("/:id", cursoController.ObterCursoPorID)
		curso.POST("/", cursoController.CriarCurso)
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
