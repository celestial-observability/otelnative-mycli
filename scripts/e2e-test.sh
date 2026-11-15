#!/usr/bin/env bash
set -uo pipefail
#
# E2E テスト
#
# バイナリをビルドして終了コード/stdout/stderrを検証する
#

#
# グローバル変数
#
RUN_ARGS=()
PASS=0
FAIL=0
STATUS=0
STDOUT=''
STDERR=''
STDOUT_TMP=$(mktemp)
STDERR_TMP=$(mktemp)
CMD_TIMEOUT_SEC=10

#
# ヘルパー関数
#
# 終了時に一時ファイルを削除する
cleanup() {
  rm -f "$STDOUT_TMP" "$STDERR_TMP"
}

trap cleanup EXIT

# テスト成功を記録して緑色で表示する
_pass() {
  local label="$1"
  PASS=$((PASS + 1))
  if [[ -t 2 ]]; then
    printf '\033[32m  PASS\033[0m  %s\n' "$label"
  else
    printf '  PASS  %s\n' "$label"
  fi
}

# テスト失敗を記録して赤色で表示する
_fail() {
  local label="$1"
  FAIL=$((FAIL + 1))
  if [[ -t 2 ]]; then
    printf '\033[31m  FAIL\033[0m  %s\n' "$label"
  else
    printf '  FAIL  %s\n' "$label"
  fi
}

# 終了コードが期待値と一致するか検証する
expect_exit_code() {
  local expected="$1"
  local label="$2"
  if [ "$STATUS" -eq "$expected" ]; then
    _pass "exit $STATUS  $label"
  else
    _fail "exit $STATUS (expected $expected)  $label"
  fi
}

# stdout に指定の文字列が含まれるか検証する
expect_stdout_contains() {
  local needle="$1"
  local label="$2"
  if printf '%s' "$STDOUT" | grep -qF "$needle"; then
    _pass "stdout '$needle' $label"
  else
    _fail "stdout missing '$needle'  $label  (actual: $(printf '%s' "$STDOUT" | head -1))"
  fi
}

# stderr が空であるか検証する
expect_stderr_empty() {
  local label="$1"
  if [ -z "$STDERR" ]; then
    _pass "stderr is empty  $label"
  else
    _fail "stderr is not empty  $label  (actual: $(printf '%s' "$STDERR" | head -1))"
  fi
}

# hi をタイムアウト付きで実行して終了コード/stdout/stderrを記録する
run_cmd() {
  perl -e "alarm ${CMD_TIMEOUT_SEC}; exec @ARGV" \
    env ./hi "${RUN_ARGS[@]+"${RUN_ARGS[@]}"}" >"$STDOUT_TMP" 2>"$STDERR_TMP"
  STATUS=$?
  STDOUT=$(cat "$STDOUT_TMP")
  STDERR=$(cat "$STDERR_TMP")
}

#
# ビルド
#
echo ''
echo '=== Build ===='

if mise run build; then
  _pass 'mise run build'
else
  if [[ -t 2 ]]; then
    printf '\033[31m ERROR\033[0m Failed build\n'
  else
    printf '  ERROR  Failed build\n'
  fi
  exit 1
fi

#
# 引数なし
#
echo ''
echo '=== no args ==='
RUN_ARGS=()
run_cmd
expect_exit_code 0 'no args'
expect_stdout_contains 'Hello, World!' "'Hello, World!' is in stdout"
expect_stderr_empty 'no args'

#
# サマリー
#
echo ''
echo '─────────────────────────────────────────────'
printf "  PASS: %d  FAIL: %d  \n" "$PASS" "$FAIL"
echo '─────────────────────────────────────────────'

# テストが1つでも失敗したら非ゼロで終了する
[ "$FAIL" -eq 0 ]
