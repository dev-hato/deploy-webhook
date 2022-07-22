# deploy-webhook
GitHubにてDocker Composeで動作するシステムのリリースを作成したらオンプレにデプロイします。

## 動作方法
### オンプレでの作業
#### デプロイしたいシステム
1. デプロイしたいシステムをcloneします。
2. オンプレ用の `docker-compose.yml` の差分である `onpremises.docker-compose.yml` を配置します。

#### deploy-webhook
1. 本リポジトリをcloneします。
    ```sh
    git clone git@github.com:dev-hato/deploy-webhook.git
    ```
2. 本リポジトリに移動します。
    ```sh
    cd deploy-webhook
    ```
3. `.env.example` をコピーして `.env` を作成します。
    ```sh
    cp .env.example .env
    ```
4. `.env` の各項目を設定します。
    * `GITHUB_WEBHOOK_SECRET`: 任意の値
    * `LOCAL_REPO_PATH`: オンプレ内のリポジトリのパス
    * `PORT`: deploy-webhookで使用するポート
5. コードをコンパイルします。
    ```sh
    go build -o deploy-webhook
    ```
6. deploy-webhookを起動します。
    ```sh
    ./deploy-webhook &
    ```
### GitHubのWebhook作成
1. デプロイしたいシステムのGitHubリポジトリ上でWebhookを作成します。
   * `Payload URL`
      * `{オンプレのURL}:{deploy-webhookで使用するポート}/webhooks`
   * `Content type`
      * `application/json`
   * `Secret`
      * deploy-webhookの `GITHUB_WEBHOOK_SECRET` の値
   * `Which events would you like to trigger this webhook?`
      * `Let me select individual events.`
      * `Releases`
