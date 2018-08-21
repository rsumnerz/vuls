
# Vuls: VULnerability Scanner

[![Slack](https://img.shields.io/badge/slack-join-blue.svg)](http://goo.gl/forms/xm5KFo35tu)

Vulnerability scanner for Linux, agentless, written in golang.

[README in English](https://github.com/future-architect/vuls/blob/master/README.md)  
Slackチームは[こちらから](http://goo.gl/forms/xm5KFo35tu)参加できます。(日本語でオッケーです)

[![asciicast](https://asciinema.org/a/bazozlxrw1wtxfu9yojyihick.png)](https://asciinema.org/a/bazozlxrw1wtxfu9yojyihick)

![Vuls-slack](img/vuls-slack-ja.png)


----

# Abstract

毎日のように発見される脆弱性の調査やソフトウェアアップデート作業は、システム管理者にとって負荷の高いタスクである。
プロダクション環境ではサービス停止リスクを避けるために、パッケージマネージャの自動更新機能を使わずに手動更新で運用するケースも多い。
だが、手動更新での運用には以下の問題がある。
- システム管理者がNVDなどで新着の脆弱性をウォッチし続けなければならない
- サーバにインストールされているソフトウェアは膨大であり、システム管理者が全てを把握するのは困難
- 新着の脆弱性がどのサーバに該当するのかといった調査コストが大きく、漏れる可能性がある


Vulsは上に挙げた手動運用での課題を解決するツールであり、以下の特徴がある。
- システムに関係ある脆弱性のみ教えてくれる
- その脆弱性に該当するサーバを教えてくれる
- 自動スキャンのため脆弱性検知の漏れを防ぐことができる
- CRONなどで定期実行、レポートすることで脆弱性の放置を防ぐことできる

![Vuls-Motivation](img/vuls-motivation.png)

----

# Main Features

- Linuxサーバに存在する脆弱性をスキャン
    - Ubuntu, Debian, CentOS, Amazon Linux, RHELに対応
    - クラウド、オンプレミス、Docker
- OSパッケージ管理対象外のミドルウェアをスキャン
    - プログラミング言語のライブラリやフレームワーク、ミドルウェアの脆弱性スキャン
    - CPEに登録されているソフトウェアが対象
- エージェントレスアーキテクチャ
    - スキャン対象サーバにSSH接続可能なマシン1台にセットアップするだけで動作
- 設定ファイルのテンプレート自動生成
    - CIDRを指定してサーバを自動検出、設定ファイルのテンプレートを生成
- EmailやSlackで通知可能（日本語でのレポートも可能）
- 付属するTerminal-Based User Interfaceビューアでは、Vim風キーバインドでスキャン結果を参照可能
reeBSDは今のところRoot権限なしでスキャン可能

----

# Usage: Scan

```
$ vuls scan -help
scan:
        scan
                [-config=/path/to/config.toml]
                [-results-dir=/path/to/results]
                [-log-dir=/path/to/log]
                [-cachedb-path=/path/to/cache.db]
                [-ssh-external]
                [-containers-only]
            -dbpath string
            /path/to/sqlite3/DBfile (default "$PWD/cve.sqlite3")
      -debug
            debug mode
      -debug-sql
            SQL debug mode
      -dump-path string
            /path/to/dump.json (default "$PWD/cve.json")
      -entire
            Fetch data for entire period.(This operation is time-consuming) (default: false)
      -month
            Fetch data in the last month (default: false)
      -week
            Fetch data in the last week. (default: false)
>>>>>>>+master
ug mode
  -http-proxy string
        http://proxy-url:port (default: empty)
  -log-dir string
        /path/to/log (default "/var/log/vuls")
  -pipe
        Use stdin via PIPE
  -results-dir string
        /path/to/results
  -skip-broken
        [For CentOS] yum update changelog with --skip-broken option
  -ssh-external
        Use external ssh command. Default: Use the Go native implementation
```

## -ssh-external option

Vulsは２種類のSSH接続方法をサポートしている。

デフォルトでは、Goのネイティブ実装 (crypto/ssh) を使ってスキャンする。
これは、SSHコマンドがインストールされていない環境でも動作する（Windowsなど）  

外部SSHコマンドを使ってスキャンするためには、`-ssh-external`を指定する。
SSH Configが使えるので、ProxyCommandを使った多段SSHなどが可能。  
CentOSでは、スキャン対象サーバの/etc/sudoersに以下を追加する必要がある(user: vuls)
```
Defaults:vuls !requiretty
```

## -ask-key-password option

| SSH key password |  -ask-key-password | |
|:-----------------|:-------------------|:----|
| empty password   |                 -  | |
| with password    |           required | or use ssh-agent |

## Example: Scan all servers defined in config file
```
$ vuls scan -ask-key-password
```
この例では、
- SSH公開鍵認証（秘密鍵パスフレーズ）を指定
- configに定義された全サーバをスキャン

## Example: Scan specific servers
```
$ vuls scan server1 server2
```
この例では、
- SSH公開鍵認証（秘密鍵パスフレーズなし）
- ノーパスワードでsudoが実行可能
- configで定義されているサーバの中の、server1, server2のみスキャン

## Example: Scan via shell instead of SSH.

ローカルホストのスキャンする場合、SSHではなく直接コマンドの発行が可能。  
config.tomlのhostに`localhost または 127.0.0.1`かつ、portに`local`を設定する必要がある。  
For more details, see [Architecture section](#architecture)

- config.toml
  ```
  [servers]

  [servers.localhost]
  host         = "localhost" # or "127.0.0.1"
  port         = "local"
  ```

### cronで動かす場合

RHEL/CentOSの場合、スキャン対象サーバの/etc/sudoersに以下を追加する必要がある。(user: vuls)
```
Defaults:vuls !requiretty
```

## Example: Scan containers (Docker/LXD)

>>>>>>> 688cfd6

    ```

- すべての期間の脆弱性情報を取得(1時間以上かかる)
    ```
    $ go-cve-dictionary fetchjvn -entire
    ```

- 直近1ヶ月間に更新された脆弱性情報を取得(1分未満)
    ```
    $ go-cve-dictionary fetchjvn -month
    ```

- 直近1週間に更新された脆弱性情報を取得(1分未満)
    ```
    $ go-cve-dictionary fetchjvn -week
    ```

- 脆弱性情報の自動アップデート  
Cronなどのジョブスケジューラを用いて実現可能。  
-week オプションを指定して夜間の日次実行を推奨。


## スキャン実行

```
$ vuls scan -lang=ja
```
Scan時にlang=jaを指定すると脆弱性レポートが日本語になる  
slack, emailは日本語対応済み TUIは日本語表示未対応

