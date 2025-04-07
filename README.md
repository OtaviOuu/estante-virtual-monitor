![Screencastfrom06-04-2025222558-ezgif com-video-to-gif-converter](https://github.com/user-attachments/assets/2f798ee8-2831-4254-b1e9-a365b855b143)

## Funcionalidades
- Web scraping do site Estante Virtual
- Coleta de informações detalhadas de livros recém-lançados
- Envio automático das informações para um canal do Telegram
- Formatação de mensagens com Markdown

## Pré-requisitos

- Go 1.16+
- Conta no Telegram
- Bot do Telegram (obtenha um token com @BotFather)
- Canal do Telegram (adicione o bot como administrador)

## Dependências

- [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery) - Para scraping de HTML
- [github.com/joho/godotenv](https://github.com/joho/godotenv) - Para gerenciamento de variáveis de ambiente

## Instalação

1. Clone o repositório:
   ```bash
   git clone https://github.com/OtaviOuu/estante-virtual-monitor
   cd estante-virtual-bot
   ```

2. Instale as dependências:
   ```bash
   go mod tidy
   ```

3. Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:
   ```
   BOT_TOKEN=seu_token_do_bot_telegram
   CHANNEL_ID=id_do_seu_canal_telegram
   ```

## Uso

Execute o programa:

```bash
go run main.go
```

Alternativamente, compile e execute o binário:

```bash
go build -o estante-bot
./estante-bot
```
