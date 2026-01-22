package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// TODO: Реализовать структуру EventLogger
// Должна содержать: канал для событий, канал завершения, файл для записи
type EventLogger struct {
	// TODO: добавить поля структуры
}

// TODO: Реализовать NewEventLogger()
// Открыть файл для логирования (создать если не существует, дописывать в конец)
// Инициализировать каналы (размер 100 для eventsChan)
// Вернуть инициализированный EventLogger
func NewEventLogger(filePath string) *EventLogger {
	// TODO: реализовать
	return nil
}

// TODO: Реализовать LogEvent()
// Отправить событие в канал eventsChan (неблокирующе)
// Если канал переполнен - залогировать warning
func (el *EventLogger) LogEvent(event string) {
	// TODO: реализовать
}

// TODO: Реализовать Start()
// Запустить worker в отдельной горутине
func (el *EventLogger) Start() {
	// TODO: реализовать
}

// TODO: Реализовать worker()
// Читать события из канала, добавлять timestamp, писать в файл
// Когда канал закроется - закрыть файл и сигнализировать done канал
func (el *EventLogger) worker() {
	// TODO: реализовать
}

// TODO: Реализовать Stop()
// Закрыть канал eventsChan, дождаться завершения worker, залогировать сообщение
func (el *EventLogger) Stop() {
	// TODO: реализовать
}
