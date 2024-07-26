package gameserver

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"server/config"
	"server/types"
	"strings"
)

type Object struct {
	uid                      int32
	name                     string
	kind                     types.ObjectKind
	meta                     string
	waypoints                [][3]float64
	resourcePath             string
	position                 []float32
	rotation                 []float32
	scale                    []float32
	isRigidbodySleepingStart bool
	isWithColorProperty      bool
	colorProperty            []float32
	variationIndex           int32
}

type LevelData struct {
	Version      uint8
	Kind         int32
	TerrainCount int32
	TerrainData  TerrainData
	ObjectsCount int32
	Objects      []Object
	Teleports    []LevelTeleport
}

type LevelTeleport struct {
	Name     string
	Target   string
	Position types.Vector3
	Rotation types.Vector3
}

type TerrainData struct {
	heightmapResolution    int32
	size                   []float32
	heightmapWidth         int32
	heightmapHeight        int32
	alphamapResolution     int32
	alphamapWidth          int32
	alphamapHeight         int32
	alphamapLayers         int32
	alphamapTextureIndices []int32
	heightmap              [][]float32
	alphamaps              [][][]float32
}

func LoadLevel() (*LevelData, error) {
	level, _ := loadLevel(config.WorldFilePath)

	byteArray := level[0]
	fmt.Println("Loading level data ...")

	gameData, err := loadLevelDataFromByteArray(byteArray)
	if err != nil {
		fmt.Printf("Error processing data: %v\n", err)
		return nil, err
	}

	fmt.Print("==============================\n")
	fmt.Printf("| Level size: %dx%d\n", int(gameData.TerrainData.size[0]), int(gameData.TerrainData.size[2]))
	fmt.Printf("| Objects count: %d\n", gameData.ObjectsCount)
	fmt.Printf("| Teleports count: %d\n", len(gameData.Teleports))
	fmt.Print("==============================\n")

	return gameData, nil
}

func loadLevel(filename string) ([][]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("LoadFromStr: error reading file or file is empty. Error: %v", err)
	}

	savedLevelStr := string(data)
	dataAsStringArray := strings.Split(savedLevelStr, "#")
	if len(dataAsStringArray) != 2 {
		return nil, fmt.Errorf("LoadFromStr: incorrect number of Base64 strings")
	}

	byteArrays := make([][]byte, 2)
	for i, str := range dataAsStringArray {
		decodedBytes, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, fmt.Errorf("LoadFromStr: base64 decoding error. Error: %v", err)
		}
		byteArrays[i] = decodedBytes
	}

	return byteArrays, nil
}

func (o *Object) isNPC() bool {
	return o.kind == types.ObjectKindNPC
}

func getColorProperty(steam *bytes.Reader) []float32 {
	var r, g, b float32

	binary.Read(steam, binary.LittleEndian, &r)
	binary.Read(steam, binary.LittleEndian, &g)
	binary.Read(steam, binary.LittleEndian, &b)

	return []float32{r, g, b}
}

func getVector(steam *bytes.Reader) []float32 {
	var x, y, z float32

	binary.Read(steam, binary.LittleEndian, &x)
	binary.Read(steam, binary.LittleEndian, &y)
	binary.Read(steam, binary.LittleEndian, &z)

	return []float32{x, y, z}
}

func getInt32(steam *bytes.Reader) int32 {
	var value int32
	binary.Read(steam, binary.LittleEndian, &value)
	return value
}

func getFloat(steam *bytes.Reader) float32 {
	var value float32
	binary.Read(steam, binary.LittleEndian, &value)
	return value
}

func getQuaternion(steam *bytes.Reader) []float32 {
	var x, y, z, w float32

	binary.Read(steam, binary.LittleEndian, &x)
	binary.Read(steam, binary.LittleEndian, &y)
	binary.Read(steam, binary.LittleEndian, &z)
	binary.Read(steam, binary.LittleEndian, &w)

	return []float32{x, y, z, w}

}

func getString(steam *bytes.Reader) string {
	length, _ := binary.ReadUvarint(steam)
	resourcePath := make([]byte, length)

	if _, err := steam.Read(resourcePath); err != nil {
		panic("Cant read getString() resource path")
	}

	return string(resourcePath)
}

func getBool(steam *bytes.Reader) bool {
	var value bool
	binary.Read(steam, binary.LittleEndian, &value)
	return value

}

func getTerrainData(stream *bytes.Reader) TerrainData {
	heightmapResolution := getInt32(stream)
	size := getVector(stream)
	heightmapWidth := getInt32(stream)
	heightmapHeight := getInt32(stream)

	heightmap := make([][]float32, heightmapWidth)
	for x := 0; x < int(heightmapWidth); x++ {
		heightmap[x] = make([]float32, heightmapHeight)
	}

	for x := 0; x < int(heightmapWidth); x++ {
		for y := 0; y < int(heightmapHeight); y++ {
			heightmap[x][y] = getFloat(stream)
		}
	}

	alphamapResolution := getInt32(stream)
	alphamapWidth := getInt32(stream)
	alphamapHeight := getInt32(stream)
	alphamapLayers := getInt32(stream)

	alphamaps := getAlphamaps(stream, alphamapWidth, alphamapHeight, alphamapLayers)

	return TerrainData{
		heightmapResolution: heightmapResolution,
		size:                size,
		heightmapWidth:      heightmapWidth,
		heightmapHeight:     heightmapHeight,

		alphamapResolution: alphamapResolution,
		alphamapWidth:      alphamapWidth,
		alphamapHeight:     alphamapHeight,
		alphamapLayers:     alphamapLayers,
		// alphamapTextureIndices: alphamapTextureIndices,
		heightmap: heightmap,
		alphamaps: alphamaps,
	}
}

func getAlphamaps(stream *bytes.Reader, alphamapWidth, alphamapHeight, alphamapLayers int32) [][][]float32 {

	alphamapTextureIndices := make([]int32, alphamapLayers)
	for i := 0; i < int(alphamapLayers); i++ {
		alphamapTextureIndices[i] = getInt32(stream)
	}

	alphamaps := make([][][]float32, alphamapWidth)
	for x := int32(0); x < alphamapWidth; x++ {
		alphamaps[x] = make([][]float32, alphamapHeight)
		for y := int32(0); y < alphamapHeight; y++ {
			alphamaps[x][y] = make([]float32, alphamapLayers)
		}
	}

	for x := int32(0); x < alphamapWidth; x++ {
		for y := int32(0); y < alphamapHeight; y++ {
			for z := int32(0); z < alphamapLayers; z++ {
				alphamaps[x][y][z] = getFloat(stream)
			}
		}
	}

	return alphamaps
}

func getObjects(stream *bytes.Reader, count int32) []Object {
	objects := make([]Object, count)

	for i := 0; i < int(count); i++ {
		var colorProperty []float32

		uid := getInt32(stream)
		resourcePath := getString(stream)
		position := getVector(stream)
		rotation := getQuaternion(stream)
		scale := getVector(stream)
		isRigidbodySleepingStart := getBool(stream)
		isWithColorProperty := getBool(stream)
		variationIndex := getInt32(stream)

		name := getString(stream)
		meta := getString(stream)
		kind := getInt32(stream)

		if isWithColorProperty {
			colorProperty = getColorProperty(stream)
		}

		waypoints := make([][3]float64, 0)
		if kind == types.ObjectKindNPC {
			waypointsCount := getInt32(stream)

			for i := 0; i < int(waypointsCount); i++ {
				waypoint := getVector(stream)
				waypoints = append(waypoints, [3]float64{float64(waypoint[0]), float64(waypoint[1]), float64(waypoint[2])})
			}
		}

		objects[i] = Object{
			uid:                      uid,
			name:                     name,
			kind:                     types.ObjectKind(kind),
			meta:                     meta,
			waypoints:                waypoints,
			resourcePath:             resourcePath,
			position:                 position,
			rotation:                 rotation,
			scale:                    scale,
			isRigidbodySleepingStart: isRigidbodySleepingStart,
			isWithColorProperty:      isWithColorProperty,
			colorProperty:            colorProperty,
			variationIndex:           variationIndex,
		}

	}

	return objects
}

func loadLevelDataFromByteArray(byteArray []byte) (*LevelData, error) {
	buf := bytes.NewReader(byteArray)

	version, _ := buf.ReadByte()
	worldKind := getInt32(buf)

	terrainCount := getInt32(buf)
	terrainData := TerrainData{}

	if terrainCount > 0 {
		terrainData = getTerrainData(buf)
	}

	objectsCount := getInt32(buf)
	objects := getObjects(buf, objectsCount)

	teleports := make([]LevelTeleport, 0)

	for _, obj := range objects {
		if obj.kind != types.ObjectKindTeleport {
			continue
		}

		teleports = append(teleports, LevelTeleport{
			Name:     obj.name,
			Target:   obj.meta,
			Position: types.Vector3{X: float64(obj.position[0]), Y: float64(obj.position[1]) + 0.6, Z: float64(obj.position[2])},
			Rotation: types.Vector3{X: float64(obj.rotation[0]), Y: float64(obj.rotation[1]), Z: float64(obj.rotation[2])},
		})

	}

	gameData := &LevelData{
		Version:      version,
		Kind:         worldKind,
		TerrainCount: terrainCount,
		TerrainData:  terrainData,
		ObjectsCount: objectsCount,
		Objects:      objects,
		Teleports:    teleports,
	}

	return gameData, nil
}
