# Tech Challenge: Autenticação de Usuário

Este repositório provisiona uma infraestrutura serverless robusta para o serviço de autenticação de usuários, integrando uma aplicação Go (Gin) como AWS Lambda, com persistência de dados no RDS (PostgreSQL) e gerenciamento de sessões no DynamoDB.

## Tecnologias Utilizadas

- **Linguagem**: Go (1.25.0)
- **Infraestrutura**: Terraform, AWS Lambda, Amazon ECR, Amazon RDS (PostgreSQL), Amazon DynamoDB
- **Ambiente Local**: LocalStack, Docker, Docker Compose
- **Monitoramento**: Datadog

## Passos para Execução e Deploy

### Pré-requisitos

- **Go**: Versão 1.25.0 ou superior.
- **Docker & Docker Compose**: Necessários para o ambiente LocalStack.
- **AWS CLI / awslocal**: Ferramentas de linha de comando para interação com AWS.
- **Terraform**: Para o deploy da infraestrutura real.

### Execução Local (via LocalStack)

1. **Suba os serviços locais**: Utilize o comando abaixo para iniciar o LocalStack e o banco de dados PostgreSQL.
   ```bash
   make local
   ```
   *Este comando automatiza o `docker-compose up`, o provisionamento inicial (IAM, DynamoDB, API Gateway) e o deploy da Lambda.*

2. **Popule o banco de dados (opcional)**: Insira usuários de teste para validar o login.
   ```bash
   make seed
   ```

3. **Teste o endpoint**: Use o comando pronto para simular uma requisição de autenticação.
   ```bash
   make curl
   ```

### Deploy Real na AWS

1. **Inicialize o Terraform**: Configure o backend (S3) e os providers.
   ```bash
   cd terraform
   terraform init
   ```

2. **Planeje as mudanças**: Verifique quais recursos serão criados ou alterados.
   ```bash
   terraform plan -var-file="suas-variaveis.tfvars"
   ```

3. **Aplique a infraestrutura**: Execute o deploy para a nuvem.
   ```bash
   terraform apply -var-file="suas-variaveis.tfvars"
   ```

*Nota: O backend utiliza S3 para persistência do estado do Terraform (tfstate), garantindo consistência em ambientes compartilhados.*

## Diagrama da Arquitetura

![Arquitetura da Base](.github/misc/lambda-authentication-architecture.png)

## APIs (Swagger/Postman)

*(em branco)*
