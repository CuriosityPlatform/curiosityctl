package servicepreparer

import (
	"bytes"
	"fmt"
	"strings"

	"curiosity/pkg/common/app/dockerclient"
)

type mysqlPreparer struct {
	client     dockerclient.Client
	dbUser     string
	dbPassword string
}

func (preparer *mysqlPreparer) Prepare(composeServiceName string) error {
	rawSQL, err := preparer.buildSQL()
	if err != nil {
		return err
	}

	sql := bytes.NewBuffer([]byte(rawSQL))

	const execCommand = "mysql -u%s -p%s"

	_, err = preparer.client.Exec(dockerclient.ExecParam{
		Service: composeServiceName,
		Command: fmt.Sprintf(execCommand, preparer.dbUser, preparer.dbPassword),
		Reader:  sql,
	})
	return err
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
