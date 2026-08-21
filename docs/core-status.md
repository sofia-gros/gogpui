# Core層 実装・スキップ状況トラッキング

本ドキュメントは、Rust版 `gpui-component/crates/ui/src` の Core層（Theme, Sizing, Styling など全般）を Go に移植する際、「スモールスタートとしてどこまで実装し、何を後回し（スキップ）にしたか」を記録するチェックリストです。

## Theme (`core/theme`)

### 構造体・機能
- [x] `Theme` 構造体 (基本構成のみ)
- [x] `ThemeMode` (Light / Dark)
- [x] `ThemeColor` (カラーパレット保持)
- [x] `ThemeFrom(ctx)` (gg.Contextからのテーマ取得)
- [ ] `registry.rs` にあるような複雑な動的テーマ登録・切り替えシステム (初期は固定のLight/Darkのみに簡略化)
- [ ] Typography Tokens (フォントサイズの動的計算システム等)
- [ ] Shadow Tokens (Shadowの詳細な定義、今回は必要最小限に留める)

### カラーパレット (`theme/color.rs` 等)
shadcn/ui 準拠の数十色あるカラーのうち、第一段階で移植したものとスキップしたもの。

**実装済み:**
- [x] `Background`, `Foreground`
- [x] `Primary`, `PrimaryForeground`
- [x] `Secondary`, `SecondaryForeground`
- [x] `Muted`, `MutedForeground`
- [x] `Border`, `Input`, `Ring`
- [x] `Destructive`, `DestructiveForeground`

**スキップ（必要に応じて後日追加）:**
- [ ] `Popover`, `PopoverForeground` (後日Overlay系コンポーネント実装時に追加)
- [ ] `Accent`, `AccentForeground`
- [ ] `Card`, `CardForeground`
- [ ] 各種ステータスカラー（Success, Warning, Info 等の詳細バリエーション）
- [ ] スクロールバー専用色 (`Scrollbar`, `ScrollbarThumb` 等)
- [ ] 各種ホバー状態の明示的カラー (ThemeColor側で持つか、描画時に半透明化するか要検討)

## Size (`core/size`)

- [x] `sizing.rs` の定数定義 (Rem 換算の `sm`, `md`, `lg` など) -> Goではピクセル等で `SizeSm`, `SizeMd` などの定数として定義。
- [ ] 複雑な `SizeTokens` 構造体としての受け渡し（最初は単純なパッケージ定数として運用）

## Style (`core/style`)

- [x] `styled.rs` 相当の機能（ggのDrawコールに対する角丸(Radius)やBorderの適用ヘルパー）
- [ ] `ElementExt` などの複雑なTrait拡張（gogpuのイミディエイトモードにはそぐわないため、ヘルパー関数として再構築）
