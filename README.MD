# Desafio Backend 
Repositório do desafio: https://github.com/PicPay/picpay-desafio-backend
Deixo claro que não estou participando do processo seletivo. Apenas utilizei esse desafio para praticar

### 🛠 Tecnologias
As seguintes ferramentas foram usadas na construção do projeto:
- [Go](https://go.dev/)
- [Keycloak](https://www.keycloak.org/)
- [MongoDB](https://www.mongodb.com/)
- [Docker](https://www.docker.com/)

## Como Executar a Aplicação
1. Clone este repositório para o seu ambiente local.
2. Navegue até a pasta do projeto através do terminal:
    ```bash
    cd pp_chlg
    ```
3. Inicie a aplicação com o seguinte comando:
    ```bash
    ./start_development_environment.sh
    ```
4. Esse script vai acessar o bash do container da aplicação, sinta-se a vontade para executar alguns dos testes que existem na camada de 'usecase', iniciar o servidor, customizar etc..

-  Não esqueça de criar o arquivo .env com as variáveis que estão listadas no arquivo de exemplo .env.example.

### Endpoints
Cadastrar Usuário:

    Método: POST
    URL: http://localhost:8000/user/signup
    Headers:
        Content-Type: application/json

Exemplo de corpo da requisição:

```json
{
  "first_name": "LOREM",
  "last_name": "IPSUM",
  "email": "lorem@ipsum.com",
  "password": "P@ssw0rd",
  "taxpayer_identification": "06947283000160", // CNPJ DO GOOGLE
  "role": "lorem_ipsum"
}
```
- Observações:
    - O atributo 'taxpayer_identification' deve ser um CNPJ ou CPF válido.
    - O atributo 'role' deve estar cadastrado no Keycloak e nas variáveis de ambiente.

Autenticar Usuário:

    Método: POST
    URL: http://localhost:8000/user/signin
    Headers:
        Content-Type: application/json

Exemplo de corpo da requisição:

```json
{
  "email": "lorem@ipsum.com",
  "password": "P@ssw0rd"
}
```

Realizar Transferência:

    Método: POST
    URL: http://localhost:8000/transfer
    Headers:
        x_access_token: Bearer {{ACCESS_TOKEN}}

Exemplo de corpo da requisição:

```json
{
  "payee": {{id do usuário que receberá a transferência}},
  "value": ? // Valor da transferência
}
```

Obter Transferências:

    Método: GET
    URL: http://localhost:8000/transfer
    Headers:
        x_access_token: Bearer {{ACCESS_TOKEN}}
        page: X // Número da página
        limit: Y // Quantidade de itens por página


