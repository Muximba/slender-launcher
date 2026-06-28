# Muximba — notas do launcher

Launcher do Muximba (MTC), baseado no slender-launcher (**Wails: Go + Svelte**). O `README.md` é o do slender original. Aqui o que é específico do Muximba.

Visão geral: `ARCHITECTURE.md` em `Muximba/distribution`.

## O que faz

Baixa 1 arquivo (o launcher), e ele baixa/atualiza o **client** do R2 comparando sha256 (baixa só o que mudou), depois dá **Play**. Também se **auto-atualiza** (compara o próprio binário com o `.sha256` remoto).

- `baseURL = https://files.muximba.dev/` (override via `base_url` no `config.toml`). Daí vem `client.<os>.json`/`assets.<os>.json` + a árvore do client + o próprio binário do launcher.
- Web service: `https://app.muximba.dev/login.php` (`status.go`) — server status + news na tela.

## Customizações Muximba

- Branding: `frontend/src/assets/images/logo-universal.png` (brasão Muximba) + `background-artwork.jpg`.
- `wails.json`: `name`/`outputfilename` = `Muximba`.
- Server status (pill online/offline + players) e News clicável (abre o site). Discord e "X players online" clicáveis com cursor de mão.
- Wails **2.12** (2.5.1 não honrava Width/Height no Linux → janela minúscula).
- `app.go`: Play roda `syscall.Exec(otclient, --battleeye)`; Update reseta contadores antes do diff.

## Build + publish

```bash
cd ~/projects/muximba/slender-launcher
export PATH="$PATH:$(go env GOPATH)/bin"
wails build                            # Linux  -> build/bin/Muximba  (precisa libgtk-3-dev + libwebkit2gtk-4.0-dev)
wails build -platform windows/amd64    # Windows -> build/bin/Muximba.exe  (cross-compile precisa mingw-w64)
cd build/bin
for f in Muximba Muximba.exe; do sha256sum "$f" | awk '{print $1}' > "$f.sha256"; done
. ~/projects/muximba/.secrets/r2.env
for f in Muximba Muximba.sha256 Muximba.exe Muximba.exe.sha256; do
  aws s3 cp "$f" "s3://$R2_BUCKET/$f" --profile r2muximba --endpoint-url "$R2_ENDPOINT"
done
```

> Cross-compile Windows precisa de `mingw-w64` instalado (hoje pode faltar na máquina). Sem ele, builda só o Linux.

## Pendências relevantes

- Display do download mostra o filename cru ("otclient") → trocar pra "MTC (Muximba Client)".
- **Prune**: o launcher só baixa o que está no manifest, não deleta o que sobrou. Remoções (ex.: cavebot/terminal) não propagam pra instalações existentes — precisa prune em `modules/`, `mods/`, `data/`.
- CI de build/release.
