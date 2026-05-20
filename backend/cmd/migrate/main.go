package main

import (
	"flag"
	"log"

	"github.com/genealogy-ma/platform/internal/infrastructure/config"
	"github.com/genealogy-ma/platform/internal/infrastructure/persistence"
)

func main() {
	// 解析命令行参数
	reset := flag.Bool("reset", false, "Reset database (drop all tables)")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := persistence.NewDatabase(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	migrator := persistence.NewMigrator(db.DB)

	// 重置数据库
	if *reset {
		log.Println("Resetting database...")
		if err := migrator.Reset(); err != nil {
			log.Fatalf("Failed to reset database: %v", err)
		}
		log.Println("Database reset successfully")
	}

	// 执行迁移
	log.Println("Running migrations...")
	if err := migrator.Run(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("All migrations completed successfully")
}
