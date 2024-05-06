package usecase

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type FaceRecognitionService struct {
	ExternalService
}

func (f *FaceRecognitionService) Execute(strategy Strategy, des any) error {
	return strategy.Execute(des)
}

type EncodeFace struct {
	FormData    *bytes.Buffer
	ContentType string
}

type VerifyFace struct {
	FormData    *bytes.Buffer
	ContentType string
}

func (e *EncodeFace) Execute(des any) error {
	// Send request to FaceRecognition to encode the face
	err := SendRequestWithFormData("https://face-reg-service-latest.onrender.com/img2encoding", e.FormData, e.ContentType, des)
	if err != nil {
		return err
	}
	return nil
}

func (v *VerifyFace) Execute(des any) error {
	// Send request to FaceRecognition to verify the face
	err := SendRequestWithFormData("https://face-reg-service-latest.onrender.com/verify", v.FormData, v.ContentType, des)
	if err != nil {
		return err
	}
	return nil
}

func SendRequestWithFormData(url string, formData *bytes.Buffer, ContentType string, des any) error {

	client := &http.Client{}
	// Send request to FaceRecognition
	req, err := http.NewRequest(http.MethodPost, url, formData)
	if err != nil {
		return err
	}
	// Set the content type
	req.Header.Set("Content-Type", ContentType)
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Get the string from the body in the key "face_encoding"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if des != nil {
		err = json.Unmarshal(body, &des)
		if err != nil {
			return err
		}
	}

	return nil
}
