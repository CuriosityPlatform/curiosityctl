package servicepreparer

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"curiosity/pkg/common/app/dockerclient"
	"curiosity/pkg/common/infrastructure/progress"
)

type mysqlPreparer struct {
	client     dockerclient.Client
	dbUser     string
	dbPassword string
}

func (preparer *mysqlPreparer) Prepare(ctx context.Context, composeServiceName string) error {
	w := progress.ContextWriter(ctx)

	rawSQL, err := preparer.buildSQL()
	if err != nil {
		return err
	}

	sql := bytes.NewBuffer([]byte(rawSQL))

	const execCommand = "mysql -u%s -p%s"

	eventID := "Migrating database"
	w.Event(progress.StartedEvent(eventID))

	_, err = preparer.client.Exec(dockerclient.ExecParam{
		Service: composeServiceName,
		Command: fmt.Sprintf(execCommand, preparer.dbUser, preparer.dbPassword),
		Reader:  sql,
	})

	if err != nil {
		w.Event(progress.ErrorEvent(eventID))
		return err
	}

	w.Event(progress.DoneEvent(eventID))
	return nil
}

func (preparer *mysqlPreparer) buildSQL() (string, error) {
	const createDatabaseSQL = `CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`

	services, err := preparer.fetchServicesToPrepare()
	if err != nil {
		return "", err
	}

	result := make([]string, 0, len(services))
	for _, service := range services {
		result = append(result, fmt.Sprintf(createDatabaseSQL, service))
	}
	return strings.Join(result, " "), nil
}

func (preparer *mysqlPreparer) fetchServicesToPrepare() ([]string, error) {
	return []string{"patcherservice", "specialservice"}, nil
}
