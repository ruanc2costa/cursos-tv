package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tvtec/controller"
	"tvtec/models"
	"tvtec/repository"
	"tvtec/services"
)

func main() {
	// Obtém a connection string do Supabase a partir da variável de ambiente DATABASE_URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Variável de ambiente DATABASE_URL não definida")
	}

	// Abre a conexão com o PostgreSQL utilizando o driver do GORM para Postgres
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao PostgreSQL: %v", err)
	}

	// Executa o AutoMigrate para criar/atualizar as tabelas
	if err := db.AutoMigrate(&models.Aluno{}, &models.Curso{}); err != nil {
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

	// Registra as rotas dos controllers
	alunoController.RegisterRoutes(router)
	cursoController.RegisterRoutes(router)

	// Obtém a porta da variável de ambiente PORT, ou usa 8080 como padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API iniciada na porta %s", port)
	router.Run(":" + port)
}
