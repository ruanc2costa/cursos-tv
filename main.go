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
	// Carrega a connection string do Supabase a partir da variável de ambiente
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Variável de ambiente DATABASE_URL não definida")
	}

	// Abre a conexão com o PostgreSQL usando o driver do GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}

	// Executa o AutoMigrate para criar/atualizar as tabelas
	if err := db.AutoMigrate(&models.Aluno{}, &models.Curso{}, &models.Inscricao{}); err != nil {
		log.Fatalf("Erro ao migrar o banco de dados: %v", err)
	}

	// Instancia os repositórios e serviços
	alunoRepo := repository.NewAlunoRepository(db)
	alunoService := service.NewAlunoService(db)
	cursoService := service.NewCursoService(db)

	// Instancia os controllers
	alunoController := controller.NewAlunoController(alunoService, alunoRepo)
	cursoController := controller.NewCursoController(cursoService)

	// Inicializa o roteador do Gin
	router := gin.Default()

	// Configura o middleware CORS para permitir requisições de qualquer origem
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ou especifique as origens permitidas, ex: "http://seu-frontend.com"
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Registra as rotas dos controllers
	alunoController.RegisterRoutes(router)
	cursoController.RegisterRoutes(router)

	// Obtém a porta definida na variável de ambiente PORT, ou usa 8080 como padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API iniciada na porta %s", port)
	router.Run(":" + port)
}
