# gogpui

`gogpui` は、Rust製 UIコンポーネントライブラリ [gpui-component](https://github.com/longbridge/gpui-component) を、`github.com/gogpu/gogpu` エコシステム (gogpu本体 + gogpu/gg) 上に Go で完全移植するプロジェクトです。

## 目標

- **完全移植**: 公開関数名・型名・メソッドシグネチャを Rust 版と 1:1 対応 (Goの命名規則の範囲内)。
- **挙動の一致**: hover / click / focus / disabled / open-close 等の状態遷移やスクロール追従などの振る舞いを Rust 版と一致させる。
- **デザインの一致**: shadcn/ui 由来のデザインや見た目を Rust 版と一致させる。

## 構成

- `/components` - 移植された UI コンポーネント群
- `/core` - コアとなるレイアウト・描画関連の基盤
- `/docs` - 移植パターンなどのドキュメント
- `/cmd/showcase` - 各コンポーネントの動作を確認できるデモ・ギャラリー

## 開発ガイドライン

開発やコントリビュートの際は、以下の点に注意してください。

1. **正典(Source of Truth)**: 
   実装の際は常に Rust 版ソースコード (`crates/ui/src`) と[ライブコンポーネントギャラリー](https://longbridge.github.io/gpui-component/)を正とし、推測での実装は避けてください。
2. **移植パターン**:
   `docs/porting-pattern.md` に記載されている確立済みの移植パターン (GPUI概念 → gogpu/gg概念の対応表) を参照してください。
3. **進捗管理**:
   コンポーネント単位の進捗や API 対応表は `PORTING_LEDGER.md` で管理されています。
4. **オーバーレイの実装**:
   gogpu/gg はイミディエイトモード GPU 描画です。ドロップダウンなどのオーバーレイの座標は、開いた瞬間のキャッシュではなく、アンカーとなるウィジェットの位置を毎フレーム再計算して描画してください。

## テストと検証

- 単体テスト (`go test ./...`) をパスすること。
- ゴールデンイメージ差分テストによる見た目の確認。
- 合成入力イベントによる hover / focus / click 等の状態遷移テスト。