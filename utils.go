/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"encoding/binary"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"io"
	"math"
	"os"
	"runtime"
	"strings"
)

const (
	BasePath = "/you_path/res/"
)

func ReadData(file, path string) []byte {
	bs, err := os.ReadFile(BasePath + file)
	HandleErr(err)
	items := strings.Split(path, "/")
	for _, item := range items {
		bs = ParseMix(bs)[item]
	}
	return bs
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadU8(reader io.Reader) uint8 {
	return ReadBytes(reader, 1)[0]
}

func ReadU16(reader io.Reader) uint16 {
	res := ReadBytes(reader, 2)
	return binary.LittleEndian.Uint16(res)
}

func ReadU32(reader io.Reader) uint32 {
	res := ReadBytes(reader, 4)
	return binary.LittleEndian.Uint32(res)
}

func ReadF32(reader io.Reader) float32 {
	temp := ReadU32(reader)
	return math.Float32frombits(temp)
}

func ReadAny[T any](reader io.Reader) *T {
	res := new(T)
	err := binary.Read(reader, binary.LittleEndian, res)
	HandleErr(err)
	return res
}

func ReadBytes(reader io.Reader, count int) []byte {
	res := make([]byte, count)
	_, err := reader.Read(res)
	HandleErr(err)
	return res
}

func ReadStr(reader io.Reader) string {
	bs := make([]byte, 0)
	for {
		temp := ReadU8(reader)
		if temp == 0 {
			return string(bs)
		} else {
			bs = append(bs, temp)
		}
	}
}

type Texture struct {
	Texture uint32
}

func (t *Texture) Bind(texture uint32) {
	gl.ActiveTexture(texture)
	gl.BindTexture(gl.TEXTURE_2D, t.Texture)
}

func LoadTexture(rgba *image.RGBA) *Texture {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Dx()), int32(rgba.Rect.Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return &Texture{
		Texture: texture,
	}
}

func ShowImage(img *image.RGBA) {
	bound := img.Bounds()
	window := NewWindow(bound.Dx(), bound.Dy(), "Test")

	shader := LoadShader("res/img")
	shader.Use()
	shader.SetI1("Texture", 0)
	texture := LoadTexture(img)
	data := []float32{
		-1.0, 1.0, 0.0, 0.0, // 左上角
		-1.0, -1.0, 0.0, 1.0, // 左下角
		1.0, -1.0, 1.0, 1.0, // 右下角
		-1.0, 1.0, 0.0, 0.0, // 左上角
		1.0, -1.0, 1.0, 1.0, // 右下角
		1.0, 1.0, 1.0, 0.0, // 右上角
	}
	vao := NewVao(data, 2, 2)

	for !window.ShouldClose() {
		gl.ClearColor(0.1, 0.1, 0.1, 0.1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		shader.Use()
		texture.Bind(gl.TEXTURE0)
		vao.Bind()
		vao.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
	glfw.Terminate()
}

func buildPoint(x int, y int, z int, scale float32, clr *Color) []float32 {
	return []float32{float32(x) * scale, float32(y) * scale, float32(z) * scale, float32(clr.Red) / 255, float32(clr.Green) / 255, float32(clr.Blue) / 255}
}

func BuildMesh(vxl *Vxl, pal []*Color) []float32 { // 3 pos  3  clr
	items := vxl.Limbs[0].Items
	scale := vxl.LimbFooters[0].Scale
	res := make([]float32, 0)
	for x := 0; x < len(items); x++ {
		for y := 0; y < len(items[x]); y++ {
			for z := 0; z < len(items[x][y]); z++ {
				if items[x][y][z] == nil {
					continue
				}
				clr := pal[items[x][y][z].Color]
				// 上
				if z+1 >= len(items[x][y]) || items[x][y][z+1] == nil {
					res = append(res, buildPoint(x, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z+1, scale, clr)...)
					res = append(res, buildPoint(x, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x, y+1, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z+1, scale, clr)...)
				}
				// 下
				if z-1 < 0 || items[x][y][z-1] == nil {
					res = append(res, buildPoint(x, y, z, scale, clr)...)
					res = append(res, buildPoint(x+1, y, z, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z, scale, clr)...)
					res = append(res, buildPoint(x, y, z, scale, clr)...)
					res = append(res, buildPoint(x, y+1, z, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z, scale, clr)...)
				}
				// 左
				if y-1 < 0 || items[x][y-1][z] == nil {
					res = append(res, buildPoint(x, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y, z, scale, clr)...)
					res = append(res, buildPoint(x, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x, y, z, scale, clr)...)
					res = append(res, buildPoint(x+1, y, z, scale, clr)...)
				}
				// 右
				if y+1 >= len(items[x]) || items[x][y+1][z] == nil {
					res = append(res, buildPoint(x, y+1, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z, scale, clr)...)
					res = append(res, buildPoint(x, y+1, z+1, scale, clr)...)
					res = append(res, buildPoint(x, y+1, z, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z, scale, clr)...)
				}
				// 前
				if x-1 < 0 || items[x-1][y][z] == nil {
					res = append(res, buildPoint(x, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x, y+1, z+1, scale, clr)...)
					res = append(res, buildPoint(x, y+1, z, scale, clr)...)
					res = append(res, buildPoint(x, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x, y, z, scale, clr)...)
					res = append(res, buildPoint(x, y+1, z, scale, clr)...)
				}
				// 后
				if x+1 >= len(items) || items[x+1][y][z] == nil {
					res = append(res, buildPoint(x+1, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z, scale, clr)...)
					res = append(res, buildPoint(x+1, y, z+1, scale, clr)...)
					res = append(res, buildPoint(x+1, y, z, scale, clr)...)
					res = append(res, buildPoint(x+1, y+1, z, scale, clr)...)
				}
			}
		}
	}
	return res
}

func NewWindow(width, height int, title string) *glfw.Window {
	runtime.LockOSThread()
	// 初始化 glfw 辅助窗口
	err := glfw.Init()
	HandleErr(err)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	HandleErr(err)
	window.MakeContextCurrent()
	// 初始化 gl
	err = gl.Init()
	HandleErr(err)
	return window
}

func loadShader(path string, shaderType uint32) uint32 {
	bs, err := os.ReadFile(path)
	HandleErr(err)
	shader := gl.CreateShader(shaderType)
	cStr, free := gl.Strs(string(bs) + "\x00") // c 字符串需要这个结束标识
	gl.ShaderSource(shader, 1, cStr, nil)
	free()
	gl.CompileShader(shader)
	// 校验错误
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		panic(fmt.Sprintf("loadShader path %s shaderType %v err %s", path, shaderType, log))
	}
	return shader
}

const (
	VertName = ".vert"
	FragName = ".frag"
)

type Shader struct {
	Program    uint32
	UniformMap map[string]int32
}

func (s *Shader) Use() {
	gl.UseProgram(s.Program)
}

func (s *Shader) getUniformLoc(name string) int32 {
	if _, ok := s.UniformMap[name]; !ok {
		s.UniformMap[name] = gl.GetUniformLocation(s.Program, gl.Str(name+"\x00")) // c 字符串需要这个结束标识
	}
	return s.UniformMap[name]
}

func (s *Shader) SetMat4(name string, mat4 mgl32.Mat4) {
	uniformLoc := s.getUniformLoc(name)
	gl.UniformMatrix4fv(uniformLoc, 1, false, &mat4[0])
}

func (s *Shader) SetF4(name string, val mgl32.Vec4) {
	uniformLoc := s.getUniformLoc(name)
	gl.Uniform4f(uniformLoc, val[0], val[1], val[2], val[3])
}

func (s *Shader) SetF3(name string, val mgl32.Vec3) {
	uniformLoc := s.getUniformLoc(name)
	gl.Uniform3f(uniformLoc, val[0], val[1], val[2])
}

func (s *Shader) SetF1(name string, val float32) {
	uniformLoc := s.getUniformLoc(name)
	gl.Uniform1f(uniformLoc, val)
}

func (s *Shader) SetI1(name string, val int32) {
	uniformLoc := s.getUniformLoc(name)
	gl.Uniform1i(uniformLoc, val)
}

func LoadShader(name string) *Shader {
	// 顶点着色器与片源着色器一定是要有的
	vertShader := loadShader(name+VertName, gl.VERTEX_SHADER)
	fragShader := loadShader(name+FragName, gl.FRAGMENT_SHADER)
	// 链接着色器
	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)
	gl.DeleteShader(vertShader)
	gl.DeleteShader(fragShader)
	// 错误检查
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		panic(fmt.Sprintf("LoadShader name %s err %s", name, log))
	}
	return &Shader{
		Program:    program,
		UniformMap: make(map[string]int32),
	}
}

func PressKey(window *glfw.Window, key glfw.Key) bool {
	return window.GetKey(key) == glfw.Press
}

func GetAxis(window *glfw.Window, min, max glfw.Key) float32 {
	if PressKey(window, min) {
		return -1
	}
	if PressKey(window, max) {
		return 1
	}
	return 0
}

var (
	VecUp    = mgl32.Vec3{0, 1, 0}
	VecRight = mgl32.Vec3{1, 0, 0}
)

type Camera struct {
	Pos, Dir   mgl32.Vec3
	dirX, dirY float32 // 不能绕Z旋转
}

func NewCamera() *Camera {
	return &Camera{Pos: mgl32.Vec3{4.5405855, -1.2734333, 4.104122}, Dir: mgl32.Vec3{-0.39915586, 0.55325437, -0.73114574}.Normalize()}
}

func (c *Camera) GetView() mgl32.Mat4 {
	return mgl32.LookAtV(c.Pos, c.Pos.Add(c.Dir), VecUp)
}

func (c *Camera) TranslateX(value float32) {
	c.Pos = c.Pos.Add(c.Dir.Cross(VecUp).Normalize().Mul(value))
}

func (c *Camera) TranslateY(value float32) {
	c.Pos = c.Pos.Add(c.Dir.Cross(VecRight).Normalize().Mul(value))
}

func (c *Camera) TranslateZ(value float32) {
	c.Pos = c.Pos.Add(c.Dir.Normalize().Mul(value))
}

func (c *Camera) RotateX(value float32) { // 左右看 沿着 Y轴
	c.Dir = mgl32.Rotate3DY(value).Mul3x1(c.Dir)
}

func (c *Camera) RotateY(value float32) { // 上下看 沿着  X轴
	c.Dir = mgl32.Rotate3DX(value).Mul3x1(c.Dir)
}

func (c *Camera) Update(window *glfw.Window) {
	offsetX := GetAxis(window, glfw.KeyA, glfw.KeyD)
	if offsetX != 0 {
		c.TranslateX(offsetX * 0.1)
	}
	offsetY := GetAxis(window, glfw.KeyE, glfw.KeyQ)
	if offsetY != 0 {
		c.TranslateY(offsetY * 0.1)
	}
	offsetZ := GetAxis(window, glfw.KeyS, glfw.KeyW)
	if offsetZ != 0 {
		c.TranslateZ(offsetZ * 0.1)
	}
	rotateX := GetAxis(window, glfw.KeyRight, glfw.KeyLeft)
	if rotateX != 0 {
		c.RotateX(rotateX * 0.01)
	}
	rotateY := GetAxis(window, glfw.KeyDown, glfw.KeyUp)
	if rotateY != 0 {
		c.RotateY(rotateY * 0.01)
	}
	if PressKey(window, glfw.KeyEnter) {
		fmt.Println(c.Pos, c.Dir)
	}
}

type Vao struct {
	Vao        uint32
	Vbo        uint32
	IndicSize  int32
	PointCount int32
}

func (v *Vao) Bind() {
	gl.BindVertexArray(v.Vao)
}

func (v *Vao) Draw() {
	gl.DrawArrays(gl.TRIANGLES, 0, v.PointCount)
}

func NewVao(data []float32, sizes ...int32) *Vao {
	// 创建对象&写入数据
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.DYNAMIC_DRAW)
	// 设置数据
	if len(sizes) == 0 {
		panic("data format not set")
	}
	sum := int32(0)
	for _, size := range sizes {
		sum += size
	}
	curr := uintptr(0)
	for i := 0; i < len(sizes); i++ {
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribPointerWithOffset(uint32(i), sizes[i], gl.FLOAT, false, sum*4, curr*4)
		curr += uintptr(sizes[i])
	}
	return &Vao{
		Vao:        vao,
		Vbo:        vbo,
		PointCount: int32(len(data)) / sum,
	}
}
