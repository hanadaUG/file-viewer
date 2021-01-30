# file-viewer
golangで実装したシンプルなディレクトリリスティング（URLで指定するとディレクトリに含まれるファイル一覧を表示するサーバー）です。
```
# 実行方法
# root オプション => 公開するフォルダを指定
# port オプション => 使用するportを指定
$ ./file-viewer \
-root ${HOME}/git/file-viewer/sample \
-port 1234
```