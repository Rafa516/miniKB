-- Executado automaticamente pela imagem oficial do Postgres na primeira
-- vez que o volume de dados é criado (scripts em /docker-entrypoint-initdb.d
-- rodam só na inicialização de um volume vazio). Cria um banco separado
-- só para os testes automatizados (go test), pra não misturar dados de
-- teste com os dados de verdade da aplicação no banco "minikb".
CREATE DATABASE minikb_test;
