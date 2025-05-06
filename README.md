# RateLimiter

Rate limiter em Go que pode ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

Este é um exemplo de rate limiter em Go que pode ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter é capaz de limitar o número de requisições com base em dois critérios:

1. **Endereço IP:** O rate limiter restringe o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
2. **Token de Acesso:** O rate limiter também pode limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. 
    O Token é informado no header no seguinte formato: `API_KEY: <TOKEN>`
3. As configurações de limite do token de acesso sobrepoem as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter utiliza as informações do token.


## Como rodar:
Utilize o comando `make up` no terminal e siga as instruções. Este comando irá subir o redis e a app e mostrará uma forma de enviar requisições para a api da app para testar.

A aplicação exibe na console log toda requisição recebida de qual IP, utilizando qual token, o limite de acordo se foi informado o token ou não, e se a requisição foi permitida. Caso atinja o limite o sistema exibe uma mensagem de erro
