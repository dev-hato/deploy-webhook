# deploy-webhook
GitHubにてDocker Composeで動作するシステムのリリースを作成したらオンプレにデプロイします。

## 動作方法
2. [オンプレ] 本リポジトリをcloneします。
    ```sh
    git clone git@github.com:dev-hato/deploy-webhook.git
    ```
3. [オンプレ] 本リポジトリに移動します。
    ```sh
    cd deploy-webhook
    ```
4. [オンプレ] `.env.example` をコピーして `.env` を作成します。
    ```sh
    cp .env.example .env
    ```
5. [オンプレ] `.env` の各項目を設定します。
    * `GITHUB_WEBHOOK_SECRET`: 任意の値
    * `LOCAL_REPO_PATH`: オンプレ内のリポジトリのパス
    * `PORT`: deploy-webhookで使用するポート
6. [オンプレ] コードをコンパイルします。
    ```sh
    go build -o deploy-webhook
    ```
6. [オンプレ] deploy-webhookを起動します。
    ```sh
    ./deploy-webhook &
    ```
7. [デプロイしたいシステム] GitHubリポジトリ上でWebhookを作成します。
   * `Payload URL`
      * `{オンプレのURL}:{deploy-webhookで使用するポート}/webhooks`
   * `Content type`
      * `application/json`
   * `Secret`
      * deploy-webhookの `GITHUB_WEBHOOK_SECRET` の値
   * `Which events would you like to trigger this webhook?`
      * `Let me select individual events.`
      * `Releases`
