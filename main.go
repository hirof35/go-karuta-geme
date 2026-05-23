package main

import (
	_ "embed" // ファイルを埋め込むための標準機能
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	screenWidth  = 1000
	screenHeight = 700
	cardWidth    = 50
	cardHeight   = 70
)

// 【ここを修正しました！】staticフォルダ内の NotoSansJP-Regular.ttf をプログラムに直接埋め込みます
//go:embed static/NotoSansJP-Regular.ttf
var defaultFontData []byte

// カルタのデータ構造
type Card struct {
	Yomifuda string  // 読み札（テキスト）
	Torifuda string  // 取り札（ひらがな1文字）
	X, Y     float64 // 画面上の配置座標
	Taken    bool    // すでに取られたかどうか
}

type Game struct {
	cards        []Card
	currentIdx   int       // 現在の出題インデックス
	questionFont *text.GoTextFace // 読み札用のフォント
	cardFont     *text.GoTextFace // 取り札用のフォント
	gameCleared  bool
	startTime    time.Time
	clearTime    float64
}

func (g *Game) Update() error {
	if g.gameCleared {
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		cx, cy := float64(mx), float64(my)

		correctCard := g.cards[g.currentIdx]

		if cx >= correctCard.X && cx <= correctCard.X+cardWidth &&
			cy >= correctCard.Y && cy <= correctCard.Y+cardHeight && !correctCard.Taken {
			
			for i := range g.cards {
				if g.cards[i].Torifuda == correctCard.Torifuda {
					g.cards[i].Taken = true
					break
				}
			}

			g.currentIdx++
			if g.currentIdx >= len(g.cards) {
				g.gameCleared = true
				g.clearTime = time.Since(g.startTime).Seconds()
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{34, 139, 34, 255}) // 緑色の畳背景

	if g.gameCleared {
		op := &text.DrawOptions{}
		op.GeoM.Translate(screenWidth/2-150, screenHeight/2-30)
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, fmt.Sprintf("全札獲得！\n\nタイム: %.2f 秒", g.clearTime), g.questionFont, op)
		return
	}

	if g.currentIdx < len(g.cards) {
		currentCard := g.cards[g.currentIdx]
		
		op := &text.DrawOptions{}
		op.GeoM.Translate(50, 50)
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, "【読み札】 "+currentCard.Yomifuda, g.questionFont, op)
		
		opCount := &text.DrawOptions{}
		opCount.GeoM.Translate(800, 50)
		opCount.ColorScale.ScaleWithColor(color.RGBA{255, 215, 0, 255})
		text.Draw(screen, fmt.Sprintf("残り: %d 枚", len(g.cards)-g.currentIdx), g.cardFont, opCount)
	}

	for _, card := range g.cards {
		if card.Taken {
			continue
		}

		cardRect := ebiten.NewImage(cardWidth, cardHeight)
		cardRect.Fill(color.White)
		
		opCard := &ebiten.DrawImageOptions{}
		opCard.GeoM.Translate(card.X, card.Y)
		screen.DrawImage(cardRect, opCard)

		borderColor := color.RGBA{50, 50, 50, 255}
		for w := 0; w < cardWidth; w++ {
			screen.Set(int(card.X)+w, int(card.Y), borderColor)
			screen.Set(int(card.X)+w, int(card.Y)+cardHeight-1, borderColor)
		}
		for h := 0; h < cardHeight; h++ {
			screen.Set(int(card.X), int(card.Y)+h, borderColor)
			screen.Set(int(card.X)+cardWidth-1, int(card.Y)+h, borderColor)
		}

		opText := &text.DrawOptions{}
		opText.GeoM.Translate(card.X+13, card.Y+15)
		opText.ColorScale.ScaleWithColor(color.Black)
		text.Draw(screen, card.Torifuda, g.cardFont, opText)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	log.Println("埋め込みフォントを読み込み中...")

	// 解凍した「NotoSansJP-Regular.ttf」からフォントソースを作成します
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(defaultFontData))
	if err != nil {
		log.Fatalf("フォントソースの作成に失敗しました: %v", err)
	}

	qFont := &text.GoTextFace{Source: fontSource, Size: 32}
	cFont := &text.GoTextFace{Source: fontSource, Size: 24}

	iroha := []struct{ yomi, tori string }{
		{"犬も歩けば棒に当たる", "い"}, {"論より証拠", "ろ"}, {"花より団子", "は"}, {"憎まれっ子世に憚る", "に"},
		{"骨折り損のくたびれ儲け", "ほ"}, {"へいたの（下手の）長談義", "へ"}, {"ちりも積もれば山となる", "ち"}, {"散りぬるを（ちりぬるを）", "り"},
		{"ちり（ぬるを）の次は「ぬ」", "ぬ"}, {"るり（瑠璃）も玻璃も照らせば光る", "る"}, {"をか（おかし）なこと", "を"}, {"わが衣手に雪は降りつつ", "わ"},
		{"カエルの面に水", "か"}, {"よのなか（世の中）は", "よ"}, {"たれ（誰）そ彼", "た"}, {"れん（連）れん", "れ"},
		{"そ（袖）ひつて", "そ"}, {"つ（月）見れば", "つ"}, {"ね（念）には念を入れよ", "ね"}, {"な（泣）き面に蜂", "な"},
		{"ら（楽）あれば苦あり", "ら"}, {"む（無理）が通れば道理引っ込む", "む"}, {"う（嘘）から出たまこと", "う"}, {"ゐ（井）の中の蛙", "ゐ"},
		{"の（喉）元過ぎれば熱さを忘れる", "の"}, {"お（終わり）良ければすべて良し", "お"}, {"く（聞）くは一時の恥", "く"}, {"や（病）は気から", "や"},
		{"ま（負）けるが勝ち", "ま"}, {"け（怪）我の功名", "け"}, {"ふ（百）聞は一見にしかず", "ふ"}, {"こ（転）ばぬ先の杖", "こ"},
		{"え（縁）は異なもの味なもの", "え"}, {"て（亭）主の好きな赤烏帽子", "て"}, {"あ（頭）隠して尻隠さず", "あ"}, {"さ（三）人寄れば文殊の知恵", "さ"},
		{"き（聞）いて極楽見て地獄", "き"}, {"ゆ（油）断大敵", "ゆ"}, {"め（目）の下の瘤", "め"}, {"み（身）から出た錆", "み"},
		{"し（知）らぬが仏", "し"}, {"ゑ（縁）に連れ", "ゑ"}, {"ひ（瓢）たんから駒", "ひ"}, {"も（餅）は餅屋", "も"},
		{"せ（背）に腹は変えられぬ", "せ"}, {"す（雀）百まで踊り忘れず", "す"},
	}

	var cards []Card
	for _, v := range iroha {
		cards = append(cards, Card{
			Yomifuda: v.yomi,
			Torifuda: v.tori,
		})
	}

	startX, startY := 60.0, 200.0
	spacingX, spacingY := 58.0, 90.0

	var positions []struct{ x, y float64 }
	for row := 0; row < 3; row++ {
		for col := 0; col < 16; col++ {
			if row*16+col >= len(cards) {
				break
			}
			x := startX + float64(col)*spacingX
			y := startY + float64(row)*spacingY
			positions = append(positions, struct{ x, y float64 }{x, y})
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	for i := range cards {
		cards[i].X = positions[i].x
		cards[i].Y = positions[i].y
	}

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	game := &Game{
		cards:        cards,
		questionFont: qFont,
		cardFont:     cFont,
		startTime:    time.Now(),
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("本格派！Go言語いろはかるた")
	
	log.Println("ゲームを起動します。")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}