package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	// Создаем новый HTTP-запрос.
	req := httptest.NewRequest("GET", "/testpath", nil)
	// Создаем ResponseRecorder для записи ответа.
	w := httptest.NewRecorder()

	// Вызываем хендлер.
	handler(w, req)

	// Получаем результат.
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	// Проверяем статус-код.
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус-код %d, получен %d", http.StatusOK, resp.StatusCode)
	}

	// Проверяем тело ответа.
	expectedBody := "Hi there, I love testpath!"
	if string(body) != expectedBody {
		t.Errorf("Ожидалось тело ответа '%s', получено '%s'", expectedBody, string(body))
	}
}
