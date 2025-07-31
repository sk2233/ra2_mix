package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"testing"
)

func TestVxl(t *testing.T) {
	window := NewWindow(1280, 720, "Test")

	shader := LoadShader("res/test")
	shader.Use()
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(16)/9, 0.1, 30.0)
	shader.SetMat4("Projection", projection)
	shader.SetMat4("Model", mgl32.Ident4())

	data := ReadData("ra2.mix", "local.mix/harv.vxl")
	vxl := ParseVxl(data)
	data = ReadData("ra2.mix", "cache.mix/uniturb.pal")
	pal := ParsePal(data)
	mesh := BuildMesh(vxl, pal)
	vao := NewVao(mesh, 3, 3)

	camera := NewCamera()
	gl.Enable(gl.DEPTH_TEST)
	for !window.ShouldClose() {
		gl.ClearColor(0.1, 0.1, 0.1, 0.1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		camera.Update(window)

		shader.Use()
		shader.SetMat4("View", camera.GetView())
		vao.Bind()
		vao.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
	glfw.Terminate()
}
