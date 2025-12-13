package monitoring

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	// "user-management/internal/domain/entities"
)

// OrderEvent representa un evento en el sistema
type OrderEvent struct {
	OrderID   int
	UserID    int
	EventType string // "created", "updated", "completed", "failed"
	Timestamp time.Time
	Data      any
}

// OrderMonitor demuestra select, context y manejo avanzado
type OrderMonitor struct {
	// Channels para diferentes tipos de eventos
	eventChan   chan OrderEvent
	alertChan   chan string
	metricsChan chan map[string]int

	// Channels de control
	shutdownChan chan struct{}
	doneChan     chan struct{}

	// Context para cancelación
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewOrderMonitor(ctx context.Context) *OrderMonitor {
	childCtx, cancel := context.WithCancel(ctx)

	return &OrderMonitor{
		eventChan:    make(chan OrderEvent, 100),
		alertChan:    make(chan string, 10),
		metricsChan:  make(chan map[string]int, 5),
		shutdownChan: make(chan struct{}),
		doneChan:     make(chan struct{}),
		ctx:          childCtx,
		cancel:       cancel,
	}
}

// Start inicia el procesamiento con select multiplexado
func (m *OrderMonitor) Start() {
	m.wg.Add(1)

	go func() {
		defer m.wg.Done()
		defer close(m.doneChan)

		log.Println("Monitor de pedidos iniciado")

		// Ticker para métricas periódicas
		metricsTicker := time.NewTicker(30 * time.Second)
		defer metricsTicker.Stop()

		for {
			select {
			// 1. Eventos de pedidos
			case event := <-m.eventChan:
				m.handleEvent(event)

			// 2. Alertas del sistema
			case alert := <-m.alertChan:
				m.handleAlert(alert)

			// 3. Métricas periódicas
			case <-metricsTicker.C:
				m.collectMetrics()

			// 4. Señal de shutdown
			case <-m.shutdownChan:
				log.Println("Recibida señal de shutdown")
				return

			// 5. Context cancelado
			case <-m.ctx.Done():
				log.Printf("Context cancelado: %v", m.ctx.Err())
				return

			// 6. Timeout para evitar bloqueos (pattern común)
			case <-time.After(5 * time.Minute):
				log.Println("Timeout de inactividad")
				m.alertChan <- "Monitor inactivo por 5 minutos"
			}
		}
	}()
}

// PublishEvent demuestra select con timeout
func (m *OrderMonitor) PublishEvent(event OrderEvent) error {
	select {
	case m.eventChan <- event:
		return nil
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("timeout publicando evento")
	case <-m.ctx.Done():
		return fmt.Errorf("context cancelado: %w", m.ctx.Err())
	}
}

// handleEvent con diferentes estrategias de timeout
func (m *OrderMonitor) handleEvent(event OrderEvent) {
	// Crear context con timeout para este procesamiento
	ctx, cancel := context.WithTimeout(m.ctx, 2*time.Second)
	defer cancel()

	// Channel para resultado
	resultChan := make(chan bool, 1)

	// Procesar en goroutine separada
	go func() {
		// Simular trabajo
		time.Sleep(time.Millisecond * 50)

		// Log basado en tipo de evento
		switch event.EventType {
		case "created":
			log.Printf("Pedido %d creado por usuario %d",
				event.OrderID, event.UserID)
		case "completed":
			log.Printf("Pedido %d completado", event.OrderID)
		case "failed":
			log.Printf("Pedido %d falló: %v",
				event.OrderID, event.Data)
			m.alertChan <- fmt.Sprintf("Pedido %d falló", event.OrderID)
		}

		resultChan <- true
	}()

	// Esperar resultado con timeout del context
	select {
	case <-resultChan:
		// Éxito
	case <-ctx.Done():
		log.Printf("Timeout procesando evento %d: %v",
			event.OrderID, ctx.Err())
	}
}

func (m *OrderMonitor) handleAlert(alert string) {
	log.Printf("ALERTA: %s", alert)
}

func (m *OrderMonitor) collectMetrics() {
	metrics := map[string]int{
		"events_processed": len(m.eventChan),
		"alerts_pending":   len(m.alertChan),
		"goroutines":       50, // Simulado
	}

	select {
	case m.metricsChan <- metrics:
		log.Println("Métricas enviadas")
	default:
		// Channel lleno, no bloquear
	}
}

// GetMetrics con timeout usando context
func (m *OrderMonitor) GetMetrics(timeout time.Duration) (map[string]int, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	select {
	case metrics := <-m.metricsChan:
		return metrics, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout obteniendo métricas: %w", ctx.Err())
	}
}

// Stop con shutdown ordenado
func (m *OrderMonitor) Stop() {
	log.Println("Deteniendo monitor...")

	// 1. Cancelar context
	m.cancel()

	// 2. Enviar señal de shutdown
	select {
	case m.shutdownChan <- struct{}{}:
	case <-time.After(100 * time.Millisecond):
		// Ya está cerrando
	}

	// 3. Esperar que termine
	select {
	case <-m.doneChan:
		log.Println("Monitor detenido correctamente")
	case <-time.After(5 * time.Second):
		log.Println("Timeout esperando monitor")
	}

	// 4. Esperar todas las goroutines
	m.wg.Wait()
}
