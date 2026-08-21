# Button コンポーネント 実装・スキップ状況トラッキング

本ドキュメントは、Rust版 `gpui-component/crates/ui/src/button/button.rs` の Button を Go に移植する際、「スモールスタートとしてどこまで実装し、何を後回し（スキップ）にしたか」を記録するチェックリストです。

## Variants
**実装済み:**
- [x] `Default`
- [x] `Primary`
- [x] `Secondary`
- [x] `Destructive` (Rust版の `Danger` 相当)
- [x] `Ghost`
- [x] `Link`

**スキップ（必要に応じて後日追加）:**
- [ ] `Info`
- [ ] `Success`
- [ ] `Warning`
- [ ] `Text`
- [ ] `Custom` (独自背景色などのカスタム指定)

## Sizes
**実装済み:**
- [x] `Sm` (Small)
- [x] `Md` (Medium / Default)
- [x] `Lg` (Large)

**スキップ:**
- [ ] `XSmall`
- [ ] 任意のピクセルサイズ指定 (`Size(Pixels)`)

## Features & Methods
**実装済み:**
- [x] `Label()` (テキストラベル指定)
- [x] `Disabled()` (非活性状態)
- [x] `OnClick()` (クリックイベント)
- [x] `Hover` 状態と `Active` (押し込み) 状態のスタイリング連動

**スキップ:**
- [ ] `Icon()` (左側のアイコン配置) ※アイコン機構自体が未整備なため今回はスキップ
- [ ] `Tooltip()` (ホバー時のツールチップ表示) ※Tooltipコンポーネント実装時に結合
- [ ] `Loading()` (ローディングインジケーターの表示)
- [ ] `Outline()` (アウトラインスタイル)
- [ ] `Compact()` (パディングを減らす)
- [ ] `DropdownCaret()` (右端のキャレットアイコン)
- [ ] `BorderCorners()`, `BorderEdges()` (特定の角だけ角丸を外す等の細かい制御)
