package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"tvtec/models/controller"
	"tvtec/repository"
	service "tvtec/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tvtec/models" // ajuste o import conforme a localização dos seus modelos
)

func main() {
	// Configuração da conexão com o MySQL (ajuste conforme necessário)
	dsn := "root:nova_senha@tcp(localhost:3306)/cadastrotv?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar ao MySQL: %v", err)
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

	// Inicia o servidor na porta 8080
	log.Println("API iniciada na porta 8080")
	router.Run(":8080")
}
