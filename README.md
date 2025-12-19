# validator-fiber

## Descrição

Este projeto é um validador para o framework web [Fiber](https://gofiber.io/). Ele fornece uma maneira simples e eficiente de validar requisições HTTP, como dados de formulários ou JSON, antes de serem processados pela sua aplicação Fiber.

## Funcionalidades

Este validador encapsula a popular biblioteca `github.com/go-playground/validator/v10` para oferecer validação de structs. As principais funcionalidades incluem:

*   **Validação de Structs**: Valida qualquer struct Go usando tags de validação (`validate:"required,email"`, `validate:"min=3"`, etc.).
*   **Retorno de Erros Formatados**: Os erros de validação são retornados em uma slice de `ErrorsHandle`, que contém informações como `FailedField` (o campo que falhou a validação), `Tag` (a tag de validação que causou o erro) e `Error` (um booleano indicando erro).

go get github.com/stolato/validator-fiber

## Instalação

Para usar este validador em seu projeto Fiber, você pode instalá-lo usando `go get`:

```bash
go get github.com/stolato/validator-fiber
```

## Uso

Para integrar o validador com sua aplicação Fiber, você pode importar o pacote e usar a função `Validator` em seus handlers.

### Estrutura de Erros

O validador retorna uma slice de `ErrorsHandle` em caso de falha na validação, com a seguinte estrutura:

```go
type ErrorsHandle struct {
	FailedField string      // Nome do campo que falhou a validação
	Tag         string      // A tag de validação (ex: "required", "min")
	Value       interface{} // O valor do campo (útil para debug)
	Error       bool        // Sempre true para indicar um erro
}
```

### Exemplo Básico

```go
// Exemplo de uso
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	validator_fiber "github.com/stolato/validator-fiber"
)

type User struct {
	Name  string `json:"name" validate:"required,min=3,max=32"`
	Email string `json:"email" validate:"required,email"`
	Age   uint   `json:"age" validate:"required,gte=18"`
}

func main() {
	app := fiber.New()

	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// Validar o struct User
		if errs := validator_fiber.Validator(user); errs != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": errs,
			})
		}

		// Se a validação passar, você pode prosseguir com a lógica de negócio
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User created successfully",
			"user":    user,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
```

## Contribuição

Contribuições são bem-vindas! Se você quiser melhorar este projeto, por favor, sinta-se à vontade para abrir uma issue ou enviar um pull request.

## Licença

Este projeto é licenciado sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
