package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Управляет асинхронной записью бизнес-событий в файл через каналы
type EventLogger struct {
	eventsChan chan string   // Канал для передачи строк-событий воркеру
	doneChan   chan struct{} // Канал для сигнализации о полном завершении работы воркера
	file       *os.File      // Дескриптор файла логов (log.txt)
}

// Открывает файл логов и инициализирует каналы
func NewEventLogger(filePath string) *EventLogger {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Fatal: failed to open or create log file %s: %v", filePath, err)
	}

	return &EventLogger{
		eventsChan: make(chan string, 100), // Буферизованный канал на 100 элементов
		doneChan:   make(chan struct{}),
		file:       file,
	}
}

// Отправляет событие в канал в неблокирующем режиме
func (el *EventLogger) LogEvent(event string) {
	select {
	case el.eventsChan <- event:
		// Событие успешно отправлено в буфер канала
	default:
		// Если буфер канала (100 элементов) переполнен, пишем предупреждение в стандартный лог,
		// чтобы не блокировать основной поток выполнения HTTP-запроса пользователя
		log.Printf("Warning: EventLogger buffer is full. Dropped log: %s", event)
	}
}

// Запускает воркера в фоновой горутине
func (el *EventLogger) Start() {
	go el.worker()
	log.Println("EventLogger async worker started successfully")
}

// Считывает сообщения из канала, имитирует задержку и записывает данные в файл
func (el *EventLogger) worker() {
	for event := range el.eventsChan {
		time.Sleep(1500 * time.Millisecond)

		// Форматируем строку лога: добавляем временную метку (Timestamp)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		logLine := fmt.Sprintf("[%s] %s\n", timestamp, event)

		// Записываем строку в файл
		if _, err := el.file.WriteString(logLine); err != nil {
			log.Printf("Error: failed to write to log file: %v", err)
		}

		// Для надежности дублируем запись в консоль
		log.Printf("[Async Log] %s", event)
	}

	// Код ниже выполняется только ПОСЛЕ закрытия канала и обработки всех накопленных логов
	log.Println("EventLogger worker is closing the log file...")
	if err := el.file.Close(); err != nil {
		log.Printf("Error: failed to close log file cleanly: %v", err)
	}

	// Сигнализируем методу Stop(), что воркер полностью завершил очистку буфера и закрыл файл
	close(el.doneChan)
}

// Корректно останавливает логгер, гарантируя запись всех оставшихся в канале сообщений
func (el *EventLogger) Stop() {
	log.Println("Stopping EventLogger... Flushing remaining logs.")

	// Закрываем входной канал событий. Новые записи через LogEvent больше не пройдут.
	close(el.eventsChan)

	// Блокируемся и ждем, пока воркер дочитает остатки из закрытого канала и закроет файл
	<-el.doneChan

	log.Println("EventLogger stopped gracefully. All logs saved.")
}
