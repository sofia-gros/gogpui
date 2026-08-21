# 移植進捗台帳 (Zed)

## 使い方
1. 最初のExploreタスクとして、AIアシスタントに `crates/ui/src`(gpui-componentのRustソース)配下の
   モジュール一覧を出力させ、下表の行を自動生成させる。手作業で覚えているコンポーネント名を
   使って埋めない(漏れ・誤記の元)。
2. 1コンポーネントの作業が完了するたびに、この台帳を必ず更新する。更新しないままタスク完了を
   報告することを禁止する(.zed/rules.md参照)。
3. 「状態」列は下記のいずれかのみを使う: `未着手` / `Plan承認待ち` / `実装中` / `検証中` / `完了` / `要修正`

| # | Rustコンポーネント (crates/ui/src) | Go対応パッケージ/型 | API 1:1対応 | ユニットテスト | ゴールデン画像diff | 状態 | 備考 |
|---|---|---|---|---|---|---|---|
| 1 | button | components/button | はい | はい | はい(Mock) | 完了 | docs/button-status.mdに一部スキップを記録 |
| 2 | select (ドロップダウン) | | | | | 未着手 | スクロール追従の期待仕様を先に `docs/reference/select.md` に記録すること |
| 3 | table (Data Table) | | | | | 未着手 | 仮想化リストあり、優先度は後半でよい |
| 4 | list | | | | | 未着手 | |
| 5 | sidebar | | | | | 未着手 | |
| 6 | dock (パネルレイアウト) | | | | | 未着手 | |
| 7 | chart | | | | | 未着手 | |
| 8 | editor (コードエディタ) | | | | | 未着手 | Tree-sitter依存、最終フェーズ推奨 |
| 9 | accordion | | | | | 未着手 | |
| 10 | alert | | | | | 未着手 | |
| 11 | avatar | | | | | 未着手 | |
| 12 | badge | | | | | 未着手 | |
| 13 | breadcrumb | | | | | 未着手 | |
| 14 | checkbox | | | | | 未着手 | |
| 15 | clipboard | | | | | 未着手 | |
| 16 | collapsible | | | | | 未着手 | |
| 17 | color_picker | | | | | 未着手 | |
| 18 | combobox | | | | | 未着手 | |
| 19 | description_list | | | | | 未着手 | |
| 20 | dialog | | | | | 未着手 | |
| 21 | form | | | | | 未着手 | |
| 22 | group_box | | | | | 未後手 | |
| 23 | highlighter | | | | | 未着手 | |
| 24 | hover_card | | | | | 未着手 | |
| 25 | icon | | | | | 未着手 | |
| 26 | input | | | | | 未着手 | |
| 27 | inspector | | | | | 未着手 | |
| 28 | kbd | | | | | 未着手 | |
| 29 | label | | | | | 未着手 | |
| 30 | link | | | | | 未着手 | |
| 31 | menu | | | | | 未着手 | |
| 32 | native_menu | | | | | 未着手 | |
| 33 | notification | | | | | 未着手 | |
| 34 | pagination | | | | | 未着手 | |
| 35 | plot | | | | | 未着手 | |
| 36 | popover | | | | | 未着手 | |
| 37 | progress | | | | | 未着手 | |
| 38 | radio | | | | | 未着手 | |
| 39 | rating | | | | | 未着手 | |
| 40 | scroll | | | | | 未着手 | |
| 41 | searchable_list | | | | | 未着手 | |
| 42 | separator | | | | | 未着手 | |
| 43 | setting | | | | | 未着手 | |
| 44 | sheet | | | | | 未着手 | |
| 45 | skeleton | | | | | 未着手 | |
| 46 | slider | | | | | 未着手 | |
| 47 | spinner | | | | | 未着手 | |
| 48 | status_bar | | | | | 未着手 | |
| 49 | stepper | | | | | 未着手 | |
| 50 | switch | | | | | 未着手 | |
| 51 | tab | | | | | 未着手 | |
| 52 | tag | | | | | 未着手 | |
| 53 | text | | | | | 未着手 | |
| 54 | theme | | | | | 未着手 | |
| 55 | time | | | | | 未着手 | |
| 56 | title_bar | | | | | 未着手 | |
| 57 | tooltip | | | | | 未着手 | |
| 58 | tree | | | | | 未着手 | |
| 59 | virtual_list | | | | | 未着手 | |
| 60 | window_border | | | | | 未着手 | |

## 移植の優先順位の考え方
- 依存関係の少ない基本コンポーネント(Button, Input, Checkbox等)から着手し、
  `docs/porting-pattern.md` のパターンを固める。
- Select / Popover / Tooltip などオーバーレイ系は、.zed/rules.mdのルール4(毎フレーム座標再計算)を
  明示的に確認しながら進める。
- Table / Editor / Chart のような複雑・大規模なものは、パターンが安定してから最後に着手する。
