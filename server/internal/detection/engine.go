package detection

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"

	"github.com/dnlmgwi/football-object-detection/internal/config"
	"github.com/dnlmgwi/football-object-detection/internal/models"
)

// Engine handles object detection using YOLO
type Engine struct {
	net              gocv.Net
	config           *config.Manager
	inputSize        int
	confidenceThresh float32
	nmsThresh        float32
	classes          []string
}

// NewEngine creates a new detection engine
func NewEngine(modelPath string, cfg *config.Manager) (*Engine, error) {
	// Load YOLO model
	net := gocv.ReadNet(modelPath, "")
	if net.Empty() {
		return nil, fmt.Errorf("failed to load model from %s", modelPath)
	}

	// Set backend and target (prefer CUDA if available)
	net.SetPreferableBackend(gocv.NetBackendDefault)
	net.SetPreferableTarget(gocv.NetTargetCPU)

	detCfg := cfg.GetConfig().Detection

	// COCO dataset classes (we care about person=0 and sports ball=32)
	classes := []string{
		"person", "bicycle", "car", "motorcycle", "airplane", "bus", "train", "truck", "boat",
		"traffic light", "fire hydrant", "stop sign", "parking meter", "bench", "bird", "cat",
		"dog", "horse", "sheep", "cow", "elephant", "bear", "zebra", "giraffe", "backpack",
		"umbrella", "handbag", "tie", "suitcase", "frisbee", "skis", "snowboard", "sports ball",
	}

	return &Engine{
		net:              net,
		config:           cfg,
		inputSize:        detCfg.InputSize,
		confidenceThresh: detCfg.ConfidenceThreshold,
		nmsThresh:        detCfg.NMSThreshold,
		classes:          classes,
	}, nil
}

// Detect performs object detection on a frame
func (e *Engine) Detect(frame gocv.Mat) ([]models.Player, *models.Ball, error) {
	if frame.Empty() {
		return nil, nil, fmt.Errorf("empty frame")
	}

	// Create blob from image
	blob := gocv.BlobFromImage(frame, 1/255.0, image.Pt(e.inputSize, e.inputSize),
		gocv.NewScalar(0, 0, 0, 0), true, false)
	defer blob.Close()

	// Set input
	e.net.SetInput(blob, "")

	// Forward pass
	prob := e.net.Forward("")
	defer prob.Close()

	// Parse detections
	players, ball := e.parseDetections(prob, frame.Cols(), frame.Rows())

	// Classify players by team
	cfg := e.config.GetConfig()
	players = e.classifyPlayers(players, frame, cfg.TeamA, cfg.TeamB, cfg.Possession.ColorTolerance)

	return players, ball, nil
}

// parseDetections parses YOLO output
func (e *Engine) parseDetections(prob gocv.Mat, frameWidth, frameHeight int) ([]models.Player, *models.Ball) {
	var players []models.Player
	var ball *models.Ball

	// YOLO output format: [batch, num_detections, 85]
	// 85 = [x, y, w, h, objectness, class1, class2, ..., class80]
	data := prob.ToBytes()
	numDetections := prob.Size()[1]

	var personBoxes []image.Rectangle
	var personConfidences []float32
	var ballBoxes []image.Rectangle
	var ballConfidences []float32

	for i := 0; i < numDetections; i++ {
		offset := i * 85 * 4 // 4 bytes per float32

		// Parse bounding box (normalized coordinates)
		cx := bytesToFloat32(data[offset : offset+4])
		cy := bytesToFloat32(data[offset+4 : offset+8])
		w := bytesToFloat32(data[offset+8 : offset+12])
		h := bytesToFloat32(data[offset+12 : offset+16])
		objectness := bytesToFloat32(data[offset+16 : offset+20])

		if objectness < e.confidenceThresh {
			continue
		}

		// Find best class
		maxConf := float32(0.0)
		classID := -1
		for j := 0; j < 80; j++ {
			classConf := bytesToFloat32(data[offset+20+j*4 : offset+24+j*4])
			if classConf > maxConf {
				maxConf = classConf
				classID = j
			}
		}

		confidence := objectness * maxConf
		if confidence < e.confidenceThresh {
			continue
		}

		// Convert to pixel coordinates
		x := int((cx - w/2) * float32(frameWidth))
		y := int((cy - h/2) * float32(frameHeight))
		width := int(w * float32(frameWidth))
		height := int(h * float32(frameHeight))

		box := image.Rect(x, y, x+width, y+height)

		// Person (class 0)
		if classID == 0 {
			personBoxes = append(personBoxes, box)
			personConfidences = append(personConfidences, confidence)
		}

		// Sports ball (class 32)
		if classID == 32 {
			ballBoxes = append(ballBoxes, box)
			ballConfidences = append(ballConfidences, confidence)
		}
	}

	// Apply Non-Maximum Suppression
	personIndices := performNMS(personBoxes, personConfidences, e.nmsThresh)
	for _, idx := range personIndices {
		box := personBoxes[idx]
		players = append(players, models.Player{
			ID:         fmt.Sprintf("p%d", len(players)),
			BBox:       models.BoundingBox{X: box.Min.X, Y: box.Min.Y, Width: box.Dx(), Height: box.Dy()},
			Team:       "neutral",
			Confidence: personConfidences[idx],
		})
	}

	ballIndices := performNMS(ballBoxes, ballConfidences, e.nmsThresh)
	if len(ballIndices) > 0 {
		idx := ballIndices[0] // Take highest confidence ball
		box := ballBoxes[idx]
		ball = &models.Ball{
			BBox:       models.BoundingBox{X: box.Min.X, Y: box.Min.Y, Width: box.Dx(), Height: box.Dy()},
			Confidence: ballConfidences[idx],
		}
	}

	return players, ball
}

// classifyPlayers classifies players into teams based on jersey color
func (e *Engine) classifyPlayers(players []models.Player, frame gocv.Mat,
	teamA, teamB config.Team, tolerance float64) []models.Player {

	for i := range players {
		// Extract player region
		bbox := players[i].BBox
		roi := frame.Region(bbox.ToRect())
		if roi.Empty() {
			continue
		}

		// Get dominant color
		dominantColor := e.getDominantColor(roi)
		roi.Close()

		// Match against team colors
		matchA := e.colorDistance(dominantColor, teamA.PrimaryColor)
		matchB := e.colorDistance(dominantColor, teamB.PrimaryColor)

		if matchA < matchB && matchA < tolerance {
			players[i].Team = teamA.ID
			players[i].ColorMatch = 1.0 - (matchA / tolerance)
		} else if matchB < matchA && matchB < tolerance {
			players[i].Team = teamB.ID
			players[i].ColorMatch = 1.0 - (matchB / tolerance)
		}
		// Otherwise remains "neutral"
	}

	return players
}

// getDominantColor extracts dominant color from image region (simplified)
func (e *Engine) getDominantColor(roi gocv.Mat) config.Color {
	// Calculate mean BGR color from ROI
	// For better accuracy, could convert to HSV first, but BGR mean works well
	meanBGR := roi.Mean()

	return config.Color{
		R: uint8(meanBGR.Val3),
		G: uint8(meanBGR.Val2),
		B: uint8(meanBGR.Val1),
	}
}

// colorDistance calculates Euclidean distance between two colors
func (e *Engine) colorDistance(c1, c2 config.Color) float64 {
	dr := float64(c1.R) - float64(c2.R)
	dg := float64(c1.G) - float64(c2.G)
	db := float64(c1.B) - float64(c2.B)
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// Close releases resources
func (e *Engine) Close() error {
	if !e.net.Empty() {
		e.net.Close()
	}
	return nil
}

// Helper functions

func bytesToFloat32(b []byte) float32 {
	if len(b) < 4 {
		return 0
	}
	bits := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	return math.Float32frombits(bits)
}

// performNMS applies Non-Maximum Suppression
func performNMS(boxes []image.Rectangle, confidences []float32, threshold float32) []int {
	if len(boxes) == 0 {
		return nil
	}

	// Convert to GoCV format
	boxRects := make([]image.Rectangle, len(boxes))
	copy(boxRects, boxes)

	// Use GoCV's NMS implementation
	indices := gocv.NMSBoxes(boxRects, confidences, 0.0, threshold)

	return indices
}

// DrawDetections draws bounding boxes on frame
func DrawDetections(frame *gocv.Mat, players []models.Player, ball *models.Ball) {
	// Draw players
	for _, player := range players {
		var col color.RGBA
		switch player.Team {
		case "teamA":
			col = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red
		case "teamB":
			col = color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue
		default:
			col = color.RGBA{R: 128, G: 128, B: 128, A: 255} // Gray
		}

		rect := player.BBox.ToRect()
		gocv.Rectangle(frame, rect, col, 2)

		label := fmt.Sprintf("%s %.2f", player.Team, player.Confidence)
		gocv.PutText(frame, label, image.Pt(rect.Min.X, rect.Min.Y-5),
			gocv.FontHersheyPlain, 1.0, col, 1)
	}

	// Draw ball
	if ball != nil {
		rect := ball.BBox.ToRect()
		gocv.Circle(frame, image.Pt(rect.Min.X+rect.Dx()/2, rect.Min.Y+rect.Dy()/2),
			rect.Dx()/2, color.RGBA{R: 0, G: 255, B: 0, A: 255}, 2)
	}
}
