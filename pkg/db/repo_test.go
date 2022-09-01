package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"poc-gnomock/pkg/db"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

var _ = Describe("Db", func() {
	It("should select all people with age over than 10", func() {
		database, container, err := Init()
		Expect(err).ToNot(HaveOccurred())
		repo, err := db.NewRepository(database)
		Expect(err).ToNot(HaveOccurred())
		people, err := repo.Select(context.Background(), 10)
		Expect(err).ToNot(HaveOccurred())
		fmt.Println(people)
		err = gnomock.Stop(container)
		Expect(err).ToNot(HaveOccurred())
	})
})

func Init() (*sql.DB, *gnomock.Container, error) {
	queries := `
		create table test(name varchar(64), age int);
		insert into test(name, age) values ('ivan', 15);
		insert into test(name, age) values ('georgi', 9);
	`
	p := postgres.Preset(
		postgres.WithUser("niki", "test"),
		postgres.WithDatabase("db"),
		postgres.WithQueries(queries), // can accept multiple queries
		// postgres.WithQueriesFile() sets queries from a file
		postgres.WithTimezone("Europe/Sofia"),
		// postgres.WithVersion() sets image version
	)

	//slow network leads to timeout (workaround is to pull the image before executing e.g docker pull postgres:11)
	container, err := gnomock.Start(p)
	// gnomock.InParallel().Start(p).Start(p).Go()
	if err != nil {
		return nil, nil, err
	}

	cfg := db.ConfigDatabase{
		Host:     container.Host,
		Port:     container.DefaultPort(),
		Username: "niki",
		Password: "test",
		Name:     "db",
	}

	database, err := db.ConnectDB(cfg)
	if err != nil {
		return nil, nil, err
	}

	return database, container, nil
}
