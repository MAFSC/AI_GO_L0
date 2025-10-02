package onnx

import (
	"context"

	ort "github.com/microsoft/onnxruntime-go"
)

type Model struct{
	env *ort.Environment
	sess *ort.Session
}

func Load(modelPath string) (*Model, error) {
	env, err := ort.NewEnvironment()
	if err != nil { return nil, err }
	sess, err := env.NewSession(modelPath)
	if err != nil { return nil, err }
	return &Model{env: env, sess: sess}, nil
}

// Predict placeholder (wire proper tensors if you export a model)
func (m *Model) Predict(ctx context.Context, features []float32) (batch int, flushMs float32, err error) {
	return 128, 35.0, nil
}
